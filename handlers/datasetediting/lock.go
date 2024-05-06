package datasetediting

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/render/flash"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

func Lock(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	dataset.Locked = true

	if validationErrs := dataset.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(c.Loc, validationErrs.(*okay.Errors)))
		views.FormErrorsDialog(c, "Unable to lock this dataset due to the following errors", errors).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdateDataset(dataset.SnapshotID, dataset, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("lock dataset: could not save the dataset:", "error", err, "dataset", dataset.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	f := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Dataset was successfully locked.</p>"))

	c.PersistFlash(w, *f)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}

func Unlock(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	dataset.Locked = false

	if validationErrs := dataset.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(c.Loc, validationErrs.(*okay.Errors)))
		views.FormErrorsDialog(c, "Unable to unlock this dataset due to the following errors", errors).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdateDataset(dataset.SnapshotID, dataset, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ShowModal(
			views.ErrorDialogWithOptions(c.Loc.Get("dataset.conflict_error_reload"), views.ErrorDialogOptions{
				RedirectURL: r.URL.Query().Get("redirect-url"),
			})).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("unlock dataset: could not save the dataset:", "error", err, "dataset", dataset.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	f := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Dataset was successfully unlocked.</p>"))

	c.PersistFlash(w, *f)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
