package commands

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(updateEmbargoes)
}

var updateEmbargoes = &cobra.Command{
	Use:   "update-embargoes",
	Short: "Update embargoes",
	Run: func(cmd *cobra.Command, args []string) {
		e := Services()
		logger := newLogger()

		updatePublicationEmbargoes(e, logger)
		updateDatasetEmbargoes(e, logger)
	},
}

func updatePublicationEmbargoes(e *backends.Services, logger *zap.SugaredLogger) {
	e.Repository.AddPublicationListener(func(p *models.Publication) {
		if err := e.PublicationSearchService.Index(p); err != nil {
			logger.Fatalf("error indexing publication %s: %v", p.ID, err)
		}
	})

	var count int = 0
	updateEmbargoErr := e.Repository.Transaction(
		context.Background(),
		func(repo backends.Repository) error {

			/*
				select live publications that have files with embargoed access
			*/
			var embargoAccessLevel string = "info:eu-repo/semantics/embargoedAccess"
			currentDateStr := time.Now().Format("2006-01-02")
			var sqlPublicationWithEmbargo string = `
			SELECT * FROM publications WHERE date_until IS NULL AND
			data->'file' IS NOT NULL AND
			EXISTS(
				SELECT 1 FROM jsonb_array_elements(data->'file') AS f
				WHERE f->>'access_level' = $1 AND
				f->>'embargo_date' <= $2
			)
			`

			publications := make([]*models.Publication, 0)
			sErr := repo.SelectPublications(
				sqlPublicationWithEmbargo,
				[]any{
					embargoAccessLevel,
					currentDateStr},
				func(publication *models.Publication) bool {
					publications = append(publications, publication)
					return true
				},
			)

			if sErr != nil {
				return sErr
			}

			for _, publication := range publications {
				/*
					clear outdated embargoes
				*/
				for _, file := range publication.File {
					if file.AccessLevel != embargoAccessLevel {
						continue
					}
					// TODO: what with empty embargo_date?
					if file.EmbargoDate == "" {
						continue
					}
					if file.EmbargoDate > currentDateStr {
						continue
					}
					file.ClearEmbargo()
				}

				publication.User = nil

				if e := repo.SavePublication(publication, nil); e != nil {
					return e
				}
				count++
			}

			return nil
		},
	)

	if updateEmbargoErr != nil {
		logger.Fatal(updateEmbargoErr)
	}

	logger.Infof("updated %d publication embargoes", count)
}

func updateDatasetEmbargoes(e *backends.Services, logger *zap.SugaredLogger) {
	e.Repository.AddDatasetListener(func(d *models.Dataset) {
		if err := e.DatasetSearchService.Index(d); err != nil {
			logger.Fatalf("error indexing dataset %s: %v", d.ID, err)
		}
	})

	var count int = 0
	updateEmbargoErr := e.Repository.Transaction(
		context.Background(),
		func(repo backends.Repository) error {
			/*
				select live datasets with embargoed access
			*/
			var embargoAccessLevel string = "info:eu-repo/semantics/embargoedAccess"
			currentDateStr := time.Now().Format("2006-01-02")
			var sqlDatasetsWithEmbargo string = `
		SELECT * FROM datasets
		WHERE date_until is null AND 
		data->>'access_level' = $1 AND
		data->>'embargo_date' <> '' AND
		data->>'embargo_date' <= $2 
		`

			datasets := make([]*models.Dataset, 0)
			sErr := repo.SelectDatasets(
				sqlDatasetsWithEmbargo,
				[]any{
					embargoAccessLevel,
					currentDateStr},
				func(dataset *models.Dataset) bool {
					datasets = append(datasets, dataset)
					return true
				},
			)

			if sErr != nil {
				return sErr
			}

			for _, dataset := range datasets {
				/*
					clear outdated embargoes
				*/
				// TODO: what with empty embargo_date?
				if dataset.EmbargoDate == "" {
					continue
				}
				if dataset.EmbargoDate > currentDateStr {
					continue
				}
				dataset.ClearEmbargo()

				if e := repo.SaveDataset(dataset, nil); e != nil {
					return e
				}
				count++
			}

			return nil
		},
	)

	if updateEmbargoErr != nil {
		logger.Fatal(updateEmbargoErr)
	}

	logger.Infof("updated %d dataset embargoes", count)
}
