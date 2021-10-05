package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/unrolled/render"
)

type DatasetProjects struct {
	engine *engine.Engine
	render *render.Render
}

func NewDatasetProjects(e *engine.Engine, r *render.Render) *DatasetProjects {
	return &DatasetProjects{
		engine: e,
		render: r,
	}
}

func (d *DatasetProjects) ListProjects(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// TODO: set constraint to research_data
	pub, err := d.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	d.render.HTML(w, 200,
		"dataset/_projects_modal",
		views.NewDatasetData(r, d.render, pub),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (d *DatasetProjects) AddToDataset(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	projectId := mux.Vars(r)["project_id"]

	// TODO: set constraint to research_data
	pub, err := d.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// TODO: get the project based on the ID from the LibreCat REST API
	project := models.PublicationProject{
		ID:   projectId,
		Name: "Ankh Morkpocian project granted by the Patrician to the University of Magic",
	}
	pub.Project = append(pub.Project, project)

	savedPub, _ := d.engine.UpdatePublication(pub)

	// TODO: error handling if project save fails

	d.render.HTML(w, 200,
		"dataset/_projects",
		views.NewDatasetData(r, d.render, savedPub),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (d *DatasetProjects) ConfirmRemoveFromDataset(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	projectId := mux.Vars(r)["project_id"]

	d.render.HTML(w, 200,
		"dataset/_projects_modal_confirm_removal",
		struct {
			ID        string
			ProjectID string
		}{
			id,
			projectId,
		},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (d *DatasetProjects) RemoveFromDataset(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	projectId := mux.Vars(r)["project_id"]

	// TODO: set constraint to research_data
	pub, err := d.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	projects := make([]models.PublicationProject, len(pub.Project))
	copy(projects, pub.Project)

	var removeKey int
	for key, project := range projects {
		if project.ID == projectId {
			removeKey = key
		}
	}

	projects = append(projects[:removeKey], projects[removeKey+1:]...)
	pub.Project = projects

	// TODO: error handling
	savedPub, _ := d.engine.UpdatePublication(pub)

	d.render.HTML(w, 200,
		"dataset/_projects",
		views.NewDatasetData(r, d.render, savedPub),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
