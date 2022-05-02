package controllers

import (
	"errors"
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

type PublicationConference struct {
	Context
}

func NewPublicationConference(c Context) *PublicationConference {
	return &PublicationConference{c}
}

func (c *PublicationConference) Show(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	c.Render.HTML(w, http.StatusOK,
		"publication/conference/_show",
		c.ViewData(r, struct {
			Publication *models.Publication
			Show        *views.ShowBuilder
		}{
			pub,
			views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationConference) Edit(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	c.Render.HTML(w, http.StatusOK,
		"publication/conference/_edit",
		c.ViewData(r, struct {
			Publication *models.Publication
			Form        *views.FormBuilder
		}{
			pub,
			views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationConference) Update(w http.ResponseWriter, r *http.Request) {
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

	savedPub := pub.Clone()
	err = c.Engine.UpdatePublication(savedPub)

	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {
		c.Render.HTML(w, http.StatusOK,
			"publication/conference/_edit",
			c.ViewData(r, struct {
				Publication *models.Publication
				Form        *views.FormBuilder
			}{
				pub,
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

	c.Render.HTML(w, http.StatusOK,
		"publication/conference/_update",
		c.ViewData(r, struct {
			Publication *models.Publication
			Show        *views.ShowBuilder
		}{
			savedPub,
			views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		},
			views.Flash{Type: "sucess", Message: "Conference updated succesfully", DismissAfter: 5 * time.Second},
		),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
