package datasetediting

import (
	"errors"
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

type BindMessage struct {
	Message string `form:"message"`
}

type BindReviewerTags struct {
	ReviewerTags []string `form:"reviewer_tags"`
}

type BindReviewerNote struct {
	ReviewerNote string `form:"reviewer_note"`
}

type YieldEditMessage struct {
	Context
	Form     *form.Form
	Conflict bool
}

type YieldEditReviewerTags struct {
	Context
	Form     *form.Form
	Conflict bool
}

type YieldEditReviewerNote struct {
	Context
	Form     *form.Form
	Conflict bool
}

func (h *Handler) EditMessage(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "show_modal", "dataset/edit_message", YieldEditMessage{
		Context: ctx,
		Form:    messageForm(ctx.User, ctx.Locale, ctx.Dataset, nil),
	})
}

func (h *Handler) UpdateMessage(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindMessage{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update dataset reviewer note: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	p := ctx.Dataset
	p.Message = b.Message

	if validationErrs := p.Validate(); validationErrs != nil {
		h.Logger.Warnw("update dataset reviewer note: could not validate message:", "errors", validationErrs, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.Layout(w, "refresh_modal", "dataset/edit_message", YieldEditMessage{
			Context:  ctx,
			Form:     messageForm(ctx.User, ctx.Locale, p, validationErrs.(validation.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "dataset/edit_message", YieldEditMessage{
			Context:  ctx,
			Form:     messageForm(ctx.User, ctx.Locale, p, nil),
			Conflict: true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("update dataset message: could not save the dataset:", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_message", ctx)
}

func (h *Handler) EditReviewerTags(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCurate() {
		render.Unauthorized(w, r)
		return
	}

	render.Layout(w, "show_modal", "dataset/edit_reviewer_tags", YieldEditReviewerTags{
		Context: ctx,
		Form:    reviewerTagsForm(ctx.User, ctx.Locale, ctx.Dataset, nil),
	})
}

func (h *Handler) UpdateReviewerTags(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCurate() {
		render.Unauthorized(w, r)
		return
	}

	b := BindReviewerTags{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update dataset reviewer tags: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	p := ctx.Dataset
	p.ReviewerTags = b.ReviewerTags

	if validationErrs := p.Validate(); validationErrs != nil {
		h.Logger.Warnw("update dataset reviewer tags: could not validate reviewer tags:", "errors", validationErrs, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.Layout(w, "refresh_modal", "dataset/edit_reviewer_tags", YieldEditReviewerTags{
			Context:  ctx,
			Form:     reviewerTagsForm(ctx.User, ctx.Locale, p, validationErrs.(validation.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "dataset/edit_reviewer_tags", YieldEditReviewerTags{
			Context:  ctx,
			Form:     reviewerTagsForm(ctx.User, ctx.Locale, p, nil),
			Conflict: true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("update dataset reviewer tags: could not save the dataset:", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_reviewer_tags", ctx)
}

func (h *Handler) EditReviewerNote(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCurate() {
		render.Unauthorized(w, r)
		return
	}

	render.Layout(w, "show_modal", "dataset/edit_reviewer_note", YieldEditReviewerNote{
		Context:  ctx,
		Form:     reviewerNoteForm(ctx.User, ctx.Locale, ctx.Dataset, nil),
		Conflict: false,
	})
}

func (h *Handler) UpdateReviewerNote(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCurate() {
		render.Unauthorized(w, r)
		return
	}

	b := BindReviewerNote{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update dataset reviewer note: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	p := ctx.Dataset
	p.ReviewerNote = b.ReviewerNote

	if validationErrs := p.Validate(); validationErrs != nil {
		h.Logger.Warnw("update dataset reviewer note: could not validate reviewer note:", "errors", validationErrs, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.Layout(w, "refresh_modal", "dataset/edit_reviewer_note", YieldEditReviewerNote{
			Context:  ctx,
			Form:     reviewerNoteForm(ctx.User, ctx.Locale, p, validationErrs.(validation.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "dataset/edit_reviewer_note", YieldEditReviewerNote{
			Context:  ctx,
			Form:     reviewerNoteForm(ctx.User, ctx.Locale, p, nil),
			Conflict: true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("update dataset reviewer note: could not save the dataset:", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_reviewer_note", ctx)
}

func messageForm(user *models.User, l *locale.Locale, p *models.Dataset, errors validation.Errors) *form.Form {
	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(l, errors)).
		AddSection(
			&form.TextArea{
				Name:  "message",
				Value: p.Message,
				Label: l.T("builder.message"),
				Cols:  9,
				Rows:  4,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					"/message",
				),
			},
		)
}

func reviewerTagsForm(user *models.User, l *locale.Locale, p *models.Dataset, errors validation.Errors) *form.Form {
	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(l, errors)).
		AddSection(
			&form.TextRepeat{
				Name:   "reviewer_tags",
				Values: p.ReviewerTags,
				Label:  l.T("builder.reviewer_tags"),
				Cols:   9,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					"/reviewer_tags",
				),
			},
		)
}

func reviewerNoteForm(user *models.User, l *locale.Locale, p *models.Dataset, errors validation.Errors) *form.Form {
	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(l, errors)).
		AddSection(
			&form.TextArea{
				Name:  "reviewer_note",
				Value: p.ReviewerNote,
				Label: l.T("builder.reviewer_note"),
				Cols:  9,
				Rows:  4,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					"/reviewer_note",
				),
			},
		)
}
