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
	Form *form.Form
}

func (h *Handler) Publish(w http.ResponseWriter, r *http.Request, ctx Context) {
	ctx.Dataset.Status = "public"

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		render.Render(w, "dataset/refresh_publish_dataset", YieldPublishDataset{
			Context: ctx,
			Form:    publishForm(ctx, validationErrs.(validation.Errors)),
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

	h.SetSessionFlash(r, w, flash.Flash{Type: "success", Body: "Dataset was succesfully published", DismissAfter: 5 * time.Second})

	destUrl := h.PathFor("dataset", "id", ctx.Dataset.ID)
	destUrl.RawQuery = r.URL.Query().Encode()

	w.Header().Set("HX-Redirect", destUrl.String())
}

func publishForm(ctx Context, errors validation.Errors) *form.Form {
	return form.New().WithTheme("default").WithErrors(localize.ValidationErrors(ctx.Locale, errors))
}
