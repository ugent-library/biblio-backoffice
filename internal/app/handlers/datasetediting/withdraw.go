package datasetediting

import (
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/flash"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
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
	if !ctx.User.CanWithdrawDataset(ctx.Dataset) {
		h.Logger.Warnw("withdraw dataset: user has no permission to withdraw", "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.Forbidden(w, r)
		return
	}

	ctx.Dataset.Status = "returned"

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(ctx.Locale, validationErrs.(validation.Errors)))
		render.Layout(w, "refresh_modal", "form_errors_dialog", struct {
			Title  string
			Errors form.Errors
		}{
			Title:  "Unable to withdraw this dataset due to the following errors",
			Errors: errors,
		})
		return
	}

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("dataset.conflict_error"))
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

	h.AddSessionFlash(r, w, *flash)

	// TODO temporary fix until we can figure out a way let ES notify this handler that it did its thing.
	// see: https://github.com/ugent-library/biblio-backend/issues/590
	time.Sleep(1250 * time.Millisecond)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
