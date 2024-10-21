package ctx

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/repositories"
	"github.com/ugent-library/httperror"
)

const DatasetKey = contextKey("dataset")

func GetDataset(r *http.Request) *models.Dataset {
	return r.Context().Value(DatasetKey).(*models.Dataset)
}

func SetDataset(repo *repositories.Repo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r)

			datasetId := chi.URLParam(r, "id")

			dataset, err := repo.GetDataset(datasetId)
			if err != nil {
				c.HandleError(w, r, err)
				return
			}

			ctx := context.WithValue(r.Context(), DatasetKey, dataset)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireViewDataset(repo *repositories.Repo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r)

			if !c.Repo.CanViewDataset(c.User, GetDataset(r)) {
				c.HandleError(w, r, httperror.Forbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireEditDataset(repo *repositories.Repo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r)

			if !c.Repo.CanEditDataset(c.User, GetDataset(r)) {
				c.HandleError(w, r, httperror.Forbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
