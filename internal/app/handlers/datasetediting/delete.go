package datasetediting

import (
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backoffice/internal/app/handlers"
	"github.com/ugent-library/biblio-backoffice/internal/render"
	"github.com/ugent-library/biblio-backoffice/internal/render/flash"
	"github.com/ugent-library/biblio-backoffice/internal/snapstore"
	"github.com/ugent-library/biblio-backoffice/models"
)

type YieldConfirmDelete struct {
	Context
	Dataset     *models.Dataset
	RedirectURL string
}

func (h *Handler) ConfirmDelete(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "show_modal", "dataset/confirm_delete", YieldConfirmDelete{
		Context:     ctx,
		Dataset:     ctx.Dataset,
		RedirectURL: r.URL.Query().Get("redirect-url"),
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanDeleteDataset(ctx.Dataset) {
		h.Logger.Warnw("delete dataset: user isn't allowed to delete dataset", "dataset", ctx.Dataset.ID, "user", ctx.User.ID, "user", ctx.User.ID)
		render.Forbidden(w, r)
		return
	}

	ctx.Dataset.Status = "deleted"

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Locale.T("dataset.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("delete dataset: Could not save the dataset:", "error", err, "identifier", ctx.Dataset.ID, "user", ctx.User.ID)
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
