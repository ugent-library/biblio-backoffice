package publicationcreating

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
)

type Handler struct {
	handlers.BaseHandler
	Repository               backends.Repository
	PublicationSearchService backends.PublicationSearchService
	PublicationSources       map[string]backends.PublicationGetter
	PublicationDecoders      map[string]backends.PublicationDecoderFactory
}

type Context struct {
	handlers.BaseContext
	Publication *models.Publication
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil {
			handlers.Unauthorized(w, r)
			return
		}

		context := Context{
			BaseContext: ctx,
		}

		if id := bind.PathValues(r).Get("id"); id != "" {
			d, err := h.Repository.GetPublication(id)
			if err != nil {
				handlers.NotFound(w, r, ctx, err)
				return
			}

			if !ctx.User.CanEditPublication(d) {
				h.Logger.Warn("create publication: user isn't allowed to edit the publication:", "errors", err, "publication", id, "user", ctx.User.ID)
				handlers.Forbidden(w, r)
				return
			}

			context.Publication = d
		}

		fn(w, r, context)
	})
}
