package mediatypes

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/views/media_types"
)

type Handler struct {
	handlers.BaseHandler
	MediaTypeSearchService backends.MediaTypeSearchService
}

func (h *Handler) Suggest(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	// TODO how can we change a param name with htmx?
	input := r.URL.Query().Get("input")
	if input == "" {
		input = "q"
	}
	query := r.URL.Query().Get(input)

	hits, err := h.MediaTypeSearchService.SuggestMediaTypes(query)
	if err != nil {
		h.Logger.Errorw("suggest mediatype: could not suggest mediatypes:", "errors", err, "query", query, "user", c.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	media_types.Suggest(c, query, hits).Render(r.Context(), w)
}
