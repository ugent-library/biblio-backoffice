package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/render"
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

type YieldNotFound struct {
	BaseContext
	PageTitle        string
	ActiveNav        string
	ErrorTitle       string
	ErrorDescription string
}

type YieldModalError struct {
	BaseContext
	ID string
}

func (h *BaseHandler) NotFound(w http.ResponseWriter, r *http.Request, ctx BaseContext) {
	w.WriteHeader(404)
	render.Layout(w, "layouts/default", "pages/notfound", YieldNotFound{
		BaseContext:      ctx,
		PageTitle:        "Biblio",
		ErrorTitle:       "This page does not exist.",
		ErrorDescription: "Your (re)search was too groundbreaking.",
	})
}

func (h *BaseHandler) UserNotFound(w http.ResponseWriter, r *http.Request, ctx BaseContext) {
	// clear session to make sure that baseHandler doesn't keep blocking other handlers
	session, err := h.SessionStore.Get(r, h.SessionName)
	if err != nil {
		h.Logger.Errorf("unable to retrieve session: %w", err)
		InternalServerError(w, r)
		return
	}
	if err := clearSession(w, r, session); err != nil {
		h.Logger.Errorf("unable to save session: %w", err)
		InternalServerError(w, r)
		return
	}

	w.WriteHeader(404)
	render.Layout(w, "layouts/default", "pages/usernotfound", YieldNotFound{
		BaseContext:      ctx,
		PageTitle:        "Biblio",
		ErrorTitle:       "Account not found",
		ErrorDescription: "Account not found",
	})
}

func clearSession(w http.ResponseWriter, r *http.Request, session *sessions.Session) error {
	delete(session.Values, UserIDKey)
	delete(session.Values, OriginalUserIDKey)
	delete(session.Values, OriginalUserRoleKey)
	delete(session.Values, UserRoleKey)
	return session.Save(r, w)
}
