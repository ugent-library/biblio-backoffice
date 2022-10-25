package datasetediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/displays"
	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type BindDetails struct {
	AccessLevel             string   `form:"access_level"`
	EmbargoDate             string   `form:"embargo_date"`
	AccessLevelAfterEmbargo string   `form:"access_level_after_embargo"`
	Format                  []string `form:"format"`
	Keyword                 []string `form:"keyword"`
	License                 string   `form:"license"`
	OtherLicense            string   `form:"other_license"`
	Publisher               string   `form:"publisher"`
	Title                   string   `form:"title"`
	URL                     string   `form:"url"`
	Year                    string   `form:"year"`
}

type YieldDetails struct {
	Context
	DisplayDetails *display.Display
}

type YieldEditDetails struct {
	Context
	Form     *form.Form
	Conflict bool
}

func (h *Handler) EditDetails(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "show_modal", "dataset/edit_details", YieldEditDetails{
		Context: ctx,
		Form:    detailsForm(ctx.Locale, ctx.Dataset, nil),
	})
}

func (h *Handler) RefreshEditFileForm(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDetails{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update dataset details: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	if b.AccessLevel != "info:eu-repo/semantics/embargoedAccess" {
		b.EmbargoDate = ""
		b.AccessLevelAfterEmbargo = ""
	}

	ctx.Dataset.AccessLevel = b.AccessLevel
	ctx.Dataset.EmbargoDate = b.EmbargoDate
	ctx.Dataset.AccessLevelAfterEmbargo = b.AccessLevelAfterEmbargo
	ctx.Dataset.Format = b.Format
	ctx.Dataset.Keyword = b.Keyword
	ctx.Dataset.License = b.License
	ctx.Dataset.OtherLicense = b.OtherLicense
	ctx.Dataset.Publisher = b.Publisher
	ctx.Dataset.Title = b.Title
	ctx.Dataset.URL = b.URL
	ctx.Dataset.Year = b.Year

	render.Layout(w, "refresh_modal", "dataset/edit_details", YieldEditDetails{
		Context:  ctx,
		Form:     detailsForm(ctx.Locale, ctx.Dataset, nil),
		Conflict: false,
	})
}

func (h *Handler) UpdateDetails(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDetails{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update dataset details: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	// @note decoding the form into a model omits empty values
	//   removing "omitempty" in the model doesn't make a difference.
	if b.AccessLevel != "info:eu-repo/semantics/embargoedAccess" {
		b.EmbargoDate = ""
		b.AccessLevelAfterEmbargo = ""
	}

	ctx.Dataset.AccessLevel = b.AccessLevel
	ctx.Dataset.EmbargoDate = b.EmbargoDate
	ctx.Dataset.AccessLevelAfterEmbargo = b.AccessLevelAfterEmbargo
	ctx.Dataset.Format = b.Format
	ctx.Dataset.Keyword = b.Keyword
	ctx.Dataset.License = b.License
	ctx.Dataset.OtherLicense = b.OtherLicense
	ctx.Dataset.Publisher = b.Publisher
	ctx.Dataset.Title = b.Title
	ctx.Dataset.URL = b.URL
	ctx.Dataset.Year = b.Year

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "dataset/edit_details", YieldEditDetails{
			Context:  ctx,
			Form:     detailsForm(ctx.Locale, ctx.Dataset, validationErrs.(validation.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "dataset/edit_details", YieldEditDetails{
			Context:  ctx,
			Form:     detailsForm(ctx.Locale, ctx.Dataset, nil),
			Conflict: true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("update dataset details: Could not save the dataset:", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_details", YieldDetails{
		Context:        ctx,
		DisplayDetails: displays.DatasetDetails(ctx.Locale, ctx.Dataset),
	})
}

func detailsForm(l *locale.Locale, d *models.Dataset, errors validation.Errors) *form.Form {
	f := form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(l, errors)).
		AddSection(
			&form.Text{
				Name:     "title",
				Value:    d.Title,
				Label:    l.T("builder.title"),
				Cols:     9,
				Error:    localize.ValidationErrorAt(l, errors, "/title"),
				Required: d.FieldIsRequired(),
			},
			&display.Text{
				Label:         l.T("builder.doi"),
				Value:         d.DOI,
				Required:      d.FieldIsRequired(),
				ValueTemplate: "format/doi",
			},
			&form.Text{
				Name:  "url",
				Value: d.URL,
				Label: l.T("builder.url"),
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/url"),
			},
		).
		AddSection(
			&form.Text{
				Name:     "publisher",
				Value:    d.Publisher,
				Label:    l.T("builder.publisher"),
				Cols:     9,
				Error:    localize.ValidationErrorAt(l, errors, "/publisher"),
				Required: d.FieldIsRequired(),
				Tooltip:  l.T("tooltip.dataset.publisher"),
			},
			&form.Text{
				Name:     "year",
				Value:    d.Year,
				Label:    l.T("builder.year"),
				Cols:     3,
				Help:     l.T("builder.year.help"),
				Error:    localize.ValidationErrorAt(l, errors, "/year"),
				Required: d.FieldIsRequired(),
			},
		).
		AddSection(
			&form.TextRepeat{
				Name:            "format",
				Values:          d.Format,
				Label:           l.T("builder.format"),
				Cols:            9,
				Error:           localize.ValidationErrorAt(l, errors, "/format"),
				Required:        d.FieldIsRequired(),
				AutocompleteURL: "suggest_media_types",
				Tooltip:         l.T("tooltip.dataset.format"),
			},
			&form.TextRepeat{
				Name:   "keyword",
				Values: d.Keyword,
				Label:  l.T("builder.keyword"),
				Cols:   9,
				Error:  localize.ValidationErrorAt(l, errors, "/keyword"),
			},
		)

	if d.License == "LicenseNotListed" {
		f.AddSection(
			&form.Select{
				Name:        "license",
				Template:    "dataset/license",
				Value:       d.License,
				Label:       l.T("builder.license"),
				Options:     localize.VocabularySelectOptions(l, "dataset_licenses"),
				Cols:        3,
				Error:       localize.ValidationErrorAt(l, errors, "/license"),
				Tooltip:     l.T("tooltip.dataset.license"),
				EmptyOption: true,
				Required:    d.FieldIsRequired(),
				Vars:        struct{ ID string }{ID: d.ID},
			},
			&form.Text{
				Name:     "other_license",
				Value:    d.OtherLicense,
				Label:    l.T("builder.other_license"),
				Cols:     9,
				Help:     l.T("builder.other_license.help"),
				Error:    localize.ValidationErrorAt(l, errors, "/other_license"),
				Required: d.FieldIsRequired(),
			},
		)
	} else {
		f.AddSection(
			&form.Select{
				Name:        "license",
				Template:    "dataset/license",
				Value:       d.License,
				Label:       l.T("builder.license"),
				Options:     localize.VocabularySelectOptions(l, "dataset_licenses"),
				Cols:        3,
				Error:       localize.ValidationErrorAt(l, errors, "/license"),
				Tooltip:     l.T("tooltip.dataset.license"),
				EmptyOption: true,
				Required:    d.FieldIsRequired(),
				Vars:        struct{ ID string }{ID: d.ID},
			},
		)
	}

	if d.AccessLevel != "info:eu-repo/semantics/embargoedAccess" {
		f.AddSection(
			&form.Select{
				Name:        "access_level",
				Template:    "dataset/access_level",
				Label:       l.T("builder.access_level"),
				Value:       d.AccessLevel,
				Options:     localize.VocabularySelectOptions(l, "dataset_access_levels"),
				Cols:        3,
				Error:       localize.ValidationErrorAt(l, errors, "/access_level"),
				Required:    d.FieldIsRequired(),
				EmptyOption: true,
				Tooltip:     l.T("tooltip.dataset.access_level"),
				Vars:        struct{ ID string }{ID: d.ID},
			},
		)
	} else {
		f.AddSection(
			&form.Select{
				Name:        "access_level",
				Template:    "dataset/access_level",
				Label:       l.T("builder.access_level"),
				Value:       d.AccessLevel,
				Options:     localize.VocabularySelectOptions(l, "dataset_access_levels"),
				Cols:        3,
				Error:       localize.ValidationErrorAt(l, errors, "/access_level"),
				Required:    d.FieldIsRequired(),
				EmptyOption: true,
				Tooltip:     l.T("tooltip.dataset.access_level"),
				Vars:        struct{ ID string }{ID: d.ID},
			},
			&form.Date{
				Name:  "embargo_date",
				Value: d.EmbargoDate,
				Label: l.T("builder.embargo_date"),
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/embargo_date"),
				// Disabled: d.AccessLevel != "info:eu-repo/semantics/embargoedAccess",
			},
			&form.Select{
				Name:        "access_level_after_embargo",
				Label:       l.T("builder.access_level_after_embargo"),
				Value:       d.AccessLevelAfterEmbargo,
				Options:     localize.VocabularySelectOptions(l, "dataset_access_levels_after_embargo"),
				Cols:        3,
				Error:       localize.ValidationErrorAt(l, errors, "/access_level_after_embargo"),
				EmptyOption: true,
				// Disabled:    d.AccessLevel != "info:eu-repo/semantics/embargoedAccess",
			},
		)
	}

	return f
}
