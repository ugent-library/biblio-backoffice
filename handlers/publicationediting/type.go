package publicationediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views/publication"
	"github.com/ugent-library/bind"
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

func (h *Handler) UpdateType(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindType{}
	if err := bind.Body(r, &b); err != nil {
		h.Logger.Warnw("update publication type: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	ctx.Publication.ChangeType(b.Type)

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		form := detailsForm(ctx.User, ctx.Loc, ctx.Publication, validationErrs.(*okay.Errors))

		// TODO: refactor to templ
		render.Layout(w, "refresh_modal", "publication/edit_details", YieldEditDetails{
			Context: ctx,
			Form:    form,
		})
		return
	}

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		// TODO: refactor to templ
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("publication.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("update publication type: Could not save the publication:", "error", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	redirectURL := h.PathFor("publication", "id", ctx.Publication.ID)
	w.Header().Set("HX-Redirect", redirectURL.String())
}
