package datasetediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/snapstore"
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

func UpdateBiblioMessage(w http.ResponseWriter, r *http.Request) {
	updateDataset(w, r, func(d *models.Dataset, b BindMessage) {
		d.Message = b.Message
	})
}

func UpdateReviewerTags(w http.ResponseWriter, r *http.Request) {
	updateDataset(w, r, func(d *models.Dataset, b BindReviewerTags) {
		d.ReviewerTags = b.ReviewerTags
	})
}

func UpdateReviewerNote(w http.ResponseWriter, r *http.Request) {
	updateDataset(w, r, func(d *models.Dataset, b BindReviewerNote) {
		d.ReviewerNote = b.ReviewerNote
	})
}

func updateDataset[TBind any](w http.ResponseWriter, r *http.Request, setter func(p *models.Dataset, b TBind)) {
	c := ctx.Get(r)

	var b TBind
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	p := ctx.GetDataset(r)
	setter(p, b)

	if validationErrs := p.Validate(); validationErrs != nil {
		datasetviews.Messages(c, datasetviews.MessagesArgs{
			Dataset:  p,
			Errors:   validationErrs.(*okay.Errors),
			Conflict: false,
		}).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		datasetviews.Messages(c, datasetviews.MessagesArgs{
			Dataset:  p,
			Conflict: true,
		}).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	datasetviews.Messages(c, datasetviews.MessagesArgs{
		Dataset: p,
	}).Render(r.Context(), w)
}
