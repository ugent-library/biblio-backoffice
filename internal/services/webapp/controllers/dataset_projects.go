package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/context"
	"github.com/unrolled/render"
)

type DatasetProjects struct {
	Base
	store                backends.Repository
	projectSearchService backends.ProjectSearchService
	projectService       backends.ProjectService
}

func NewDatasetProjects(base Base, store backends.Repository, projectSearchService backends.ProjectSearchService,
	projectService backends.ProjectService) *DatasetProjects {
	return &DatasetProjects{
		Base:                 base,
		store:                store,
		projectSearchService: projectSearchService,
		projectService:       projectService,
	}
}

func (c *DatasetProjects) Choose(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	// Get 20 random projects (no search, init state)
	hits, _ := c.projectSearchService.SuggestProjects("")

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
	hits, _ := c.projectSearchService.SuggestProjects(query[0])

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

	project, err := c.projectService.GetProject(projectId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	dataset.Project = append(dataset.Project, models.DatasetProject{
		ID:   projectId,
		Name: project.Title,
	})

	savedDataset := dataset.Clone()
	err = c.store.SaveDataset(savedDataset)

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
	err := c.store.SaveDataset(dataset)

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
