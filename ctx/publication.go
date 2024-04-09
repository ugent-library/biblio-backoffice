package ctx

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/repositories"
)

const PublicationKey = contextKey("publication")

func GetPublication(r *http.Request) *models.Publication {
	return r.Context().Value(PublicationKey).(*models.Publication)
}

func SetPublication(repo *repositories.Repo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r)

			publicationId := chi.URLParam(r, "id")

			publication, err := repo.GetPublication(publicationId)
			if err != nil {
				c.HandleError(w, r, err)
				return
			}

			ctx := context.WithValue(r.Context(), PublicationKey, publication)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
