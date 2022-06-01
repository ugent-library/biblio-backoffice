package dataset

import (
	"fmt"
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
		Form:    abstractForm(ctx, BindAbstract{}, len(ctx.Dataset.Abstract), nil),
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
			Form:    abstractForm(ctx, bind, len(d.Abstract)-1, validationErrors),
		})
		return
	}

	if render.Must(w, err) {
		ctx.RenderYield(w, "dataset/create_abstract", YieldAbstract{
			Dataset: d,
		})
	}
}

func abstractForm(ctx EditContext, bind BindAbstract, index int, errors validation.Errors) *render.Form {
	formErrors := make([]string, len(errors))
	for i, e := range errors {
		formErrors[i] = ctx.Locale.TS("validation", e.Code)
	}

	var (
		textErr string
		langErr string
	)
	if e := errors.At(fmt.Sprintf("/abstract/%d/text", index)); e != nil {
		textErr = ctx.Locale.TS("validation", e.Code)
	}
	if e := errors.At(fmt.Sprintf("/abstract/%d/lang", index)); e != nil {
		langErr = ctx.Locale.TS("validation", e.Code)
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
				Value:       bind.Text,
				Label:       ctx.Locale.T("builder.abstract.text"),
				Cols:        12,
				Rows:        6,
				Placeholder: ctx.Locale.T("builder.abstract.text.placeholder"),
				Error:       textErr,
			},
			&render.Select{
				Name:    "lang",
				Value:   bind.Lang,
				Label:   ctx.Locale.T("builder.abstract.lang"),
				Options: langOpts,
				Cols:    12,
				Error:   langErr,
			},
		},
	}
}
