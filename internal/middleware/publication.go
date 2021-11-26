package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/engine"
)

func SetPublication(e *engine.Engine) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pub, err := e.GetPublication(mux.Vars(r)["id"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			c := context.WithPublication(r.Context(), pub)
			next.ServeHTTP(w, r.WithContext(c))
		})
	}
}

func RequireCanViewPublication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		pub := context.GetPublication(c)
		user := context.GetUser(c)

		if user.CanViewPublication(pub) {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
}

func RequireCanEditPublication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		pub := context.GetPublication(c)
		user := context.GetUser(c)

		if user.CanEditPublication(pub) {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
}

func RequireCanPublishPublication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		pub := context.GetPublication(c)
		user := context.GetUser(c)

		if user.CanPublishPublication(pub) {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
}

func RequireCanDeletePublication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		pub := context.GetPublication(c)
		user := context.GetUser(c)

		if user.CanDeletePublication(pub) {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
}
