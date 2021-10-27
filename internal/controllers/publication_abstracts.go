package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type PublicationAbstracts struct {
	Context
}

func NewPublicationAbstracts(c Context) *PublicationAbstracts {
	return &PublicationAbstracts{c}
}

// Show the "Add abstract" modal
func (c *PublicationAbstracts) AddAbstract(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	abstract := &models.Text{}

	w.Header().Set("HX-Trigger", "PublicationAddAbstract")
	w.Header().Set("HX-Trigger-After-Swap", "PublicationAddAbstractAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "PublicationAddAbstractAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		"publication/abstracts/_form",
		views.NewData(c.Render, r, struct {
			PublicationID string
			Abstract      *models.Text
			Form          *views.FormBuilder
			Vocabularies  map[string][]string
		}{
			id,
			abstract,
			views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
			c.Engine.Vocabularies(),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// Save an abstract to Librecat
func (c *PublicationAbstracts) CreateAbstract(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := c.Engine.GetPublication(id)
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

	abstract := &models.Text{}

	if err := forms.Decode(abstract, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	abstracts := make([]models.Text, len(pub.Abstract))
	copy(abstracts, pub.Abstract)

	abstracts = append(abstracts, *abstract)
	pub.Abstract = abstracts

	savedPub, err := c.Engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, 200,
			"publication/abstracts/_form",
			views.NewData(c.Render, r, struct {
				PublicationID string
				Abstract      *models.Text
				Form          *views.FormBuilder
				Vocabularies  map[string][]string
			}{
				savedPub.ID,
				abstract,
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

	w.Header().Set("HX-Trigger", "PublicationCreateAbstract")
	w.Header().Set("HX-Trigger-After-Swap", "PublicationCreateAbstractAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "PublicationCreateAbstractAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		"publication/abstracts/_table_body",
		views.NewData(c.Render, r, struct {
			Publication *models.Publication
		}{
			savedPub,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// Show the "Edit abstract" modal
func (c *PublicationAbstracts) EditAbstract(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub, err := c.Engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	abstract := &pub.Abstract[rowDelta]

	w.Header().Set("HX-Trigger", "PublicationAddAbstract")
	w.Header().Set("HX-Trigger-After-Swap", "PublicationAddAbstractAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "PublicationAddAbstractAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		"publication/abstracts/_form_edit",
		views.NewData(c.Render, r, struct {
			PublicationID string
			Delta         string
			Abstract      *models.Text
			Form          *views.FormBuilder
			Vocabularies  map[string][]string
		}{
			id,
			muxRowDelta,
			abstract,
			views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
			c.Engine.Vocabularies(),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Save the updated abstract to Librecat
func (c *PublicationAbstracts) UpdateAbstract(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub, err := c.Engine.GetPublication(id)
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

	abstract := &models.Text{}

	if err := forms.Decode(abstract, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	abstracts := make([]models.Text, len(pub.Abstract))
	copy(abstracts, pub.Abstract)

	abstracts[rowDelta] = *abstract
	pub.Abstract = abstracts

	savedPub, err := c.Engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, 200,
			"publication/abstracts/_form_edit",
			views.NewData(c.Render, r, struct {
				PublicationID string
				Delta         string
				Abstract      *models.Text
				Form          *views.FormBuilder
				Vocabularies  map[string][]string
			}{
				savedPub.ID,
				strconv.Itoa(rowDelta),
				abstract,
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

	w.Header().Set("HX-Trigger", "PublicationUpdateAbstract")
	w.Header().Set("HX-Trigger-After-Swap", "PublicationUpdateAbstractAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "PublicationUpdateAbstractAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		"publication/abstracts/_table_body",
		views.NewData(c.Render, r, struct {
			Publication *models.Publication
		}{
			savedPub,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Show the "Confirm remove" modal
func (c *PublicationAbstracts) ConfirmRemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	w.Header().Set("HX-Trigger", "PublicationConfirmRemove")
	w.Header().Set("HX-Trigger-After-Swap", "PublicationConfirmRemoveAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "PublicationConfirmRemoveAfterSettle")

	c.Render.HTML(w, 200,
		"publication/_abstracts_modal_confirm_removal",
		struct {
			ID  string
			Key string
		}{
			id,
			muxRowDelta,
		},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Remove an abstract from Librecat
func (c *PublicationAbstracts) RemoveAbstract(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub, err := c.Engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	abstracts := make([]models.Text, len(pub.Abstract))
	copy(abstracts, pub.Abstract)

	abstracts = append(abstracts[:rowDelta], abstracts[rowDelta+1:]...)
	pub.Abstract = abstracts

	// TODO: error handling
	c.Engine.UpdatePublication(pub)

	w.Header().Set("HX-Trigger", "PublicationRemoveAbstract")
	w.Header().Set("HX-Trigger-After-Swap", "PublicationRemoveAbstractAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "PublicationRemoveAbstractAfterSettle")

	// Empty content, denotes we deleted the record
	fmt.Fprintf(w, "")
}
