package datasetediting

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/localize"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type BindAbstract struct {
	Position int    `path:"position"`
	Text     string `form:"text"`
	Lang     string `form:"lang"`
}

type YieldAbstracts struct {
	Context
}
type YieldAddAbstract struct {
	Context
	Form *form.Form
}
type YieldEditAbstract struct {
	Context
	Position int
	Form     *form.Form
}
type YieldDeleteAbstract struct {
	Context
	Position int
}

func (h *Handler) AddAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	form := abstractForm(ctx, BindAbstract{Position: len(ctx.Dataset.Abstract)}, nil)

	render.Render(w, "dataset/add_abstract", YieldAddAbstract{
		Context: ctx,
		Form:    form,
	})
}

func (h *Handler) CreateAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAbstract{Position: len(ctx.Dataset.Abstract)}
	if render.BadRequest(w, bind.RequestForm(r, &b)) {
		return
	}

	ctx.Dataset.Abstract = append(ctx.Dataset.Abstract, models.Text{Text: b.Text, Lang: b.Lang})

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		render.Render(w, "dataset/refresh_add_abstract", YieldAddAbstract{
			Context: ctx,
			Form:    abstractForm(ctx, b, validationErrs.(validation.Errors)),
		})
		return
	}

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Render(w, "error_dialog", ctx.T("dataset.conflict_error"))
		return
	}

	if !render.InternalServerError(w, err) {
		render.Render(w, "dataset/refresh_abstracts", YieldAbstracts{
			Context: ctx,
		})
	}
}

func (h *Handler) EditAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
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
		Context:  ctx,
		Position: b.Position,
		Form:     abstractForm(ctx, b, nil),
	})
}

func (h *Handler) UpdateAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
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
			Context:  ctx,
			Position: b.Position,
			Form:     form,
		})
		return
	}

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Render(w, "error_dialog", ctx.T("dataset.conflict_error"))
		return
	}

	if !render.InternalServerError(w, err) {
		render.Render(w, "dataset/refresh_abstracts", YieldAbstracts{
			Context: ctx,
		})
	}
}

func (h *Handler) ConfirmDeleteAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindAbstract
	if render.BadRequest(w, bind.RequestPath(r, &b)) {
		return
	}

	if _, err := ctx.Dataset.GetAbstract(b.Position); render.BadRequest(w, err) {
		return
	}

	render.Render(w, "dataset/confirm_delete_abstract", YieldDeleteAbstract{
		Context:  ctx,
		Position: b.Position,
	})
}

func (h *Handler) DeleteAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindAbstract
	if render.BadRequest(w, bind.Request(r, &b)) {
		return
	}

	if render.BadRequest(w, ctx.Dataset.RemoveAbstract(b.Position)) {
		return
	}

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Render(w, "error_dialog", ctx.T("dataset.conflict_error"))
		return
	}

	if !render.InternalServerError(w, err) {
		render.Render(w, "dataset/refresh_abstracts", YieldAbstracts{
			Context: ctx,
		})
	}
}

func abstractForm(ctx Context, b BindAbstract, errors validation.Errors) *form.Form {
	return &form.Form{
		Theme:  "default",
		Errors: localize.ValidationErrors(ctx.Locale, errors),
		Fields: []form.Field{
			&form.TextArea{
				Name:        "text",
				Value:       b.Text,
				Label:       ctx.T("builder.abstract.text"),
				Cols:        12,
				Rows:        6,
				Placeholder: ctx.T("builder.abstract.text.placeholder"),
				Error:       localize.ValidationErrorAt(ctx.Locale, errors, fmt.Sprintf("/abstract/%d/text", b.Position)),
			},
			&form.Select{
				Name:    "lang",
				Value:   b.Lang,
				Label:   ctx.T("builder.abstract.lang"),
				Options: localize.LanguageSelectOptions(ctx.Locale),
				Cols:    12,
				Error:   localize.ValidationErrorAt(ctx.Locale, errors, fmt.Sprintf("/abstract/%d/lang", b.Position)),
			},
		},
	}
}
