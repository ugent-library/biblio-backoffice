package controllers

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
)

type MediaTypes struct {
	Base
	mediaTypeSearchService backends.MediaTypeSearchService
}

func NewMediaTypes(base Base, mediaTypeSearchService backends.MediaTypeSearchService) *MediaTypes {
	return &MediaTypes{
		Base:                   base,
		mediaTypeSearchService: mediaTypeSearchService,
	}
}

// TODO how can we change a param name with htmx?
func (c *MediaTypes) Choose(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")
	if input == "" {
		input = "q"
	}
	q := r.URL.Query().Get(input)

	suggestions, _ := c.mediaTypeSearchService.SuggestMediaTypes(q)

	c.Render.HTML(w, http.StatusOK, "media_types/_choose", c.ViewData(r, struct {
		Suggestions []models.Completion
	}{
		suggestions,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
