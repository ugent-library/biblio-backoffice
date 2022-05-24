package controllers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/unrolled/render"
)

type PublicationDepartments struct {
	Base
	store                     backends.Store
	organizationSearchService backends.OrganizationSearchService
	organizationService       backends.OrganizationService
}

func NewPublicationDepartments(base Base, store backends.Store,
	organizationSearchService backends.OrganizationSearchService,
	organizationService backends.OrganizationService) *PublicationDepartments {
	return &PublicationDepartments{
		Base:                      base,
		store:                     store,
		organizationSearchService: organizationSearchService,
		organizationService:       organizationService,
	}
}

func (c *PublicationDepartments) List(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	// Get 20 random departments (no search, init state)
	hits, _ := c.organizationSearchService.SuggestOrganizations("")

	c.Render.HTML(w, http.StatusOK, "publication/departments/_modal", c.ViewData(r, struct {
		Publication *models.Publication
		Hits        []models.Completion
	}{
		pub,
		hits,
	}),
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
	hits, _ := c.organizationSearchService.SuggestOrganizations(query[0])

	c.Render.HTML(w, http.StatusOK, "publication/departments/_modal_hits", c.ViewData(r, struct {
		Publication *models.Publication
		Hits        []models.Completion
	}{
		pub,
		hits,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDepartments) Add(w http.ResponseWriter, r *http.Request) {
	departmentID := mux.Vars(r)["department_id"]
	// because mux var {department_id} is not properly decoded
	// i.e. LW06%2A remains LW06%2A instead of LW06*
	departmentID, _ = url.PathUnescape(departmentID)

	pub := context.GetPublication(r.Context())

	org, err := c.organizationService.GetOrganization(departmentID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	pubDepartment := models.PublicationDepartment{
		ID: org.ID,
	}
	for _, o := range org.Tree {
		pubDepartment.Tree = append(pubDepartment.Tree, models.PublicationDepartmentRef{ID: o.ID})
	}
	pub.Department = append(pub.Department, pubDepartment)
	savedPub := pub.Clone()

	if err := c.store.UpdatePublication(savedPub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/departments/_show", c.ViewData(r, struct {
		Publication *models.Publication
	}{
		savedPub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDepartments) ConfirmRemove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	departmentId := mux.Vars(r)["department_id"]
	// because mux var {department_id} is not properly decoded
	// i.e. LW06%2A remains LW06%2A instead of LW06*
	departmentId, _ = url.PathUnescape(departmentId)

	c.Render.HTML(w, http.StatusOK, "publication/departments/_modal_confirm_removal", c.ViewData(r, struct {
		ID           string
		DepartmentID string
	}{
		id,
		departmentId,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDepartments) Remove(w http.ResponseWriter, r *http.Request) {
	departmentId := mux.Vars(r)["department_id"]
	// because mux var {department_id} is not properly decoded
	// i.e. LW06%2A remains LW06%2A instead of LW06*
	departmentId, _ = url.PathUnescape(departmentId)

	pub := context.GetPublication(r.Context())

	departments := make([]models.PublicationDepartment, len(pub.Department))
	copy(departments, pub.Department)

	var removeKey int = -1
	for key, department := range departments {
		if department.ID == departmentId {
			removeKey = key
		}
	}
	if removeKey >= 0 {
		departments = append(departments[:removeKey], departments[removeKey+1:]...)
	}
	pub.Department = departments

	savedPub := pub.Clone()
	err := c.store.UpdatePublication(pub)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/departments/_show", c.ViewData(r, struct {
		Publication *models.Publication
	}{
		savedPub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
