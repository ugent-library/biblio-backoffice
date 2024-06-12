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

type bindCandidateRecord struct {
	ID string `path:"id" form:"id"`
}

func CandidateRecords(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	countRecs, err := c.Repo.CountCandidateRecords(r.Context())
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	recs, err := c.Repo.GetCandidateRecords(r.Context(), searchArgs.Offset(), searchArgs.Limit())
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	searchHits := &models.SearchHits{
		Pagination: pagination.Pagination{
			Offset: searchArgs.Offset(),
			Limit:  searchArgs.Limit(),
			Total:  countRecs,
		},
	}
	candidaterecordviews.List(c, searchArgs, searchHits, recs).Render(r.Context(), w)
}

func CandidateRecordPreview(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if !c.User.CanCurate() {
		c.HandleError(w, r, httperror.Unauthorized)
		return
	}

	b := bindCandidateRecord{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	rec, err := c.Repo.GetCandidateRecord(r.Context(), b.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.ShowModal(candidaterecordviews.Preview(c, rec)).Render(r.Context(), w)
}

func CandidateRecordsIcon(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	countRecs, err := c.Repo.CountCandidateRecords(r.Context())
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	views.CandidateRecordsIcon(c, countRecs > 0).Render(r.Context(), w)
}

func ConfirmRejectCandidateRecord(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if !c.User.CanCurate() {
		c.HandleError(w, r, httperror.Unauthorized)
		return
	}

	b := bindCandidateRecord{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	rec, err := c.Repo.GetCandidateRecord(r.Context(), b.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.ShowModal(candidaterecordviews.ConfirmHide(c, rec)).Render(r.Context(), w)
}

func RejectCandidateRecord(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if !c.User.CanCurate() {
		c.HandleError(w, r, httperror.Unauthorized)
		return
	}

	b := bindCandidateRecord{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	err := c.Repo.RejectCandidateRecord(r.Context(), b.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	f := flash.SimpleFlash().
		WithLevel("success").
		WithBody("<p>Candidate record was successfully deleted.</p>")

	c.PersistFlash(w, *f)

	w.Header().Set("HX-Redirect", c.URLTo("candidate_records").String())
}

func ImportCandidateRecord(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if !c.User.CanCurate() {
		c.HandleError(w, r, httperror.Unauthorized)
		return
	}

	b := bindCandidateRecord{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	pubID, err := c.Repo.ImportCandidateRecordAsPublication(r.Context(), b.ID, c.User)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	f := flash.SimpleFlash().
		WithLevel("success").
		WithBody("<p>Suggestion was successfully imported!</p>")
	c.PersistFlash(w, *f)

	w.Header().Set("HX-Redirect", c.URLTo("publication", "id", pubID).String())
}
