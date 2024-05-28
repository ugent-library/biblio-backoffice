package publicationediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	"github.com/ugent-library/biblio-backoffice/views/flash"
	"github.com/ugent-library/okay"
)

func Lock(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	publication.Locked = true

	if validationErrs := publication.Validate(); validationErrs != nil {
		errors := localize.ValidationErrors(c.Loc, validationErrs.(*okay.Errors))
		w.Header().Add("HX-Retarget", "#modals")
		w.Header().Add("HX-Reswap", "innerHTML")
		views.ShowModal(views.FormErrorsDialog("Unable to lock this publication due to the following errors", errors)).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(publication.SnapshotID, publication, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	f := flash.SimpleFlash().
		WithLevel("success").
		WithBody("<p>Publication was successfully locked.</p>")

	c.PersistFlash(w, *f)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}

func Unlock(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	publication.Locked = false

	if validationErrs := publication.Validate(); validationErrs != nil {
		errors := localize.ValidationErrors(c.Loc, validationErrs.(*okay.Errors))
		w.Header().Add("HX-Retarget", "#modals")
		w.Header().Add("HX-Reswap", "innerHTML")
		views.ShowModal(views.FormErrorsDialog("Unable to unlock this publication due to the following errors", errors)).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(publication.SnapshotID, publication, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ShowModal(
			views.ErrorDialogWithOptions(c.Loc.Get("publication.conflict_error"), views.ErrorDialogOptions{
				RedirectURL: r.URL.Query().Get("redirect-url"),
			})).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	f := flash.SimpleFlash().
		WithLevel("success").
		WithBody("<p>Publication was successfully unlocked.</p>")

	c.PersistFlash(w, *f)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
