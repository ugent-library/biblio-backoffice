package publicationediting

import (
	"errors"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	"github.com/ugent-library/biblio-backoffice/views/flash"
	"github.com/ugent-library/httperror"
)

func ConfirmDelete(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	views.ConfirmDeleteDialog(views.ConfirmDeleteDialogArgs{
		Context:    c,
		Question:   "Are you sure you want to delete this publication?",
		DeleteUrl:  c.PathTo("publication_delete", "id", publication.ID, "redirect-url", r.URL.Query().Get("redirect-url")),
		SnapshotID: publication.SnapshotID,
	}).Render(r.Context(), w)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	if !c.Repo.CanDeletePublication(c.User, publication) {
		c.HandleError(w, r, httperror.Forbidden)
		return
	}

	publication.Status = "deleted"

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), publication, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(
			views.ErrorDialogWithOptions(c.Loc.Get("publication.conflict_error_reload"), views.ErrorDialogOptions{
				RedirectURL: r.URL.Query().Get("redirect-url"),
			})).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	flash := flash.SimpleFlash().
		WithLevel("success").
		WithBody("<p>Publication was successfully deleted.</p>")

	c.PersistFlash(w, *flash)

	// TODO temporary fix until we can figure out a way let ES notify this handler that it did its thing.
	// see: https://github.com/ugent-library/biblio-backoffice/issues/590
	time.Sleep(1250 * time.Millisecond)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
