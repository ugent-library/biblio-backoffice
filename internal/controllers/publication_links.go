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
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type PublicationLinks struct {
	Context
}

func NewPublicationLinks(c Context) *PublicationLinks {
	return &PublicationLinks{c}
}

// Show the "Add link" modal
func (c *PublicationLinks) AddLink(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	link := &models.PublicationLink{}

	c.Render.HTML(w, http.StatusOK, "publication/links/_form", views.NewData(c.Render, r, struct {
		PublicationID string
		Link          *models.PublicationLink
		Form          *views.FormBuilder
		Vocabularies  map[string][]string
	}{
		id,
		link,
		views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
		c.Engine.Vocabularies(),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// Save a link to Librecat
func (c *PublicationLinks) CreateLink(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	link := &models.PublicationLink{}

	if err := forms.Decode(link, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	links := make([]models.PublicationLink, len(pub.Link))
	copy(links, pub.Link)

	links = append(links, *link)
	pub.Link = links

	savedPub, err := c.Engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK, "publication/links/_form", views.NewData(c.Render, r, struct {
			PublicationID string
			Link          *models.PublicationLink
			Form          *views.FormBuilder
			Vocabularies  map[string][]string
		}{
			savedPub.ID,
			link,
			views.NewFormBuilder(c.Render, locale.Get(r.Context()), formErrors),
			c.Engine.Vocabularies(),
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "PublicationCreateLink")

	c.Render.HTML(w, http.StatusOK, "publication/links/_table_body", views.NewData(c.Render, r, struct {
		Publication *models.Publication
	}{
		savedPub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// Show the "Edit link" modal
func (c *PublicationLinks) EditLink(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	link := &pub.Link[rowDelta]

	c.Render.HTML(w, http.StatusOK, "publication/links/_form_edit", views.NewData(c.Render, r, struct {
		PublicationID string
		Delta         string
		Link          *models.PublicationLink
		Form          *views.FormBuilder
		Vocabularies  map[string][]string
	}{
		pub.ID,
		muxRowDelta,
		link,
		views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
		c.Engine.Vocabularies(),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Save the updated link to Librecat
func (c *PublicationLinks) UpdateLink(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	link := &models.PublicationLink{}

	if err := forms.Decode(link, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	links := make([]models.PublicationLink, len(pub.Link))
	copy(links, pub.Link)

	links[rowDelta] = *link
	pub.Link = links

	log.Println(links)

	savedPub, err := c.Engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK, "publication/links/_form_edit", views.NewData(c.Render, r, struct {
			PublicationID string
			Delta         string
			Link          *models.PublicationLink
			Form          *views.FormBuilder
			Vocabularies  map[string][]string
		}{
			savedPub.ID,
			strconv.Itoa(rowDelta),
			link,
			views.NewFormBuilder(c.Render, locale.Get(r.Context()), formErrors),
			c.Engine.Vocabularies(),
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "PublicationUpdateLink")

	c.Render.HTML(w, http.StatusOK,
		"publication/links/_table_body",
		views.NewData(c.Render, r, struct {
			Publication *models.Publication
		}{
			savedPub,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Show the "Confirm remove" modal
func (c *PublicationLinks) ConfirmRemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	c.Render.HTML(w, http.StatusOK, "publication/links/_modal_confirm_removal", views.NewData(c.Render, r, struct {
		ID  string
		Key string
	}{
		id,
		muxRowDelta,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Remove a link from Librecat
func (c *PublicationLinks) RemoveLink(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	links := make([]models.PublicationLink, len(pub.Link))
	copy(links, pub.Link)

	links = append(links[:rowDelta], links[rowDelta+1:]...)
	pub.Link = links

	// TODO: error handling
	c.Engine.UpdatePublication(pub)

	w.Header().Set("HX-Trigger", "PublicationRemoveLink")

	// Empty content, denotes we deleted the record
	fmt.Fprintf(w, "")
}
