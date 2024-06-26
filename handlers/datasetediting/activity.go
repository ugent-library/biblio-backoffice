package datasetediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
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
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	d := ctx.GetDataset(r)
	d.Message = b.Message

	if validationErrs := d.Validate(); validationErrs != nil {
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
		c.HandleError(w, r, err)
		return
	}

	views.CloseModalAndReplace(datasetviews.MessageBodySelector, datasetviews.MessageBody(c, d)).Render(r.Context(), w)
}

func EditReviewerTags(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	views.ShowModal(datasetviews.EditReviewerTagsDialog(c, datasetviews.EditReviewerTagsDialogArgs{
		Dataset: d,
	})).Render(r.Context(), w)
}

func UpdateReviewerTags(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindReviewerTags{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	d := ctx.GetDataset(r)
	d.ReviewerTags = b.ReviewerTags

	if validationErrs := d.Validate(); validationErrs != nil {
		views.ReplaceModal(datasetviews.EditReviewerTagsDialog(c, datasetviews.EditReviewerTagsDialogArgs{
			Dataset:  d,
			Errors:   validationErrs.(*okay.Errors),
			Conflict: false,
		})).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), d, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(datasetviews.EditReviewerTagsDialog(c, datasetviews.EditReviewerTagsDialogArgs{
			Dataset:  d,
			Conflict: true,
		})).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.CloseModalAndReplace(datasetviews.ReviewerTagsSelector, datasetviews.ReviewerTagsBody(c, d)).Render(r.Context(), w)
}

func EditReviewerNote(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	views.ShowModal(datasetviews.EditReviewerNoteDialog(c, datasetviews.EditReviewerNoteDialogArgs{
		Dataset: d,
	})).Render(r.Context(), w)
}

func UpdateReviewerNote(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindReviewerNote{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	d := ctx.GetDataset(r)
	d.ReviewerNote = b.ReviewerNote

	if validationErrs := d.Validate(); validationErrs != nil {
		views.ReplaceModal(datasetviews.EditReviewerNoteDialog(c, datasetviews.EditReviewerNoteDialogArgs{
			Dataset:  d,
			Errors:   validationErrs.(*okay.Errors),
			Conflict: false,
		}))
		return
	}

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), d, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(datasetviews.EditReviewerNoteDialog(c, datasetviews.EditReviewerNoteDialogArgs{
			Dataset:  d,
			Conflict: true,
		})).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.CloseModalAndReplace(datasetviews.ReviewerNoteSelector, datasetviews.ReviewerNoteBody(c, d)).Render(r.Context(), w)
}
