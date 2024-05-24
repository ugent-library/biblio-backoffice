package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/views"
)

func UserNotFound(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	// clear session to make sure that baseHandler doesn't keep blocking other handlers
	session, err := c.SessionStore.Get(r, c.SessionName)
	if err != nil {
		c.Log.Errorf("unable to retrieve session: %w", err)
		InternalServerError(w, r)
		return
	}
	if err := clearSession(w, r, session); err != nil {
		c.Log.Errorf("unable to save session: %w", err)
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

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	w.WriteHeader(500)
	views.InternalServerError(c).Render(r.Context(), w)
}

func clearSession(w http.ResponseWriter, r *http.Request, session *sessions.Session) error {
	delete(session.Values, ctx.UserIDKey)
	delete(session.Values, ctx.OriginalUserIDKey)
	delete(session.Values, ctx.OriginalUserRoleKey)
	delete(session.Values, ctx.UserRoleKey)
	return session.Save(r, w)
}
