package publicationediting

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/render"
)

type BindType struct {
	Type string `form:"type"`
}

type YieldUpdateType struct {
	Context
	Type string
}

func (h *Handler) ConfirmUpdateType(w http.ResponseWriter, r *http.Request, ctx Context) {
	// TODO validate type
	render.Layout(w, "show_modal", "publication/confirm_update_type", YieldUpdateType{
		Context: ctx,
		Type:    r.URL.Query().Get("type"),
	})
}

func (h *Handler) UpdateType(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindType{}
	if err := bind.RequestForm(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	ctx.Publication.ChangeType(b.Type)

	redirectURL := h.PathFor("publication", "id", ctx.Publication.ID)
	w.Header().Set("HX-Redirect", redirectURL.String())
}
