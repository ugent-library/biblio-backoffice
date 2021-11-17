package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/unrolled/render"
)

type PublicationDepartments struct {
	Context
}

func NewPublicationDepartments(c Context) *PublicationDepartments {
	return &PublicationDepartments{c}
}

func (c *PublicationDepartments) ListDepartments(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	// Get 20 random departments (no search, init state)
	hits, _ := c.Engine.SuggestDepartments("")

	c.Render.HTML(w, http.StatusOK, "publication/departments/_modal", struct {
		Publication *models.Publication
		Hits        []models.Completion
	}{
		pub,
		hits,
	},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDepartments) ActiveSearch(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get 20 results from the search query
	query := r.Form["search"]
	hits, _ := c.Engine.SuggestDepartments(query[0])

	c.Render.HTML(w, http.StatusOK, "publication/departments/_modal_hits", struct {
		Publication *models.Publication
		Hits        []models.Completion
	}{
		pub,
		hits,
	},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDepartments) AddToPublication(w http.ResponseWriter, r *http.Request) {
	departmentId := mux.Vars(r)["department_id"]

	pub := context.GetPublication(r.Context())

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
	savedPub, _ := c.Engine.UpdatePublication(pub)

	// TODO: error handling if department save fails

	c.Render.HTML(w, http.StatusOK, "publication/departments/_show", views.NewData(c.Render, r, struct {
		Publication *models.Publication
	}{
		savedPub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDepartments) ConfirmRemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	departmentId := mux.Vars(r)["department_id"]

	c.Render.HTML(w, http.StatusOK, "publication/departments/_modal_confirm_removal", struct {
		ID           string
		DepartmentID string
	}{
		id,
		departmentId,
	},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDepartments) RemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	departmentId := mux.Vars(r)["department_id"]

	pub := context.GetPublication(r.Context())

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
	savedPub, _ := c.Engine.UpdatePublication(pub)

	c.Render.HTML(w, http.StatusOK, "publication/departments/_show", views.NewData(c.Render, r, struct {
		Publication *models.Publication
	}{
		savedPub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
