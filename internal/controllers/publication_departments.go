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

type PublicationDepartments struct {
	engine *engine.Engine
	render *render.Render
}

func NewPublicationDepartments(e *engine.Engine, r *render.Render) *PublicationDepartments {
	return &PublicationDepartments{
		engine: e,
		render: r,
	}
}

func (p *PublicationDepartments) ListDepartments(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := p.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get 20 random departments (no search, init state)
	hits, _ := p.engine.SuggestDepartments("")

	p.render.HTML(w, 200,
		"publication/_departments_modal",
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

func (p *PublicationDepartments) ActiveSearch(w http.ResponseWriter, r *http.Request) {
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
	hits, _ := p.engine.SuggestDepartments(query[0])

	p.render.HTML(w, 200,
		"publication/_departments_modal_hits",
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

func (p *PublicationDepartments) AddToPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	departmentId := mux.Vars(r)["department_id"]

	pub, err := p.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// department, err := p.engine.GetDepartment(departmentId)
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, err.Error(), http.StatusNotFound)
	// 	return
	// }

	publicationDepartment := models.PublicationDepartment{
		ID: departmentId,
	}
	pub.Department = append(pub.Department, publicationDepartment)
	savedPub, _ := p.engine.UpdatePublication(pub)

	// TODO: error handling if department save fails

	p.render.HTML(w, 200,
		"publication/_departments",
		views.NewPublicationData(r, p.render, savedPub),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationDepartments) ConfirmRemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	departmentId := mux.Vars(r)["department_id"]

	p.render.HTML(w, 200,
		"publication/_departments_modal_confirm_removal",
		struct {
			ID           string
			DepartmentID string
		}{
			id,
			departmentId,
		},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationDepartments) RemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	departmentId := mux.Vars(r)["department_id"]

	pub, err := p.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	departments := make([]models.PublicationDepartment, len(pub.Department))
	copy(departments, pub.Department)

	var removeKey int
	for key, department := range departments {
		if department.ID == departmentId {
			removeKey = key
		}
	}

	departments = append(departments[:removeKey], departments[removeKey+1:]...)
	pub.Department = departments

	// TODO: error handling
	savedPub, _ := p.engine.UpdatePublication(pub)

	p.render.HTML(w, 200,
		"publication/_departments",
		views.NewPublicationData(r, p.render, savedPub),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
