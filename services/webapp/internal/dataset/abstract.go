package dataset

import (
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/localize"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type BindAbstract struct {
	Position int    `path:"position"`
	Text     string `form:"text"`
	Lang     string `form:"lang"`
}

type YieldNewAbstract struct {
	Dataset *models.Dataset
	Form    *render.Form
}

type YieldAbstract struct {
	Dataset  *models.Dataset
	Position int
	Form     *render.Form
}

func (c *Controller) AddAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	ctx.RenderYield(w, "dataset/add_abstract", YieldNewAbstract{
		Dataset: ctx.Dataset,
		Form:    abstractForm(ctx, BindAbstract{Position: len(ctx.Dataset.Abstract)}, nil),
	})
}

func (c *Controller) CreateAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	var bind BindAbstract
	if render.BadRequest(w, render.BindForm(r, &bind)) {
		return
	}

	d := ctx.Dataset
	d.SnapshotID = r.Header.Get("If-Match")
	d.Abstract = append(d.Abstract, models.Text{Text: bind.Text, Lang: bind.Lang})
	err := c.store.UpdateDataset(d)

	if validationErrors := validation.From(err); validationErrors != nil {
		ctx.RenderYield(w, "dataset/create_abstract_failed", YieldNewAbstract{
			Dataset: d,
			Form:    abstractForm(ctx, bind, validationErrors),
		})
		return
	}

	if !render.InternalServerError(w, err) {
		ctx.RenderYield(w, "dataset/create_abstract", YieldNewAbstract{
			Dataset: d,
		})
	}
}

func (c *Controller) EditAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	var bind BindAbstract
	if render.BadRequest(w, render.BindPath(r, &bind)) {
		return
	}

	if bind.Position >= len(ctx.Dataset.Abstract) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	bind.Lang = ctx.Dataset.Abstract[bind.Position].Lang
	bind.Text = ctx.Dataset.Abstract[bind.Position].Text

	ctx.RenderYield(w, "dataset/edit_abstract", YieldAbstract{
		Dataset:  ctx.Dataset,
		Form:     abstractForm(ctx, bind, nil),
		Position: bind.Position,
	})
}

func (c *Controller) UpdateAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	var bind BindAbstract
	if render.BadRequest(w, render.Bind(r, &bind)) {
		return
	}

	if bind.Position >= len(ctx.Dataset.Abstract) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	d := ctx.Dataset
	d.SnapshotID = r.Header.Get("If-Match")
	d.Abstract[bind.Position].Lang = bind.Lang
	d.Abstract[bind.Position].Text = bind.Text
	err := c.store.UpdateDataset(d)

	if validationErrors := validation.From(err); validationErrors != nil {
		ctx.RenderYield(w, "dataset/update_abstract_failed", YieldAbstract{
			Dataset:  d,
			Position: bind.Position,
			Form:     abstractForm(ctx, bind, validationErrors),
		})
		return
	}

	if !render.InternalServerError(w, err) {
		ctx.RenderYield(w, "dataset/update_abstract", YieldAbstract{
			Dataset:  d,
			Position: bind.Position,
		})
	}
}

func abstractForm(ctx EditContext, bind BindAbstract, errors validation.Errors) *render.Form {
	return &render.Form{
		Errors: localize.ValidationErrors(ctx.Locale, errors),
		Fields: []render.FormField{
			&render.TextArea{
				Name:        "text",
				Value:       bind.Text,
				Label:       ctx.Locale.T("builder.abstract.text"),
				Cols:        12,
				Rows:        6,
				Placeholder: ctx.Locale.T("builder.abstract.text.placeholder"),
				Error:       localize.ValidationErrorAt(ctx.Locale, errors, fmt.Sprintf("/abstract/%d/text", bind.Position)),
			},
			&render.Select{
				Name:    "lang",
				Value:   bind.Lang,
				Label:   ctx.Locale.T("builder.abstract.lang"),
				Options: localize.LanguageSelectOptions(ctx.Locale),
				Cols:    12,
				Error:   localize.ValidationErrorAt(ctx.Locale, errors, fmt.Sprintf("/abstract/%d/lang", bind.Position)),
			},
		},
	}
}
