package candidaterecords

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/pagination"
	"github.com/ugent-library/biblio-backoffice/views"
	candidaterecordviews "github.com/ugent-library/biblio-backoffice/views/candidaterecord"
	"github.com/ugent-library/biblio-backoffice/views/flash"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
)

func CandidateRecords(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	var total int
	var recs []*models.CandidateRecord
	var err error

	if c.ProxiedPerson != nil {
		total, recs, err = c.Repo.GetCandidateRecordsByPersonID(r.Context(), c.ProxiedPerson.ID, searchArgs.Offset(), searchArgs.Limit(), false)
	} else if c.UserRole == "curator" {
		total, recs, err = c.Repo.GetCandidateRecords(r.Context(), searchArgs.Offset(), searchArgs.Limit())
	} else {
		total, recs, err = c.Repo.GetCandidateRecordsByPersonID(r.Context(), c.User.ID, searchArgs.Offset(), searchArgs.Limit(), false)
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	searchHits := &models.SearchHits{
		Pagination: pagination.Pagination{
			Offset: searchArgs.Offset(),
			Limit:  searchArgs.Limit(),
			Total:  total,
		},
	}

	candidaterecordviews.List(c, searchArgs, searchHits, recs).Render(r.Context(), w)
}

func CandidateRecordPreview(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	rec := ctx.GetCandidateRecord(r)

	views.ShowModal(candidaterecordviews.Preview(c, rec)).Render(r.Context(), w)
}

func ConfirmRejectCandidateRecord(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	rec := ctx.GetCandidateRecord(r)

	views.ShowModal(candidaterecordviews.ConfirmHide(c, rec)).Render(r.Context(), w)
}

func RejectCandidateRecord(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	rec := ctx.GetCandidateRecord(r)

	err := c.Repo.RejectCandidateRecord(r.Context(), rec.ID, c.User)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	f := flash.SimpleFlash().
		WithLevel("success").
		WithBody("<p>Candidate record was successfully rejected.</p>")
	c.Flash = append(c.Flash, *f)

	rec, err = c.Repo.GetCandidateRecord(r.Context(), rec.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.Cat(
		candidaterecordviews.ListItem(c, rec),
		views.CloseModal(),
		views.Replace("#flash-messages", views.FlashMessages(c)),
	).Render(r.Context(), w)
}

func ImportCandidateRecord(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	rec := ctx.GetCandidateRecord(r)

	pubID, err := c.Repo.ImportCandidateRecordAsPublication(r.Context(), rec.ID, c.User)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	f := flash.SimpleFlash().
		WithLevel("success").
		WithBody("<p>Suggestion was successfully imported!</p>")
	c.PersistFlash(w, *f)

	w.Header().Set("HX-Redirect", c.PathTo("publication", "id", pubID).String())
}

func RestoreRejectedCandidateRecord(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	rec := ctx.GetCandidateRecord(r)

	err := c.Repo.RestoreCandidateRecord(r.Context(), rec.ID, c.User)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	f := flash.SimpleFlash().
		WithLevel("success").
		WithBody("<p>Candidate record was successfully restored.</p>")
	c.Flash = append(c.Flash, *f)

	rec, err = c.Repo.GetCandidateRecord(r.Context(), rec.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.Cat(
		candidaterecordviews.ListItem(c, rec),
		views.Replace("#flash-messages", views.FlashMessages(c)),
	).Render(r.Context(), w)
}
