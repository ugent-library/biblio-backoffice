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
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type PublicationAuthors struct {
	engine *engine.Engine
	render *render.Render
}

func NewPublicationAuthors(e *engine.Engine, r *render.Render) *PublicationAuthors {
	return &PublicationAuthors{
		engine: e,
		render: r,
	}
}

func (p *PublicationAuthors) List(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := p.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("HX-Trigger", "ITList")
	w.Header().Set("HX-Trigger-After-Swap", "ITListAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITListAfterSettle")

	p.render.HTML(w, 200,
		"publication/authors/_default_table_body",
		views.NewContributorData(r, p.render, pub, nil, "0"),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationAuthors) AddRow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	rowDelta++

	muxRowDelta = strconv.Itoa(rowDelta)

	// Skeleton to make the render fields happy
	author := &models.PublicationContributor{}

	w.Header().Set("HX-Trigger", "ITAddRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITAddRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITAddRowAfterSettle")

	p.render.HTML(w, http.StatusOK,
		"publication/authors/_default_form",
		views.NewContributorForm(r, p.render, id, author, muxRowDelta, nil),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationAuthors) ShiftRow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	// Note: we don't increment the delta in this method!

	// Skeleton to make the render fields happy
	author := &models.PublicationContributor{}

	w.Header().Set("HX-Trigger", "ITAddRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITAddRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITAddRowAfterSettle")

	p.render.HTML(w, http.StatusOK,
		"publication/authors/_default_form",
		views.NewContributorForm(r, p.render, id, author, muxRowDelta, nil),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationAuthors) CancelAddRow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Trigger", "ITCancelAddRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITCancelAddRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITCancelAddRowAfterSettle")

	// Empty content, denotes we deleted the row
	fmt.Fprintf(w, "")
}

func (p *PublicationAuthors) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub, err := p.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	author := &models.PublicationContributor{}

	if err := forms.Decode(author, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	placeholder := models.PublicationContributor{}

	authors := make([]models.PublicationContributor, len(pub.Author))
	copy(authors, pub.Author)

	authors = append(authors, placeholder)
	copy(authors[rowDelta+1:], authors[rowDelta:])
	authors[rowDelta] = *author
	pub.Author = authors

	savedPub, err := p.engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		p.render.HTML(w, 200,
			fmt.Sprintf("publication/authors/_%s_form", pub.Type),
			views.NewContributorForm(r, p.render, savedPub.ID, author, muxRowDelta, formErrors),
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

	p.render.HTML(w, http.StatusOK,
		"publication/authors/_default_row",
		views.NewContributorData(r, p.render, savedPub, savedAuthor, muxRowDelta),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationAuthors) EditRow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub, err := p.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Skeleton to make the render fields happy
	author := &pub.Author[rowDelta]

	w.Header().Set("HX-Trigger", "ITEditRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITEditRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITEditRowAfterSettle")

	p.render.HTML(w, http.StatusOK,
		"publication/authors/_default_form_edit",
		views.NewContributorForm(r, p.render, id, author, muxRowDelta, nil),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationAuthors) CancelEditRow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub, err := p.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	author := &pub.Author[rowDelta]

	w.Header().Set("HX-Trigger", "ITCancelEditRow")
	w.Header().Set("HX-Trigger-After-Swap", "ITCancelEditRowAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITCancelEditRowAfterSettle")

	p.render.HTML(w, http.StatusOK,
		"publication/authors/_default_row",
		views.NewContributorData(r, p.render, pub, author, muxRowDelta),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationAuthors) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub, err := p.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	author := &models.PublicationContributor{}

	if err := forms.Decode(author, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authors := make([]models.PublicationContributor, len(pub.Author))
	copy(authors, pub.Author)

	authors[rowDelta] = *author
	pub.Author = authors

	savedPub, err := p.engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		p.render.HTML(w, 200,
			fmt.Sprintf("publication/authors/_%s_edit_form", pub.Type),
			views.NewContributorForm(r, p.render, savedPub.ID, author, muxRowDelta, formErrors),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Use the SavedAuthor since Librecat returns Author.FullName
	savedAuthor := &savedPub.Author[rowDelta]

	w.Header().Set("HX-Trigger", "ITUpdateItem")
	w.Header().Set("HX-Trigger-After-Swap", "ITUpdateItemAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITUpdateItemAfterSettle")

	p.render.HTML(w, http.StatusOK,
		"publication/authors/_default_row",
		views.NewContributorData(r, p.render, savedPub, savedAuthor, muxRowDelta),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationAuthors) ConfirmRemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	w.Header().Set("HX-Trigger", "ITConfirmRemoveFromPublication")
	w.Header().Set("HX-Trigger-After-Swap", "ITConfirmRemoveFromPublicationAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITConfirmRemoveFromPublicationAfterSettle")

	p.render.HTML(w, 200,
		"publication/_authors_modal_confirm_removal",
		struct {
			ID          string
			AuthorDelta string
		}{
			id,
			muxRowDelta,
		},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationAuthors) RemoveAuthor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub, err := p.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	authors := make([]models.PublicationContributor, len(pub.Author))
	copy(authors, pub.Author)

	authors = append(authors[:rowDelta], authors[rowDelta+1:]...)
	pub.Author = authors

	// TODO: error handling
	p.engine.UpdatePublication(pub)

	w.Header().Set("HX-Trigger", "ITRemoveItem")
	w.Header().Set("HX-Trigger-After-Swap", "ITRemoveItemAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITRemoveItemAfterSettle")

	// Empty content, denotes we deleted the record
	fmt.Fprintf(w, "")
}

func (p *PublicationAuthors) OrderAuthors(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	muxStart := mux.Vars(r)["start"]
	start, _ := strconv.Atoi(muxStart)

	muxEnd := mux.Vars(r)["end"]
	end, _ := strconv.Atoi(muxEnd)

	pub, err := p.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	author := &pub.Author[start]

	// Remove the author
	authors := make([]models.PublicationContributor, len(pub.Author))
	copy(authors, pub.Author)
	authors = append(authors[:start], authors[start+1:]...)
	pub.Author = authors

	// Re-insert the author at the new position
	placeholder := models.PublicationContributor{}
	authors = append(authors, placeholder)
	copy(authors[end+1:], authors[end:])
	authors[end] = *author

	// Save everything
	pub.Author = authors

	savedPub, err := p.engine.UpdatePublication(pub)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "ITOrderAuthors")
	w.Header().Set("HX-Trigger-After-Swap", "ITOrderAuthorsAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITOrderAuthorsAfterSettle")

	p.render.HTML(w, 200,
		"publication/authors/_default_table_body",
		views.NewContributorData(r, p.render, savedPub, nil, "0"),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
