package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/unrolled/render"
)

type PublicationProjects struct {
	Context
}

func NewPublicationProjects(c Context) *PublicationProjects {
	return &PublicationProjects{c}
}

func (c *PublicationProjects) ListProjects(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	// Get 20 random projects (no search, init state)
	hits, _ := c.Engine.SuggestProjects("")

	c.Render.HTML(w, http.StatusOK, "publication/projects/_modal", views.NewData(c.Render, r, struct {
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
	hits, _ := c.Engine.SuggestProjects(query[0])

	c.Render.HTML(w, http.StatusOK, "publication/projects/_modal_hits", views.NewData(c.Render, r, struct {
		Publication *models.Publication
		Hits        []models.Completion
	}{
		pub,
		hits,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationProjects) AddToPublication(w http.ResponseWriter, r *http.Request) {
	projectId := mux.Vars(r)["project_id"]

	pub := context.GetPublication(r.Context())

	project, err := c.Engine.GetProject(projectId)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// TODO: get the project based on the ID from the LibreCat REST API
	publicationProject := models.PublicationProject{
		ID:   projectId,
		Name: project.Name,
	}
	pub.Project = append(pub.Project, publicationProject)

	savedPub, _ := c.Engine.UpdatePublication(pub)

	// TODO: error handling if project save fails

	c.Render.HTML(w, http.StatusOK, "publication/projects/_show", views.NewData(c.Render, r, struct {
		Publication *models.Publication
	}{
		savedPub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationProjects) ConfirmRemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	projectId := mux.Vars(r)["project_id"]

	c.Render.HTML(w, http.StatusOK, "publication/projects/_modal_confirm_removal", views.NewData(c.Render, r, struct {
		ID        string
		ProjectID string
	}{
		id,
		projectId,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationProjects) RemoveFromPublication(w http.ResponseWriter, r *http.Request) {
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
	savedPub, _ := c.Engine.UpdatePublication(pub)

	c.Render.HTML(w, http.StatusOK, "publication/projects/_show", views.NewData(c.Render, r, struct {
		Publication *models.Publication
	}{
		savedPub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
