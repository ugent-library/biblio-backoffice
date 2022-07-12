package datasetediting

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/bind"
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

type BindDeleteAbstract struct {
	Position int `path:"position"`
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

	render.Layout(w, "show_modal", "dataset/add_abstract", YieldAddAbstract{
		Context: ctx,
		Form:    form,
	})
}

func (h *Handler) CreateAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAbstract{Position: len(ctx.Dataset.Abstract)}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	ctx.Dataset.Abstract = append(ctx.Dataset.Abstract, models.Text{Text: b.Text, Lang: b.Lang})

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "dataset/add_abstract", YieldAddAbstract{
			Context: ctx,
			Form:    abstractForm(ctx, b, validationErrs.(validation.Errors)),
		})
		return
	}

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("dataset.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_abstracts", YieldAbstracts{
		Context: ctx,
	})
}

func (h *Handler) EditAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAbstract{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	a, err := ctx.Dataset.GetAbstract(b.Position)
	if err != nil {
		render.BadRequest(w, r, err)
		return
	}

	b.Lang = a.Lang
	b.Text = a.Text

	render.Layout(w, "show_modal", "dataset/edit_abstract", YieldEditAbstract{
		Context:  ctx,
		Position: b.Position,
		Form:     abstractForm(ctx, b, nil),
	})
}

func (h *Handler) UpdateAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAbstract{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	a := models.Text{Text: b.Text, Lang: b.Lang}
	if err := ctx.Dataset.SetAbstract(b.Position, a); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		form := abstractForm(ctx, b, validationErrs.(validation.Errors))

		render.Layout(w, "refresh_modal", "dataset/edit_abstract", YieldEditAbstract{
			Context:  ctx,
			Position: b.Position,
			Form:     form,
		})
		return
	}

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "show_modal", "error_dialog", ctx.Locale.T("dataset.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_abstracts", YieldAbstracts{
		Context: ctx,
	})
}

func (h *Handler) ConfirmDeleteAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteAbstract
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.Layout(w, "show_modal", "dataset/confirm_delete_abstract", YieldDeleteAbstract{
		Context:  ctx,
		Position: b.Position,
	})
}

func (h *Handler) DeleteAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteAbstract
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	if err := ctx.Dataset.RemoveAbstract(b.Position); err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "show_modal", "error_dialog", ctx.Locale.T("dataset.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_abstracts", YieldAbstracts{
		Context: ctx,
	})
}

func abstractForm(ctx Context, b BindAbstract, errors validation.Errors) *form.Form {
	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(ctx.Locale, errors)).
		AddSection(
			&form.TextArea{
				Name:        "text",
				Value:       b.Text,
				Label:       ctx.Locale.T("builder.abstract.text"),
				Cols:        12,
				Rows:        6,
				Placeholder: ctx.Locale.T("builder.abstract.text.placeholder"),
				Error:       localize.ValidationErrorAt(ctx.Locale, errors, fmt.Sprintf("/abstract/%d/text", b.Position)),
			},
			&form.Select{
				Name:    "lang",
				Value:   b.Lang,
				Label:   ctx.Locale.T("builder.abstract.lang"),
				Options: localize.LanguageSelectOptions(ctx.Locale),
				Cols:    12,
				Error:   localize.ValidationErrorAt(ctx.Locale, errors, fmt.Sprintf("/abstract/%d/lang", b.Position)),
			},
		)
}
