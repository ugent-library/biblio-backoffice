package publicationediting

import (
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/flash"
	"github.com/ugent-library/biblio-backoffice/snapstore"
)

type YieldConfirmDelete struct {
	Context
	Publication *models.Publication
	RedirectURL string
}

func (h *Handler) ConfirmDelete(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "show_modal", "publication/confirm_delete", YieldConfirmDelete{
		Context:     ctx,
		Publication: ctx.Publication,
		RedirectURL: r.URL.Query().Get("redirect-url"),
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanDeletePublication(ctx.Publication) {
		h.Logger.Warnw("delete publication: user is unauthorized", "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.Forbidden(w, r)
		return
	}

	ctx.Publication.Status = "deleted"

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message:     ctx.Loc.Get("publication.conflict_error_reload"),
			RedirectURL: r.URL.Query().Get("redirect-url"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("delete publication: Could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	flash := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Publication was successfully deleted.</p>"))

	h.AddFlash(r, w, *flash)

	// TODO temporary fix until we can figure out a way let ES notify this handler that it did its thing.
	// see: https://github.com/ugent-library/biblio-backoffice/issues/590
	time.Sleep(1250 * time.Millisecond)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
