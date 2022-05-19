package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/unrolled/render"
)

type PublicationContributors struct {
	Base
	store               backends.Store
	personSearchService backends.PersonSearchService
	personService       backends.PersonService
}

func NewPublicationContributors(base Base, store backends.Store, personSearchService backends.PersonSearchService,
	personService backends.PersonService) *PublicationContributors {
	return &PublicationContributors{
		Base:                base,
		store:               store,
		personSearchService: personSearchService,
		personService:       personService,
	}
}

func (c *PublicationContributors) Add(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())
	role := mux.Vars(r)["role"]
	position := len(pub.Contributors(role))
	q := r.URL.Query()

	c.Render.HTML(w, http.StatusOK, "publication/contributors/_add", c.ViewData(r, struct {
		Role        string
		Publication *models.Publication
		Contributor *models.Contributor
		Position    int
		Form        *views.FormBuilder
	}{
		role,
		pub,
		&models.Contributor{
			CreditRole: q["credit_role"],
			FirstName:  q.Get("first_name"),
			LastName:   q.Get("last_name"),
		},
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
		user, err := c.personService.GetPerson(id)
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

	savedPub := pub.Clone()
	err := c.store.UpdatePublication(savedPub)

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
		user, err := c.personService.GetPerson(id)
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

	savedPub := pub.Clone()
	err := c.store.UpdatePublication(savedPub)

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

	if err := c.store.UpdatePublication(pub); err != nil {
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

	suggestions, _ := c.personSearchService.SuggestPeople(firstName + " " + lastName)

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
	position, _ := strconv.Atoi(mux.Vars(r)["position"])

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contributor := &models.Contributor{
		CreditRole: r.Form["credit_role"],
		FirstName:  r.FormValue("first_name"),
		LastName:   r.FormValue("last_name"),
	}

	var tmpl string
	if len(pub.Contributors(role)) > position {
		tmpl = "publication/contributors/_edit"
	} else {
		tmpl = "publication/contributors/_add"
	}

	c.Render.HTML(w, http.StatusOK, tmpl, c.ViewData(r, struct {
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
	position, _ := strconv.Atoi(mux.Vars(r)["position"])

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	person, err := c.personService.GetPerson(r.FormValue("id"))
	if err != nil || person == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contributor := &models.Contributor{
		CreditRole: r.Form["credit_role"],
		FirstName:  person.FirstName,
		ID:         person.ID,
		LastName:   person.LastName,
	}

	var tmpl string
	if len(pub.Contributors(role)) > position {
		tmpl = "publication/contributors/_edit"
	} else {
		tmpl = "publication/contributors/_add"
	}

	c.Render.HTML(w, http.StatusOK, tmpl, c.ViewData(r, struct {
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

	if err := c.store.UpdatePublication(publication); err != nil {
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
