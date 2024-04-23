package publicationediting

import (
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/flash"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
)

type YieldConfirmDelete struct {
	Context
	Publication *models.Publication
	RedirectURL string
}

func (h *Handler) ConfirmDelete(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to delete this publication?",
		DeleteUrl:  c.PathTo("publication_delete", "id", publication.ID, "redirect-url", r.URL.Query().Get("redirect-url")),
		SnapshotID: publication.SnapshotID,
	}).Render(r.Context(), w)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	if !c.User.CanDeletePublication(publication) {
		h.Logger.Warnw("delete publication: user is unauthorized", "publication", publication.ID, "user", c.User.ID)
		render.Forbidden(w, r)
		return
	}

	publication.Status = "deleted"

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), publication, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		// TODO: refactor to templ
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message:     c.Loc.Get("publication.conflict_error_reload"),
			RedirectURL: r.URL.Query().Get("redirect-url"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("delete publication: Could not save the publication:", "errors", err, "publication", publication.ID, "user", c.User.ID)
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
