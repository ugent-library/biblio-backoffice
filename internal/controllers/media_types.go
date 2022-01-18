package controllers

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/unrolled/render"
)

type MediaTypes struct {
	Context
}

func NewMediaTypes(c Context) *MediaTypes {
	return &MediaTypes{c}
}

// TODO how can we change a param name with htmx?
func (c *MediaTypes) Choose(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")
	if input == "" {
		input = "q"
	}
	q := r.URL.Query().Get(input)

	suggestions, _ := c.Engine.SuggestMediaTypes(q)

	c.Render.HTML(w, http.StatusOK, "media_types/_choose", views.NewData(c.Render, r, struct {
		Suggestions []models.Completion
	}{
		suggestions,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
