package datasetcreating

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
	Repository           backends.Repository
	DatasetSearchService backends.DatasetSearchService
	DatasetSources       map[string]backends.DatasetGetter
}

type Context struct {
	handlers.BaseContext
	SearchArgs *models.SearchArgs
	Dataset    *models.Dataset
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil {
			render.Unauthorized(w, r)
			return
		}

		// TODO Needed because called in dataset/show_nav which is reused in add_description
		searchArgs := models.NewSearchArgs()
		if err := bind.Request(r, searchArgs); err != nil {
			render.BadRequest(w, r, err)
			return
		}

		context := Context{
			BaseContext: ctx,
			SearchArgs:  searchArgs,
		}

		if id := bind.PathValues(r).Get("id"); id != "" {
			d, err := h.Repository.GetDataset(id)
			if err != nil {
				render.InternalServerError(w, r, err)
				return
			}

			if !ctx.User.CanEditDataset(d) {
				render.Forbidden(w, r)
				return
			}

			context.Dataset = d
		}

		fn(w, r, context)
	})
}
