package mediatypes

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
)

type Handler struct {
	handlers.BaseHandler
	MediaTypeSearchService backends.MediaTypeSearchService
}

type Context struct {
	handlers.BaseContext
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil {
			h.Logger.Warnw("mediatypes: user is not authorized to access this resource:", "user", ctx.User.ID)
			render.Unauthorized(w, r)
			return
		}

		fn(w, r, Context{
			BaseContext: ctx,
		})
	})
}

type YieldSuggest struct {
	Context
	Hits  []models.Completion
	Query string
}

func (h *Handler) Suggest(w http.ResponseWriter, r *http.Request, ctx Context) {
	// TODO how can we change a param name with htmx?
	input := r.URL.Query().Get("input")
	if input == "" {
		input = "q"
	}
	q := r.URL.Query().Get(input)

	hits, err := h.MediaTypeSearchService.SuggestMediaTypes(q)
	if err != nil {
		h.Logger.Errorw("suggest mediatype: could not suggest mediatypes:", "errors", err, "query", q, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "media_types/suggest", YieldSuggest{
		Context: ctx,
		Hits:    hits,
		Query:   q,
	})
}
