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

type PublicationProjects struct {
	engine *engine.Engine
	render *render.Render
}

func NewPublicationProjects(e *engine.Engine, r *render.Render) *PublicationProjects {
	return &PublicationProjects{
		engine: e,
		render: r,
	}
}

func (p *PublicationProjects) ListProjects(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := p.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get 20 random projects (no search, init state)
	hits, _ := p.engine.SuggestProjects("")

	p.render.HTML(w, 200,
		"publication/_projects_modal",
		struct {
			Publication *models.Publication
			Hits        []models.Completion
		}{
			pub,
			hits,
		},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationProjects) ActiveSearch(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := p.engine.GetPublication(id)
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
	hits, _ := p.engine.SuggestProjects(query[0])

	p.render.HTML(w, 200,
		"publication/_projects_modal_hits",
		struct {
			Publication *models.Publication
			Hits        []models.Completion
		}{
			pub,
			hits,
		},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationProjects) AddToPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	projectId := mux.Vars(r)["project_id"]

	pub, err := p.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	project, err := p.engine.GetProject(projectId)
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

	savedPub, _ := p.engine.UpdatePublication(pub)

	// TODO: error handling if project save fails

	p.render.HTML(w, 200,
		"publication/_projects",
		views.NewPublicationData(r, p.render, savedPub),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationProjects) ConfirmRemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	projectId := mux.Vars(r)["project_id"]

	p.render.HTML(w, 200,
		"publication/_projects_modal_confirm_removal",
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

func (p *PublicationProjects) RemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	projectId := mux.Vars(r)["project_id"]

	// TODO: set constraint to research_data
	pub, err := p.engine.GetPublication(id)
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
	savedPub, _ := p.engine.UpdatePublication(pub)

	p.render.HTML(w, 200,
		"publication/_projects",
		views.NewPublicationData(r, p.render, savedPub),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
