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
	render.Layout(w, "show_modal", "dataset/edit_details", YieldEditDetails{
		Context: ctx,
		Form:    detailsForm(ctx.Locale, ctx.Dataset, nil),
	})
}

func (h *Handler) EditDetailsAccessLevel(w http.ResponseWriter, r *http.Request, ctx Context) {
	// Clear embargo and embargoTo fields if access level is not embargo
	//   TODO Disabled per https://github.com/ugent-library/biblio-backend/issues/217
	//
	//   Another issue: the old JS also temporary stored the data in these fields if
	//   access level changed from embargo to something else. The data would be restored
	//   into the form fields again if embargo level is chosen again. This feature isn't
	//   implemented in this solution since state isn't kept across HTTP requests.
	//
	dataset := ctx.Dataset
	if dataset.AccessLevel != "info:eu-repo/semantics/embargoedAccess" {
		dataset.Embargo = ""
		dataset.EmbargoTo = ""
	}

	render.Layout(w, "refresh_modal", "dataset/edit_details", YieldEditDetails{
		Context: ctx,
		Form:    detailsForm(ctx.Locale, dataset, nil),
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
		form := detailsForm(ctx.Locale, ctx.Dataset, validationErrs.(validation.Errors))

		render.Layout(w, "refresh_modal", "dataset/edit_details", YieldEditDetails{
			Context: ctx,
			Form:    form,
		})
		return
	}

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("dataset.conflict_error"))
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

func detailsForm(l *locale.Locale, dataset *models.Dataset, errors validation.Errors) *form.Form {
	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(l, errors)).
		AddSection(
			&form.Text{
				Name:     "title",
				Value:    dataset.Title,
				Label:    l.T("builder.title"),
				Cols:     9,
				Error:    localize.ValidationErrorAt(l, errors, "/title"),
				Required: true,
			},
			&display.Text{
				Label:         l.T("builder.doi"),
				Value:         dataset.DOI,
				Required:      true,
				ValueTemplate: "format/doi",
			},
			&form.Text{
				Name:  "url",
				Value: dataset.URL,
				Label: l.T("builder.url"),
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/url"),
			},
		).
		AddSection(
			&form.Text{
				Name:     "publisher",
				Value:    dataset.Publisher,
				Label:    l.T("builder.publisher"),
				Cols:     9,
				Error:    localize.ValidationErrorAt(l, errors, "/publisher"),
				Required: true,
				Tooltip:  l.T("tooltip.dataset.publisher"),
			},
			&form.Text{
				Name:     "year",
				Value:    dataset.Year,
				Label:    l.T("builder.year"),
				Cols:     3,
				Help:     l.T("builder.year.help"),
				Error:    localize.ValidationErrorAt(l, errors, "/year"),
				Required: true,
			},
		).
		AddSection(
			&form.TextRepeat{
				Name:            "format",
				Values:          dataset.Format,
				Label:           l.T("builder.format"),
				Cols:            9,
				Error:           localize.ValidationErrorAt(l, errors, "/format"),
				Required:        true,
				AutocompleteURL: "suggest_media_types",
				Tooltip:         l.T("tooltip.dataset.format"),
			},
			&form.TextRepeat{
				Name:   "keyword",
				Values: dataset.Keyword,
				Label:  l.T("builder.keyword"),
				Cols:   9,
				Error:  localize.ValidationErrorAt(l, errors, "/keyword"),
			},
		).
		AddSection(
			&form.Select{
				Name:        "license",
				Value:       dataset.License,
				Label:       l.T("builder.license"),
				Options:     localize.VocabularySelectOptions(l, "licenses"),
				Cols:        3,
				Error:       localize.ValidationErrorAt(l, errors, "/license"),
				Tooltip:     l.T("tooltip.dataset.license"),
				EmptyOption: true,
				Required:    true,
			},
			&form.Text{
				Name:     "other_license",
				Value:    dataset.OtherLicense,
				Label:    l.T("builder.other_license"),
				Cols:     9,
				Help:     l.T("builder.other_license.help"),
				Error:    localize.ValidationErrorAt(l, errors, "/other_license"),
				Required: true,
			},
			&form.Select{
				//TODO: closes modal because controller reuses full edit view
				//Template:    "dataset/access_level",
				Name:        "access_level",
				Label:       l.T("builder.access_level"),
				Value:       dataset.AccessLevel,
				Options:     localize.VocabularySelectOptions(l, "access_levels"),
				Cols:        3,
				Error:       localize.ValidationErrorAt(l, errors, "/access_level"),
				Required:    true,
				EmptyOption: true,
				Tooltip:     l.T("tooltip.dataset.access_level"),
				Vars:        struct{ ID string }{ID: dataset.ID},
			},
			&form.Date{
				Name:     "embargo",
				Value:    dataset.Embargo,
				Label:    l.T("builder.embargo"),
				Cols:     3,
				Error:    localize.ValidationErrorAt(l, errors, "/embargo"),
				Disabled: dataset.AccessLevel != "info:eu-repo/semantics/embargoedAccess",
			},
			&form.Select{
				Name:        "embargo_to",
				Label:       l.T("builder.embargo_to"),
				Value:       dataset.EmbargoTo,
				Options:     localize.VocabularySelectOptions(l, "access_levels"),
				Cols:        3,
				Error:       localize.ValidationErrorAt(l, errors, "/embargo_to"),
				EmptyOption: true,
				Disabled:    dataset.AccessLevel != "info:eu-repo/semantics/embargoedAccess",
			},
		)
}
