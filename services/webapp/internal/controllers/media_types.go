package controllers

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
)

type MediaTypes struct {
	Base
}

func NewMediaTypes(c Base) *MediaTypes {
	return &MediaTypes{c}
}

// TODO how can we change a param name with htmx?
func (c *MediaTypes) Choose(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")
	if input == "" {
		input = "q"
	}
	q := r.URL.Query().Get(input)

	suggestions, _ := c.Services.SuggestMediaTypes(q)

	c.Render.HTML(w, http.StatusOK, "media_types/_choose", c.ViewData(r, struct {
		Suggestions []models.Completion
	}{
		suggestions,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
