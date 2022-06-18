package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/context"
	"github.com/unrolled/render"
)

type PublicationProjects struct {
	Base
	store                backends.Repository
	projectSearchService backends.ProjectSearchService
	projectService       backends.ProjectService
}

func NewPublicationProjects(base Base, store backends.Repository,
	projectSearchService backends.ProjectSearchService,
	projectSerive backends.ProjectService) *PublicationProjects {
	return &PublicationProjects{
		Base:                 base,
		store:                store,
		projectSearchService: projectSearchService,
		projectService:       projectSerive,
	}
}

func (c *PublicationProjects) List(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	// Get 20 random projects (no search, init state)
	hits, _ := c.projectSearchService.SuggestProjects("")

	c.Render.HTML(w, http.StatusOK, "publication/projects/_modal", c.ViewData(r, struct {
		Publication *models.Publication
		Hits        []models.Completion
	}{
		pub,
		hits,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationProjects) ActiveSearch(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get 20 results from the search query
	query := r.Form["search"]
	hits, _ := c.projectSearchService.SuggestProjects(query[0])

	c.Render.HTML(w, http.StatusOK, "publication/projects/_modal_hits", c.ViewData(r, struct {
		Publication *models.Publication
		Hits        []models.Completion
	}{
		pub,
		hits,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationProjects) Add(w http.ResponseWriter, r *http.Request) {
	projectId := mux.Vars(r)["project_id"]

	pub := context.GetPublication(r.Context())

	project, err := c.projectService.GetProject(projectId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	pub.Project = append(pub.Project, models.PublicationProject{
		ID:   projectId,
		Name: project.Title,
	})

	savedPub := pub.Clone()
	err = c.store.SavePublication(savedPub)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/projects/_show", c.ViewData(r, struct {
		Publication *models.Publication
	}{
		savedPub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationProjects) ConfirmRemove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	projectId := mux.Vars(r)["project_id"]

	c.Render.HTML(w, http.StatusOK, "publication/projects/_modal_confirm_removal", c.ViewData(r, struct {
		ID        string
		ProjectID string
	}{
		id,
		projectId,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationProjects) Remove(w http.ResponseWriter, r *http.Request) {
	projectId := mux.Vars(r)["project_id"]

	pub := context.GetPublication(r.Context())

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
	savedPub := pub.Clone()
	err := c.store.SavePublication(savedPub)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/projects/_show", c.ViewData(r, struct {
		Publication *models.Publication
	}{
		savedPub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
