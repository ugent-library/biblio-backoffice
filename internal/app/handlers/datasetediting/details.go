package datasetediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/displays"
	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type BindDetails struct {
	AccessLevel  string   `form:"access_level"`
	Embargo      string   `form:"embargo"`
	EmbargoTo    string   `form:"embargo_to"`
	Format       []string `form:"format"`
	Keyword      []string `form:"keyword"`
	License      string   `form:"license"`
	OtherLicense string   `form:"other_license"`
	Publisher    string   `form:"publisher"`
	Title        string   `form:"title"`
	URL          string   `form:"url"`
	Year         string   `form:"year"`
}

type YieldDetails struct {
	Context
	DisplayDetails *display.Display
}

type YieldEditDetails struct {
	Context
	Form *form.Form
}

func (h *Handler) EditDetails(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDetails{}
	b.AccessLevel = ctx.Dataset.AccessLevel
	b.Embargo = ctx.Dataset.Embargo
	b.EmbargoTo = ctx.Dataset.EmbargoTo
	b.Format = ctx.Dataset.Format
	b.Keyword = ctx.Dataset.Keyword
	b.License = ctx.Dataset.License
	b.OtherLicense = ctx.Dataset.OtherLicense
	b.Publisher = ctx.Dataset.Publisher
	b.Title = ctx.Dataset.Title
	b.URL = ctx.Dataset.URL
	b.Year = ctx.Dataset.Year

	render.Render(w, "dataset/edit_details", YieldEditDetails{
		Context: ctx,
		Form:    detailsForm(ctx, b, nil),
	})
}

