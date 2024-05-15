package datasetediting

import (
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/render/flash"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	"github.com/ugent-library/httperror"
)

func ConfirmDelete(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to delete this dataset?",
		DeleteUrl:  c.PathTo("dataset_delete", "id", dataset.ID, "redirect-url", r.URL.Query().Get("redirect-url")),
		SnapshotID: dataset.SnapshotID,
	}).Render(r.Context(), w)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	if !c.User.CanDeleteDataset(dataset) {
		c.Log.Warnw("delete dataset: user isn't allowed to delete dataset", "dataset", dataset.ID, "user", c.User.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.Forbidden)
		return
	}

	dataset.Status = "deleted"

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), dataset, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("delete dataset: Could not save the dataset:", "error", err, "identifier", dataset.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	flash := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Dataset was successfully deleted.</p>"))

	c.PersistFlash(w, *flash)

	// TODO temporary fix until we can figure out a way let ES notify this handler that it did its thing.
	// see: https://github.com/ugent-library/biblio-backoffice/issues/590
	time.Sleep(1250 * time.Millisecond)

	redirectURL := r.URL.Query().Get("redirect-url")
	if redirectURL == "" {
		redirectURL = c.PathTo("datasets").String()
	}
	w.Header().Set("HX-Redirect", redirectURL)
}
