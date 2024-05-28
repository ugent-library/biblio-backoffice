package mediatypes

import (
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/views/media_types"
	"github.com/ugent-library/httperror"
)

func Suggest(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	// TODO how can we change a param name with htmx?
	input := r.URL.Query().Get("input")
	if input == "" {
		input = "q"
	}
	query := r.URL.Query().Get(input)

	hits, err := c.Services.MediaTypeSearchService.SuggestMediaTypes(query)
	if err != nil {
		c.HandleError(w, r, httperror.InternalServerError.Wrap(fmt.Errorf("could not suggest mediatypes: %w", err)))
		return
	}

	media_types.Suggest(c, query, hits).Render(r.Context(), w)
}
