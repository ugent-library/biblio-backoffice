package controllers

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
)

type Licenses struct {
	Base
	licenseSearchService backends.LicenseSearchService
}

func NewLicenses(base Base, licenseSearchService backends.LicenseSearchService) *Licenses {
	return &Licenses{
		Base:                 base,
		licenseSearchService: licenseSearchService,
	}
}

// TODO how can we change a param name with htmx?
func (c *Licenses) Choose(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")
	if input == "" {
		input = "q"
	}
	q := r.URL.Query().Get(input)

	suggestions, _ := c.licenseSearchService.SuggestLicenses(q)

	c.Render.HTML(w, http.StatusOK, "licenses/_choose", c.ViewData(r, struct {
		Suggestions []models.Completion
	}{
		suggestions,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
