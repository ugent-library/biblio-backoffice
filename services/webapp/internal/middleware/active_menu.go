package middleware

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
)

func SetActiveMenu(menu string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := context.WithActiveMenu(r.Context(), menu)
			next.ServeHTTP(w, r.WithContext(c))
		})
	}
}
