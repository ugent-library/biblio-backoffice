package frontoffice

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
)

type updateResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Count   int    `json:"count"`
}

func (h *Handler) PublicationUpdateEmbargo(w http.ResponseWriter, r *http.Request) {

	var count int = 0
	updateEmbargoErr := h.Repository.Transaction(
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
				if e := repo.SavePublication(publication, nil); e != nil {
					return e
				}
				count++
			}

			return nil
		},
	)

	updateResponse := &updateResponse{}

	if updateEmbargoErr == nil {
		updateResponse.Count = count
		updateResponse.Status = "ok"
	} else {
		updateResponse.Status = "error"
		updateResponse.Message = updateEmbargoErr.Error()
	}

	j, jErr := json.Marshal(updateResponse)

	if updateEmbargoErr != nil {
		render.InternalServerError(w, r, jErr)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(j)
}

func (h *Handler) DatasetUpdateEmbargo(w http.ResponseWriter, r *http.Request) {

	var count int = 0
	updateEmbargoErr := h.Repository.Transaction(
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

	updateResponse := &updateResponse{}

	if updateEmbargoErr == nil {
		updateResponse.Count = count
		updateResponse.Status = "ok"
	} else {
		updateResponse.Status = "error"
		updateResponse.Message = updateEmbargoErr.Error()
	}

	j, jErr := json.Marshal(updateResponse)

	if updateEmbargoErr != nil {
		render.InternalServerError(w, r, jErr)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(j)
}
