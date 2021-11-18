// Publication Authors controller
//
// Manages the listing of authors on the Publication detail page.
//
// HTMX Custom Events:
//  See: https://htmx.org/headers/hx-trigger/
//
// 	ITList
//		The table listing is being refreshed.
//  ITListAfterSwap
//      The table listing is being refreshed, trigger on htmx:AfterSwap
// 	ITAddRow
//		A row w/ inline-edit form for a new author is being added
// 	ITCancelAddRow
//		A row w/ inline-edit form for a new author is being cancelled
// 	ITCreateItem
//		A new author has been added to the publication
// 	ITEditRow
//		A row w/ inline-edit form for an existing author is inserted
// 	ITCancelEditRow
//		A row w/ inline-edit form for an existing author is being cancelled
// 	ITUpdateItem
//		An existing author has been updated
//  ITConfirmRemoveFromPublication
//      The confirmation pop-up for removing an author is being shown.
// 	ITRemoveItem
//		An existing author has been removed

package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type PublicationAuthors struct {
	Context
}

func NewPublicationAuthors(c Context) *PublicationAuthors {
	return &PublicationAuthors{c}
}

func (c *PublicationAuthors) List(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	w.Header().Set("HX-Trigger", "ITList")
	w.Header().Set("HX-Trigger-After-Swap", "ITListAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITListAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		"publication/authors/_default_table_body",
		views.NewData(c.Render, r, views.NewContributorData(c.Render, pub, nil, "0")),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationAuthors) AddRow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	rowDelta++

	muxRowDelta = strconv.Itoa(rowDelta)

	// Skeleton to make the render fields happy
	author := &models.Contributor{}

	w.Header().Set("HX-Trigger", "ITAddRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITAddRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITAddRowAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		"publication/authors/_default_form",
		views.NewData(c.Render, r, views.NewContributorForm(c.Render, id, author, muxRowDelta, nil)),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationAuthors) ShiftRow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	// Note: we don't increment the delta in this method!

	// Skeleton to make the render fields happy
	author := &models.Contributor{}

	w.Header().Set("HX-Trigger", "ITAddRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITAddRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITAddRowAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		"publication/authors/_default_form",
		views.NewData(c.Render, r, views.NewContributorForm(c.Render, id, author, muxRowDelta, nil)),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationAuthors) CancelAddRow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Trigger", "ITCancelAddRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITCancelAddRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITCancelAddRowAfterSettle")

	// Empty content, denotes we deleted the row
	fmt.Fprintf(w, "")
}

