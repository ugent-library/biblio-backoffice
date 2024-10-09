package publicationediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/snapstore"
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

func UpdateBiblioMessage(w http.ResponseWriter, r *http.Request) {
	updatePublication(w, r, func(p *models.Publication, b BindMessage) {
		p.Message = b.Message
	})
}

func UpdateReviewerTags(w http.ResponseWriter, r *http.Request) {
	updatePublication(w, r, func(p *models.Publication, b BindReviewerTags) {
		p.ReviewerTags = b.ReviewerTags
	})
}

func UpdateReviewerNote(w http.ResponseWriter, r *http.Request) {
	updatePublication(w, r, func(p *models.Publication, b BindReviewerNote) {
		p.ReviewerNote = b.ReviewerNote
	})
}

func updatePublication[TBind any](w http.ResponseWriter, r *http.Request, setter func(p *models.Publication, b TBind)) {
	c := ctx.Get(r)

	var b TBind
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	p := ctx.GetPublication(r)
	setter(p, b)

	if validationErrs := p.Validate(); validationErrs != nil {
		publicationviews.Messages(c, publicationviews.MessagesArgs{
			Publication: p,
			Errors:      validationErrs.(*okay.Errors),
			Conflict:    false,
		}).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		publicationviews.Messages(c, publicationviews.MessagesArgs{
			Publication: p,
			Conflict:    true,
		}).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	publicationviews.Messages(c, publicationviews.MessagesArgs{
		Publication: p,
	}).Render(r.Context(), w)
}
