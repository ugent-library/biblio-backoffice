package publicationediting

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/bind"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/repositories"
)

type Handler struct {
	handlers.BaseHandler
	Repo                      *repositories.Repo
	ProjectSearchService      backends.ProjectSearchService
	ProjectService            backends.ProjectService
	PersonService             backends.PersonService
	PersonSearchService       backends.PersonSearchService
	OrganizationSearchService backends.OrganizationSearchService
	OrganizationService       backends.OrganizationService
	DatasetSearchIndex        backends.DatasetIndex
	FileStore                 backends.FileStore
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
		pub, err := h.Repo.GetPublication(id)
		if err != nil {
			if err == models.ErrNotFound {
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
