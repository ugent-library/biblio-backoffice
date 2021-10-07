package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type DatasetDetails struct {
	engine *engine.Engine
	render *render.Render
}

func NewDatasetDetails(e *engine.Engine, r *render.Render) *DatasetDetails {
	return &DatasetDetails{
		engine: e,
		render: r,
	}
}

func (d *DatasetDetails) Show(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// TODO: set constriant to research_data
	pub, err := d.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	d.render.HTML(w, 200,
		"dataset/_details",
		views.NewDatasetData(r, d.render, pub),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (d *DatasetDetails) OpenForm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// TODO: set constriant to research_data
	pub, err := d.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	d.render.HTML(w, 200,
		"dataset/_details_edit_form",
		views.NewDatasetForm(r, d.render, pub, nil),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (d *DatasetDetails) SaveForm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// TODO: set constriant to research_data
	pub, err := d.engine.GetPublication(id)
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

	// TODO: set constriant to research_data
	savedPub, err := d.engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		d.render.HTML(w, 200,
			"dataset/_details_edit_form",
			views.NewDatasetForm(r, d.render, pub, formErrors),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	}

	d.render.HTML(w, 200,
		"dataset/_details_edit_submit",
		views.NewDatasetData(r, d.render, savedPub),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
