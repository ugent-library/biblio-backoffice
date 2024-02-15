package suggestions

import (
	"html/template"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/pagination"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/flash"
	"github.com/ugent-library/biblio-backoffice/views"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
)

type bindSuggestion struct {
	ID string `path:"id" form:"id"`
}

func Suggestions(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if !c.User.CanCurate() {
		c.HandleError(w, r, httperror.Unauthorized)
		return
	}

	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		c.Log.Warnw("could not bind search arguments", "errors", err, "request", r, "user", c.User.ID)
		render.BadRequest(w, r, err)
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
	views.Suggestions(c, searchArgs, searchHits, recs).Render(r.Context(), w)
}

func SuggestionsIcon(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	countRecs, err := c.Repo.CountCandidateRecords(r.Context())
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.SuggestionsIcon(c, countRecs > 0).Render(r.Context(), w)
}

func ConfirmDeleteSuggestion(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if !c.User.CanCurate() {
		c.HandleError(w, r, httperror.Unauthorized)
		return
	}

	b := bindSuggestion{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete suggestion: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	rec, err := c.Repo.GetCandidateRecord(r.Context(), b.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.ConfirmDeleteSuggestion(c, rec).Render(r.Context(), w)
}

func DeleteSuggestion(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if !c.User.CanCurate() {
		c.HandleError(w, r, httperror.Unauthorized)
		return
	}

	b := bindSuggestion{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("delete suggestion: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	err := c.Repo.DeleteCandidateRecord(r.Context(), b.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	f := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Suggestion was successfully deleted.</p>"))

	c.PersistFlash(w, *f)

	w.Header().Set("HX-Redirect", c.URLTo("suggestions").String())
}

func ImportSuggestion(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if !c.User.CanCurate() {
		c.HandleError(w, r, httperror.Unauthorized)
		return
	}

	b := bindSuggestion{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("import suggestion: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	rec, err := c.Repo.GetCandidateRecord(r.Context(), b.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	var pubID string
	if pubID, err = c.Repo.ImportCandidateRecordAsPublication(r.Context(), rec, c.User); err != nil {
		c.HandleError(w, r, err)
		return
	}

	f := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Suggestion was successfully imported!</p>"))
	c.PersistFlash(w, *f)

	w.Header().Set("HX-Redirect", c.URLTo("publication", "id", pubID).String())
}