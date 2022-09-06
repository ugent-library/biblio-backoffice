package publicationediting

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
	Form *form.Form
}

type YieldEditReviewerTags struct {
	Context
	Form *form.Form
}

type YieldEditReviewerNote struct {
	Context
	Form *form.Form
}

func (h *Handler) EditMessage(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "show_modal", "publication/edit_message", YieldEditMessage{
		Context: ctx,
		Form:    messageForm(ctx.User, ctx.Locale, ctx.Publication, nil),
	})
}

func (h *Handler) UpdateMessage(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindMessage{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update publication reviewer note: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	p := ctx.Publication
	p.Message = b.Message

	if validationErrs := p.Validate(); validationErrs != nil {
		h.Logger.Warnw("update publication reviewer note: could not validate message:", "errors", validationErrs, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		form := messageForm(ctx.User, ctx.Locale, p, validationErrs.(validation.Errors))

		render.Layout(w, "refresh_modal", "publication/edit_message", YieldEditMessage{
			Context: ctx,
			Form:    form,
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
		h.Logger.Errorf("update publication message: could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_message", ctx)
}

func (h *Handler) EditReviewerTags(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCuratePublications() {
		render.Unauthorized(w, r)
		return
	}

	render.Layout(w, "show_modal", "publication/edit_reviewer_tags", YieldEditReviewerTags{
		Context: ctx,
		Form:    reviewerTagsForm(ctx.User, ctx.Locale, ctx.Publication, nil),
	})
}

func (h *Handler) UpdateReviewerTags(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCuratePublications() {
		render.Unauthorized(w, r)
		return
	}

	b := BindReviewerTags{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update publication reviewer tags: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	p := ctx.Publication
	p.ReviewerTags = b.ReviewerTags

	if validationErrs := p.Validate(); validationErrs != nil {
		h.Logger.Warnw("update publication reviewer tags: could not validate reviewer tags:", "errors", validationErrs, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		form := reviewerTagsForm(ctx.User, ctx.Locale, p, validationErrs.(validation.Errors))

		render.Layout(w, "refresh_modal", "publication/edit_reviewer_tags", YieldEditReviewerTags{
			Context: ctx,
			Form:    form,
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
		h.Logger.Errorf("update publication reviewer tags: could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_reviewer_tags", ctx)
}

func (h *Handler) EditReviewerNote(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCuratePublications() {
		render.Unauthorized(w, r)
		return
	}

	render.Layout(w, "show_modal", "publication/edit_reviewer_note", YieldEditReviewerNote{
		Context: ctx,
		Form:    reviewerNoteForm(ctx.User, ctx.Locale, ctx.Publication, nil),
	})
}

func (h *Handler) UpdateReviewerNote(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCuratePublications() {
		render.Unauthorized(w, r)
		return
	}

	b := BindReviewerNote{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update publication reviewer note: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	p := ctx.Publication
	p.ReviewerNote = b.ReviewerNote

	if validationErrs := p.Validate(); validationErrs != nil {
		h.Logger.Warnw("update publication reviewer note: could not validate reviewer note:", "errors", validationErrs, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		form := reviewerNoteForm(ctx.User, ctx.Locale, p, validationErrs.(validation.Errors))

		render.Layout(w, "refresh_modal", "publication/edit_reviewer_note", YieldEditReviewerNote{
			Context: ctx,
			Form:    form,
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
		h.Logger.Errorf("update publication reviewer note: could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_reviewer_note", ctx)
}

func messageForm(user *models.User, l *locale.Locale, p *models.Publication, errors validation.Errors) *form.Form {
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

func reviewerTagsForm(user *models.User, l *locale.Locale, p *models.Publication, errors validation.Errors) *form.Form {
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

func reviewerNoteForm(user *models.User, l *locale.Locale, p *models.Publication, errors validation.Errors) *form.Form {
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
