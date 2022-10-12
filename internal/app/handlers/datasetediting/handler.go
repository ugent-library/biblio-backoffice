package datasetediting

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/backends"
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
	PublicationSearchService  backends.PublicationSearchService
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

		id := bind.PathValues(r).Get("id")
		d, err := h.Repository.GetDataset(id)
		if err != nil {
			render.NotFound(w, r, err)
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
