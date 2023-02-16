package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(createHandles)
}

var createHandles = &cobra.Command{
	Use:   "create-handles",
	Short: "Create handles",
	Run: func(cmd *cobra.Command, args []string) {
		e := Services()
		logger := newLogger()
		handleService := e.HandleService

		if handleService == nil {
			logger.Fatal("handle server updates are not enabled")
		}

		createPublicationHandles(e, logger, handleService)
		createDatasetHandles(e, logger, handleService)
	},
}

func createPublicationHandles(e *backends.Services, logger *zap.SugaredLogger, handleService backends.HandleService) {
	e.Repository.AddPublicationListener(func(p *models.Publication) {
		if p.DateUntil == nil {
			if err := e.PublicationSearchService.Index(p); err != nil {
				logger.Fatalf("error indexing publication %s: %v", p.ID, err)
			}
		}
	})

	var count int = 0
	createHandlesErr := e.Repository.Transaction(
		context.Background(),
		func(repo backends.Repository) error {

			publications := make([]*models.Publication, 0)
			sql := `
			SELECT * FROM publications WHERE date_until IS NULL AND
			data->>'status' = 'public' AND
			NOT data ? 'handle'
			`

			selectErr := repo.SelectPublications(
				sql,
				[]any{},
				func(p *models.Publication) bool {
					publications = append(publications, p)
					return true
				},
			)

			if selectErr != nil {
				return selectErr
			}

			for _, p := range publications {
				h, hErr := handleService.UpsertHandle(p.ID)
				if hErr != nil {
					return fmt.Errorf(
						"error adding handle for publication %s: %s",
						p.ID,
						hErr,
					)
				} else if !h.IsSuccess() {
					return fmt.Errorf(
						"error adding handle for publication %s: %s",
						p.ID,
						h.Message,
					)
				} else {
					logger.Infof(
						"added handle url %s to publication %s",
						h.GetFullHandleURL(),
						p.ID,
					)
					p.Handle = h.GetFullHandleURL()
					p.User = nil

					if e := repo.SavePublication(p, nil); e != nil {
						return e
					}
					count++
				}
			}

			return nil
		},
	)

	if createHandlesErr != nil {
		logger.Fatal(createHandlesErr)
	}

	logger.Infof("created %d publication handles", count)
}

func createDatasetHandles(e *backends.Services, logger *zap.SugaredLogger, handleService backends.HandleService) {
	e.Repository.AddDatasetListener(func(d *models.Dataset) {
		if d.DateUntil == nil {
			if err := e.DatasetSearchService.Index(d); err != nil {
				logger.Fatalf("error indexing dataset %s: %v", d.ID, err)
			}
		}
	})

	var count int = 0
	createHandlesErr := e.Repository.Transaction(
		context.Background(),
		func(repo backends.Repository) error {

			datasets := make([]*models.Dataset, 0)
			sql := `
			SELECT * FROM datasets WHERE date_until IS NULL AND
			data->>'status' = 'public' AND
			NOT data ? 'handle'
			`

			selectErr := repo.SelectDatasets(
				sql,
				[]any{},
				func(p *models.Dataset) bool {
					datasets = append(datasets, p)
					return true
				},
			)

			if selectErr != nil {
				return selectErr
			}

			for _, d := range datasets {
				h, hErr := handleService.UpsertHandle(d.ID)
				if hErr != nil {
					return fmt.Errorf(
						"error adding handle for dataset %s: %s",
						d.ID,
						hErr,
					)
				} else if !h.IsSuccess() {
					return fmt.Errorf(
						"error adding handle for dataset %s: %s",
						d.ID,
						h.Message,
					)
				} else {
					logger.Infof(
						"added handle url %s to dataset %s",
						h.GetFullHandleURL(),
						d.ID,
					)
					d.Handle = h.GetFullHandleURL()
					d.User = nil

					if e := repo.SaveDataset(d, nil); e != nil {
						return e
					}
					count++
				}
			}

			return nil
		},
	)

	if createHandlesErr != nil {
		logger.Fatal(createHandlesErr)
	}

	logger.Infof("created %d dataset handles", count)
}
