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
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type PublicationAbstracts struct {
	engine *engine.Engine
	render *render.Render
}

func NewPublicationAbstracts(e *engine.Engine, r *render.Render) *PublicationAbstracts {
	return &PublicationAbstracts{
		engine: e,
		render: r,
	}
}

// Show the "Add abstract" modal
func (p *PublicationAbstracts) AddAbstract(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	abstract := &models.Text{}

	w.Header().Set("HX-Trigger", "PublicationAddAbstract")
	w.Header().Set("HX-Trigger-After-Swap", "PublicationAddAbstractAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "PublicationAddAbstractAfterSettle")

	p.render.HTML(w, http.StatusOK,
		"publication/abstracts/_form",
		struct {
			views.Data
			PublicationID string
			Abstract      *models.Text
			Form          *views.FormBuilder
			Vocabularies  map[string][]string
		}{
			views.NewData(p.render, r),
			id,
			abstract,
			views.NewFormBuilder(p.render, locale.Get(r.Context()), nil),
			p.engine.Vocabularies(),
		},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// Save an abstract to Librecat
func (p *PublicationAbstracts) CreateAbstract(w http.ResponseWriter, r *http.Request) {
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

	abstract := &models.Text{}

	if err := forms.Decode(abstract, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	abstracts := make([]models.Text, len(pub.Abstract))
	copy(abstracts, pub.Abstract)

	abstracts = append(abstracts, *abstract)
	pub.Abstract = abstracts

	savedPub, err := p.engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		p.render.HTML(w, 200,
			"publication/abstracts/_form",
			struct {
				views.Data
				PublicationID string
				Abstract      *models.Text
				Form          *views.FormBuilder
				Vocabularies  map[string][]string
			}{
				views.NewData(p.render, r),
				savedPub.ID,
				abstract,
				views.NewFormBuilder(p.render, locale.Get(r.Context()), formErrors),
				p.engine.Vocabularies(),
			},
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

	p.render.HTML(w, http.StatusOK,
		"publication/abstracts/_table_body",
		struct {
			views.Data
			Publication *models.Publication
		}{
			views.NewData(p.render, r),
			savedPub,
		},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// Show the "Edit abstract" modal
func (p *PublicationAbstracts) EditAbstract(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub, err := p.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	abstract := &pub.Abstract[rowDelta]

	w.Header().Set("HX-Trigger", "PublicationAddAbstract")
	w.Header().Set("HX-Trigger-After-Swap", "PublicationAddAbstractAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "PublicationAddAbstractAfterSettle")

	p.render.HTML(w, http.StatusOK,
		"publication/abstracts/_form_edit",
		struct {
			views.Data
			PublicationID string
			Delta         string
			Abstract      *models.Text
			Form          *views.FormBuilder
			Vocabularies  map[string][]string
		}{
			views.NewData(p.render, r),
			id,
			muxRowDelta,
			abstract,
			views.NewFormBuilder(p.render, locale.Get(r.Context()), nil),
			p.engine.Vocabularies(),
		},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Save the updated abstract to Librecat
func (p *PublicationAbstracts) UpdateAbstract(w http.ResponseWriter, r *http.Request) {
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

	abstract := &models.Text{}

	if err := forms.Decode(abstract, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	abstracts := make([]models.Text, len(pub.Abstract))
	copy(abstracts, pub.Abstract)

	abstracts[rowDelta] = *abstract
	pub.Abstract = abstracts

	savedPub, err := p.engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		p.render.HTML(w, 200,
			"publication/abstracts/_form_edit",
			struct {
				views.Data
				PublicationID string
				Delta         string
				Abstract      *models.Text
				Form          *views.FormBuilder
				Vocabularies  map[string][]string
			}{
				views.NewData(p.render, r),
				savedPub.ID,
				strconv.Itoa(rowDelta),
				abstract,
				views.NewFormBuilder(p.render, locale.Get(r.Context()), formErrors),
				p.engine.Vocabularies(),
			},
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

	p.render.HTML(w, http.StatusOK,
		"publication/abstracts/_table_body",
		struct {
			views.Data
			Publication *models.Publication
		}{
			views.NewData(p.render, r),
			savedPub,
		},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Show the "Confirm remove" modal
func (p *PublicationAbstracts) ConfirmRemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	w.Header().Set("HX-Trigger", "PublicationConfirmRemove")
	w.Header().Set("HX-Trigger-After-Swap", "PublicationConfirmRemoveAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "PublicationConfirmRemoveAfterSettle")

	p.render.HTML(w, 200,
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
func (p *PublicationAbstracts) RemoveAbstract(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub, err := p.engine.GetPublication(id)
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
	p.engine.UpdatePublication(pub)

	w.Header().Set("HX-Trigger", "PublicationRemoveAbstract")
	w.Header().Set("HX-Trigger-After-Swap", "PublicationRemoveAbstractAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "PublicationRemoveAbstractAfterSettle")

	// Empty content, denotes we deleted the record
	fmt.Fprintf(w, "")
}
