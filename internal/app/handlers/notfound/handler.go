package notfound

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/render"
)

type Handler struct {
	handlers.BaseHandler
}

type Context struct {
	handlers.BaseContext
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		fn(w, r, Context{
			BaseContext: ctx,
		})
	})
}

type YieldNotFound struct {
	Context
	PageTitle string
	ActiveNav string
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.NotFoundLayout(w, "layouts/default", "pages/notfound", YieldNotFound{
		Context:   ctx,
		PageTitle: "Biblio",
	})
}
