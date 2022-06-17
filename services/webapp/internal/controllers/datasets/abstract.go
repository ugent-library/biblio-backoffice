package datasets

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/localize"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type BindAbstract struct {
	Position int    `path:"position"`
	Text     string `form:"text"`
	Lang     string `form:"lang"`
}

type YieldAbstracts struct {
	Ctx EditContext
}
type YieldAddAbstract struct {
	Ctx  EditContext
	Form *render.Form
}
type YieldEditAbstract struct {
	Ctx      EditContext
	Position int
	Form     *render.Form
}
type YieldDeleteAbstract struct {
	Ctx      EditContext
	Position int
}

func (c *Controller) AddAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	form := abstractForm(ctx, BindAbstract{Position: len(ctx.Dataset.Abstract)}, nil)

	render.Render(w, "dataset/add_abstract", YieldAddAbstract{
		Ctx:  ctx,
		Form: form,
	})
}

func (c *Controller) CreateAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	b := BindAbstract{Position: len(ctx.Dataset.Abstract)}
	if render.BadRequest(w, bind.RequestForm(r, &b)) {
		return
	}

	ctx.Dataset.Abstract = append(ctx.Dataset.Abstract, models.Text{Text: b.Text, Lang: b.Lang})

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		render.Render(w, "dataset/refresh_add_abstract", YieldAddAbstract{
			Ctx:  ctx,
			Form: abstractForm(ctx, b, validationErrs.(validation.Errors)),
		})
		return
	}

	err := c.store.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Render(w, "error_dialog", ctx.Locale.T("dataset.conflict_error"))
		return
	}

	if !render.InternalServerError(w, err) {
		render.Render(w, "dataset/refresh_abstracts", YieldAbstracts{
			Ctx: ctx,
		})
	}
}

func (c *Controller) EditAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	b := BindAbstract{}
	if render.BadRequest(w, bind.RequestPath(r, &b)) {
		return
	}

	a, err := ctx.Dataset.GetAbstract(b.Position)
	if render.BadRequest(w, err) {
		return
	}

	b.Lang = a.Lang
	b.Text = a.Text

	render.Render(w, "dataset/edit_abstract", YieldEditAbstract{
		Ctx:      ctx,
		Position: b.Position,
		Form:     abstractForm(ctx, b, nil),
	})
}

func (c *Controller) UpdateAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	b := BindAbstract{}
	if render.BadRequest(w, bind.Request(r, &b)) {
		return
	}

	a := models.Text{Text: b.Text, Lang: b.Lang}
	if render.BadRequest(w, ctx.Dataset.SetAbstract(b.Position, a)) {
		return
	}

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		form := abstractForm(ctx, b, validationErrs.(validation.Errors))

		render.Render(w, "dataset/refresh_edit_abstract", YieldEditAbstract{
			Ctx:      ctx,
			Position: b.Position,
			Form:     form,
		})
		return
	}

	err := c.store.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Render(w, "error_dialog", ctx.Locale.T("dataset.conflict_error"))
		return
	}

	if !render.InternalServerError(w, err) {
		render.Render(w, "dataset/refresh_abstracts", YieldAbstracts{
			Ctx: ctx,
		})
	}
}

func (c *Controller) ConfirmDeleteAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	var b BindAbstract
	if render.BadRequest(w, bind.RequestPath(r, &b)) {
		return
	}

	if _, err := ctx.Dataset.GetAbstract(b.Position); render.BadRequest(w, err) {
		return
	}

	render.Render(w, "dataset/confirm_delete_abstract", YieldDeleteAbstract{
		Ctx:      ctx,
		Position: b.Position,
	})
}

func (c *Controller) DeleteAbstract(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	var b BindAbstract
	if render.BadRequest(w, bind.Request(r, &b)) {
		return
	}

	if render.BadRequest(w, ctx.Dataset.RemoveAbstract(b.Position)) {
		return
	}

	err := c.store.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Render(w, "error_dialog", ctx.Locale.T("dataset.conflict_error"))
		return
	}

	if !render.InternalServerError(w, err) {
		render.Render(w, "dataset/refresh_abstracts", YieldAbstracts{
			Ctx: ctx,
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
