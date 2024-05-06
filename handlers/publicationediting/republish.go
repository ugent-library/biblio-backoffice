package publicationediting

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/flash"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	"github.com/ugent-library/biblio-backoffice/views/publication"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

func ConfirmRepublish(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	redirectUrl := r.URL.Query().Get("redirect-url")

	publication.ConfirmRepublish(c, ctx.GetPublication(r), redirectUrl).Render(r.Context(), w)
}

func Republish(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	if !c.User.CanPublishPublication(publication) {
		c.Log.Warnw("republish publication: user has no permission to republish", "user", c.User.ID, "publication", publication.ID)
		c.HandleError(w, r, httperror.Forbidden)
		return
	}

	publication.Status = "public"

	if validationErrs := publication.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(c.Loc, validationErrs.(*okay.Errors)))
		render.Layout(w, "refresh_modal", "form_errors_dialog", struct {
			Title  string
			Errors form.Errors
		}{
			Title:  "Unable to republish this publication due to the following errors",
			Errors: errors,
		})
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), publication, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(
			views.ErrorDialogWithOptions(c.Loc.Get("publication.conflict_error"), views.ErrorDialogOptions{
				RedirectURL: r.URL.Query().Get("redirect-url"),
			})).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("republish publication: could not save the publication:", "error", err, "publication", publication.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	flash := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Publication was successfully republished.</p>"))

	c.PersistFlash(w, *flash)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
