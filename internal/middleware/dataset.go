package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/engine"
)

func SetDataset(e *engine.Engine) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dataset, err := e.GetDataset(mux.Vars(r)["id"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			user := context.GetUser(r.Context())

			if user.CanEditDataset(dataset) {
				c := context.WithDataset(r.Context(), dataset)
				next.ServeHTTP(w, r.WithContext(c))
				return
			}

			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		})
	}
}
