package datasetcreating

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backend/internal/app/displays"
	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/render/flash"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/ulid"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type BindImport struct {
	Source     string `form:"source"`
	Identifier string `form:"identifier"`
}

type YieldAdd struct {
	Context
	PageTitle           string
	Step                int
	Source              string
	Identifier          string
	Dataset             *models.Dataset
	DuplicateDataset    bool
	DatasetPublications []*models.Publication
	ActiveNav           string
	SubNavs             []string // needed to render show_description
	ActiveSubNav        string   // needed to render show_description
	RedirectURL         string   // needed to render show_description
	DisplayDetails      *display.Display
	Errors              *YieldValidationErrors
}

type YieldValidationErrors struct {
	Title  string
	Errors form.Errors
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "layouts/default", "dataset/pages/add", YieldAdd{
		Context:   ctx,
		PageTitle: "Add - Datasets - Biblio",
		Step:      1,
		ActiveNav: "datasets",
	})
}

func (h *Handler) ConfirmImport(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindImport{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("create dataset: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	// check for duplicates
	if b.Source == "datacite" {
		args := models.NewSearchArgs().WithFilter("doi", b.Identifier)

		existing, err := h.DatasetSearchService.Search(args)

		if err != nil {
			h.Logger.Errorw("create dataset: could not execute search", "error", err)
			render.InternalServerError(w, r, err)
			return
		}

		if existing.Total > 0 {
			render.Layout(w, "layouts/default", "dataset/pages/add", YieldAdd{
				Context:          ctx,
				PageTitle:        "Add - Datasets - Biblio",
				Step:             1,
				ActiveNav:        "datasets",
				Source:           b.Source,
				Identifier:       b.Identifier,
				Dataset:          existing.Hits[0],
				DuplicateDataset: true,
			})
			return
		}
	}

	h.AddImport(w, r, ctx)
}

func (h *Handler) AddImport(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindImport{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("create dataset: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	d, err := h.fetchDatasetByIdentifier(b.Source, b.Identifier)
	if err != nil {
		h.Logger.Warnw("create dataset: could not fetch dataset:", "error", err, "identifier", b.Identifier)
		flash := flash.Flash{
			Type:         "error",
			Body:         template.HTML(ctx.Locale.TS("dataset.single_import", "import_by_id.import_failed")),
			DismissAfter: 5 * time.Second,
		}

		ctx.Flash = append(ctx.Flash, flash)

		render.Layout(w, "layouts/default", "dataset/pages/add", YieldAdd{
			Context:    ctx,
			PageTitle:  "Add - Datasets - Biblio",
			Step:       1,
			ActiveNav:  "datasets",
			Source:     b.Source,
			Identifier: b.Identifier,
		})
		return
	}

	d.ID = ulid.MustGenerate()
	d.CreatorID = ctx.User.ID
	d.UserID = ctx.User.ID
	d.Status = "private"

	if validationErrs := d.Validate(); validationErrs != nil {
		h.Logger.Warnw("create dataset: could not validate dataset:", "errors", validationErrs, "identifier", b.Identifier)
		errors := form.Errors(localize.ValidationErrors(ctx.Locale, validationErrs.(validation.Errors)))
		render.Layout(w, "layouts/default", "dataset/pages/add", YieldAdd{
			Context:    ctx,
			PageTitle:  "Add - Datasets - Biblio",
			Step:       1,
			ActiveNav:  "datasets",
			Source:     b.Source,
			Identifier: b.Identifier,
			Errors: &YieldValidationErrors{
				Title:  "Unable to import this dataset due to the following errors",
				Errors: errors,
			},
		})
		return
	}

	err = h.Repository.SaveDataset(d)

	if err != nil {
		h.Logger.Warnw("create dataset: could not save dataset:", "error", err, "identifier", b.Identifier)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "layouts/default", "dataset/pages/add_description", YieldAdd{
		Context:        ctx,
		PageTitle:      "Add - Datasets - Biblio",
		Step:           2,
		ActiveNav:      "datasets",
		SubNavs:        []string{"description", "contributors", "publications"},
		ActiveSubNav:   "description",
		Dataset:        d,
		DisplayDetails: displays.DatasetDetails(ctx.Locale, d),
	})
}

func (h *Handler) AddDescription(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "layouts/default", "dataset/pages/add_description", YieldAdd{
		Context:        ctx,
		PageTitle:      "Add - Datasets - Biblio",
		Step:           2,
		ActiveNav:      "datasets",
		SubNavs:        []string{"description", "contributors", "publications"},
		ActiveSubNav:   "description",
		Dataset:        ctx.Dataset,
		DisplayDetails: displays.DatasetDetails(ctx.Locale, ctx.Dataset),
	})
}

func (h *Handler) AddConfirm(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "layouts/default", "dataset/pages/add_confirm", YieldAdd{
		Context:   ctx,
		PageTitle: "Add - Datasets - Biblio",
		Step:      3,
		ActiveNav: "datasets",
		Dataset:   ctx.Dataset,
	})
}

func (h *Handler) AddPublish(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanPublishDataset(ctx.Dataset) {
		render.Forbidden(w, r)
		return
	}

	ctx.Dataset.Status = "public"

	if err := ctx.Dataset.Validate(); err != nil {
		h.Logger.Warnw("create dataset: could not validate dataset:", "errors", err, "identifier", ctx.Dataset.ID)
		errors := form.Errors(localize.ValidationErrors(ctx.Locale, err.(validation.Errors)))
		render.Layout(w, "show_modal", "form_errors_dialog", struct {
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
		h.Logger.Warnf("create dataset: snapstore detected a conflicting dataset:", "errors", errors.As(err, &conflict), "identifier", ctx.Dataset.ID)
		render.Layout(w, "show_modal", "error_dialog", ctx.Locale.T("dataset.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		h.Logger.Warnf("create dataset: Could not save the dataset:", "error", err, "identifier", ctx.Dataset.ID)
		return
	}

	redirectURL := h.PathFor("dataset_add_finish", "id", ctx.Dataset.ID)
	redirectURL.RawQuery = r.URL.Query().Encode()

	w.Header().Set("HX-Redirect", redirectURL.String())
}

func (h *Handler) AddFinish(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "layouts/default", "dataset/pages/add_finish", YieldAdd{
		Context:   ctx,
		PageTitle: "Add - Datasets - Biblio",
		Step:      4,
		ActiveNav: "datasets",
		Dataset:   ctx.Dataset,
	})
}

func (h *Handler) fetchDatasetByIdentifier(source, identifier string) (*models.Dataset, error) {
	s, ok := h.DatasetSources[source]

	if !ok {
		return nil, fmt.Errorf("unkown dataset source: %s", source)
	}

	d, err := s.GetDataset(identifier)
	if err != nil {
		return nil, err
	}

	return d, nil
}