func (c *PublicationAuthors) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	author := &models.Contributor{}

	id := r.Form["ID"]

	if id[0] != "" {
		log.Println("Submitted ugent author")
		log.Println(id)
		log.Println(id[0])
		// Submitted an UGent author

		// Check if the user really exists
		user, err := c.Engine.GetPerson(id[0])
		log.Println(user)
		log.Println(err)
		if err != nil {
			// TODO: throw appropriate error
			return
		}

		// Use the registered values from the user, avoid relying on user input.
		author.ID = user.ID
		author.FirstName = user.FirstName
		author.LastName = user.LastName
		log.Println(author)
	} else {
		// Submitted an external member

		if err := forms.Decode(author, r.Form); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	placeholder := models.Contributor{}

	authors := make([]models.Contributor, len(pub.Author))
	copy(authors, pub.Author)

	authors = append(authors, placeholder)
	copy(authors[rowDelta+1:], authors[rowDelta:])
	authors[rowDelta] = *author
	pub.Author = authors

	savedPub, err := c.Engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK,
			"publication/authors/_default_form",
			views.NewData(c.Render, r, views.NewContributorForm(c.Render, savedPub.ID, author, muxRowDelta, formErrors)),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Use the SavedAuthor since Librecat returns Author.FullName
	savedAuthor := &savedPub.Author[rowDelta]

	w.Header().Set("HX-Trigger", "ITCreateItem")
	w.Header().Set("HX-Trigger-After-Swap", "ITCreateItemAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITCreateItemAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		"publication/authors/_default_row",
		views.NewData(c.Render, r, views.NewContributorData(c.Render, savedPub, savedAuthor, muxRowDelta)),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationAuthors) EditRow(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	// Skeleton to make the render fields happy
	author := &pub.Author[rowDelta]

	w.Header().Set("HX-Trigger", "ITEditRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITEditRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITEditRowAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		"publication/authors/_default_form_edit",
		views.NewData(c.Render, r, views.NewContributorForm(c.Render, pub.ID, author, muxRowDelta, nil)),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationAuthors) CancelEditRow(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	author := &pub.Author[rowDelta]

	w.Header().Set("HX-Trigger", "ITCancelEditRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITCancelEditRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITCancelEditRowAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		"publication/authors/_default_row",
		views.NewData(c.Render, r, views.NewContributorData(c.Render, pub, author, muxRowDelta)),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationAuthors) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	author := &models.Contributor{}

	id := r.Form["ID"]

	if id[0] != "" {
		log.Println("Submitted ugent author")
		log.Println(id)
		log.Println(id[0])
		// Submitted an UGent author

		// Check if the user really exists
		user, err := c.Engine.GetPerson(id[0])
		log.Println(user)
		log.Println(err)
		if err != nil {
			// TODO: throw appropriate error
			return
		}

		// Use the registered values from the user, avoid relying on user input.
		author.ID = user.ID
		author.FirstName = user.FirstName
		author.LastName = user.LastName
		log.Println(author)
	} else {
		// Submitted an external member

		if err := forms.Decode(author, r.Form); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	authors := make([]models.Contributor, len(pub.Author))
	copy(authors, pub.Author)

	authors[rowDelta] = *author
	pub.Author = authors

	savedPub, err := c.Engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK,
			fmt.Sprintf("publication/authors/_%s_edit_form", pub.Type),
			views.NewData(c.Render, r, views.NewContributorForm(c.Render, savedPub.ID, author, muxRowDelta, formErrors)),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	} else if err != nil {
		// TODO: throw appropriate error
		return
	}

	// Use the SavedAuthor since Librecat returns Author.FullName
	savedAuthor := &savedPub.Author[rowDelta]

	w.Header().Set("HX-Trigger", "ITUpdateItem")
	w.Header().Set("HX-Trigger-After-Swap", "ITUpdateItemAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITUpdateItemAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		"publication/authors/_default_row",
		views.NewData(c.Render, r, views.NewContributorData(c.Render, savedPub, savedAuthor, muxRowDelta)),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationAuthors) ConfirmRemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	w.Header().Set("HX-Trigger", "ITConfirmRemoveFromPublication")
	w.Header().Set("HX-Trigger-After-Swap", "ITConfirmRemoveFromPublicationAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITConfirmRemoveFromPublicationAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		"publication/_authors_modal_confirm_removal",
		views.NewData(c.Render, r, struct {
			ID          string
			AuthorDelta string
		}{
			id,
			muxRowDelta,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationAuthors) RemoveAuthor(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	authors := make([]models.Contributor, len(pub.Author))
	copy(authors, pub.Author)

	authors = append(authors[:rowDelta], authors[rowDelta+1:]...)
	pub.Author = authors

	// TODO: error handling
	c.Engine.UpdatePublication(pub)

	w.Header().Set("HX-Trigger", "ITRemoveItem")
	w.Header().Set("HX-Trigger-After-Swap", "ITRemoveItemAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITRemoveItemAfterSettle")

	// Empty content, denotes we deleted the record
	fmt.Fprintf(w, "")
}

// @todo
//   Temporarily disabling dragging / re-ordering authors. It's a complex feature which
//   might introduce complex bugs. May re-enable this later again when there's a real need
//   for this feature.
//
// func (c *PublicationAuthors) OrderAuthors(w http.ResponseWriter, r *http.Request) {
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

func (c *PublicationAuthors) PromoteSearchAuthor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	person := &struct {
		FirstName string
		LastName  string
	}{}

	if err := forms.Decode(person, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := person.FirstName + " " + person.LastName
	people, _ := c.Engine.SuggestPersons(q)

	// pub, err := c.Engine.GetPublication(id)
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, err.Error(), http.StatusNotFound)
	// 	return
	// }

	length := strconv.Itoa(len(people))

	w.Header().Set("HX-Trigger", "ITPromoteModal")
	w.Header().Set("HX-Trigger-After-Swap", "ITPromoteModalAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITPromoteModalAfterSettle")

	c.Render.HTML(w, 200,
		"publication/_authors_modal_promote_author",
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
