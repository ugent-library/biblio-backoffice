package datasetediting

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/internal/app/handlers"
	"github.com/ugent-library/biblio-backoffice/internal/app/localize"
	"github.com/ugent-library/biblio-backoffice/internal/render"
	"github.com/ugent-library/biblio-backoffice/internal/render/flash"
	"github.com/ugent-library/biblio-backoffice/internal/render/form"
	"github.com/ugent-library/biblio-backoffice/internal/snapstore"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
)

type YieldLock struct {
	Context
	RedirectURL string
}

func (h *Handler) Lock(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCurate() {
		h.Logger.Warnw("lock dataset: user has no permission to lock", "user", ctx.User.ID, "dataset", ctx.Dataset.ID)
		render.Forbidden(w, r)
		return
	}

	ctx.Dataset.Locked = true

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(ctx.Locale, validationErrs.(validation.Errors)))
		render.Layout(w, "show_modal", "form_errors_dialog", struct {
			Title  string
			Errors form.Errors
		}{
			Title:  "Unable to lock this dataset due to the following errors",
			Errors: errors,
		})
		return
	}

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Locale.T("dataset.conflict_error"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("lock dataset: could not save the dataset:", "error", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	flash := flash.SimpleFlash().
		WithLevel("error").
		WithBody(template.HTML("<p>Dataset was successfully locked.</p>"))

	h.AddSessionFlash(r, w, *flash)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}

func (h *Handler) Unlock(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCurate() {
		h.Logger.Warnw("unlock dataset: user has no permission to lock", "user", ctx.User.ID, "dataset", ctx.Dataset.ID)
		render.Forbidden(w, r)
		return
	}

	ctx.Dataset.Locked = false

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(ctx.Locale, validationErrs.(validation.Errors)))
		render.Layout(w, "show_modal", "form_errors_dialog", struct {
			Title  string
			Errors form.Errors
		}{
			Title:  "Unable to unlock this dataset due to the following errors",
			Errors: errors,
		})
		return
	}

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message:     ctx.Locale.T("dataset.conflict_error_reload"),
			RedirectURL: r.URL.Query().Get("redirect-url"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("unlock dataset: could not save the dataset:", "error", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	flash := flash.SimpleFlash().
		WithLevel("error").
		WithBody(template.HTML("<p>Dataset was successfully unlocked.</p>"))

	h.AddSessionFlash(r, w, *flash)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
