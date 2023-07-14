package dashboard

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/internal/app/handlers"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/bind"
	"github.com/ugent-library/biblio-backoffice/internal/render"
)

type Handler struct {
	handlers.BaseHandler
	SearchService          backends.SearchService
	DatasetSearchIndex     backends.DatasetIndex
	PublicationSearchIndex backends.PublicationIndex
}

type Context struct {
	handlers.BaseContext
	Type string
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil || !ctx.User.CanViewDashboard() {
			render.Unauthorized(w, r)
			return
		}

		context := Context{
			BaseContext: ctx,
			Type:        bind.PathValues(r).Get("type"),
		}

		fn(w, r, context)
	})
}
