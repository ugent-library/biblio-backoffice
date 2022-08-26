package datasetediting

import (
	"errors"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/flash"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
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
		h.Logger.Warnw("delete dataset: user is unauthorized", "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.Forbidden(w, r)
		return
	}

	ctx.Dataset.Status = "deleted"

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		h.Logger.Warnf("delete dataset: snapstore detected a conflicting dataset:", "errors", errors.As(err, &conflict), "identifier", ctx.Dataset.ID)
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("dataset.conflict_error"))
		return
	}

	if err != nil {
		h.Logger.Errorf("delete dataset: Could not save the dataset:", "error", err, "identifier", ctx.Dataset.ID)
		render.InternalServerError(w, r, err)
		return
	}

	h.AddSessionFlash(r, w, flash.Flash{
		Type:         "success",
		Body:         "Dataset was succesfully deleted",
		DismissAfter: 5 * time.Second,
	})

	w.Header().Set("HX-Redirect", r.URL.Query().Get("redirect-url"))
}
