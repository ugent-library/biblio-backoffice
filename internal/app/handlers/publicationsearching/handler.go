package publicationsearching

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/internal/app/handlers"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/bind"
	"github.com/ugent-library/biblio-backoffice/internal/render"
	"github.com/ugent-library/biblio-backoffice/models"
)

type Handler struct {
	handlers.BaseHandler
	PublicationSearchIndex backends.PublicationIndex
	FileStore              backends.FileStore
}

type Context struct {
	handlers.BaseContext
	SearchArgs *models.SearchArgs
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil {
			render.Unauthorized(w, r)
			return
		}

		searchArgs := models.NewSearchArgs()
		if err := bind.Request(r, searchArgs); err != nil {
			h.Logger.Warnw("publication search: could not bind search arguments", "errors", err, "request", r, "user", ctx.User.ID)
			render.BadRequest(w, r, err)
			return
		}
		searchArgs.Cleanup()

		fn(w, r, Context{
			BaseContext: ctx,
			SearchArgs:  searchArgs,
		})
	})
}
