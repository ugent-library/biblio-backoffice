package controllers

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
)

type Licenses struct {
	Context
}

func NewLicenses(c Context) *Licenses {
	return &Licenses{c}
}

// TODO how can we change a param name with htmx?
func (c *Licenses) Choose(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")
	if input == "" {
		input = "q"
	}
	q := r.URL.Query().Get(input)

	suggestions, _ := c.Engine.SuggestLicenses(q)

	c.Render.HTML(w, http.StatusOK, "licenses/_choose", c.ViewData(r, struct {
		Suggestions []models.Completion
	}{
		suggestions,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
