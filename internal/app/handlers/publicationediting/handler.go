package publicationediting

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/backends/filestore"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
)

type Handler struct {
	handlers.BaseHandler
	Repository                backends.Repository
	ProjectSearchService      backends.ProjectSearchService
	ProjectService            backends.ProjectService
	PersonService             backends.PersonService
	PersonSearchService       backends.PersonSearchService
	OrganizationSearchService backends.OrganizationSearchService
	OrganizationService       backends.OrganizationService
	DatasetSearchService      backends.DatasetSearchService
	FileStore                 *filestore.Store
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

		id := bind.PathValues(r).Get("id")
		pub, err := h.Repository.GetPublication(id)
		if err != nil {
			if err == backends.ErrNotFound {
				h.Logger.Warn("edit publication: could not find publication with id:", "errors", err, "id", id, "user", ctx.User.ID)
				handlers.NotFound(w, r, ctx, err)
			} else {
				h.Logger.Error(
					"edit publication: unexpected error when retrieving publication with id:",
					"errors", err,
					"id", id,
					"user", ctx.User.ID,
				)
				handlers.InternalServerError(w, r, err)
			}
			return
		}

		if !ctx.User.CanEditPublication(pub) {
			h.Logger.Warn("edit publication: user isn't allowed to edit the publication:", "errors", err, "publication", id, "user", ctx.User.ID)
			handlers.Forbidden(w, r)
			return
		}

		fn(w, r, Context{
			BaseContext: ctx,
			Publication: pub,
		})
	})
}
