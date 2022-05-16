package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/unrolled/render"
)

type PublicationLaySummaries struct {
	Context
}

func NewPublicationLaySummaries(c Context) *PublicationLaySummaries {
	return &PublicationLaySummaries{c}
}

func (c *PublicationLaySummaries) Add(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	c.Render.HTML(w, http.StatusOK, "publication/lay_summaries/_form", c.ViewData(r, struct {
		PublicationID string
		LaySummary    *models.Text
		Form          *views.FormBuilder
	}{
		PublicationID: id,
		LaySummary:    &models.Text{},
		Form:          views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// Save a lay summary to Librecat
func (c *PublicationLaySummaries) Create(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lay_summary := &models.Text{}

	if err := DecodeForm(lay_summary, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lay_summaries := make([]models.Text, len(pub.LaySummary))
	copy(lay_summaries, pub.LaySummary)

	lay_summaries = append(lay_summaries, *lay_summary)
	pub.LaySummary = lay_summaries

	savedPub := pub.Clone()
	err = c.Engine.Store.StorePublication(savedPub)

	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {
		c.Render.HTML(w, http.StatusOK, "publication/lay_summaries/_form", c.ViewData(r, struct {
			PublicationID string
			LaySummary    *models.Text
			Form          *views.FormBuilder
		}{
			PublicationID: savedPub.ID,
			LaySummary:    lay_summary,
			Form:          views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), validationErrors),
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/*
		target is modal, so send back empty answer (that closes modal),
		and update list out of band
		TODO: if you send back multiple hx-swab-oob (e.g. by sending a Flash message)
			  then the second hx-swap-oob attribute is ignored (where you're list is),
			  and so your list is inserted in the .. modal.
	*/
	c.Render.HTML(w, http.StatusOK,
		"publication/lay_summaries/_created",
		c.ViewData(r, struct {
			Publication *models.Publication
		}{
			Publication: savedPub,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// Show the "Edit lay summary" modal
func (c *PublicationLaySummaries) Edit(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	lay_summary := &pub.LaySummary[rowDelta]

	c.Render.HTML(w, http.StatusOK, "publication/lay_summaries/_form_edit", c.ViewData(r, struct {
		PublicationID string
		Delta         string
		LaySummary    *models.Text
		Form          *views.FormBuilder
	}{
		PublicationID: pub.ID,
		Delta:         muxRowDelta,
		LaySummary:    lay_summary,
		Form:          views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Save the updated lay summary to Librecat
func (c *PublicationLaySummaries) Update(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lay_summary := &models.Text{}

	if err := DecodeForm(lay_summary, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lay_summaries := make([]models.Text, len(pub.LaySummary))
	copy(lay_summaries, pub.LaySummary)

	lay_summaries[rowDelta] = *lay_summary
	pub.LaySummary = lay_summaries

	savedPub := pub.Clone()
	err = c.Engine.Store.StorePublication(savedPub)

	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {
		c.Render.HTML(w, http.StatusOK,
			"publication/lay_summaries/_form_edit",
			c.ViewData(r, struct {
				PublicationID string
				Delta         string
				LaySummary    *models.Text
				Form          *views.FormBuilder
			}{
				PublicationID: savedPub.ID,
				Delta:         strconv.Itoa(rowDelta),
				LaySummary:    lay_summary,
				Form:          views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), validationErrors),
			}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/lay_summaries/_updated", c.ViewData(r, struct {
		Publication *models.Publication
	}{
		Publication: savedPub,
	},
	),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Show the "Confirm remove" modal
func (c *PublicationLaySummaries) ConfirmRemove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	c.Render.HTML(w, http.StatusOK, "publication/lay_summaries/_modal_confirm_removal", c.ViewData(r, struct {
		ID  string
		Key string
	}{
		ID:  id,
		Key: muxRowDelta,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Remove a lay summary from Librecat
func (c *PublicationLaySummaries) Remove(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	lay_summaries := make([]models.Text, len(pub.LaySummary))
	copy(lay_summaries, pub.LaySummary)

	lay_summaries = append(lay_summaries[:rowDelta], lay_summaries[rowDelta+1:]...)
	pub.LaySummary = lay_summaries

	// TODO: error handling
	savedPub := pub.Clone()
	c.Engine.Store.StorePublication(savedPub)

	c.Render.HTML(w, http.StatusOK, "publication/lay_summaries/_deleted", c.ViewData(r, struct {
		Publication *models.Publication
	}{
		Publication: savedPub,
	},
	),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
