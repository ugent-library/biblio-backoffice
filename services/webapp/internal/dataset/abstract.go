package dataset

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
)

type BindAbstract struct {
	Text string `form:"text"`
	Lang string `form:"lang"`
}

type YieldAbstract struct {
	Dataset *models.Dataset
	Form    *render.Form
}

func (c *Controller) AddAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	ctx.RenderYield(w, "dataset/add_abstract", YieldAbstract{
		Dataset: ctx.Dataset,
		Form:    abstractForm(ctx, nil),
	})
}

func (c *Controller) CreateAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	var bind BindAbstract
	if !render.MustBindForm(w, r, &bind) {
		return
	}

	d := ctx.Dataset.Clone()
	d.SnapshotID = r.Header.Get("If-Match")
	d.Abstract = append(d.Abstract, models.Text{Text: bind.Text, Lang: bind.Lang})
	err := c.store.UpdateDataset(d)

	if validationErrors := validation.From(err); validationErrors != nil {
		ctx.RenderYield(w, "dataset/create_abstract_failed", YieldAbstract{
			Dataset: d,
			Form:    abstractForm(ctx, validationErrors),
		})
		return
	}

	if render.Must(w, err) {
		ctx.RenderYield(w, "dataset/create_abstract", YieldAbstract{
			Dataset: d,
		})
	}
}

func abstractForm(ctx EditContext, errors validation.Errors) *render.Form {
	formErrors := make([]string, len(errors))
	for i, err := range errors {
		formErrors[i] = err.Error()
	}

	langOpts := []render.SelectOption{}
	for _, lang := range vocabularies.Map["language_codes"] {
		langOpts = append(langOpts, render.SelectOption{
			Value: lang,
			Label: ctx.Locale.LanguageName(lang),
		})
	}

	return &render.Form{
		Errors: formErrors,
		Fields: []render.FormField{
			&render.TextArea{
				Name:        "text",
				Label:       ctx.Locale.T("builder.abstract.text"),
				Cols:        12,
				Rows:        6,
				Placeholder: ctx.Locale.T("builder.abstract.text.placeholder"),
			},
			&render.Select{
				Name:    "lang",
				Label:   ctx.Locale.T("builder.abstract.lang"),
				Options: langOpts,
				Cols:    12,
			},
		},
	}
}
