package controllers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
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

	c.Render.HTML(w, http.StatusOK, "publication/details/_show", c.ViewData(r, struct {
		Publication  *models.Publication
		Show         *views.ShowBuilder
		Vocabularies map[string][]string
	}{
		pub,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		c.Engine.Vocabularies(),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDetails) Edit(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	c.Render.HTML(w, http.StatusOK, "publication/details/_edit", c.ViewData(r, struct {
		Publication  *models.Publication
		Show         *views.ShowBuilder
		Form         *views.FormBuilder
		Vocabularies map[string][]string
	}{
		pub,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
		c.Engine.Vocabularies(),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDetails) Update(w http.ResponseWriter, r *http.Request) {
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

	// TODO handle checkbox boolean values elegantly
	if r.FormValue("extern") != "true" {
		pub.Extern = false
	}

	savedPub, err := c.Engine.UpdatePublication(pub)

	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {
		log.Printf("%+v", validationErrors)
		c.Render.HTML(w, http.StatusOK, "publication/details/_edit", c.ViewData(r, struct {
			Publication  *models.Publication
			Show         *views.ShowBuilder
			Form         *views.FormBuilder
			Vocabularies map[string][]string
		}{
			pub,
			views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
			views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), validationErrors),
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

	c.Render.HTML(w, http.StatusOK, "publication/details/_update", c.ViewData(r, struct {
		Publication  *models.Publication
		Show         *views.ShowBuilder
		Vocabularies map[string][]string
	}{
		savedPub,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		c.Engine.Vocabularies(),
	},
		views.Flash{Type: "success", Message: "Details updated succesfully", DismissAfter: 5 * time.Second},
	),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
