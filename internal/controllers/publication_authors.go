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

func (p *PublicationAuthors) AddAuthorToTable(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxAuthorDelta := mux.Vars(r)["author_delta"]

	pub, err := p.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Convert string to int
	authorDelta, _ := strconv.Atoi(muxAuthorDelta)
	// Create an empty author
	author := models.PublicationContributor{}

	// Deep copy authors from pub.Author
	authors := make([]models.PublicationContributor, len(pub.Author))
	copy(authors, pub.Author)

	// Intersperse a new item in an slice
	// 1. Append an empty author to the end fo the slice, creating room in the slice.
	authors = append(authors, author)
	// 2. Shift all authors one element to the right, the last item will be overwritten.
	//    Element at autohrDelta+1 will be duplicated, with the element at authorDelta
	copy(authors[authorDelta+1:], authors[authorDelta:])
	// 3. Overwrite the element at authorDelta+1 with the actual element we want show
	//    at that position
	authors[authorDelta+1] = author

	pub.Author = authors

	p.render.HTML(w, 200,
		fmt.Sprintf("publication/authors/_%s_form", pub.Type),
		views.NewContributorForm(r, p.render, pub, authorDelta+1, nil),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationAuthors) CancelAddAuthorToTable(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := p.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	p.render.HTML(w, 200,
		fmt.Sprintf("publication/authors/_%s", pub.Type),
		views.NewContributorData(r, p.render, pub),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationAuthors) SaveAuthorToPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

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
	formAuthorDelta := r.Form.Get("delta")

	if err := forms.Decode(author, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(r.Form)
	log.Println(author)

	authorDelta, _ := strconv.Atoi(formAuthorDelta)
	placeholder := models.PublicationContributor{}

	authors := make([]models.PublicationContributor, len(pub.Author))
	copy(authors, pub.Author)

	authors = append(authors, placeholder)
	copy(authors[authorDelta+1:], authors[authorDelta:])
	authors[authorDelta] = *author
	pub.Author = authors

	savedPub, err := p.engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		p.render.HTML(w, 200,
			fmt.Sprintf("publication/authors/_%s_form", pub.Type),
			views.NewContributorForm(r, p.render, pub, authorDelta, formErrors),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.render.HTML(w, 200,
		fmt.Sprintf("publication/authors/_%s_form_submit", savedPub.Type),
		views.NewContributorData(r, p.render, savedPub),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationAuthors) ConfirmRemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	authorDelta := mux.Vars(r)["author_delta"]

	p.render.HTML(w, 200,
		"publication/_authors_modal_confirm_removal",
		struct {
			ID          string
			AuthorDelta string
		}{
			id,
			authorDelta,
		},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (p *PublicationAuthors) RemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxAuthorDelta := mux.Vars(r)["author_delta"]
	authorDelta, _ := strconv.Atoi(muxAuthorDelta)
	pub, err := p.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	authors := make([]models.PublicationContributor, len(pub.Author))
	copy(authors, pub.Author)

	authors = append(authors[:authorDelta], authors[authorDelta+1:]...)
	pub.Author = authors

	// TODO: error handling
	savedPub, _ := p.engine.UpdatePublication(pub)

	p.render.HTML(w, 200,
		fmt.Sprintf("publication/authors/_%s_form_submit", pub.Type),
		views.NewPublicationData(r, p.render, savedPub),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// func (p *PublicationAuthors) ActiveSearch(w http.ResponseWriter, r *http.Request) {
// 	id := mux.Vars(r)["id"]

// 	pub, err := p.engine.GetPublication(id)
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	err = r.ParseForm()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Get 20 results from the search query
// 	query := r.Form["search"]
// 	hits, _ := p.engine.SuggestDepartments(query[0])

// 	p.render.HTML(w, 200,
// 		"publication/_departments_modal_hits",
// 		struct {
// 			Publication *models.Publication
// 			Hits        []models.Completion
// 		}{
// 			pub,
// 			hits,
// 		},
// 		render.HTMLOptions{Layout: "layouts/htmx"},
// 	)
// }
