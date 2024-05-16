package datasetediting

import (
	"errors"
	"net/http"

	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	datasetviews "github.com/ugent-library/biblio-backoffice/views/dataset"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
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

func EditMessage(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	views.ShowModal(datasetviews.EditMessageDialog(c, datasetviews.EditMessageDialogArgs{
		Dataset: d,
	})).Render(r.Context(), w)
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindMessage{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("update dataset reviewer note: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	d := ctx.GetDataset(r)
	d.Message = b.Message

	if validationErrs := d.Validate(); validationErrs != nil {
		c.Log.Warnw("update dataset reviewer note: could not validate message:", "errors", validationErrs, "dataset", d.ID, "user", c.User.ID)
		views.ReplaceModal(datasetviews.EditMessageDialog(c, datasetviews.EditMessageDialogArgs{
			Dataset:  d,
			Errors:   validationErrs.(*okay.Errors),
			Conflict: false,
		})).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), d, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(datasetviews.EditMessageDialog(c, datasetviews.EditMessageDialogArgs{
			Dataset:  d,
			Conflict: true,
		})).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("update dataset message: could not save the dataset:", "errors", err, "dataset", d.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(datasetviews.MessageBodySelector, datasetviews.MessageBody(c, d)).Render(r.Context(), w)
}

func (h *Handler) EditReviewerTags(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCurate() {
		render.Unauthorized(w, r)
		return
	}

	render.Layout(w, "show_modal", "dataset/edit_reviewer_tags", YieldEditReviewerTags{
		Context: ctx,
		Form:    reviewerTagsForm(ctx.User, ctx.Loc, ctx.Dataset, nil),
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
			Form:     reviewerTagsForm(ctx.User, ctx.Loc, p, validationErrs.(*okay.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "dataset/edit_reviewer_tags", YieldEditReviewerTags{
			Context:  ctx,
			Form:     reviewerTagsForm(ctx.User, ctx.Loc, p, nil),
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
		Form:     reviewerNoteForm(ctx.User, ctx.Loc, ctx.Dataset, nil),
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
			Form:     reviewerNoteForm(ctx.User, ctx.Loc, p, validationErrs.(*okay.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "dataset/edit_reviewer_note", YieldEditReviewerNote{
			Context:  ctx,
			Form:     reviewerNoteForm(ctx.User, ctx.Loc, p, nil),
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

func reviewerTagsForm(user *models.Person, loc *gotext.Locale, p *models.Dataset, errors *okay.Errors) *form.Form {
	return form.New().
		WithTheme("cols").
		WithErrors(localize.ValidationErrors(loc, errors)).
		AddSection(
			&form.TextRepeat{
				Name:   "reviewer_tags",
				Values: p.ReviewerTags,
				Label:  loc.Get("builder.reviewer_tags"),
				Cols:   9,
				Error: localize.ValidationErrorAt(
					loc,
					errors,
					"/reviewer_tags",
				),
			},
		)
}

func reviewerNoteForm(user *models.Person, loc *gotext.Locale, p *models.Dataset, errors *okay.Errors) *form.Form {
	return form.New().
		WithTheme("cols").
		WithErrors(localize.ValidationErrors(loc, errors)).
		AddSection(
			&form.TextArea{
				Name:  "reviewer_note",
				Value: p.ReviewerNote,
				Label: loc.Get("builder.reviewer_note"),
				Cols:  9,
				Rows:  4,
				Error: localize.ValidationErrorAt(
					loc,
					errors,
					"/reviewer_note",
				),
			},
		)
}
