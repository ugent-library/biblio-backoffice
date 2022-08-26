package publicationediting

import (
	"errors"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/flash"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type YieldPublish struct {
	Context
	RedirectURL string
}

func (h *Handler) ConfirmPublish(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "show_modal", "publication/confirm_publish", YieldPublish{
		Context:     ctx,
		RedirectURL: r.URL.Query().Get("redirect-url"),
	})
}

func (h *Handler) Publish(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanPublishPublication(ctx.Publication) {
		h.Logger.Warnw("publish dataset: user has no permission to publish", "user", ctx.User.ID, "dataset", ctx.Publication.ID)
		render.Forbidden(w, r)
		return
	}

	ctx.Publication.Status = "public"

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		h.Logger.Warnw("publish dataset: could not validate dataset:", "errors", validationErrs, "identifier", ctx.Publication.ID)
		errors := form.Errors(localize.ValidationErrors(ctx.Locale, validationErrs.(validation.Errors)))
		render.Layout(w, "refresh_modal", "form_errors_dialog", struct {
			Title  string
			Errors form.Errors
		}{
			Title:  "Unable to publish this publication due to the following errors",
			Errors: errors,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		h.Logger.Warnf("publish dataset: snapstore detected a conflicting dataset:", "errors", errors.As(err, &conflict), "identifier", ctx.Publication.ID)
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		h.Logger.Errorf("publish dataset: could not save the dataset:", "error", err, "identifier", ctx.Publication.ID)
		render.InternalServerError(w, r, err)
		return
	}

	h.AddSessionFlash(r, w, flash.Flash{
		Type:         "success",
		Body:         "Publication was succesfully published",
		DismissAfter: 5 * time.Second,
	})

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
