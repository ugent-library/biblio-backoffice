package handlers

import (
	"net/http"
	"strings"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/handlers/authenticating"
	"github.com/ugent-library/biblio-backoffice/views"
	"github.com/ugent-library/htmx"
)

func UserNotFound(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	// clear session to make sure that baseHandler doesn't keep blocking other handlers
	session, err := c.SessionStore.Get(r, c.SessionName)
	if err != nil {
		c.Log.Error("unable to retrieve session", "error", err)
		InternalServerError(w, r)
		return
	}
	if err := authenticating.ClearSession(w, r, session); err != nil {
		c.Log.Error("unable to save session", "error", err)
		InternalServerError(w, r)
		return
	}

	w.WriteHeader(404)
	views.UserNotFound(c).Render(r.Context(), w)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	w.WriteHeader(404)
	views.NotFound(c).Render(r.Context(), w)
}

func Unauthorized(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && strings.Contains(r.Header.Get("Accept"), "text/html") && !htmx.Request(r) {
		c := ctx.Get(r)
		http.Redirect(w, r, c.PathTo("login", "destination", r.URL.String()).String(), http.StatusSeeOther)
		return
	}

	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	w.WriteHeader(500)
	views.InternalServerError(c).Render(r.Context(), w)
}
