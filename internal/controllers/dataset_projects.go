package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/unrolled/render"
)

type DatasetProjects struct {
	Context
}

func NewDatasetProjects(c Context) *DatasetProjects {
	return &DatasetProjects{c}
}

func (c *DatasetProjects) ListProjects(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// TODO: set constraint to research_data
	pub, err := c.Engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get 20 random projects (no search, init state)
	hits, _ := c.Engine.SuggestProjects("")

	c.Render.HTML(w, 200,
		"dataset/_projects_modal",
		views.NewData(c.Render, r, struct {
			Dataset *models.Publication
			Hits    []models.Completion
		}{
			pub,
			hits,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetProjects) ActiveSearch(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// TODO: set constraint to research_data
	pub, err := c.Engine.GetPublication(id)
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
	hits, _ := c.Engine.SuggestProjects(query[0])

	c.Render.HTML(w, 200,
		"dataset/_projects_modal_hits",
		views.NewData(c.Render, r, struct {
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

	dataset, err := c.Engine.GetDataset(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	project, err := c.Engine.GetProject(projectId)
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

	savedDataset, _ := c.Engine.UpdateDataset(dataset)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, 200,
		"dataset/_projects",
		views.NewData(c.Render, r, struct {
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

	c.Render.HTML(w, 200,
		"dataset/_projects_modal_confirm_removal",
		views.NewData(c.Render, r, struct {
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

	dataset, err := c.Engine.GetDataset(id)
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
	savedDataset, _ := c.Engine.UpdateDataset(dataset)

	c.Render.HTML(w, 200,
		"dataset/_projects",
		views.NewData(c.Render, r, struct {
			Dataset *models.Dataset
		}{
			savedDataset,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
