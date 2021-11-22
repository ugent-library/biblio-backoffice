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
	ctype := mux.Vars(r)["type"]
	pub := context.GetPublication(r.Context())

	w.Header().Set("HX-Trigger", "ITList")
	w.Header().Set("HX-Trigger-After-Swap", "ITListAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITListAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		fmt.Sprintf("publication/%s/_default_table_body", ctype),
		views.NewData(c.Render, r, struct {
			Publication *models.Publication
			Show        *views.ShowBuilder
			Author      *models.Contributor
			Key         string
		}{
			pub,
			views.NewShowBuilder(c.Render, locale.Get(r.Context())),
			nil,
			"0",
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) AddRow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ctype := mux.Vars(r)["type"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	rowDelta++

	muxRowDelta = strconv.Itoa(rowDelta)

	// Skeleton to make the render fields happy
	contributor := &models.Contributor{}

	w.Header().Set("HX-Trigger", "ITAddRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITAddRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITAddRowAfterSettle")

	c.Render.HTML(w, http.StatusOK, fmt.Sprintf("publication/%s/_default_form", ctype), views.NewData(c.Render, r, struct {
		Author *models.Contributor
		Form   *views.FormBuilder
		ID     string
		Key    string
	}{
		contributor,
		views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
		id,
		muxRowDelta,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) ShiftRow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ctype := mux.Vars(r)["type"]
	muxRowDelta := mux.Vars(r)["delta"]

	// Note: we don't increment the delta in this method!

	// Skeleton to make the render fields happy
	contributor := &models.Contributor{}

	w.Header().Set("HX-Trigger", "ITAddRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITAddRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITAddRowAfterSettle")

	c.Render.HTML(w, http.StatusOK, fmt.Sprintf("publication/%s/_default_form", ctype), views.NewData(c.Render, r, struct {
		Author *models.Contributor
		Form   *views.FormBuilder
		ID     string
		Key    string
	}{
		contributor,
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
	ctype := mux.Vars(r)["type"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contributor := &models.Contributor{}

	// Add the contributor to the publication
	switch ctype {
	case "authors":
		// Authors can be "UGent Author" / "External member".
		// Separate form processing if a submitted author is an "UGent Author"
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

		c.Engine.AddAuthorToPublication(pub, contributor, rowDelta)
	default:
		// Throw an error, unkown type
		return
	}

	savedPub, err := c.Engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK, fmt.Sprintf("publication/%s/_default_form", ctype), views.NewData(c.Render, r, struct {
			Author *models.Contributor
			Form   *views.FormBuilder
			ID     string
			Key    string
		}{
			contributor,
			views.NewFormBuilder(c.Render, locale.Get(r.Context()), formErrors),
			savedPub.ID,
			muxRowDelta,
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Use the SavedContributor since Librecat returns contributor.FullName
	var savedContributor *models.Contributor
	switch ctype {
	case "authors":
		savedContributor = c.Engine.GetAuthorFromPublication(savedPub, rowDelta)
	default:
		// @todo: Throw an error, unkown type
	}

	w.Header().Set("HX-Trigger", "ITCreateItem")
	w.Header().Set("HX-Trigger-After-Swap", "ITCreateItemAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITCreateItemAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		fmt.Sprintf("publication/%s/_default_row", ctype),
		views.NewData(c.Render, r, struct {
			Publication *models.Publication
			Show        *views.ShowBuilder
			Author      *models.Contributor
			Key         string
		}{
			savedPub,
			views.NewShowBuilder(c.Render, locale.Get(r.Context())),
			savedContributor,
			muxRowDelta,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) EditRow(w http.ResponseWriter, r *http.Request) {
	ctype := mux.Vars(r)["type"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	var contributor *models.Contributor
	switch ctype {
	case "authors":
		contributor = c.Engine.GetAuthorFromPublication(pub, rowDelta)
	default:
		// @todo: Throw an error, unkown type
		return

	}

	w.Header().Set("HX-Trigger", "ITEditRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITEditRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITEditRowAfterSettle")

	c.Render.HTML(w, http.StatusOK, fmt.Sprintf("publication/%s/_default_form_edit", ctype), views.NewData(c.Render, r, struct {
		Author *models.Contributor
		Form   *views.FormBuilder
		ID     string
		Key    string
	}{
		contributor,
		views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
		pub.ID,
		muxRowDelta,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) CancelEditRow(w http.ResponseWriter, r *http.Request) {
	ctype := mux.Vars(r)["type"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	var contributor *models.Contributor
	switch ctype {
	case "authors":
		contributor = c.Engine.GetAuthorFromPublication(pub, rowDelta)
	default:
		// @todo: Throw an error, unkown type
		return
	}

	w.Header().Set("HX-Trigger", "ITCancelEditRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITCancelEditRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITCancelEditRowAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		fmt.Sprintf("publication/%s/_default_row", ctype),
		views.NewData(c.Render, r, struct {
			render      *render.Render
			Publication *models.Publication
			Author      *models.Contributor
			Key         string
		}{
			c.Render,
			pub,
			contributor,
			muxRowDelta,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) UpdateContributor(w http.ResponseWriter, r *http.Request) {
	ctype := mux.Vars(r)["type"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contributor := &models.Contributor{}

	// Update the contributor
	switch ctype {
	case "authors":
		// Authors can be "UGent Author" / "External member".
		// Separate form processing if a submitted author is an "UGent Author"
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
		c.Engine.UpdateAuthorOnPublication(pub, contributor, rowDelta)
	default:
		return
	}

	savedPub, err := c.Engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK, fmt.Sprintf("publication/%s/%s_edit_form", ctype, savedPub.Type), views.NewData(c.Render, r, struct {
			Author *models.Contributor
			Form   *views.FormBuilder
			ID     string
			Key    string
		}{
			contributor,
			views.NewFormBuilder(c.Render, locale.Get(r.Context()), formErrors),
			savedPub.ID,
			muxRowDelta,
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)
		return
	} else if err != nil {
		// @todo: throw appropriate error if saving the publication fails
		return
	}

	// Use the SavedContributor since Librecat returns contributor.FullName
	var savedContributor *models.Contributor
	switch ctype {
	case "authors":
		savedContributor = c.Engine.GetAuthorFromPublication(savedPub, rowDelta)
	default:
		// @todo: Throw an error, unkown type
		return
	}

	w.Header().Set("HX-Trigger", "ITUpdateItem")
	w.Header().Set("HX-Trigger-After-Swap", "ITUpdateItemAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITUpdateItemAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		fmt.Sprintf("publication/%s/_default_row", ctype),
		views.NewData(c.Render, r, struct {
			Publication *models.Publication
			Show        *views.ShowBuilder
			Author      *models.Contributor
			Key         string
		}{
			savedPub,
			views.NewShowBuilder(c.Render, locale.Get(r.Context())),
			savedContributor,
			muxRowDelta,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) ConfirmRemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	ctype := mux.Vars(r)["type"]
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	w.Header().Set("HX-Trigger", "ITConfirmRemoveFromPublication")
	w.Header().Set("HX-Trigger-After-Swap", "ITConfirmRemoveFromPublicationAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITConfirmRemoveFromPublicationAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		fmt.Sprintf("publication/%s/_modal_confirm_removal", ctype),
		views.NewData(c.Render, r, struct {
			ID               string
			ContributorDelta string
		}{
			id,
			muxRowDelta,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationContributors) RemoveContributor(w http.ResponseWriter, r *http.Request) {
	ctype := mux.Vars(r)["type"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	switch ctype {
	case "authors":
		c.Engine.RemoveAuthorFromPublication(pub, rowDelta)
	default:
		return
	}

	// @todo: error handling
	c.Engine.UpdatePublication(pub)

	w.Header().Set("HX-Trigger", "ITRemoveItem")
	w.Header().Set("HX-Trigger-After-Swap", "ITRemoveItemAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITRemoveItemAfterSettle")

	// Empty content, denotes we deleted the record
	fmt.Fprintf(w, "")
}

func (c *PublicationContributors) PromoteSearchContributor(w http.ResponseWriter, r *http.Request) {
	ctype := mux.Vars(r)["type"]
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

	c.Render.HTML(w, 200,
		fmt.Sprintf("publication/%s/_modal_promote_contributor", ctype),
		views.NewData(c.Render, r, struct {
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
