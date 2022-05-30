package dataset

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/vocabularies"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/render"
)

type AddAbstractData struct {
	ID         string
	SnapshotID string
	Form       *render.Form
}

func (c *Controller) AddAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	langOpts := []render.SelectOption{}
	for _, lang := range vocabularies.Map["language_codes"] {
		langOpts = append(langOpts, render.SelectOption{
			Value: lang,
			Label: ctx.Locale.LanguageName(lang),
		})
	}

	form := &render.Form{
		Fields: []render.FormField{
			&render.TextArea{
				Name:        "text",
				Label:       ctx.Locale.T("builder.abstract", "text"),
				Cols:        12,
				Rows:        6,
				Placeholder: ctx.Locale.T("builder.abstract", "text.placeholder"),
			},
			&render.Select{
				Name:    "lang",
				Label:   ctx.Locale.T("builder.abstract", "lang"),
				Options: langOpts,
				Cols:    12,
			},
		},
	}

	c.addAbstractPartial.Render(w, AddAbstractData{
		ID:         ctx.Dataset.ID,
		SnapshotID: ctx.Dataset.SnapshotID,
		Form:       form,
	})
}
