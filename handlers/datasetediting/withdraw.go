package datasetediting

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	datasetviews "github.com/ugent-library/biblio-backoffice/views/dataset"
	"github.com/ugent-library/biblio-backoffice/views/flash"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

func ConfirmWithdraw(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	redirectUrl := r.URL.Query().Get("redirect-url")

	datasetviews.ConfirmWithdraw(c, ctx.GetDataset(r), redirectUrl).Render(r.Context(), w)
}

func Withdraw(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	if !c.User.CanWithdrawDataset(dataset) {
		c.HandleError(w, r, httperror.Forbidden)
		return
	}

	dataset.Status = "returned"

	if validationErrs := dataset.Validate(); validationErrs != nil {
		errors := localize.ValidationErrors(c.Loc, validationErrs.(*okay.Errors))
		views.ReplaceModal(views.FormErrorsDialog("Unable to withdraw this dataset due to the following errors", errors)).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), dataset, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(
			views.ErrorDialogWithOptions(c.Loc.Get("dataset.conflict_error"), views.ErrorDialogOptions{
				RedirectURL: r.URL.Query().Get("redirect-url"),
			})).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, httperror.InternalServerError.Wrap(fmt.Errorf("could not save the dataset: %w", err)))
		return
	}

	flash := flash.SimpleFlash().
		WithLevel("success").
		WithBody("<p>Dataset was successfully withdrawn.</p>")

	c.PersistFlash(w, *flash)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
