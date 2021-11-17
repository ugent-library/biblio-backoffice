package controllers

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type PublicationDetails struct {
	Context
}

func NewPublicationDetails(c Context) *PublicationDetails {
	return &PublicationDetails{c}
}

func (c *PublicationDetails) Show(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	c.Render.HTML(w, http.StatusOK, "publication/details/_show", views.NewData(c.Render, r, struct {
		Publication  *models.Publication
		Show         *views.ShowBuilder
		Vocabularies map[string][]string
	}{
		pub,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
		c.Engine.Vocabularies(),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDetails) OpenForm(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	c.Render.HTML(w, http.StatusOK, "publication/details/_edit", views.NewData(c.Render, r, struct {
		Publication  *models.Publication
		Form         *views.FormBuilder
		Vocabularies map[string][]string
	}{
		pub,
		views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
		c.Engine.Vocabularies(),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDetails) SaveForm(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := forms.Decode(pub, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savedPub, err := c.Engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK, "publication/details/_edit", views.NewData(c.Render, r, struct {
			Publication  *models.Publication
			Form         *views.FormBuilder
			Vocabularies map[string][]string
		}{
			pub,
			views.NewFormBuilder(c.Render, locale.Get(r.Context()), formErrors),
			c.Engine.Vocabularies(),
		},
			views.Flash{Type: "error", Message: "There are some problems with your input"},
		),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/details/_update", views.NewData(c.Render, r, struct {
		Publication  *models.Publication
		Show         *views.ShowBuilder
		Vocabularies map[string][]string
	}{
		savedPub,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
		c.Engine.Vocabularies(),
	},
		views.Flash{Type: "success", Message: "Details updated succesfully"},
	),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
