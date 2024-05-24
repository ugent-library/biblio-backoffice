package publicationediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	publicationviews "github.com/ugent-library/biblio-backoffice/views/publication"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

type BindType struct {
	Type string `form:"type"`
}

func ConfirmUpdateType(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	// TODO validate type
	publicationviews.ConfirmUpdateType(c, ctx.GetPublication(r), r.URL.Query().Get("type")).Render(r.Context(), w)
}

func UpdateType(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindType{}
	if err := bind.Body(r, &b); err != nil {
		c.Log.Warnw("update publication type: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	p.ChangeType(b.Type)

	if validationErrs := p.Validate(); validationErrs != nil {
		views.ReplaceModal(publicationviews.EditDetailsDialog(c, p, false, validationErrs.(*okay.Errors))).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("update publication type: Could not save the publication:", "error", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	redirectURL := c.PathTo("publication", "id", p.ID)
	w.Header().Set("HX-Redirect", redirectURL.String())
}
