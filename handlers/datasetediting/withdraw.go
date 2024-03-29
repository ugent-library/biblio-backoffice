package datasetediting

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/flash"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/validation"
)

type YieldWithdraw struct {
	Context
	RedirectURL string
}

func (h *Handler) ConfirmWithdraw(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "show_modal", "dataset/confirm_withdraw", YieldWithdraw{
		Context:     ctx,
		RedirectURL: r.URL.Query().Get("redirect-url"),
	})
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanEditDataset(ctx.Dataset) {
		h.Logger.Warnw("withdraw dataset: user has no permission to withdraw", "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.Forbidden(w, r)
		return
	}

	ctx.Dataset.Status = "returned"

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(ctx.Loc, validationErrs.(validation.Errors)))
		render.Layout(w, "refresh_modal", "form_errors_dialog", struct {
			Title  string
			Errors form.Errors
		}{
			Title:  "Unable to withdraw this dataset due to the following errors",
			Errors: errors,
		})
		return
	}

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message:     ctx.Loc.Get("dataset.conflict_error"),
			RedirectURL: r.URL.Query().Get("redirect-url"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("withdraw dataset: could not save the dataset:", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	flash := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Dataset was successfully withdrawn.</p>"))

	h.AddFlash(r, w, *flash)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
