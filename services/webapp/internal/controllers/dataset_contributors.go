package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
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

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_add", c.ViewData(r, struct {
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
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
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
		c.Render.HTML(w, http.StatusOK, "dataset/contributors/_add", c.ViewData(r, struct {
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
			views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), formErrors),
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savedContributor := savedDataset.Contributors(role)[position]

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_insert_row", c.ViewData(r, struct {
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

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_edit", c.ViewData(r, struct {
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
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
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
		c.Render.HTML(w, http.StatusOK, "dataset/contributors/_edit", c.ViewData(r, struct {
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
			views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), formErrors),
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savedContributor := savedDataset.Contributors(role)[position]

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_update_row", c.ViewData(r, struct {
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

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_confirm_remove", c.ViewData(r, struct {
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

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_table", c.ViewData(r, struct {
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

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_choose", c.ViewData(r, struct {
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

func (c *DatasetContributors) Demote(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())
	role := mux.Vars(r)["role"]
	positionVar := mux.Vars(r)["position"]
	position, _ := strconv.Atoi(positionVar)

	contributor := dataset.Contributors(role)[position]

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := forms.Decode(contributor, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Demote contributor from "UGent member" to "External member"
	// We do this by resetting the ID field to an empty string
	contributor.ID = ""

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_edit", c.ViewData(r, struct {
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
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetContributors) Promote(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())
	role := mux.Vars(r)["role"]
	positionVar := mux.Vars(r)["position"]
	position, _ := strconv.Atoi(positionVar)

	contributor := dataset.Contributors(role)[position]

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Promoting the user from "external member" to "UGent member"
	// The form contains an ID field value pushed from the "choose modal"
	// This value gets relayed to the edit form via the Contributor model.
	// UGent FirstName and LastName are interspersed into the Contributor model
	// as well.
	if err := forms.Decode(contributor, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_edit", c.ViewData(r, struct {
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
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetContributors) Order(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())
	role := mux.Vars(r)["role"]
	contributors := dataset.Contributors(role)
	newContributors := make([]*models.Contributor, len(contributors))

	r.ParseForm()

	for i, v := range r.Form["position"] {
		pos, _ := strconv.Atoi(v)
		newContributors[i] = contributors[pos]
	}

	dataset.SetContributors(role, newContributors)

	if _, err := c.Engine.UpdateDataset(dataset); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/contributors/_table", c.ViewData(r, struct {
		Role    string
		Dataset *models.Dataset
	}{
		role,
		dataset,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
