package dataset

import (
	"errors"
	"log"
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
	c.abstractView.Render(w, "add-abstract", YieldAbstract{
		Context: ctx,
		Form:    c.abstractForm(ctx),
	})
}

func (c *Controller) CreateAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	r.ParseForm()

	var bind BindAbstract
	if err := render.Bind(&bind, r.Form); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	d := ctx.Dataset.Clone()
	d.Abstract = append(d.Abstract, models.Text{Text: bind.Text, Lang: bind.Lang})

	err := c.store.UpdateDataset(d)

	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {
		c.abstractView.Render(w, "create-abstract-failed", YieldAbstract{
			Context: ctx,
			Form:    c.abstractForm(ctx),
		})
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	c.abstractView.Render(w, "create-abstract", YieldAbstract{
		Context: ctx,
	})
}

func (c *Controller) abstractForm(ctx Context) *render.Form {
	langOpts := []render.SelectOption{}
	for _, lang := range vocabularies.Map["language_codes"] {
		langOpts = append(langOpts, render.SelectOption{
			Value: lang,
			Label: ctx.Locale.LanguageName(lang),
		})
	}

	return &render.Form{
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
