package controllers

import (
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
	// claims := &oidc.Claims{}
	// err := c.oidcClient.Exchange(r.URL.Query().Get("code"), claims)
	// if err != nil {
	// 	log.Printf("oidc error: %s", err)
	// 	// TODO
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// log.Printf("user claims: %+v", claims)

	// session, _ := c.sessionStore.Get(r, c.sessionName)
	// if err != nil {
	// 	log.Printf("session error: %s", err)
	// 	// TODO
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// user, err := c.engine.GetUserByEmail(claims.Email)
	// if err != nil {
	// 	log.Printf("get user error: %s", err)
	// 	// TODO
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// log.Printf("user: %+v", user)

	// session.Values["user_id"] = user["_id"].(string)
	// session.Save(r, w)

	// http.Redirect(w, r, "/publication", http.StatusFound)
}

// func (c *Auth) Login(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, c.oidcClient.AuthorizationURL(), http.StatusFound)
// }

func (c *Auth) Login(w http.ResponseWriter, r *http.Request) {
	session, err := c.sessionStore.Get(r, c.sessionName)
	if err != nil {
		log.Printf("session error: %s", err)
		// TODO
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := c.engine.GetUserByUsername(r.URL.Query().Get("username"))
	if err != nil {
		log.Printf("get user error: %s", err)
		// TODO
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("user: %+v", user)

	session.Values["user_id"] = user.ID()
	session.Save(r, w)

	http.Redirect(w, r, "/publication", http.StatusFound)
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
