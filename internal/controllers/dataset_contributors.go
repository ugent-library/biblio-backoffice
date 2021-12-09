package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type DatasetContributors struct {
	Context
}

func NewDatasetContributors(c Context) *DatasetContributors {
	return &DatasetContributors{c}
}

func (c *DatasetContributors) Add(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())
	role := mux.Vars(r)["role"]
	positionVar := r.URL.Query().Get("position")
	contributors := dataset.Contributors(role)
	position := len(contributors)
	if positionVar != "" {
		position, _ = strconv.Atoi(positionVar)
	}

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_add", views.NewData(c.Render, r, struct {
		Role        string
		Dataset     *models.Dataset
		Contributor *models.Contributor
		Position    int
		Form        *views.FormBuilder
	}{
		role,
		dataset,
		&models.Contributor{},
		position,
		views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetContributors) Create(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())
	role := mux.Vars(r)["role"]
	r.ParseForm()
	positionVar := r.FormValue("position")
	contributors := dataset.Contributors(role)
	position := len(contributors)
	if positionVar != "" {
		position, _ = strconv.Atoi(positionVar)
	}

	contributor := &models.Contributor{}

	id := r.FormValue("id")
	if id != "" {
		// Check if the user really exists
		user, err := c.Engine.GetPerson(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		contributor.ID = user.ID
		contributor.CreditRole = r.Form["credit_role"]
		contributor.FirstName = user.FirstName
		contributor.LastName = user.LastName
	} else {
		contributor.FirstName = r.FormValue("first_name")
		contributor.LastName = r.FormValue("last_name")
	}

	dataset.AddContributor(role, position, contributor)

	savedDataset, err := c.Engine.UpdateDataset(dataset)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK, "dataset/contributors/_add", views.NewData(c.Render, r, struct {
			Role        string
			Dataset     *models.Dataset
			Contributor *models.Contributor
			Position    int
			Form        *views.FormBuilder
		}{
			role,
			dataset,
			contributor,
			position,
			views.NewFormBuilder(c.Render, locale.Get(r.Context()), formErrors),
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savedContributor := savedDataset.Contributors(role)[position]

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_insert_row", views.NewData(c.Render, r, struct {
		Role        string
		Dataset     *models.Dataset
		Contributor *models.Contributor
		Position    int
	}{
		role,
		savedDataset,
		savedContributor,
		position,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetContributors) Edit(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())
	role := mux.Vars(r)["role"]
	positionVar := mux.Vars(r)["position"]
	position, _ := strconv.Atoi(positionVar)

	contributor := dataset.Contributors(role)[position]

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_edit", views.NewData(c.Render, r, struct {
		Role        string
		Dataset     *models.Dataset
		Contributor *models.Contributor
		Position    int
		Form        *views.FormBuilder
	}{
		role,
		dataset,
		contributor,
		position,
		views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetContributors) Update(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())
	role := mux.Vars(r)["role"]
	positionVar := mux.Vars(r)["position"]
	position, _ := strconv.Atoi(positionVar)

	contributor := &models.Contributor{}

	r.ParseForm()

	id := r.FormValue("id")
	if id != "" {
		// Check if the user really exists
		user, err := c.Engine.GetPerson(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		contributor.ID = user.ID
		contributor.CreditRole = r.Form["credit_role"]
		contributor.FirstName = user.FirstName
		contributor.LastName = user.LastName
	} else {
		contributor.FirstName = r.FormValue("first_name")
		contributor.LastName = r.FormValue("last_name")
	}

	dataset.Contributors(role)[position] = contributor

	savedDataset, err := c.Engine.UpdateDataset(dataset)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK, "dataset/contributors/_edit", views.NewData(c.Render, r, struct {
			Role        string
			Dataset     *models.Dataset
			Contributor *models.Contributor
			Position    int
			Form        *views.FormBuilder
		}{
			role,
			dataset,
			contributor,
			position,
			views.NewFormBuilder(c.Render, locale.Get(r.Context()), formErrors),
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savedContributor := savedDataset.Contributors(role)[position]

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_update_row", views.NewData(c.Render, r, struct {
		Role        string
		Dataset     *models.Dataset
		Contributor *models.Contributor
		Position    int
	}{
		role,
		savedDataset,
		savedContributor,
		position,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetContributors) ConfirmRemove(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())
	role := mux.Vars(r)["role"]
	positionVar := mux.Vars(r)["position"]
	position, _ := strconv.Atoi(positionVar)

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_confirm_remove", views.NewData(c.Render, r, struct {
		Role     string
		Dataset  *models.Dataset
		Position int
	}{
		role,
		dataset,
		position,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetContributors) Remove(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())
	role := mux.Vars(r)["role"]
	positionVar := mux.Vars(r)["position"]
	position, _ := strconv.Atoi(positionVar)

	dataset.RemoveContributor(role, position)

	if _, err := c.Engine.UpdateDataset(dataset); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_table", views.NewData(c.Render, r, struct {
		Role    string
		Dataset *models.Dataset
	}{
		role,
		dataset,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetContributors) Choose(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())
	role := mux.Vars(r)["role"]
	positionVar := mux.Vars(r)["position"]
	position, _ := strconv.Atoi(positionVar)

	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")

	suggestions, _ := c.Engine.SuggestPeople(firstName + " " + lastName)

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_choose", views.NewData(c.Render, r, struct {
		Role        string
		Dataset     *models.Dataset
		Position    int
		Suggestions []models.Person
	}{
		role,
		dataset,
		position,
		suggestions,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
