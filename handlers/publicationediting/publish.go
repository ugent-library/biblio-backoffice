package publicationediting

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
	if !ctx.User.CanEditPublication(ctx.Publication) {
		h.Logger.Warnw("publish dataset: user has no permission to publish", "user", ctx.User.ID, "publication", ctx.Publication.ID)
		render.Forbidden(w, r)
		return
	}

	ctx.Publication.Status = "public"

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
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

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message:     ctx.Locale.T("publication.conflict_error"),
			RedirectURL: r.URL.Query().Get("redirect-url"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("publish dataset: could not save the dataset:", "error", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	flash := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Publication was successfully published.</p>"))

	h.AddFlash(r, w, *flash)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
