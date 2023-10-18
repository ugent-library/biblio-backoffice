package datasetediting

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/bind"
	"github.com/ugent-library/biblio-backoffice/displays"
	"github.com/ugent-library/biblio-backoffice/locale"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/display"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/validation"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
)

type BindDetails struct {
	AccessLevel             string   `form:"access_level"`
	AccessLevelAfterEmbargo string   `form:"access_level_after_embargo"`
	EmbargoDate             string   `form:"embargo_date"`
	Format                  []string `form:"format"`
	Identifier              string   `form:"identifier"`
	IdentifierType          string   `form:"identifier_type"`
	Keyword                 []string `form:"keyword"`
	Language                []string `form:"language"`
	License                 string   `form:"license"`
	OtherLicense            string   `form:"other_license"`
	Publisher               string   `form:"publisher"`
	Title                   string   `form:"title"`
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
	ctx.Dataset.AccessLevelAfterEmbargo = b.AccessLevelAfterEmbargo
	ctx.Dataset.EmbargoDate = b.EmbargoDate
	ctx.Dataset.Format = b.Format
	ctx.Dataset.Identifiers = models.Values{b.IdentifierType: []string{b.Identifier}}
	ctx.Dataset.Keyword = b.Keyword
	ctx.Dataset.Language = b.Language
	ctx.Dataset.License = b.License
	ctx.Dataset.OtherLicense = b.OtherLicense
	ctx.Dataset.Publisher = b.Publisher
	ctx.Dataset.Title = b.Title
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
	ctx.Dataset.Language = b.Language
	ctx.Dataset.Keyword = b.Keyword
	ctx.Dataset.Identifiers = models.Values{b.IdentifierType: []string{b.Identifier}}
	ctx.Dataset.License = b.License
	ctx.Dataset.OtherLicense = b.OtherLicense
	ctx.Dataset.Publisher = b.Publisher
	ctx.Dataset.Title = b.Title
	ctx.Dataset.Year = b.Year

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "dataset/edit_details", YieldEditDetails{
			Context:  ctx,
			Form:     detailsForm(ctx.Locale, ctx.Dataset, validationErrs.(validation.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

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
		DisplayDetails: displays.DatasetDetails(ctx.User, ctx.Locale, ctx.Dataset),
	})
}

func detailsForm(l *locale.Locale, d *models.Dataset, errors validation.Errors) *form.Form {
	if d.Keyword == nil {
		d.Keyword = []string{}
	}
	keywordsJSON, _ := json.Marshal(d.Keyword)

	var identifierType, identifier string
	for _, key := range vocabularies.Map["dataset_identifier_types"] {
		if val := d.Identifiers.Get(key); val != "" {
			identifierType = key
			identifier = val
			break
		}
	}

	identifierTypeOptions := make([]form.SelectOption, len(vocabularies.Map["dataset_identifier_types"]))
	for i, v := range vocabularies.Map["dataset_identifier_types"] {
		identifierTypeOptions[i].Label = l.TS("identifier", v)
		identifierTypeOptions[i].Value = v
	}

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
				Required: true,
			},
			&form.Select{
				Name:        "identifier_type",
				Value:       identifierType,
				Label:       l.T("builder.identifier_type"),
				Options:     identifierTypeOptions,
				Cols:        3,
				Help:        template.HTML(l.T("builder.identifier_type.help")),
				Error:       localize.ValidationErrorAt(l, errors, "/identifier"),
				EmptyOption: true,
				Required:    true,
			},
			&form.Text{
				Name:     "identifier",
				Value:    identifier,
				Required: true,
				Label:    l.T("builder.identifier"),
				Cols:     3,
				Help:     template.HTML(l.T("builder.identifier.help")),
				Error:    localize.ValidationErrorAt(l, errors, "/identifier"),
				Tooltip:  l.T("tooltip.dataset.identifier"),
			},
		).
		AddSection(
			&form.SelectRepeat{
				Name:        "language",
				Label:       l.T("builder.language"),
				Options:     localize.LanguageSelectOptions(l),
				Values:      d.Language,
				EmptyOption: true,
				Cols:        9,
				Error:       localize.ValidationErrorAt(l, errors, "/language"),
			},
			&form.Text{
				Name:     "year",
				Value:    d.Year,
				Label:    l.T("builder.year"),
				Cols:     3,
				Help:     template.HTML(l.T("builder.year.help")),
				Error:    localize.ValidationErrorAt(l, errors, "/year"),
				Required: true,
			},
			&form.Text{
				Name:     "publisher",
				Value:    d.Publisher,
				Label:    l.T("builder.publisher"),
				Cols:     9,
				Error:    localize.ValidationErrorAt(l, errors, "/publisher"),
				Required: true,
				Tooltip:  l.T("tooltip.dataset.publisher"),
			},
		).
		AddSection(
			&form.TextRepeat{
				Name:            "format",
				Values:          d.Format,
				Label:           l.T("builder.format"),
				Cols:            9,
				Error:           localize.ValidationErrorAt(l, errors, "/format"),
				Required:        true,
				AutocompleteURL: "suggest_media_types",
				Tooltip:         l.T("tooltip.dataset.format"),
			},
			&form.Text{
				Name:     "keyword",
				Template: "tags",
				Value:    string(keywordsJSON), // TODO just pass the object itself
				Label:    l.T("builder.keyword"),
				Cols:     9,
				Error:    localize.ValidationErrorAt(l, errors, "/keyword"),
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
				Required:    true,
				Vars:        struct{ ID string }{ID: d.ID},
			},
			&form.Text{
				Name:     "other_license",
				Value:    d.OtherLicense,
				Label:    l.T("builder.other_license"),
				Cols:     9,
				Help:     template.HTML(l.T("builder.other_license.help")),
				Error:    localize.ValidationErrorAt(l, errors, "/other_license"),
				Required: true,
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
				Required:    true,
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
				Required:    true,
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
				Required:    true,
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
