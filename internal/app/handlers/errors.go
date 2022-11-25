package handlers

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/render"
)

type YieldNotFound struct {
	BaseContext
	PageTitle        string
	ActiveNav        string
	ErrorTitle       string
	ErrorDescription string
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
