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

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_add", views.NewData(c.Render, r, struct {
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
		views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
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

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK, "publication/contributors/_add", views.NewData(c.Render, r, struct {
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
			views.NewFormBuilder(c.Render, locale.Get(r.Context()), formErrors),
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savedContributor := savedPub.Contributors(role)[position]

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_insert_row", views.NewData(c.Render, r, struct {
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

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_edit", views.NewData(c.Render, r, struct {
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
		views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
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

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK, "publication/contributors/_edit", views.NewData(c.Render, r, struct {
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
			views.NewFormBuilder(c.Render, locale.Get(r.Context()), formErrors),
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savedContributor := savedPub.Contributors(role)[position]

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_update_row", views.NewData(c.Render, r, struct {
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

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_confirm_remove", views.NewData(c.Render, r, struct {
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

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_table", views.NewData(c.Render, r, struct {
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

	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")

	suggestions, _ := c.Engine.SuggestPeople(firstName + " " + lastName)

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_choose", views.NewData(c.Render, r, struct {
		Role        string
		Publication *models.Publication
		Position    int
		Suggestions []models.Person
	}{
		role,
		pub,
		position,
		suggestions,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
