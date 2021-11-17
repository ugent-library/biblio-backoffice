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
		views.NewData(c.Render, r, struct {
			Publication *models.Publication
			Show        *views.ShowBuilder
		}{
			pub,
			views.NewShowBuilder(c.Render, locale.Get(r.Context())),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationConference) OpenForm(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	c.Render.HTML(w, http.StatusOK,
		"publication/conference/_edit",
		views.NewData(c.Render, r, struct {
			Publication *models.Publication
			Form        *views.FormBuilder
		}{
			pub,
			views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationConference) SaveForm(w http.ResponseWriter, r *http.Request) {
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
		c.Render.HTML(w, http.StatusOK,
			"publication/conference/_edit",
			views.NewData(c.Render, r, struct {
				Publication *models.Publication
				Form        *views.FormBuilder
			}{
				pub,
				views.NewFormBuilder(c.Render, locale.Get(r.Context()), formErrors),
			}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK,
		"publication/conference/_update",
		views.NewData(c.Render, r, struct {
			Publication *models.Publication
			Show        *views.ShowBuilder
		}{
			savedPub,
			views.NewShowBuilder(c.Render, locale.Get(r.Context())),
		},
			views.Flash{Type: "sucess", Message: "Conference updated succesfully"},
		),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
