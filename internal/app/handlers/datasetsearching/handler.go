package datasetsearching

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
	DatasetSearchService backends.DatasetSearchService
}

type Context struct {
	handlers.BaseContext
	SearchArgs *models.SearchArgs
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil {
			render.Unauthorized(w, r)
			return
		}

		searchArgs := models.NewSearchArgs()
		if err := bind.Request(r, searchArgs); err != nil {
			render.BadRequest(w, r, err)
			return
		}

		fn(w, r, Context{
			BaseContext: ctx,
			SearchArgs:  searchArgs,
		})
	})
}
