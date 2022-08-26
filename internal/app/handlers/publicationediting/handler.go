package publicationediting

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/backends/filestore"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
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
			render.Unauthorized(w, r)
			return
		}

		id := bind.PathValues(r).Get("id")
		pub, err := h.Repository.GetPublication(id)
		if err != nil {
			h.Logger.Warn("edit dataset: could not find dataset with id:", "error", err, "id", id)
			render.InternalServerError(w, r, err)
			return
		}

		if !ctx.User.CanEditPublication(pub) {
			h.Logger.Warn("edit dataset: user isn't allowed to edit the dataset:", "error", err, "id", id, "user", ctx.User.ID)
			render.Forbidden(w, r)
			return
		}

		fn(w, r, Context{
			BaseContext: ctx,
			Publication: pub,
		})
	})
}
