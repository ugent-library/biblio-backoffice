package publicationediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	"github.com/ugent-library/biblio-backoffice/views/flash"
	"github.com/ugent-library/biblio-backoffice/views/publication"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

func ConfirmPublish(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	redirectUrl := r.URL.Query().Get("redirect-url")

	views.ShowModal(publication.ConfirmPublish(c, ctx.GetPublication(r), redirectUrl)).Render(r.Context(), w)
}

func Publish(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	if !c.User.CanEditPublication(publication) {
		c.HandleError(w, r, httperror.Forbidden)
		return
	}

	publication.Status = "public"

	if validationErrs := publication.Validate(); validationErrs != nil {
		errors := localize.ValidationErrors(c.Loc, validationErrs.(*okay.Errors))
		views.ReplaceModal(views.FormErrorsDialog("Unable to publish this publication due to the following errors", errors)).Render(r.Context(), w)
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
		c.HandleError(w, r, err)
		return
	}

	flash := flash.SimpleFlash().
		WithLevel("success").
		WithBody("<p>Publication was successfully published.</p>")

	c.PersistFlash(w, *flash)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
