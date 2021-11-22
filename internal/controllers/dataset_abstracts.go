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

type DatasetAbstracts struct {
	Context
}

func NewDatasetAbstracts(c Context) *DatasetAbstracts {
	return &DatasetAbstracts{c}
}

// Show the "Add abstract" modal
func (c *DatasetAbstracts) AddAbstract(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	abstract := &models.Text{}

	w.Header().Set("HX-Trigger", "DatasetAddAbstract")
	w.Header().Set("HX-Trigger-After-Swap", "DatasetAddAbstractAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "DatasetAddAbstractAfterSettle")

	c.Render.HTML(w, http.StatusOK, "dataset/abstracts/_form", views.NewData(c.Render, r, struct {
		DatasetID    string
		Abstract     *models.Text
		Form         *views.FormBuilder
		Vocabularies map[string][]string
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
func (c *DatasetAbstracts) CreateAbstract(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	abstract := &models.Text{}

	if err := forms.Decode(abstract, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	abstracts := make([]models.Text, len(dataset.Abstract))
	copy(abstracts, dataset.Abstract)

	abstracts = append(abstracts, *abstract)
	dataset.Abstract = abstracts

	savedDataset, err := c.Engine.UpdateDataset(dataset)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK, "dataset/abstracts/_form", views.NewData(c.Render, r, struct {
			DatasetID    string
			Abstract     *models.Text
			Form         *views.FormBuilder
			Vocabularies map[string][]string
		}{
			savedDataset.ID,
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

	w.Header().Set("HX-Trigger", "DatasetCreateAbstract")
	w.Header().Set("HX-Trigger-After-Swap", "DatasetCreateAbstractAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "DatasetCreateAbstractAfterSettle")

	c.Render.HTML(w, http.StatusOK,
		"dataset/abstracts/_table_body",
		views.NewData(c.Render, r, struct {
			Dataset *models.Dataset
		}{
			savedDataset,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// Show the "Edit abstract" modal
func (c *DatasetAbstracts) EditAbstract(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	dataset := context.GetDataset(r.Context())

	abstract := &dataset.Abstract[rowDelta]

	w.Header().Set("HX-Trigger", "DatasetAddAbstract")
	w.Header().Set("HX-Trigger-After-Swap", "DatasetAddAbstractAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "DatasetAddAbstractAfterSettle")

	c.Render.HTML(w, http.StatusOK, "dataset/abstracts/_form_edit", views.NewData(c.Render, r, struct {
		DatasetID    string
		Delta        string
		Abstract     *models.Text
		Form         *views.FormBuilder
		Vocabularies map[string][]string
	}{
		dataset.ID,
		muxRowDelta,
		abstract,
		views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
		c.Engine.Vocabularies(),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Save the updated abstract to Librecat
func (c *DatasetAbstracts) UpdateAbstract(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	dataset := context.GetDataset(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	abstract := &models.Text{}

	if err := forms.Decode(abstract, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	abstracts := make([]models.Text, len(dataset.Abstract))
	copy(abstracts, dataset.Abstract)

	abstracts[rowDelta] = *abstract
	dataset.Abstract = abstracts

	savedDataset, err := c.Engine.UpdateDataset(dataset)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK,
			"dataset/abstracts/_form_edit",
			views.NewData(c.Render, r, struct {
				DatasetID    string
				Delta        string
				Abstract     *models.Text
				Form         *views.FormBuilder
				Vocabularies map[string][]string
			}{
				savedDataset.ID,
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

	w.Header().Set("HX-Trigger", "DatasetUpdateAbstract")
	w.Header().Set("HX-Trigger-After-Swap", "DatasetUpdateAbstractAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "DatasetUpdateAbstractAfterSettle")

	c.Render.HTML(w, http.StatusOK, "dataset/abstracts/_table_body", views.NewData(c.Render, r, struct {
		Dataset *models.Dataset
	}{
		savedDataset,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Show the "Confirm remove" modal
func (c *DatasetAbstracts) ConfirmRemoveFromDataset(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	w.Header().Set("HX-Trigger", "DatasetConfirmRemove")
	w.Header().Set("HX-Trigger-After-Swap", "DatasetConfirmRemoveAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "DatasetConfirmRemoveAfterSettle")

	c.Render.HTML(w, http.StatusOK, "dataset/abstracts/_modal_confirm_removal", views.NewData(c.Render, r, struct {
		ID  string
		Key string
	}{
		id,
		muxRowDelta,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Remove an abstract from Librecat
func (c *DatasetAbstracts) RemoveAbstract(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetDataset(r.Context())

	abstracts := make([]models.Text, len(pub.Abstract))
	copy(abstracts, pub.Abstract)

	abstracts = append(abstracts[:rowDelta], abstracts[rowDelta+1:]...)
	pub.Abstract = abstracts

	// TODO: error handling
	c.Engine.UpdateDataset(pub)

	w.Header().Set("HX-Trigger", "DatasetRemoveAbstract")
	w.Header().Set("HX-Trigger-After-Swap", "DatasetRemoveAbstractAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "DatasetRemoveAbstractAfterSettle")

	// Empty content, denotes we deleted the record
	fmt.Fprintf(w, "")
}
