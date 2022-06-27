package datasetediting

import (
	"errors"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backend/internal/localize"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/flash"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type YieldPublishDataset struct {
	Context
	Errors form.Errors
}

func (h *Handler) Publish(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanPublishDataset(ctx.Dataset) {
		render.Forbidden(w, r)
		return
	}

	ctx.Dataset.Status = "public"

	if err := ctx.Dataset.Validate(); err != nil {
		errors := form.Errors(localize.ValidationErrors(ctx.Locale, err.(validation.Errors)))
		render.Render(w, "dataset/refresh_publish_dataset", YieldPublishDataset{
			Context: ctx,
			Errors:  errors,
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

	h.AddSessionFlash(r, w, flash.Flash{Type: "success", Body: "Dataset was succesfully published", DismissAfter: 5 * time.Second})

	destUrl := h.PathFor("dataset", "id", ctx.Dataset.ID)
	destUrl.RawQuery = r.URL.Query().Encode()

	w.Header().Set("HX-Redirect", destUrl.String())
}
