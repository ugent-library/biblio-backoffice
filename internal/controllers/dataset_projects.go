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

	// Get 20 random projects (no search, init state)
	hits, _ := d.engine.SuggestProjects("")

	d.render.HTML(w, 200,
		"dataset/_projects_modal",
		struct {
			Dataset *models.Publication
			Hits    []models.Completion
		}{
			pub,
			hits,
		},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (d *DatasetProjects) ActiveSearch(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// TODO: set constraint to research_data
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

	// Get 20 results from the search query
	query := r.Form["search"]
	hits, _ := d.engine.SuggestProjects(query[0])

	d.render.HTML(w, 200,
		"dataset/_projects_modal_hits",
		struct {
			Dataset *models.Publication
			Hits    []models.Completion
		}{
			pub,
			hits,
		},
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

	project, err := d.engine.GetProject(projectId)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	publicationProject := models.PublicationProject{
		ID:   projectId,
		Name: project.Name,
	}
	pub.Project = append(pub.Project, publicationProject)

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
