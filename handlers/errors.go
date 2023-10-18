package handlers

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/views"
)

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

func (h *BaseHandler) InternalServerError(w http.ResponseWriter, r *http.Request, ctx BaseContext) {
	w.WriteHeader(500)
	render.Layout(w, "layouts/default", "pages/internalerror", YieldNotFound{
		BaseContext:      ctx,
		PageTitle:        "Biblio",
		ErrorTitle:       "Something went wrong.",
		ErrorDescription: "Your (re)search was too groundbreaking.",
	})
}

func (h *BaseHandler) ErrorModal(w http.ResponseWriter, r *http.Request, errID string, ctx BaseContext) {
	render.Layout(w, "show_modal", "modals/error", YieldModalError{
		BaseContext: ctx,
		ID:          errID,
	})
}
