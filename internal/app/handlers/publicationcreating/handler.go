package publicationcreating

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
	Repository          backends.Repository
	SearchService       backends.SearchService
	PublicationSources  map[string]backends.PublicationGetter
	PublicationDecoders map[string]backends.PublicationDecoderFactory
	OrganizationService backends.OrganizationService
}

type Context struct {
	handlers.BaseContext
	Publication *models.Publication
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil {
			render.Unauthorized(w, r)
			return
		}

		context := Context{
			BaseContext: ctx,
		}

		if id := bind.PathValues(r).Get("id"); id != "" {
			d, err := h.Repository.GetPublication(id)
			if err != nil {
				render.NotFound(w, r, err)
				return
			}

			if !ctx.User.CanEditPublication(d) {
				h.Logger.Warn("create publication: user isn't allowed to edit the publication:", "errors", err, "publication", id, "user", ctx.User.ID)
				render.Forbidden(w, r)
				return
			}

			context.Publication = d
		}

		fn(w, r, context)
	})
}
