package dashboard

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/bind"
)

type Handler struct {
	handlers.BaseHandler
	DatasetSearchService     backends.DatasetSearchService
	PublicationSearchService backends.PublicationSearchService
}

type Context struct {
	handlers.BaseContext
	Type string
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil || !ctx.User.CanViewDashboard() {
			handlers.Unauthorized(w, r)
			return
		}

		context := Context{
			BaseContext: ctx,
			Type:        bind.PathValues(r).Get("type"),
		}

		fn(w, r, context)
	})
}
