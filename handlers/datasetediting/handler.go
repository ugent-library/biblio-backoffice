package datasetediting

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/repositories"
	"github.com/ugent-library/bind"
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
	PublicationSearchIndex    backends.PublicationIndex
}

type Context struct {
	handlers.BaseContext
	Dataset *models.Dataset
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil {
			render.Unauthorized(w, r)
			return
		}

		id := bind.PathValue(r, "id")
		d, err := h.Repo.GetDataset(id)
		if err != nil {
			if err == models.ErrNotFound {
				render.NotFound(w, r, err)
			} else {
				render.InternalServerError(w, r, err)
			}
			return
		}

		if !ctx.User.CanEditDataset(d) {
			h.Logger.Warn("edit dataset: user isn't allowed to edit the dataset:", "error", err, "dataset", id, "user", ctx.User.ID)
			render.Forbidden(w, r)
			return
		}

		fn(w, r, Context{
			BaseContext: ctx,
			Dataset:     d,
		})
	})
}
