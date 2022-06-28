package datasetediting

import (
	"errors"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/flash"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
)

type YieldConfirmDelete struct {
	Context
	Dataset    *models.Dataset
	SearchArgs *models.SearchArgs
}

func (h *Handler) ConfirmDelete(w http.ResponseWriter, r *http.Request, ctx Context) {
	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.Render(w, "dataset/confirm_delete", YieldConfirmDelete{
		Context:    ctx,
		Dataset:    ctx.Dataset,
		SearchArgs: searchArgs,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanDeleteDataset(ctx.Dataset) {
		render.Forbidden(w, r)
		return
	}

	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	ctx.Dataset.Status = "deleted"

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Render(w, "error_dialog", ctx.T("dataset.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	h.AddSessionFlash(r, w, flash.Flash{
		Type:         "success",
		Body:         "Dataset was succesfully deleted",
		DismissAfter: 5 * time.Second,
	})

	redirectURL := h.PathFor("datasets")
	redirectURL.RawQuery = r.URL.RawQuery
	w.Header().Set("HX-Redirect", redirectURL.String())
}
