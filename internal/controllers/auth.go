package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/go-oidc/oidc"
)

type Auth struct {
	engine       *engine.Engine
	sessionName  string
	sessionStore sessions.Store
	oidcClient   *oidc.Client
}

func NewAuth(e *engine.Engine, sessionName string, sessionStore sessions.Store, oidcClient *oidc.Client) *Auth {
	return &Auth{
		engine:       e,
		sessionName:  sessionName,
		sessionStore: sessionStore,
		oidcClient:   oidcClient,
	}
}

func (c *Auth) Callback(w http.ResponseWriter, r *http.Request) {
	claims := &oidc.Claims{}
	err := c.oidcClient.Exchange(r.URL.Query().Get("code"), claims)
	if err != nil {
		log.Printf("oidc error: %s", err)
		// TODO
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("user claims: %+v", claims)

	session, _ := c.sessionStore.Get(r, c.sessionName)
	if err != nil {
		log.Printf("session error: %s", err)
		// TODO
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := c.engine.GetUserByEmail(claims.Email)
	if err != nil {
		log.Printf("get user error: %s", err)
		// TODO
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("user: %+v", user)

	session.Values["user_id"] = user["_id"].(string)
	session.Save(r, w)

	http.Redirect(w, r, "/publication", http.StatusFound)
}

func (c *Auth) Login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, c.oidcClient.AuthorizationURL(), http.StatusFound)
}

func (c *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := c.sessionStore.Get(r, c.sessionName)
	if err != nil {
		// TODO
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	delete(session.Values, "user_id")
	session.Save(r, w)
}

var UserKey = &key{"User"}

type key struct {
	name string
}

func (c *key) String() string {
	return c.name
}

func HasUser(r *http.Request) bool {
	_, ok := r.Context().Value(UserKey).(string)
	return ok
}

func GetUser(r *http.Request) string {
	return r.Context().Value(UserKey).(string)
}

func SetUser(sessionStore sessions.Store, sessionName string) func(next http.Handler) http.Handler {
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

			c := context.WithValue(r.Context(), UserKey, userID)
			next.ServeHTTP(w, r.WithContext(c))
		})
	}
}

func RequireUser() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !HasUser(r) {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
