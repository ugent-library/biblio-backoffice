package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type PublicationDetails struct {
	engine *engine.Engine
	render *render.Render
}

func NewPublicationDetails(e *engine.Engine, r *render.Render) *PublicationDetails {
	return &PublicationDetails{
		engine: e,
		render: r,
	}
}

func (c *PublicationDetails) Show(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.render.HTML(w, 200,
		"publication/details/_show",
		views.NewData(c.render, r, struct {
			Publication  *models.Publication
			Show         *views.ShowBuilder
			Vocabularies map[string][]string
		}{
			pub,
			views.NewShowBuilder(c.render, locale.Get(r.Context())),
			c.engine.Vocabularies(),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDetails) OpenForm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.render.HTML(w, 200,
		"publication/details/_edit",
		views.NewData(c.render, r, struct {
			Publication  *models.Publication
			Form         *views.FormBuilder
			Vocabularies map[string][]string
		}{
			pub,
			views.NewFormBuilder(c.render, locale.Get(r.Context()), nil),
			c.engine.Vocabularies(),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDetails) SaveForm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := forms.Decode(pub, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savedPub, err := c.engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.render.HTML(w, 200,
			"publication/details/_edit",
			views.NewData(c.render, r, struct {
				Publication  *models.Publication
				Form         *views.FormBuilder
				Vocabularies map[string][]string
			}{
				pub,
				views.NewFormBuilder(c.render, locale.Get(r.Context()), formErrors),
				c.engine.Vocabularies(),
			}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.render.HTML(w, 200,
		"publication/details/_update",
		views.NewData(c.render, r, struct {
			Publication  *models.Publication
			Show         *views.ShowBuilder
			Vocabularies map[string][]string
		}{
			savedPub,
			views.NewShowBuilder(c.render, locale.Get(r.Context())),
			c.engine.Vocabularies(),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
