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
		views.NewData(d.render, r, struct {
			Dataset *models.Publication
			Hits    []models.Completion
		}{
			pub,
			hits,
		}),
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
		views.NewData(d.render, r, struct {
			Dataset *models.Publication
			Hits    []models.Completion
		}{
			pub,
			hits,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetProjects) AddToDataset(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	projectId := mux.Vars(r)["project_id"]

	dataset, err := c.engine.GetDataset(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	project, err := c.engine.GetProject(projectId)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	publicationProject := models.PublicationProject{
		ID:   projectId,
		Name: project.Name,
	}
	dataset.Project = append(dataset.Project, publicationProject)

	savedDataset, _ := c.engine.UpdateDataset(dataset)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.render.HTML(w, 200,
		"dataset/_projects",
		views.NewData(c.render, r, struct {
			Dataset *models.Dataset
		}{
			savedDataset,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetProjects) ConfirmRemoveFromDataset(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	projectId := mux.Vars(r)["project_id"]

	c.render.HTML(w, 200,
		"dataset/_projects_modal_confirm_removal",
		views.NewData(c.render, r, struct {
			ID        string
			ProjectID string
		}{
			id,
			projectId,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetProjects) RemoveFromDataset(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	projectId := mux.Vars(r)["project_id"]

	dataset, err := c.engine.GetDataset(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	projects := make([]models.PublicationProject, len(dataset.Project))
	copy(projects, dataset.Project)

	var removeKey int
	for key, project := range projects {
		if project.ID == projectId {
			removeKey = key
		}
	}

	projects = append(projects[:removeKey], projects[removeKey+1:]...)
	dataset.Project = projects

	// TODO: error handling
	savedDataset, _ := c.engine.UpdateDataset(dataset)

	c.render.HTML(w, 200,
		"dataset/_projects",
		views.NewData(c.render, r, struct {
			Dataset *models.Dataset
		}{
			savedDataset,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
