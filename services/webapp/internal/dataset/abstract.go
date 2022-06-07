package dataset

import (
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/bind"
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

type YieldAbstract struct {
	Dataset  *models.Dataset
	Position int
	Form     *render.Form
}

func (c *Controller) AddAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	ctx.RenderYield(w, "dataset/add_abstract", YieldAbstract{
		Dataset: ctx.Dataset,
		Form:    abstractForm(ctx, BindAbstract{Position: len(ctx.Dataset.Abstract)}, nil),
	})
}

func (c *Controller) CreateAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	var b BindAbstract
	if render.BadRequest(w, bind.RequestForm(r, &b)) {
		return
	}

	d := ctx.Dataset
	d.SnapshotID = r.Header.Get("If-Match")
	d.Abstract = append(d.Abstract, models.Text{Text: b.Text, Lang: b.Lang})

	if validationErrs := d.Validate(); validationErrs != nil {
		ctx.RenderYield(w, "dataset/create_abstract_failed", YieldAbstract{
			Dataset: d,
			Form:    abstractForm(ctx, b, validationErrs),
		})
		return
	}

	err := c.store.UpdateDataset(d)
	// TODO handle conflict errors

	if !render.InternalServerError(w, err) {
		ctx.RenderYield(w, "dataset/create_abstract", YieldAbstract{
			Dataset: d,
		})
	}
}

func (c *Controller) EditAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	var b BindAbstract
	if render.BadRequest(w, bind.RequestPath(r, &b)) {
		return
	}

	if b.Position >= len(ctx.Dataset.Abstract) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	b.Lang = ctx.Dataset.Abstract[b.Position].Lang
	b.Text = ctx.Dataset.Abstract[b.Position].Text

	ctx.RenderYield(w, "dataset/edit_abstract", YieldAbstract{
		Dataset:  ctx.Dataset,
		Form:     abstractForm(ctx, b, nil),
		Position: b.Position,
	})
}

func (c *Controller) UpdateAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	var b BindAbstract
	if render.BadRequest(w, bind.Request(r, &b)) {
		return
	}

	if b.Position >= len(ctx.Dataset.Abstract) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	d := ctx.Dataset
	d.SnapshotID = r.Header.Get("If-Match")
	d.Abstract[b.Position].Lang = b.Lang
	d.Abstract[b.Position].Text = b.Text

	if validationErrs := d.Validate(); validationErrs != nil {
		ctx.RenderYield(w, "dataset/update_abstract_failed", YieldAbstract{
			Dataset:  d,
			Position: b.Position,
			Form:     abstractForm(ctx, b, validationErrs),
		})
		return
	}

	err := c.store.UpdateDataset(d)
	// TODO handle conflict errors

	if !render.InternalServerError(w, err) {
		ctx.RenderYield(w, "dataset/update_abstract", YieldAbstract{
			Dataset:  d,
			Position: b.Position,
		})
	}
}

func (c *Controller) ConfirmDeleteAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	var b BindAbstract
	if render.BadRequest(w, bind.RequestPath(r, &b)) {
		return
	}

	if b.Position >= len(ctx.Dataset.Abstract) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx.RenderYield(w, "dataset/confirm_delete_abstract", YieldAbstract{
		Dataset:  ctx.Dataset,
		Position: b.Position,
	})
}

func (c *Controller) DeleteAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	var b BindAbstract
	if render.BadRequest(w, bind.Request(r, &b)) {
		return
	}

	if b.Position >= len(ctx.Dataset.Abstract) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	d := ctx.Dataset
	d.SnapshotID = r.Header.Get("If-Match")
	d.Abstract = append(d.Abstract[:b.Position], d.Abstract[b.Position+1:]...)

	err := c.store.UpdateDataset(d)
	// TODO handle validation errors
	// TOOD handle conflict errors

	if !render.InternalServerError(w, err) {
		ctx.RenderYield(w, "dataset/delete_abstract", YieldAbstract{
			Dataset: d,
		})
	}
}

func abstractForm(ctx EditContext, b BindAbstract, errors validation.Errors) *render.Form {
	return &render.Form{
		Errors: localize.ValidationErrors(ctx.Locale, errors),
		Fields: []render.FormField{
			&render.TextArea{
				Name:        "text",
				Value:       b.Text,
				Label:       ctx.Locale.T("builder.abstract.text"),
				Cols:        12,
				Rows:        6,
				Placeholder: ctx.Locale.T("builder.abstract.text.placeholder"),
				Error:       localize.ValidationErrorAt(ctx.Locale, errors, fmt.Sprintf("/abstract/%d/text", b.Position)),
			},
			&render.Select{
				Name:    "lang",
				Value:   b.Lang,
				Label:   ctx.Locale.T("builder.abstract.lang"),
				Options: localize.LanguageSelectOptions(ctx.Locale),
				Cols:    12,
				Error:   localize.ValidationErrorAt(ctx.Locale, errors, fmt.Sprintf("/abstract/%d/lang", b.Position)),
			},
		},
	}
}
