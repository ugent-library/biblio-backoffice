package datasetcreating

import (
	"errors"
	"fmt"
	"html/template"
	"log"
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

type YieldAddDataset struct {
	Context
	PageTitle           string
	Step                int
	Source              string
	Identifier          string
	Dataset             *models.Dataset
	DuplicateDataset    bool
	DatasetPublications []*models.Publication
	ActiveNav           string
	ActiveSubNav        string
	DisplayDetails      *display.Display
	Errors              *YieldValidationErrors
}

type YieldValidationErrors struct {
	Title  string
	Errors form.Errors
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Wrap(w, "layouts/default", "dataset/add", YieldAddDataset{
		Context:   ctx,
		PageTitle: "Add - Datasets - Biblio",
		Step:      1,
		ActiveNav: "datasets",
		Errors:    nil,
	})
}

func (h *Handler) ConfirmImportDataset(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindImport{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	// Check for duplicates
	if b.Source == "datacite" {
		args := models.NewSearchArgs().
			WithFilter("doi", b.Identifier)
			// WithFilter("status", "public")

		if existing, _ := h.DatasetSearchService.Search(args); existing.Total > 0 {
			render.Wrap(w, "layouts/default", "dataset/add", YieldAddDataset{
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
		render.BadRequest(w, r, err)
		return
	}

	d, err := h.FetchDatasetByIdentifier(b.Source, b.Identifier)
	if err != nil {
		log.Println(err)
		flash := flash.Flash{
			Type:         "error",
			Body:         template.HTML(ctx.TS("dataset.single_import", "import_by_id.import_failed")),
			DismissAfter: 5 * time.Second,
		}

		ctx.Flash = append(ctx.Flash, flash)

		render.Wrap(w, "layouts/default", "dataset/add", YieldAddDataset{
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

	// TODO This block implements validation on import. This implies that a record imported
	//   from a source (e.g. DataCite) should at least have the minimum set of fields as
	//   required by Biblio. There's no guarantee that this is the case, though. If a record
	//   lacks crucial Biblio fields in DataCite, it's impossible to import in Biblio until fixed in
	//   DataCite.
	if validationErrs := d.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(ctx.Locale, err.(validation.Errors)))
		render.Wrap(w, "layouts/default", "dataset/add", YieldAddDataset{
			Context:    ctx,
			PageTitle:  "Add - Datasets - Biblio",
			Step:       1,
			ActiveNav:  "datasets",
			Source:     b.Source,
			Identifier: b.Identifier,
			Errors: &YieldValidationErrors{
				Title:  "Unable to publish this dataset due to the following errors",
				Errors: errors,
			},
		})
		return
	}
	err = h.Repository.SaveDataset(d)

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Wrap(w, "layouts/default", "dataset/add_description", YieldAddDataset{
		Context:        ctx,
		PageTitle:      "Add - Datasets - Biblio",
		Step:           2,
		ActiveNav:      "datasets",
		ActiveSubNav:   "description",
		Dataset:        d,
		DisplayDetails: displays.DatasetDetails(ctx.Locale, d),
	})
}

func (h *Handler) AddDescription(w http.ResponseWriter, r *http.Request, ctx Context) {
	// TODO bind and validate
	activeSubNav := "description"
	if r.URL.Query().Get("show") == "contributors" || r.URL.Query().Get("show") == "description" {
		activeSubNav = r.URL.Query().Get("show")
	}

	render.Wrap(w, "layouts/default", "dataset/add_description", YieldAddDataset{
		Context:        ctx,
		PageTitle:      "Add - Datasets - Biblio",
		Step:           2,
		ActiveNav:      "datasets",
		ActiveSubNav:   activeSubNav,
		Dataset:        ctx.Dataset,
		DisplayDetails: displays.DatasetDetails(ctx.Locale, ctx.Dataset),
	})
}

func (h *Handler) AddConfirm(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Wrap(w, "layouts/default", "dataset/add_confirm_import", YieldAddDataset{
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

	redirectURL := h.PathFor("dataset_add_finish", "id", ctx.Dataset.ID)
	redirectURL.RawQuery = r.URL.Query().Encode()

	w.Header().Set("HX-Redirect", redirectURL.String())
}

func (h *Handler) AddFinish(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Wrap(w, "layouts/default", "dataset/add_finish", YieldAddDataset{
		Context:   ctx,
		PageTitle: "Add - Datasets - Biblio",
		Step:      4,
		ActiveNav: "datasets",
		Dataset:   ctx.Dataset,
	})
}

func (h *Handler) FetchDatasetByIdentifier(source, identifier string) (*models.Dataset, error) {
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
