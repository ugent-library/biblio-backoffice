package publicationediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	"github.com/ugent-library/biblio-backoffice/views/publication"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

type BindType struct {
	Type string `form:"type"`
}

type YieldUpdateType struct {
	Context
	Type string
}

func ConfirmUpdateType(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	// TODO validate type
	publication.ConfirmUpdateType(c, ctx.GetPublication(r), r.URL.Query().Get("type")).Render(r.Context(), w)
}

func UpdateType(w http.ResponseWriter, r *http.Request, legacyContext Context) {
	c := ctx.Get(r)

	b := BindType{}
	if err := bind.Body(r, &b); err != nil {
		c.Log.Warnw("update publication type: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	legacyContext.Publication.ChangeType(b.Type)

	if validationErrs := legacyContext.Publication.Validate(); validationErrs != nil {
		form := detailsForm(c.User, c.Loc, legacyContext.Publication, validationErrs.(*okay.Errors))

		// TODO: refactor to templ
		render.Layout(w, "refresh_modal", "publication/edit_details", YieldEditDetails{
			Context: legacyContext,
			Form:    form,
		})
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), legacyContext.Publication, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("update publication type: Could not save the publication:", "error", err, "publication", legacyContext.Publication.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	redirectURL := c.PathTo("publication", "id", legacyContext.Publication.ID)
	w.Header().Set("HX-Redirect", redirectURL.String())
}
