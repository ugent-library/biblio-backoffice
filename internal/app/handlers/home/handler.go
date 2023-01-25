package home

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/internal/app/handlers"
	"github.com/ugent-library/biblio-backoffice/internal/render"
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

type YieldHome struct {
	Context
	PageTitle string
	ActiveNav string
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "layouts/default", "pages/home", YieldHome{
		Context:   ctx,
		PageTitle: "Biblio",
	})
}
