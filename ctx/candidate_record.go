package ctx

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/repositories"
	"github.com/ugent-library/httperror"
)

const CandidateRecordKey = contextKey("publication")

func GetCandidateRecord(r *http.Request) *models.CandidateRecord {
	return r.Context().Value(CandidateRecordKey).(*models.CandidateRecord)
}

func SetCandidateRecord(repo *repositories.Repo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r)

			id := chi.URLParam(r, "id")

			rec, err := repo.GetCandidateRecord(r.Context(), id)
			if err != nil {
				c.HandleError(w, r, err)
				return
			}

			ctx := context.WithValue(r.Context(), CandidateRecordKey, rec)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireViewCandidateRecord(repo *repositories.Repo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r)

			if !repo.CanViewPublication(c.User, GetCandidateRecord(r).Publication) {
				c.HandleError(w, r, httperror.Forbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireEditCandidateRecord(repo *repositories.Repo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r)

			if !repo.CanEditPublication(c.User, GetCandidateRecord(r).Publication) {
				c.HandleError(w, r, httperror.Forbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
