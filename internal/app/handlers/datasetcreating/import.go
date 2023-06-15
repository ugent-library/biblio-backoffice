package datasetcreating

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/oklog/ulid/v2"
	"github.com/ugent-library/biblio-backoffice/internal/app/displays"
	"github.com/ugent-library/biblio-backoffice/internal/app/handlers"
	"github.com/ugent-library/biblio-backoffice/internal/app/localize"
	"github.com/ugent-library/biblio-backoffice/internal/bind"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/render"
	"github.com/ugent-library/biblio-backoffice/internal/render/display"
	"github.com/ugent-library/biblio-backoffice/internal/render/flash"
	"github.com/ugent-library/biblio-backoffice/internal/render/form"
	"github.com/ugent-library/biblio-backoffice/internal/snapstore"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
)

type BindAdd struct {
	Method string `form:"method"`
}

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
	b := BindAdd{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("add dataset: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	tmpl := ""
	switch b.Method {
	case "identifier":
		tmpl = "dataset/pages/add_identifier"
	case "manual":
		h.ConfirmImport(w, r, ctx)
		return
	default:
		tmpl = "dataset/pages/add"
	}

	render.Layout(w, "layouts/default", tmpl, YieldAdd{
		Context:   ctx,
		PageTitle: "Add - Datasets - Biblio",
		Step:      1,
		ActiveNav: "datasets",
	})
}

func (h *Handler) ConfirmImport(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindImport{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("confirm import dataset: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	// check for duplicates
	if b.Source == "datacite" {
		args := models.NewSearchArgs().WithFilter("identifier_values", strings.ToLower(b.Identifier)).WithFilter("status", "public")

		existing, err := h.DatasetSearchIndex.Search(args)

		if err != nil {
			h.Logger.Errorw("confirm import dataset: could not execute search", "errors", err, "user", ctx.User.ID)
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
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("add import dataset: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	var (
		d   *models.Dataset
		err error
	)

	if b.Identifier != "" {
		d, err = h.fetchDatasetByIdentifier(b.Source, b.Identifier)
		if err != nil {
			flash := flash.SimpleFlash().
				WithLevel("error").
				WithTitle("Failed to save draft").
				WithBody(template.HTML(ctx.Locale.TS("dataset.single_import", "import_by_id.import_failed")))

			ctx.Flash = append(ctx.Flash, *flash)

			render.Layout(w, "layouts/default", "dataset/pages/add_identifier", YieldAdd{
				Context:    ctx,
				PageTitle:  "Add - Datasets - Biblio",
				Step:       1,
				ActiveNav:  "datasets",
				Source:     b.Source,
				Identifier: b.Identifier,
			})
			return
		}
	} else {
		// or start with empty dataset
		d = &models.Dataset{}
	}

	d.ID = ulid.Make().String()
	d.CreatorID = ctx.User.ID
	d.Creator = &ctx.User.Person
	d.UserID = ctx.User.ID
	d.User = &ctx.User.Person
	d.Status = "private"

	if len(ctx.User.Affiliations) > 0 {
		d.AddOrganization(ctx.User.Affiliations[0].Organization)
	}

	if validationErrs := d.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(ctx.Locale, validationErrs.(validation.Errors)))
		render.Layout(w, "layouts/default", "dataset/pages/add_identifier", YieldAdd{
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

	err = h.Repository.SaveDataset(d, ctx.User)

	if err != nil {
		h.Logger.Warnw("add import dataset: could not save dataset:", "errors", err, "dataset", b.Identifier, "user", ctx.User.ID)
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
		DisplayDetails: displays.DatasetDetails(ctx.User, ctx.Locale, d),
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
		DisplayDetails: displays.DatasetDetails(ctx.User, ctx.Locale, ctx.Dataset),
	})
}

func (h *Handler) AddSaveDraft(w http.ResponseWriter, r *http.Request, ctx Context) {
	flash := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Dataset successfully saved as a draft.</p>"))

	h.AddFlash(r, w, *flash)

	redirectURL := h.PathFor("datasets")
	w.Header().Set("HX-Redirect", redirectURL.String())
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
	if !ctx.User.CanEditDataset(ctx.Dataset) {
		h.Logger.Warn("publish dataset: user isn't allowed to edit the dataset:", "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.Forbidden(w, r)
		return
	}

	ctx.Dataset.Status = "public"

	if err := ctx.Dataset.Validate(); err != nil {
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

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Locale.T("dataset.conflict_error_reload"),
		})
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
