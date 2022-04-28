package controllers

import (
	"errors"
	"fmt"
	"log"
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

type PublicationLinks struct {
	Context
}

func NewPublicationLinks(c Context) *PublicationLinks {
	return &PublicationLinks{c}
}

// Show the "Add link" modal
func (c *PublicationLinks) Add(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	link := &models.PublicationLink{}

	c.Render.HTML(w, http.StatusOK, "publication/links/_form", c.ViewData(r, struct {
		PublicationID string
		Link          *models.PublicationLink
		Form          *views.FormBuilder
	}{
		id,
		link,
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// Save a link to Librecat
func (c *PublicationLinks) Create(w http.ResponseWriter, r *http.Request) {
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

	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {
		c.Render.HTML(w, http.StatusOK, "publication/links/_form", c.ViewData(r, struct {
			PublicationID string
			Link          *models.PublicationLink
			Form          *views.FormBuilder
		}{
			savedPub.ID,
			link,
			views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), validationErrors),
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "PublicationCreateLink")

	c.Render.HTML(w, http.StatusOK, "publication/links/_table_body", c.ViewData(r, struct {
		Publication *models.Publication
	}{
		savedPub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// Show the "Edit link" modal
func (c *PublicationLinks) Edit(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	link := &pub.Link[rowDelta]

	c.Render.HTML(w, http.StatusOK, "publication/links/_form_edit", c.ViewData(r, struct {
		PublicationID string
		Delta         string
		Link          *models.PublicationLink
		Form          *views.FormBuilder
	}{
		pub.ID,
		muxRowDelta,
		link,
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Save the updated link to Librecat
func (c *PublicationLinks) Update(w http.ResponseWriter, r *http.Request) {
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

	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {
		c.Render.HTML(w, http.StatusOK, "publication/links/_form_edit", c.ViewData(r, struct {
			PublicationID string
			Delta         string
			Link          *models.PublicationLink
			Form          *views.FormBuilder
		}{
			savedPub.ID,
			strconv.Itoa(rowDelta),
			link,
			views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), validationErrors),
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
		c.ViewData(r, struct {
			Publication *models.Publication
		}{
			savedPub,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Show the "Confirm remove" modal
func (c *PublicationLinks) ConfirmRemove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	c.Render.HTML(w, http.StatusOK, "publication/links/_modal_confirm_removal", c.ViewData(r, struct {
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
func (c *PublicationLinks) Remove(w http.ResponseWriter, r *http.Request) {
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
