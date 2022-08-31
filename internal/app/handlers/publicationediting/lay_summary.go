package publicationediting

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type BindLaySummary struct {
	LaySummaryID string `path:"lay_summary_id"`
	Text         string `form:"text"`
	Lang         string `form:"lang"`
}

type BindDeleteLaySummary struct {
	LaySummaryID string `path:"lay_summary_id"`
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
	LaySummaryID string
	Form         *form.Form
}
type YieldDeleteLaySummary struct {
	Context
	LaySummaryID string
}

func (h *Handler) AddLaySummary(w http.ResponseWriter, r *http.Request, ctx Context) {
	form := laySummaryForm(ctx.Locale, ctx.Publication, &models.Text{}, nil)

	render.Layout(w, "show_modal", "publication/add_lay_summary", YieldAddLaySummary{
		Context: ctx,
		Form:    form,
	})
}

func (h *Handler) CreateLaySummary(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLaySummary{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("create publication lay summary: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	laySummary := models.Text{
		Lang: b.Lang,
		Text: b.Text,
	}
	ctx.Publication.AddLaySummary(&laySummary)

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		h.Logger.Warnw("create publication lay summary: could not validate contributor:", "errors", validationErrs, "identifier", ctx.Publication.ID)
		render.Layout(w, "refresh_modal", "publication/add_lay_summary", YieldAddLaySummary{
			Context: ctx,
			Form:    laySummaryForm(ctx.Locale, ctx.Publication, &laySummary, validationErrs.(validation.Errors)),
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		h.Logger.Warnf("create publication lay summary: snapstore detected a conflicting publication:", "errors", errors.As(err, &conflict), "identifier", ctx.Publication.ID)
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		h.Logger.Errorf("create publication lay summary: Could not save the publication:", "error", err, "identifier", ctx.Publication.ID)
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
		h.Logger.Warnw("edit publication lay summary: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	laySummary := ctx.Publication.GetLaySummary(b.LaySummaryID)
	if laySummary == nil {
		render.BadRequest(
			w,
			r,
			fmt.Errorf("no lay summary found for %s in publication %s", b.LaySummaryID, ctx.Publication.ID),
		)
		return
	}

	render.Layout(w, "show_modal", "publication/edit_lay_summary", YieldEditLaySummary{
		Context:      ctx,
		LaySummaryID: b.LaySummaryID,
		Form:         laySummaryForm(ctx.Locale, ctx.Publication, laySummary, nil),
	})
}

func (h *Handler) UpdateLaySummary(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLaySummary{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("update publication lay summary: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	/*
		TODO: this may throw a "bad request"
		when it should be throwing a conflict error.
		Always throw a conflict error (no one just "knows" an invalid id)?
	*/
	laySummary := ctx.Publication.GetLaySummary(b.LaySummaryID)
	if laySummary == nil {
		h.Logger.Warnw("update publication lay summary: could not get lay summary", "laysummary", b.LaySummaryID, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.BadRequest(
			w,
			r,
			fmt.Errorf("no lay summary found for %s in publication %s", b.LaySummaryID, ctx.Publication.ID),
		)
		return
	}
	laySummary.Text = b.Text
	laySummary.Lang = b.Lang

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		form := laySummaryForm(ctx.Locale, ctx.Publication, laySummary, validationErrs.(validation.Errors))

		render.Layout(w, "refresh_modal", "publication/edit_lay_summary", YieldEditLaySummary{
			Context:      ctx,
			LaySummaryID: b.LaySummaryID,
			Form:         form,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		h.Logger.Errorf("update publication details: Could not save the publication:", "error", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
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
		h.Logger.Warnw("confirm delete publication lay summary: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	render.Layout(w, "show_modal", "publication/confirm_delete_lay_summary", YieldDeleteLaySummary{
		Context:      ctx,
		LaySummaryID: b.LaySummaryID,
	})
}

func (h *Handler) DeleteLaySummary(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteLaySummary
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete publication lay summary: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	/*
		Note: ignore fact that lay summary is already removed:
		conflicting resolving will solve this
	*/
	ctx.Publication.RemoveLaySummary(b.LaySummaryID)

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		h.Logger.Errorf("delete publication lay summary: Could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_lay_summaries", YieldLaySummaries{
		Context: ctx,
	})
}

func laySummaryForm(l *locale.Locale, publication *models.Publication, laySummary *models.Text, errors validation.Errors) *form.Form {
	idx := -1
	for i, ls := range publication.LaySummary {
		if ls.ID == laySummary.ID {
			idx = i
			break
		}
	}
	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(l, errors)).
		AddSection(
			&form.TextArea{
				Name:  "text",
				Value: laySummary.Text,
				Label: l.T("builder.lay_summary.text"),
				Cols:  12,
				Rows:  6,
				Error: localize.ValidationErrorAt(l, errors, fmt.Sprintf("/lay_summary/%d/text", idx)),
			},
			&form.Select{
				Name:    "lang",
				Value:   laySummary.Lang,
				Label:   l.T("builder.lay_summary.lang"),
				Options: localize.LanguageSelectOptions(l),
				Cols:    12,
				Error:   localize.ValidationErrorAt(l, errors, fmt.Sprintf("/lay_summary/%d/lang", idx)),
			},
		)
}
