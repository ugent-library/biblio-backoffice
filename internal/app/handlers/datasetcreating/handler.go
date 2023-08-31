package datasetcreating

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
	DatasetSearchIndex  backends.DatasetIndex
	DatasetSources      map[string]backends.DatasetGetter
	OrganizationService backends.OrganizationService
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

		context := Context{
			BaseContext: ctx,
		}

		if id := bind.PathValues(r).Get("id"); id != "" {
			d, err := h.Repository.GetDataset(id)
			if err != nil {
				render.NotFound(w, r, err)
				return
			}

			if !ctx.User.CanEditDataset(d) {
				h.Logger.Warn("create dataset: user isn't allowed to edit the dataset:", "error", err, "dataset", id, "user", ctx.User.ID)
				render.Forbidden(w, r)
				return
			}

			context.Dataset = d
		}

		fn(w, r, context)
	})
}
