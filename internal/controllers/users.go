package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/unrolled/render"
)

type Users struct {
	engine       *engine.Engine
	render       *render.Render
	sessionName  string
	sessionStore sessions.Store
	router       *mux.Router
}

func NewUsers(e *engine.Engine, r *render.Render, sessionName string, sessionStore sessions.Store, router *mux.Router) *Users {
	return &Users{
		engine:       e,
		render:       r,
		router:       router,
		sessionName:  sessionName,
		sessionStore: sessionStore,
	}
}

func (c *Users) ImpersonateChoose(w http.ResponseWriter, r *http.Request) {
	c.render.HTML(w, 200, "user/_impersonate_choose",
		views.NewData(c.render, r, nil),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *Users) Impersonate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")

	user, err := c.engine.GetUserByUsername(username)
	if err != nil {
		log.Printf("impersonate get user error: %s", err)
		// TODO
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !context.GetUser(r.Context()).CanImpersonateUser() {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	session, _ := c.sessionStore.Get(r, c.sessionName)
	session.Values["original_user_id"] = context.GetUser(r.Context()).ID
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	redirectURL, _ := c.router.Get("home").URLPath()
	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}

func (c *Users) ImpersonateRemove(w http.ResponseWriter, r *http.Request) {
	session, _ := c.sessionStore.Get(r, c.sessionName)
	if origUserID := session.Values["original_user_id"]; origUserID != nil {
		delete(session.Values, "original_user_id")
		session.Values["user_id"] = origUserID
		session.Save(r, w)
	}

	redirectURL, _ := c.router.Get("home").URLPath()
	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}
