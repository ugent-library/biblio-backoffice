package datasetediting

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/flash"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views/dataset"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

func ConfirmRepublish(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	redirectUrl := r.URL.Query().Get("redirect-url")

	dataset.ConfirmRepublish(c, ctx.GetDataset(r), redirectUrl).Render(r.Context(), w)
}

func Republish(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	if !c.User.CanPublishDataset(dataset) {
		c.Log.Warnw("republish dataset: user has no permission to republish", "dataset", dataset.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.Forbidden)
		return
	}

	dataset.Status = "public"

	if validationErrs := dataset.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(c.Loc, validationErrs.(*okay.Errors)))
		render.Layout(w, "refresh_modal", "form_errors_dialog", struct {
			Title  string
			Errors form.Errors
		}{
			Title:  "Unable to republish this dataset due to the following errors",
			Errors: errors,
		})
		return
	}

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), dataset, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message:     c.Loc.Get("dataset.conflict_error"),
			RedirectURL: r.URL.Query().Get("redirect-url"),
		})
		return
	}

	if err != nil {
		c.Log.Errorf("republish dataset: could not save the dataset:", "errors", err, "dataset", dataset.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	flash := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Dataset was successfully republished.</p>"))

	c.PersistFlash(w, *flash)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
