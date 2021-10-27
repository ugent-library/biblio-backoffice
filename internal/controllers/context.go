package controllers

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/go-locale/locale"
	"github.com/unrolled/render"
)

type Context struct {
	Engine       *engine.Engine
	BaseURL      *url.URL
	Router       *mux.Router
	Render       *render.Render
	Localizer    *locale.Localizer
	SessionName  string
	SessionStore sessions.Store
}

func (c *Context) Session(r *http.Request) (*sessions.Session, error) {
	return c.SessionStore.Get(r, c.SessionName)
}
