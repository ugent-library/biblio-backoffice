package publicationediting

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
	"github.com/ugent-library/biblio-backoffice/views/publication"
	"github.com/ugent-library/okay"
)

func (h *Handler) ConfirmWithdraw(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	redirectUrl := r.URL.Query().Get("redirect-url")

	publication.ConfirmWithdraw(c, ctx.GetPublication(r), redirectUrl).Render(r.Context(), w)
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	if !c.User.CanWithdrawPublication(publication) {
		h.Logger.Warnw("witdraw publication: user has no permission to withdraw", "user", c.User.ID, "publication", publication.ID)
		render.Forbidden(w, r)
		return
	}

	publication.Status = "returned"

	if validationErrs := publication.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(c.Loc, validationErrs.(*okay.Errors)))
		render.Layout(w, "refresh_modal", "form_errors_dialog", struct {
			Title  string
			Errors form.Errors
		}{
			Title:  "Unable to withdraw this publication due to the following errors",
			Errors: errors,
		})
		return
	}

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), publication, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message:     c.Loc.Get("publication.conflict_error"),
			RedirectURL: r.URL.Query().Get("redirect-url"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("withdraw publication: could not save the publication:", "error", err, "publication", publication.ID, "user", c.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	flash := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Publication was successfully withdrawn.</p>"))

	h.AddFlash(r, w, *flash)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
