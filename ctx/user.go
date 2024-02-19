package ctx

import (
	"net/http"

	"github.com/ugent-library/httperror"
)

func RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := Get(r)

		if c.User == nil {
			c.HandleError(w, r, httperror.Unauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RequireUserRole(userRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r)

			if c.UserRole != userRole {
				c.HandleError(w, r, httperror.Unauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
