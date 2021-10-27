package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type DatasetDetails struct {
	Context
}

func NewDatasetDetails(c Context) *DatasetDetails {
	return &DatasetDetails{c}
}

func (c *DatasetDetails) Show(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// TODO: set constraint to research_data
	dataset, err := c.Engine.GetDataset(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.Render.HTML(w, 200, "dataset/details/_show",
		views.NewData(c.Render, r, struct {
			Dataset *models.Dataset
			Show    *views.ShowBuilder
		}{
			dataset,
			views.NewShowBuilder(c.Render, locale.Get(r.Context())),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetDetails) OpenForm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	dataset, err := c.Engine.GetDataset(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.Render.HTML(w, http.StatusOK,
		"dataset/details/_edit",
		views.NewData(c.Render, r, struct {
			Dataset *models.Dataset
			Form    *views.FormBuilder
		}{
			dataset,
			views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetDetails) SaveForm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// TODO: set constriant to research_data
	dataset, err := c.Engine.GetDataset(id)
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

	if err := forms.Decode(dataset, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savedDataset, err := c.Engine.UpdateDataset(dataset)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, 200,
			"dataset/details/_edit",
			views.NewData(c.Render, r, struct {
				Dataset *models.Dataset
				Form    *views.FormBuilder
			}{
				dataset,
				views.NewFormBuilder(c.Render, locale.Get(r.Context()), formErrors),
			}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	}

	c.Render.HTML(w, 200,
		"dataset/details/_update",
		views.NewData(c.Render, r, struct {
			Dataset *models.Dataset
			Show    *views.ShowBuilder
		}{
			savedDataset,
			views.NewShowBuilder(c.Render, locale.Get(r.Context())),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
