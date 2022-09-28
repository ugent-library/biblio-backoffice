package publicationediting

import (
	"errors"
	"html/template"
	"net/http"

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
	render.Layout(w, "show_modal", "publication/confirm_withdraw", YieldWithdraw{
		Context:     ctx,
		RedirectURL: r.URL.Query().Get("redirect-url"),
	})
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanWithdrawPublication(ctx.Publication) {
		h.Logger.Warnw("witdraw publication: user has no permission to withdraw", "user", ctx.User.ID, "publication", ctx.Publication.ID)
		render.Forbidden(w, r)
		return
	}

	ctx.Publication.Status = "returned"

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(ctx.Locale, validationErrs.(validation.Errors)))
		render.Layout(w, "refresh_modal", "form_errors_dialog", struct {
			Title  string
			Errors form.Errors
		}{
			Title:  "Unable to withdraw this publication due to the following errors",
			Errors: errors,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		h.Logger.Errorf("withdraw publication: could not save the publication:", "error", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	flash := flash.SimpleFlash().
		WithLevel("error").
		WithBody(template.HTML("<p>Publication was successfully witdrawn.</p>"))

	h.AddSessionFlash(r, w, *flash)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