func (h *Handler) EditDetailsAccessLevel(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDetails{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	// Clear embargo and embargoTo fields if access level is not embargo
	//   TODO Disabled per https://github.com/ugent-library/biblio-backend/issues/217
	//
	//   Another issue: the old JS also temporary stored the data in these fields if
	//   access level changed from embargo to something else. The data would be restored
	//   into the form fields again if embargo level is chosen again. This feature isn't
	//   implemented in this solution since state isn't kept across HTTP requests.
	//
	if b.AccessLevel != "info:eu-repo/semantics/embargoedAccess" {
		b.Embargo = ""
		b.EmbargoTo = ""
	}

	render.Render(w, "dataset/edit_details", YieldEditDetails{
		Context: ctx,
		Form:    detailsForm(ctx, b, nil),
	})
}

func (h *Handler) UpdateDetails(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDetails{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	// @note decoding the form into a model omits empty values
	//   removing "omitempty" in the model doesn't make a difference.
	if b.AccessLevel != "info:eu-repo/semantics/embargoedAccess" {
		b.Embargo = ""
		b.EmbargoTo = ""
	}

	ctx.Dataset.AccessLevel = b.AccessLevel
	ctx.Dataset.Embargo = b.Embargo
	ctx.Dataset.EmbargoTo = b.EmbargoTo
	ctx.Dataset.Format = b.Format
	ctx.Dataset.Keyword = b.Keyword
	ctx.Dataset.License = b.License
	ctx.Dataset.OtherLicense = b.OtherLicense
	ctx.Dataset.Publisher = b.Publisher
	ctx.Dataset.Title = b.Title
	ctx.Dataset.URL = b.URL
	ctx.Dataset.Year = b.Year

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		form := detailsForm(ctx, b, validationErrs.(validation.Errors))

		render.Render(w, "dataset/refresh_edit_details", YieldEditDetails{
			Context: ctx,
			Form:    form,
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

	render.Render(w, "dataset/refresh_details", YieldDetails{
		Context:        ctx,
		DisplayDetails: displays.DatasetDetails(ctx.Locale, ctx.Dataset),
	})
}

func detailsForm(ctx Context, b BindDetails, errors validation.Errors) *form.Form {
	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(ctx.Locale, errors)).
		AddSection(
			&form.Text{
				Name:        "title",
				Value:       b.Title,
				Label:       ctx.T("builder.title"),
				Cols:        9,
				Placeholder: ctx.T("builder.details.title.placeholder"),
				Error:       localize.ValidationErrorAt(ctx.Locale, errors, "/title"),
				Required:    true,
			},
			&display.Text{
				Label:         ctx.T("builder.doi"),
				Value:         ctx.Dataset.DOI,
				Required:      true,
				ValueTemplate: "format/doi",
			},
			&form.Text{
				Name:  "url",
				Value: b.URL,
				Label: ctx.T("builder.url"),
				Cols:  3,
				Error: localize.ValidationErrorAt(ctx.Locale, errors, "/url"),
			},
		).
		AddSection(
			&form.Text{
				Name:        "publisher",
				Value:       b.Publisher,
				Label:       ctx.T("builder.publisher"),
				Cols:        9,
				Placeholder: ctx.T("builder.details.publisher.placeholder"),
				Error:       localize.ValidationErrorAt(ctx.Locale, errors, "/publisher"),
				Required:    true,
				Tooltip:     ctx.T("tooltip.dataset.publisher"),
			},
			&form.Text{
				Name:        "year",
				Value:       b.Year,
				Label:       ctx.T("builder.year"),
				Cols:        3,
				Placeholder: ctx.T("builder.year.placeholder"),
				Error:       localize.ValidationErrorAt(ctx.Locale, errors, "/year"),
				Required:    true,
			},
		).
		AddSection(
			&form.TextRepeat{
				Name:            "format",
				Values:          b.Format,
				Label:           ctx.T("builder.format"),
				Cols:            9,
				Error:           localize.ValidationErrorAt(ctx.Locale, errors, "/format"),
				Required:        true,
				AutocompleteURL: "suggest_media_types",
				Tooltip:         ctx.T("tooltip.dataset.format"),
			},
			&form.TextRepeat{
				Name:   "keyword",
				Values: b.Keyword,
				Label:  ctx.T("builder.keyword"),
				Cols:   9,
				Error:  localize.ValidationErrorAt(ctx.Locale, errors, "/keyword"),
			},
		).
		AddSection(
			&form.Select{
				Name:        "license",
				Value:       b.License,
				Label:       ctx.T("builder.license"),
				Options:     localize.VocabularySelectOptions(ctx.Locale, "cc_licenses"),
				Cols:        3,
				Error:       localize.ValidationErrorAt(ctx.Locale, errors, "/license"),
				Tooltip:     ctx.T("tooltip.dataset.license"),
				EmptyOption: true,
				Required:    true,
			},
			&form.Text{
				Name:        "other_license",
				Value:       b.OtherLicense,
				Label:       ctx.T("builder.other_license"),
				Cols:        9,
				Placeholder: "e.g. https://creativecommons.org/licenses/by/4.0/",
				Error:       localize.ValidationErrorAt(ctx.Locale, errors, "/other_license"),
				Required:    true,
			},
			&form.Select{
				Template:    "dataset/access_level",
				Name:        "access_level",
				Label:       ctx.T("builder.access_level"),
				Value:       b.AccessLevel,
				Options:     localize.VocabularySelectOptions(ctx.Locale, "access_levels"),
				Cols:        3,
				Error:       localize.ValidationErrorAt(ctx.Locale, errors, "/access_level"),
				Required:    true,
				EmptyOption: true,
				Tooltip:     ctx.T("tooltip.dataset.access_level"),
				Vars:        struct{ ID string }{ID: ctx.Dataset.ID},
			},
			&form.Date{
				Name:     "embargo",
				Value:    b.Embargo,
				Label:    ctx.T("builder.embargo"),
				Cols:     3,
				Error:    localize.ValidationErrorAt(ctx.Locale, errors, "/embargo"),
				Disabled: b.AccessLevel != "info:eu-repo/semantics/embargoedAccess",
			},
			&form.Select{
				Name:        "embargo_to",
				Label:       ctx.T("builder.embargo_to"),
				Value:       b.EmbargoTo,
				Options:     localize.VocabularySelectOptions(ctx.Locale, "access_levels"),
				Cols:        3,
				Error:       localize.ValidationErrorAt(ctx.Locale, errors, "/embargo_to"),
				EmptyOption: true,
				Disabled:    b.AccessLevel != "info:eu-repo/semantics/embargoedAccess",
			},
		)
}
