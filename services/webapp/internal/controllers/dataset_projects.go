package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/unrolled/render"
)

type DatasetProjects struct {
	Base
}

func NewDatasetProjects(c Base) *DatasetProjects {
	return &DatasetProjects{c}
}

func (c *DatasetProjects) Choose(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	// Get 20 random projects (no search, init state)
	hits, _ := c.Engine.SuggestProjects("")

	c.Render.HTML(w, http.StatusOK, "dataset/projects/_modal", c.ViewData(r, struct {
		Dataset *models.Dataset
		Hits    []models.Completion
	}{
		dataset,
		hits,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetProjects) ActiveSearch(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get 20 results from the search query
	query := r.Form["search"]
	hits, _ := c.Engine.SuggestProjects(query[0])

	c.Render.HTML(w, http.StatusOK, "dataset/projects/_modal_hits", c.ViewData(r, struct {
		Dataset *models.Dataset
		Hits    []models.Completion
	}{
		dataset,
		hits,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetProjects) Add(w http.ResponseWriter, r *http.Request) {
	projectId := mux.Vars(r)["project_id"]

	dataset := context.GetDataset(r.Context())

	project, err := c.Engine.GetProject(projectId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	dataset.Project = append(dataset.Project, models.DatasetProject{
		ID:   projectId,
		Name: project.Title,
	})

	savedDataset := dataset.Clone()
	err = c.Engine.Store.UpdateDataset(dataset)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/projects/_show", c.ViewData(r, struct {
		Dataset *models.Dataset
	}{
		savedDataset,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetProjects) ConfirmRemove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	projectId := mux.Vars(r)["project_id"]

	c.Render.HTML(w, http.StatusOK, "dataset/projects/_modal_confirm_removal", c.ViewData(r, struct {
		ID        string
		ProjectID string
	}{
		id,
		projectId,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetProjects) Remove(w http.ResponseWriter, r *http.Request) {
	projectId := mux.Vars(r)["project_id"]

	dataset := context.GetDataset(r.Context())

	projects := make([]models.DatasetProject, len(dataset.Project))
	copy(projects, dataset.Project)

	var removeKey int
	for key, project := range projects {
		if project.ID == projectId {
			removeKey = key
		}
	}

	projects = append(projects[:removeKey], projects[removeKey+1:]...)
	dataset.Project = projects

	savedDataset := dataset.Clone()
	err := c.Engine.Store.UpdateDataset(dataset)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/projects/_show", c.ViewData(r, struct {
		Dataset *models.Dataset
	}{
		savedDataset,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
