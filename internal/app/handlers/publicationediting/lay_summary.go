package publicationediting

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

type BindLaySummary struct {
	Position int    `path:"position"`
	Text     string `form:"text"`
	Lang     string `form:"lang"`
}

type BindDeleteLaySummary struct {
	Position int `path:"position"`
}

type YieldLaySummaries struct {
	Context
}
type YieldAddLaySummary struct {
	Context
	Form *form.Form
}
type YieldEditLaySummary struct {
	Context
	Position int
	Form     *form.Form
}
type YieldDeleteLaySummary struct {
	Context
	Position int
}

func (h *Handler) AddLaySummary(w http.ResponseWriter, r *http.Request, ctx Context) {
	form := laySummaryForm(ctx, BindLaySummary{Position: len(ctx.Publication.LaySummary)}, nil)

	render.Layout(w, "show_modal", "publication/add_lay_summary", YieldAddLaySummary{
		Context: ctx,
		Form:    form,
	})
}

func (h *Handler) CreateLaySummary(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLaySummary{Position: len(ctx.Publication.LaySummary)}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	ctx.Publication.LaySummary = append(ctx.Publication.LaySummary, models.Text{Text: b.Text, Lang: b.Lang})

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "publication/add_lay_summary", YieldAddLaySummary{
			Context: ctx,
			Form:    laySummaryForm(ctx, b, validationErrs.(validation.Errors)),
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_lay_summaries", YieldLaySummaries{
		Context: ctx,
	})
}

func (h *Handler) EditLaySummary(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLaySummary{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	a, err := ctx.Publication.GetLaySummary(b.Position)
	if err != nil {
		render.BadRequest(w, r, err)
		return
	}

	b.Lang = a.Lang
	b.Text = a.Text

	render.Layout(w, "show_modal", "publication/edit_lay_summary", YieldEditLaySummary{
		Context:  ctx,
		Position: b.Position,
		Form:     laySummaryForm(ctx, b, nil),
	})
}

func (h *Handler) UpdateLaySummary(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLaySummary{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	a := models.Text{Text: b.Text, Lang: b.Lang}
	if err := ctx.Publication.SetLaySummary(b.Position, a); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		form := laySummaryForm(ctx, b, validationErrs.(validation.Errors))

		render.Layout(w, "refresh_modal", "publication/edit_lay_summary", YieldEditLaySummary{
			Context:  ctx,
			Position: b.Position,
			Form:     form,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_lay_summaries", YieldLaySummaries{
		Context: ctx,
	})
}

func (h *Handler) ConfirmDeleteLaySummary(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteLaySummary
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.Layout(w, "show_modal", "publication/confirm_delete_lay_summary", YieldDeleteLaySummary{
		Context:  ctx,
		Position: b.Position,
	})
}

func (h *Handler) DeleteLaySummary(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteLaySummary
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	if err := ctx.Publication.RemoveLaySummary(b.Position); err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_lay_summaries", YieldLaySummaries{
		Context: ctx,
	})
}

func laySummaryForm(ctx Context, b BindLaySummary, errors validation.Errors) *form.Form {
	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(ctx.Locale, errors)).
		AddSection(
			&form.TextArea{
				Name:        "text",
				Value:       b.Text,
				Label:       ctx.T("builder.lay_summary.text"),
				Cols:        12,
				Rows:        6,
				Placeholder: ctx.T("builder.lay_summary.text.placeholder"),
				Error:       localize.ValidationErrorAt(ctx.Locale, errors, fmt.Sprintf("/lay_summary/%d/text", b.Position)),
			},
			&form.Select{
				Name:    "lang",
				Value:   b.Lang,
				Label:   ctx.T("builder.lay_summary.lang"),
				Options: localize.LanguageSelectOptions(ctx.Locale),
				Cols:    12,
				Error:   localize.ValidationErrorAt(ctx.Locale, errors, fmt.Sprintf("/lay_summary/%d/lang", b.Position)),
			},
		)
}
