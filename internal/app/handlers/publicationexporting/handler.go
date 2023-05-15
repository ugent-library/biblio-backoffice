package publicationexporting

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/internal/app/handlers"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/bind"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/render"
)

type Handler struct {
	handlers.BaseHandler
	PublicationSearchService backends.PublicationSearchService
	PublicationListExporters map[string]backends.PublicationListExporterFactory
}

type Context struct {
	handlers.BaseContext
	SearchArgs *models.SearchArgs
	ExportArgs *ExportArgs
}

type ExportArgs struct {
	Format string `path:"format,omitempty"`
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

		exportArgs := &ExportArgs{}
		if err := bind.Request(r, exportArgs); err != nil {
			h.Logger.Warnw("publication search: could not bind export arguments", "errors", err, "request", r, "user", ctx.User.ID)
			render.BadRequest(w, r, err)
			return
		}
		if exportArgs.Format == "" {
			exportArgs.Format = "xlsx"
		}

		fn(w, r, Context{
			BaseContext: ctx,
			SearchArgs:  searchArgs,
			ExportArgs:  exportArgs,
		})
	})
}
