package controllers

import (
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backend/internal/forms"
	"github.com/ugent-library/biblio-backend/internal/jsonapi"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/unrolled/render"
)

type PublicationAdditionalInfo struct {
	Context
}

func NewPublicationAdditionalInfo(c Context) *PublicationAdditionalInfo {
	return &PublicationAdditionalInfo{c}
}

func (c *PublicationAdditionalInfo) Show(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	c.Render.HTML(w, http.StatusOK,
		"publication/additional_info/_show",
		c.ViewData(r, struct {
			Publication  *models.Publication
			Show         *views.ShowBuilder
			Vocabularies map[string][]string
		}{
			pub,
			views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
			c.Engine.Vocabularies(),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationAdditionalInfo) Edit(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	c.Render.HTML(w, http.StatusOK,
		"publication/additional_info/_edit",
		c.ViewData(r, struct {
			Publication  *models.Publication
			Form         *views.FormBuilder
			Vocabularies map[string][]string
		}{
			pub,
			views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
			c.Engine.Vocabularies(),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationAdditionalInfo) Update(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := forms.Decode(pub, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savedPub, err := c.Engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK,
			"publication/additional_info/_edit",
			c.ViewData(r, struct {
				Publication  *models.Publication
				Form         *views.FormBuilder
				Vocabularies map[string][]string
			}{
				pub,
				views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), formErrors),
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

	c.Render.HTML(w, http.StatusOK,
		"publication/additional_info/_update",
		c.ViewData(r, struct {
			Publication  *models.Publication
			Show         *views.ShowBuilder
			Vocabularies map[string][]string
		}{
			savedPub,
			views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
			c.Engine.Vocabularies(),
		},
			views.Flash{Type: "success", Message: "Additional info updated succesfully", DismissAfter: 5 * time.Second},
		),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
