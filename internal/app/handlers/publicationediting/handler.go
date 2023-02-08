package publicationediting

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/internal/app/handlers"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/backends/filestore"
	"github.com/ugent-library/biblio-backoffice/internal/bind"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/render"
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
	MaxFileSize               int
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

		id := bind.PathValues(r).Get("id")
		pub, err := h.Repository.GetPublication(id)
		if err != nil {
			if err == backends.ErrNotFound {
				h.Logger.Warn("edit publication: could not find publication with id:", "errors", err, "id", id, "user", ctx.User.ID)
				render.NotFound(w, r, err)
			} else {
				h.Logger.Error(
					"edit publication: unexpected error when retrieving publication with id:",
					"errors", err,
					"id", id,
					"user", ctx.User.ID,
				)
				render.InternalServerError(w, r, err)
			}
			return
		}

		if !ctx.User.CanEditPublication(pub) {
			h.Logger.Warn("edit publication: user isn't allowed to edit the publication:", "errors", err, "publication", id, "user", ctx.User.ID)
			render.Forbidden(w, r)
			return
		}

		fn(w, r, Context{
			BaseContext: ctx,
			Publication: pub,
		})
	})
}
