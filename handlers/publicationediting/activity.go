package publicationediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
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
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	p := ctx.GetPublication(r)
	p.Message = b.Message

	if validationErrs := p.Validate(); validationErrs != nil {
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
		c.HandleError(w, r, err)
		return
	}

	views.CloseModalAndReplace(publicationviews.MessageBodySelector, publicationviews.MessageBody(c, p)).Render(r.Context(), w)
}

func EditReviewerTags(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	views.ShowModal(publicationviews.EditReviewerTagsDialog(c, publicationviews.EditReviewerTagsDialogArgs{
		Publication: p,
	})).Render(r.Context(), w)
}

func UpdateReviewerTags(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindReviewerTags{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	p := ctx.GetPublication(r)
	p.ReviewerTags = b.ReviewerTags

	if validationErrs := p.Validate(); validationErrs != nil {
		views.ReplaceModal(publicationviews.EditReviewerTagsDialog(c, publicationviews.EditReviewerTagsDialogArgs{
			Publication: p,
			Errors:      validationErrs.(*okay.Errors),
			Conflict:    false,
		})).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(publicationviews.EditReviewerTagsDialog(c, publicationviews.EditReviewerTagsDialogArgs{
			Publication: p,
			Conflict:    true,
		})).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.CloseModalAndReplace(publicationviews.ReviewerTagsBodySelector, publicationviews.ReviewerTagsBody(c, p)).Render(r.Context(), w)
}

func EditReviewerNote(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	views.ShowModal(publicationviews.EditReviewerNoteDialog(c, publicationviews.EditReviewerNoteDialogArgs{
		Publication: p,
	})).Render(r.Context(), w)
}

func UpdateReviewerNote(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindReviewerNote{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	p := ctx.GetPublication(r)
	p.ReviewerNote = b.ReviewerNote

	if validationErrs := p.Validate(); validationErrs != nil {
		views.ReplaceModal(publicationviews.EditReviewerNoteDialog(c, publicationviews.EditReviewerNoteDialogArgs{
			Publication: p,
			Errors:      validationErrs.(*okay.Errors),
			Conflict:    false,
		}))
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(publicationviews.EditReviewerNoteDialog(c, publicationviews.EditReviewerNoteDialogArgs{
			Publication: p,
			Conflict:    true,
		})).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.CloseModalAndReplace(publicationviews.ReviewerNoteBodySelector, publicationviews.ReviewerNoteBody(c, p)).Render(r.Context(), w)
}
