package controllers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
	"github.com/unrolled/render"
)

type PublicationDetails struct {
	Base
	store backends.Store
}

func NewPublicationDetails(base Base, store backends.Store) *PublicationDetails {
	return &PublicationDetails{
		Base:  base,
		store: store,
	}
}

func (c *PublicationDetails) Show(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	c.Render.HTML(w, http.StatusOK, "publication/details/_show", c.ViewData(r, struct {
		Publication *models.Publication
		Show        *views.ShowBuilder
	}{
		pub,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDetails) Edit(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	c.Render.HTML(w, http.StatusOK, "publication/details/_edit", c.ViewData(r, struct {
		Publication *models.Publication
		Show        *views.ShowBuilder
		Form        *views.FormBuilder
	}{
		pub,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
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
	if err := DecodeForm(pub, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO handle checkbox boolean values elegantly
	if r.FormValue("extern") != "true" {
		pub.Extern = false
	}

	savedPub := pub.Clone()
	err = c.store.UpdatePublication(pub)

	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {
		log.Printf("%+v", validationErrors)
		c.Render.HTML(w, http.StatusOK, "publication/details/_edit", c.ViewData(r, struct {
			Publication *models.Publication
			Show        *views.ShowBuilder
			Form        *views.FormBuilder
		}{
			pub,
			views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
			views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), validationErrors),
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
		Publication *models.Publication
		Show        *views.ShowBuilder
	}{
		savedPub,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
	},
		views.Flash{Type: "success", Message: "Details updated succesfully", DismissAfter: 5 * time.Second},
	),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
