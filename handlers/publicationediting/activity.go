package publicationediting

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
	publicationviews "github.com/ugent-library/biblio-backoffice/views/publication"
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
	p := ctx.GetPublication(r)

	views.ShowModal(publicationviews.EditMessageDialog(c, publicationviews.EditMessageDialogArgs{
		Publication: p,
	})).Render(r.Context(), w)
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindMessage{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("update publication reviewer note: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	p := ctx.GetPublication(r)
	p.Message = b.Message

	if validationErrs := p.Validate(); validationErrs != nil {
		c.Log.Warnw("update publication reviewer note: could not validate message:", "errors", validationErrs, "publication", p.ID, "user", c.User.ID)
		views.ReplaceModal(publicationviews.EditMessageDialog(c, publicationviews.EditMessageDialogArgs{
			Publication: p,
			Errors:      validationErrs.(*okay.Errors),
			Conflict:    false,
		})).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(publicationviews.EditMessageDialog(c, publicationviews.EditMessageDialogArgs{
			Publication: p,
			Conflict:    true,
		})).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("update publication message: could not save the publication:", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(publicationviews.MessageBodySelector, publicationviews.MessageBody(c, p)).Render(r.Context(), w)
}

func (h *Handler) EditReviewerTags(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "show_modal", "publication/edit_reviewer_tags", YieldEditReviewerTags{
		Context:  ctx,
		Form:     reviewerTagsForm(ctx.User, ctx.Loc, ctx.Publication, nil),
		Conflict: false,
	})
}

func (h *Handler) UpdateReviewerTags(w http.ResponseWriter, r *http.Request, ctx Context) {
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
		render.Layout(w, "refresh_modal", "publication/edit_reviewer_tags", YieldEditReviewerTags{
			Context:  ctx,
			Form:     reviewerTagsForm(ctx.User, ctx.Loc, p, validationErrs.(*okay.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "publication/edit_reviewer_tags", YieldEditReviewerTags{
			Context:  ctx,
			Form:     reviewerTagsForm(ctx.User, ctx.Loc, p, nil),
			Conflict: true,
		})
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
	render.Layout(w, "show_modal", "publication/edit_reviewer_note", YieldEditReviewerNote{
		Context:  ctx,
		Form:     reviewerNoteForm(ctx.User, ctx.Loc, ctx.Publication, nil),
		Conflict: false,
	})
}

func (h *Handler) UpdateReviewerNote(w http.ResponseWriter, r *http.Request, ctx Context) {

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
		render.Layout(w, "refresh_modal", "publication/edit_reviewer_note", YieldEditReviewerNote{
			Context:  ctx,
			Form:     reviewerNoteForm(ctx.User, ctx.Loc, p, validationErrs.(*okay.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "publication/edit_reviewer_note", YieldEditReviewerNote{
			Context:  ctx,
			Form:     reviewerNoteForm(ctx.User, ctx.Loc, p, nil),
			Conflict: true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("update publication reviewer note: could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_reviewer_note", ctx)
}

func messageForm(_ *models.Person, loc *gotext.Locale, p *models.Publication, errors *okay.Errors) *form.Form {
	return form.New().
		WithTheme("cols").
		WithErrors(localize.ValidationErrors(loc, errors)).
		AddSection(
			&form.TextArea{
				Name:  "message",
				Value: p.Message,
				Label: loc.Get("builder.message"),
				Cols:  9,
				Rows:  10,
				Error: localize.ValidationErrorAt(
					loc,
					errors,
					"/message",
				),
			},
		)
}

func reviewerTagsForm(_ *models.Person, loc *gotext.Locale, p *models.Publication, errors *okay.Errors) *form.Form {
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

func reviewerNoteForm(_ *models.Person, loc *gotext.Locale, p *models.Publication, errors *okay.Errors) *form.Form {
	return form.New().
		WithTheme("cols").
		WithErrors(localize.ValidationErrors(loc, errors)).
		AddSection(
			&form.TextArea{
				Name:  "reviewer_note",
				Value: p.ReviewerNote,
				Label: loc.Get("builder.reviewer_note"),
				Cols:  9,
				Rows:  10,
				Error: localize.ValidationErrorAt(
					loc,
					errors,
					"/reviewer_note",
				),
			},
		)
}
