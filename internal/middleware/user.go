package middleware

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/engine"
)

func SetUser(e *engine.Engine, sessionName string, sessionStore sessions.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := sessionStore.Get(r, sessionName)
			if err != nil {
				// TODO
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			userID := session.Values["user_id"]
			if userID == nil {
				next.ServeHTTP(w, r)
				return
			}

			user, err := e.GetUser(userID.(string))
			if err != nil {
				log.Printf("get user error: %s", err)
				// TODO
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			c := context.WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(c))
		})
	}
}

func RequireUser(redirectURL string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if context.GetUser(r.Context()) == nil {
				http.Redirect(w, r, redirectURL, http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
