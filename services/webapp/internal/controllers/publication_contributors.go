package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
	"github.com/unrolled/render"
)

type PublicationContributors struct {
	Context
}

func NewPublicationContributors(c Context) *PublicationContributors {
	return &PublicationContributors{c}
}

func (c *PublicationContributors) Add(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())
	role := mux.Vars(r)["role"]
	positionVar := r.URL.Query().Get("position")
	contributors := pub.Contributors(role)
	position := len(contributors)
	if positionVar != "" {
		position, _ = strconv.Atoi(positionVar)
	}

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_add", c.ViewData(r, struct {
		Role        string
		Publication *models.Publication
		Contributor *models.Contributor
		Position    int
		Form        *views.FormBuilder
	}{
		role,
		pub,
		&models.Contributor{},
		position,
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) Create(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())
	role := mux.Vars(r)["role"]
	r.ParseForm()
	positionVar := r.FormValue("position")
	contributors := pub.Contributors(role)
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
		contributor.CreditRole = r.Form["credit_role"]
		contributor.FirstName = r.FormValue("first_name")
		contributor.LastName = r.FormValue("last_name")
	}

	pub.AddContributor(role, position, contributor)

	savedPub, err := c.Engine.UpdatePublication(pub)

	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {
		c.Render.HTML(w, http.StatusOK, "publication/contributors/_add", c.ViewData(r, struct {
			Role        string
			Publication *models.Publication
			Contributor *models.Contributor
			Position    int
			Form        *views.FormBuilder
		}{
			role,
			pub,
			contributor,
			position,
			views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), validationErrors),
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savedContributor := savedPub.Contributors(role)[position]

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_insert_row", c.ViewData(r, struct {
		Role        string
		Publication *models.Publication
		Contributor *models.Contributor
		Position    int
	}{
		role,
		savedPub,
		savedContributor,
		position,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) Edit(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())
	role := mux.Vars(r)["role"]
	positionVar := mux.Vars(r)["position"]
	position, _ := strconv.Atoi(positionVar)

	contributor := pub.Contributors(role)[position]

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_edit", c.ViewData(r, struct {
		Role        string
		Publication *models.Publication
		Contributor *models.Contributor
		Position    int
		Form        *views.FormBuilder
	}{
		role,
		pub,
		contributor,
		position,
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) Update(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())
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
		contributor.CreditRole = r.Form["credit_role"]
		contributor.FirstName = r.FormValue("first_name")
		contributor.LastName = r.FormValue("last_name")
	}

	pub.Contributors(role)[position] = contributor

	savedPub, err := c.Engine.UpdatePublication(pub)

	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {
		c.Render.HTML(w, http.StatusOK, "publication/contributors/_edit", c.ViewData(r, struct {
			Role        string
			Publication *models.Publication
			Contributor *models.Contributor
			Position    int
			Form        *views.FormBuilder
		}{
			role,
			pub,
			contributor,
			position,
			views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), validationErrors),
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savedContributor := savedPub.Contributors(role)[position]

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_update_row", c.ViewData(r, struct {
		Role        string
		Publication *models.Publication
		Contributor *models.Contributor
		Position    int
	}{
		role,
		savedPub,
		savedContributor,
		position,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) ConfirmRemove(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())
	role := mux.Vars(r)["role"]
	positionVar := mux.Vars(r)["position"]
	position, _ := strconv.Atoi(positionVar)

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_confirm_remove", c.ViewData(r, struct {
		Role        string
		Publication *models.Publication
		Position    int
	}{
		role,
		pub,
		position,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) Remove(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())
	role := mux.Vars(r)["role"]
	positionVar := mux.Vars(r)["position"]
	position, _ := strconv.Atoi(positionVar)

	pub.RemoveContributor(role, position)

	if _, err := c.Engine.UpdatePublication(pub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_table", c.ViewData(r, struct {
		Role        string
		Publication *models.Publication
	}{
		role,
		pub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) Choose(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())
	role := mux.Vars(r)["role"]
	positionVar := mux.Vars(r)["position"]
	position, _ := strconv.Atoi(positionVar)

	// firstName := r.URL.Query().Get("first_name")
	// lastName := r.URL.Query().Get("last_name")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	creditRole := r.Form["credit_role"]

	suggestions, _ := c.Engine.SuggestPeople(firstName + " " + lastName)

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_choose", c.ViewData(r, struct {
		Role        string
		CreditRole  []string
		Publication *models.Publication
		Position    int
		Suggestions []models.Person
	}{
		role,
		creditRole,
		pub,
		position,
		suggestions,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) Demote(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())
	role := mux.Vars(r)["role"]
	positionVar := mux.Vars(r)["position"]
	position, _ := strconv.Atoi(positionVar)

	contributor := pub.Contributors(role)[position]

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

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_edit", c.ViewData(r, struct {
		Role        string
		Publication *models.Publication
		Contributor *models.Contributor
		Position    int
		Form        *views.FormBuilder
	}{
		role,
		pub,
		contributor,
		position,
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) Promote(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())
	role := mux.Vars(r)["role"]
	positionVar := mux.Vars(r)["position"]
	position, _ := strconv.Atoi(positionVar)

	contributor := pub.Contributors(role)[position]

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

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_edit", c.ViewData(r, struct {
		Role        string
		Publication *models.Publication
		Contributor *models.Contributor
		Position    int
		Form        *views.FormBuilder
	}{
		role,
		pub,
		contributor,
		position,
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) Order(w http.ResponseWriter, r *http.Request) {
	publication := context.GetPublication(r.Context())
	role := mux.Vars(r)["role"]
	contributors := publication.Contributors(role)
	newContributors := make([]*models.Contributor, len(contributors))

	r.ParseForm()

	for i, v := range r.Form["position"] {
		pos, _ := strconv.Atoi(v)
		newContributors[i] = contributors[pos]
	}

	publication.SetContributors(role, newContributors)

	if _, err := c.Engine.UpdatePublication(publication); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_table", c.ViewData(r, struct {
		Role        string
		Publication *models.Publication
	}{
		role,
		publication,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
