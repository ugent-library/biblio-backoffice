package publicationediting

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
	publication := ctx.GetPublication(r)

	publication.Locked = true

	if validationErrs := publication.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(c.Loc, validationErrs.(*okay.Errors)))
		views.FormErrorsDialog(c, "Unable to lock this publication due to the following errors", errors).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(publication.SnapshotID, publication, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"), "")).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("lock publication: could not save the publication:", "error", err, "publication", publication.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	f := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Publication was successfully locked.</p>"))

	c.PersistFlash(w, *f)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}

func Unlock(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	publication.Locked = false

	if validationErrs := publication.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(c.Loc, validationErrs.(*okay.Errors)))
		views.FormErrorsDialog(c, "Unable to unlock this publication due to the following errors", errors).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(publication.SnapshotID, publication, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error"), r.URL.Query().Get("redirect-url"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("unlock publication: could not save the publication:", "error", err, "publication", publication.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	f := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Publication was successfully unlocked.</p>"))

	c.PersistFlash(w, *f)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
