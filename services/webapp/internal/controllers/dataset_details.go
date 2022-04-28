package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/unrolled/render"
)

type DatasetDetails struct {
	Context
}

func NewDatasetDetails(c Context) *DatasetDetails {
	return &DatasetDetails{c}
}

func (c *DatasetDetails) Show(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	c.Render.HTML(w, http.StatusOK, "dataset/details/_show", c.ViewData(r, struct {
		Dataset *models.Dataset
		Show    *views.ShowBuilder
	}{
		dataset,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetDetails) Edit(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	c.Render.HTML(w, http.StatusOK, "dataset/details/_edit", c.ViewData(r, struct {
		Dataset      *models.Dataset
		Show         *views.ShowBuilder
		Form         *views.FormBuilder
		Vocabularies map[string][]string
	}{
		dataset,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
		c.Engine.Vocabularies(),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetDetails) AccessLevel(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := DecodeForm(dataset, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Clear embargo and embargoTo fields if access level is not embargo
	//   @todo Disabled per https://github.com/ugent-library/biblio-backend/issues/217
	//
	//   Another issue: the old JS also temporary stored the data in these fields if
	//   access level changed from embargo to something else. The data would be restored
	//   into the form fields again if embargo level is chosen again. This feature isn't
	//   implemented in this solution since state isn't kept across HTTP requests.
	//
	// if dataset.AccessLevel != "info:eu-repo/semantics/embargoedAccess" {
	// 	dataset.Embargo = ""
	// 	dataset.EmbargoTo = ""
	// }

	c.Render.HTML(w, http.StatusOK, "dataset/details/_edit", c.ViewData(r, struct {
		Dataset      *models.Dataset
		Show         *views.ShowBuilder
		Form         *views.FormBuilder
		Vocabularies map[string][]string
	}{
		dataset,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
		c.Engine.Vocabularies(),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetDetails) Update(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := DecodeForm(dataset, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savedDataset, err := c.Engine.UpdateDataset(dataset)

	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {
		c.Render.HTML(w, http.StatusOK, "dataset/details/_edit", c.ViewData(r, struct {
			Dataset      *models.Dataset
			Show         *views.ShowBuilder
			Form         *views.FormBuilder
			Vocabularies map[string][]string
		}{
			dataset,
			views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
			views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), validationErrors),
			c.Engine.Vocabularies(),
		},
			views.Flash{Type: "error", Message: "There are some problems with your input"},
		),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/details/_update", c.ViewData(r, struct {
		Dataset *models.Dataset
		Show    *views.ShowBuilder
	}{
		savedDataset,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
	},
		views.Flash{Type: "success", Message: "Details updated succesfully", DismissAfter: 5 * time.Second},
	),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
