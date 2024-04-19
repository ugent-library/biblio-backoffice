package datasetediting

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
	"github.com/ugent-library/biblio-backoffice/views/dataset"
)

type YieldConfirmDelete struct {
	Context
	Dataset     *models.Dataset
	RedirectURL string
}

func (h *Handler) ConfirmDelete(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	dataset.ConfirmDelete(c, ctx.GetDataset(r), r.URL.Query().Get("redirect-url")).Render(r.Context(), w)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	if !c.User.CanDeleteDataset(dataset) {
		h.Logger.Warnw("delete dataset: user isn't allowed to delete dataset", "dataset", dataset.ID, "user", c.User.ID, "user", c.User.ID)
		render.Forbidden(w, r)
		return
	}

	dataset.Status = "deleted"

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), dataset, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		// TODO: refactor to templ
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: c.Loc.Get("dataset.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("delete dataset: Could not save the dataset:", "error", err, "identifier", dataset.ID, "user", c.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	flash := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Dataset was successfully deleted.</p>"))

	h.AddFlash(r, w, *flash)

	// TODO temporary fix until we can figure out a way let ES notify this handler that it did its thing.
	// see: https://github.com/ugent-library/biblio-backoffice/issues/590
	time.Sleep(1250 * time.Millisecond)

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
