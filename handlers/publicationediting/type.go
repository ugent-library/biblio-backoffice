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

func (h *Handler) ConfirmUpdateType(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	// TODO validate type
	publication.ConfirmUpdateType(c, ctx.GetPublication(r), r.URL.Query().Get("type")).Render(r.Context(), w)
}

func (h *Handler) UpdateType(w http.ResponseWriter, r *http.Request, legacyContext Context) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	b := BindType{}
	if err := bind.Body(r, &b); err != nil {
		h.Logger.Warnw("update publication type: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	publication.ChangeType(b.Type)

	// TODO This breaks converting the record type.
	//  e.g. moving from "conference" to "dissertation" is not possible. Dissertation requires a
	//  a "supervisor". This block will revert the record back to conference. The "supervisors" table
	//  won't become available and the user won't be able to set a supervisor, therefor being unable
	//  to satisfy the validation rule. (circular logic)
	if validationErrs := publication.Validate(); validationErrs != nil {
		form := detailsForm(c.User, c.Loc, publication, validationErrs.(*okay.Errors))

		// TODO: refactor to templ
		render.Layout(w, "refresh_modal", "publication/edit_details", YieldEditDetails{
			Context: legacyContext,
			Form:    form,
		})
		return
	}

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), publication, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		// TODO: refactor to templ
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: c.Loc.Get("publication.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("update publication type: Could not save the publication:", "error", err, "publication", publication.ID, "user", c.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	redirectURL := h.PathFor("publication", "id", publication.ID)
	w.Header().Set("HX-Redirect", redirectURL.String())
}
