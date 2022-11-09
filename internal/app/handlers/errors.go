package handlers

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/render"
)

type YieldNotFound struct {
	BaseContext
	PageTitle string
	ActiveNav string
}

func (h *BaseHandler) NotFound(w http.ResponseWriter, r *http.Request, ctx BaseContext) {
	w.WriteHeader(404)
	render.Layout(w, "layouts/default", "pages/notfound", YieldNotFound{
		BaseContext: ctx,
		PageTitle:   "Biblio",
	})
}
