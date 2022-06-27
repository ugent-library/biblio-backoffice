package datasetediting

import (
	"errors"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/localize"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/flash"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type YieldPublishDataset struct {
	Context
	SearchArgs *models.SearchArgs
}

func (h *Handler) ConfirmPublishDataset(w http.ResponseWriter, r *http.Request, ctx Context) {
	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.Render(w, "dataset/confirm_publish_dataset", YieldPublishDataset{
		Context:    ctx,
		SearchArgs: searchArgs,
	})
}

func (h *Handler) Publish(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanPublishDataset(ctx.Dataset) {
		render.Forbidden(w, r)
		return
	}

	ctx.Dataset.Status = "public"

	if err := ctx.Dataset.Validate(); err != nil {
		errors := form.Errors(localize.ValidationErrors(ctx.Locale, err.(validation.Errors)))
		render.Render(w, "form_errors_dialog", struct {
			Title  string
			Errors form.Errors
		}{
			Title:  "Unable to publish this dataset due to the following errors",
			Errors: errors,
		})
		return
	}

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
		Body:         "Dataset was succesfully published",
		DismissAfter: 5 * time.Second,
	})

	redirectURL := h.PathFor("dataset", "id", ctx.Dataset.ID)
	redirectURL.RawQuery = r.URL.Query().Encode()

	w.Header().Set("HX-Redirect", redirectURL.String())
}
