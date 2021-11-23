// Publication Contributors controller
//
// Manages the listing of contributors (authors, editors,...) on the Publication detail page.
//
// HTMX Custom Events:
//  See: https://htmx.org/headers/hx-trigger/
//
// 	ITList
//		The table listing is being refreshed.
//  ITListAfterSwap
//      The table listing is being refreshed, trigger on htmx:AfterSwap
// 	ITAddRow
//		A row w/ inline-edit form for a new contributor is being added
// 	ITCancelAddRow
//		A row w/ inline-edit form for a new contributor is being cancelled
// 	ITCreateItem
//		A new contributor has been added to the publication
// 	ITEditRow
//		A row w/ inline-edit form for an existing contributor is inserted
// 	ITCancelEditRow
//		A row w/ inline-edit form for an existing contributor is being cancelled
// 	ITUpdateItem
//		An existing contributor has been updated
//  ITConfirmRemoveFromPublication
//      The confirmation pop-up for removing an contributor is being shown.
// 	ITRemoveItem
//		An existing contributor has been removed

package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type PublicationContributors struct {
	Context
}

func NewPublicationContributors(c Context) *PublicationContributors {
	return &PublicationContributors{c}
}

func (c *PublicationContributors) List(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	w.Header().Set("HX-Trigger", "ITList")
	w.Header().Set("HX-Trigger-After-Swap", "ITListAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITListAfterSettle")

	c.Render.HTML(w, http.StatusOK, "contributors/_table_body", views.NewData(c.Render, r, struct {
		Role        string
		Publication *models.Publication
		Show        *views.ShowBuilder
	}{
		mux.Vars(r)["role"],
		pub,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) AddRow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	rowDelta++

	muxRowDelta = strconv.Itoa(rowDelta)

	w.Header().Set("HX-Trigger", "ITAddRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITAddRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITAddRowAfterSettle")

	c.Render.HTML(w, http.StatusOK, "contributors/_form", views.NewData(c.Render, r, struct {
		Role        string
		Contributor *models.Contributor
		Form        *views.FormBuilder
		ID          string
		Key         string
	}{
		mux.Vars(r)["role"],
		&models.Contributor{},
		views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
		id,
		muxRowDelta,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) ShiftRow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	// Note: we don't increment the delta in this method!

	w.Header().Set("HX-Trigger", "ITAddRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITAddRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITAddRowAfterSettle")

	c.Render.HTML(w, http.StatusOK, "contributors/_form", views.NewData(c.Render, r, struct {
		Role        string
		Contributor *models.Contributor
		Form        *views.FormBuilder
		ID          string
		Key         string
	}{
		mux.Vars(r)["role"],
		&models.Contributor{},
		views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
		id,
		muxRowDelta,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) CancelAddRow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Trigger", "ITCancelAddRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITCancelAddRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITCancelAddRowAfterSettle")

	// Empty content, denotes we deleted the row
	fmt.Fprintf(w, "")
}

func (c *PublicationContributors) CreateContributor(w http.ResponseWriter, r *http.Request) {
	role := mux.Vars(r)["role"]
	muxDelta := mux.Vars(r)["delta"]
	delta, _ := strconv.Atoi(muxDelta)

	pub := context.GetPublication(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contributor := &models.Contributor{}

	id := r.Form["ID"]
	if id[0] != "" {
		// Check if the user really exists
		user, err := c.Engine.GetPerson(id[0])
		if err != nil {
			// @todo: throw an error
			return
		}
		contributor.ID = user.ID
		contributor.CreditRole = r.Form["credit_role"]
		contributor.FirstName = user.FirstName
		contributor.LastName = user.LastName
	} else {
		if err := forms.Decode(contributor, r.Form); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	pub.AddContributor(role, delta, contributor)

	savedPub, err := c.Engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK, "contributors/_form", views.NewData(c.Render, r, struct {
			Role        string
			Contributor *models.Contributor
			Form        *views.FormBuilder
			ID          string
			Key         string
		}{
			role,
			contributor,
			views.NewFormBuilder(c.Render, locale.Get(r.Context()), formErrors),
			savedPub.ID,
			muxDelta,
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Use the SavedContributor since Librecat returns contributor.FullName
	savedContributor := savedPub.Contributors(role)[delta]

	w.Header().Set("HX-Trigger", "ITCreateItem")
	w.Header().Set("HX-Trigger-After-Swap", "ITCreateItemAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITCreateItemAfterSettle")

	c.Render.HTML(w, http.StatusOK, "contributors/_row", views.NewData(c.Render, r, struct {
		Role        string
		Publication *models.Publication
		Show        *views.ShowBuilder
		Contributor *models.Contributor
		Key         string
	}{
		role,
		savedPub,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
		savedContributor,
		muxDelta,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) EditRow(w http.ResponseWriter, r *http.Request) {
	role := mux.Vars(r)["role"]
	muxDelta := mux.Vars(r)["delta"]
	delta, _ := strconv.Atoi(muxDelta)

	pub := context.GetPublication(r.Context())

	contributor := pub.Contributors(role)[delta]

	w.Header().Set("HX-Trigger", "ITEditRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITEditRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITEditRowAfterSettle")

	c.Render.HTML(w, http.StatusOK, "contributors/_form_edit", views.NewData(c.Render, r, struct {
		Role        string
		Contributor *models.Contributor
		Form        *views.FormBuilder
		ID          string
		Key         string
	}{
		role,
		contributor,
		views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
		pub.ID,
		muxDelta,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) CancelEditRow(w http.ResponseWriter, r *http.Request) {
	role := mux.Vars(r)["role"]
	muxDelta := mux.Vars(r)["delta"]
	delta, _ := strconv.Atoi(muxDelta)

	pub := context.GetPublication(r.Context())

	contributor := pub.Contributors(role)[delta]

	w.Header().Set("HX-Trigger", "ITCancelEditRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITCancelEditRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITCancelEditRowAfterSettle")

	c.Render.HTML(w, http.StatusOK, "contributors/_row", views.NewData(c.Render, r, struct {
		Role        string
		Publication *models.Publication
		Contributor *models.Contributor
		Key         string
	}{
		role,
		pub,
		contributor,
		muxDelta,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) UpdateContributor(w http.ResponseWriter, r *http.Request) {
	role := mux.Vars(r)["role"]
	muxDelta := mux.Vars(r)["delta"]
	delta, _ := strconv.Atoi(muxDelta)

	pub := context.GetPublication(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contributor := &models.Contributor{}

	id := r.Form["ID"]
	if id[0] != "" {
		// Check if the user really exists
		user, err := c.Engine.GetPerson(id[0])
		if err != nil {
			// TODO throw an error
			return
		}
		contributor.ID = user.ID
		contributor.CreditRole = r.Form["credit_role"]
		contributor.FirstName = user.FirstName
		contributor.LastName = user.LastName
	} else {
		if err := forms.Decode(contributor, r.Form); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	pub.Contributors(role)[delta] = contributor

	savedPub, err := c.Engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK, "contributors/_form_edit", views.NewData(c.Render, r, struct {
			Role        string
			Contributor *models.Contributor
			Form        *views.FormBuilder
			ID          string
			Key         string
		}{
			role,
			contributor,
			views.NewFormBuilder(c.Render, locale.Get(r.Context()), formErrors),
			savedPub.ID,
			muxDelta,
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)
		return
	} else if err != nil {
		// @todo: throw appropriate error if saving the publication fails
		return
	}

	// Use the SavedContributor since Librecat returns contributor.FullName
	savedContributor := savedPub.Contributors(role)[delta]

	w.Header().Set("HX-Trigger", "ITUpdateItem")
	w.Header().Set("HX-Trigger-After-Swap", "ITUpdateItemAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITUpdateItemAfterSettle")

	c.Render.HTML(w, http.StatusOK, "contributors/_row", views.NewData(c.Render, r, struct {
		Role        string
		Publication *models.Publication
		Show        *views.ShowBuilder
		Contributor *models.Contributor
		Key         string
	}{
		role,
		savedPub,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
		savedContributor,
		muxDelta,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) ConfirmRemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	w.Header().Set("HX-Trigger", "ITConfirmRemoveFromPublication")
	w.Header().Set("HX-Trigger-After-Swap", "ITConfirmRemoveFromPublicationAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITConfirmRemoveFromPublicationAfterSettle")

	c.Render.HTML(w, http.StatusOK, "contributors/_modal_confirm_removal", views.NewData(c.Render, r, struct {
		Role             string
		ID               string
		ContributorDelta string
	}{
		mux.Vars(r)["role"],
		id,
		muxRowDelta,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) RemoveContributor(w http.ResponseWriter, r *http.Request) {
	role := mux.Vars(r)["role"]
	delta, _ := strconv.Atoi(mux.Vars(r)["delta"])

	pub := context.GetPublication(r.Context())
	pub.RemoveContributor(role, delta)
	// @todo: error handling
	c.Engine.UpdatePublication(pub)

	w.Header().Set("HX-Trigger", "ITRemoveItem")
	w.Header().Set("HX-Trigger-After-Swap", "ITRemoveItemAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITRemoveItemAfterSettle")

	// Empty content, denotes we deleted the record
	fmt.Fprintf(w, "")
}

func (c *PublicationContributors) PromoteSearchContributor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	person := &models.Person{}

	if err := forms.Decode(person, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := person.FirstName + " " + person.LastName
	people, _ := c.Engine.SuggestPersons(q)

	length := strconv.Itoa(len(people))

	w.Header().Set("HX-Trigger", "ITPromoteModal")
	w.Header().Set("HX-Trigger-After-Swap", "ITPromoteModalAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITPromoteModalAfterSettle")

	c.Render.HTML(w, 200, "contributors/_modal_promote_contributor", views.NewData(c.Render, r, struct {
		ID       string
		People   []models.Person
		Length   string
		RowDelta string
	}{
		id,
		people,
		length,
		muxRowDelta,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// @todo
//   Temporarily disabling dragging / re-ordering authors. It's a complex feature which
//   might introduce complex bugs. May re-enable this later again when there's a real need
//   for this feature.
//
// func (c *PublicationContributors) OrderAuthors(w http.ResponseWriter, r *http.Request) {
// 	muxStart := mux.Vars(r)["start"]
// 	start, _ := strconv.Atoi(muxStart)

// 	muxEnd := mux.Vars(r)["end"]
// 	end, _ := strconv.Atoi(muxEnd)

// 	pub := context.GetPublication(r.Context())

// 	author := &pub.Author[start]

// 	// Remove the author
// 	authors := make([]models.Contributor, len(pub.Author))
// 	copy(authors, pub.Author)
// 	authors = append(authors[:start], authors[start+1:]...)
// 	pub.Author = authors

// 	// Re-insert the author at the new position
// 	placeholder := models.Contributor{}
// 	authors = append(authors, placeholder)
// 	copy(authors[end+1:], authors[end:])
// 	authors[end] = *author

// 	// Save everything
// 	pub.Author = authors

// 	savedPub, err := c.Engine.UpdatePublication(pub)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("HX-Trigger", "ITOrderAuthors")
// 	w.Header().Set("HX-Trigger-After-Swap", "ITOrderAuthorsAfterSwap")
// 	w.Header().Set("HX-Trigger-After-Settle", "ITOrderAuthorsAfterSettle")

// 	c.Render.HTML(w, http.StatusOK,
// 		"publication/authors/_default_table_body",
// 		views.NewData(c.Render, r, views.NewContributorData(c.Render, savedPub, nil, "0")),
// 		render.HTMLOptions{Layout: "layouts/htmx"},
// 	)
// }
