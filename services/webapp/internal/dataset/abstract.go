package dataset

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/render"
)

type YieldAbstract struct {
	Context Context
	Form    *render.Form
}

type BindAbstract struct {
	Text string `form:"text"`
	Lang string `form:"lang"`
}

func (c *Controller) AddAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Render(w, "dataset/add_abstract", YieldAbstract{
		Context: ctx,
		Form:    c.abstractForm(ctx, nil),
	})
}

func (c *Controller) CreateAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	var bind BindAbstract
	if !render.MustBindForm(w, r, &bind) {
		return
	}

	d := ctx.Dataset.Clone()
	d.Abstract = append(d.Abstract, models.Text{Text: bind.Text, Lang: bind.Lang})
	err := c.store.UpdateDataset(d)

	if err := validation.As(err); err != nil {
		render.Render(w, "dataset/create_abstract_failed", YieldAbstract{
			Context: ctx,
			Form:    c.abstractForm(ctx, err),
		})
		return
	}

	render.MustRender(w, "dataset/create_abstract", YieldAbstract{
		Context: ctx.WithDataset(d),
	}, err)
}

func (c *Controller) abstractForm(ctx Context, errors validation.Errors) *render.Form {
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
}
