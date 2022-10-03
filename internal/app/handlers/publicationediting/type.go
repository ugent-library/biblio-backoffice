package publicationediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type BindType struct {
	Type string `form:"type"`
}

type YieldUpdateType struct {
	Context
	Type string
}

func (h *Handler) ConfirmUpdateType(w http.ResponseWriter, r *http.Request, ctx Context) {
	// TODO validate type
	render.Layout(w, "show_modal", "publication/confirm_update_type", YieldUpdateType{
		Context: ctx,
		Type:    r.URL.Query().Get("type"),
	})
}

func (h *Handler) UpdateType(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindType{}
	if err := bind.RequestForm(r, &b); err != nil {
		h.Logger.Warnw("update publication type: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	ctx.Publication.ChangeType(b.Type)

	// TODO This breaks converting the record type.
	//  e.g. moving from "conference" to "dissertation" is not possible. Dissertation requires a
	//  a "supervisor". This block will revert the record back to conference. The "supervisors" table
	//  won't become available and the user won't be able to set a supervisor, therefor being unable
	//  to satisfy the validation rule. (circular logic)
	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		form := detailsForm(ctx.User, ctx.Locale, ctx.Publication, validationErrs.(validation.Errors))

		render.Layout(w, "refresh_modal", "publication/edit_details", YieldEditDetails{
			Context: ctx,
			Form:    form,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Locale.T("publication.conflict_error_reload"),
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
