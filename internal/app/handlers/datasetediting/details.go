package datasetediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/localize"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/fields"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type BindDetails struct {
	ID           string   `form:"-"`
	AccessLevel  string   `form:"access_level"`
	DOI          string   `form:"-"`
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
	DetailFields []*fields.Fields
}

type YieldEditDetails struct {
	Context
	Form *form.Form
}

func (h *Handler) EditDetails(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDetails{}
	if err := bind.RequestPath(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	b.ID = ctx.Dataset.ID
	b.AccessLevel = ctx.Dataset.AccessLevel
	b.DOI = ctx.Dataset.DOI
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
	if err := bind.RequestForm(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	// Band-aid: omit empty values from fields with repeating values (text repeat, list repeat)
	//   @note: This should be part of "form:,omitEmpty" in the Dataset struct. However,
	//   go-playground/form doesn't support omitEmpty on lists or nested form structures
	//   (slices, maps,...)
	omitEmpty := func(keywords []string) []string {
		var tmp []string
		for _, str := range keywords {
			if str != "" {
				tmp = append(tmp, str)
			}
		}

		return tmp
	}

	b.Keyword = omitEmpty(b.Keyword)
	b.Format = omitEmpty(b.Format)
	b.ID = ctx.Dataset.ID

	// Clear embargo and embargoTo fields if access level is not embargo
	//   @todo Disabled per https://github.com/ugent-library/biblio-backend/issues/217
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

func (h *Handler) SaveDetails(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDetails{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	// @note decoding the form into a model omits empty values
	//   removing "omitempty" in the model doesn't make a difference.
	if b.AccessLevel != "info:eu-repo/semantics/embargoedAccess" {
		b.Embargo = ""
		b.EmbargoTo = ""
	}

	// Band-aid: omit empty values from fields with repeating values (text repeat, list repeat)
	//   @note: This should be part of "form:,omitEmpty" in the Dataset struct. However,
	//   go-playground/form doesn't support omitEmpty on lists or nested form structures
	//   (slices, maps,...)
	omitEmpty := func(keywords []string) []string {
		var tmp []string
		for _, str := range keywords {
			if str != "" {
				tmp = append(tmp, str)
			}
		}

		return tmp
	}

	b.Keyword = omitEmpty(b.Keyword)
	b.Format = omitEmpty(b.Format)

	ctx.Dataset.AccessLevel = b.AccessLevel
	ctx.Dataset.DOI = b.DOI
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
		Context:      ctx,
		DetailFields: detailFields(ctx),
	})
}

func detailsForm(ctx Context, b BindDetails, errors validation.Errors) *form.Form {
	detailsForm := form.NewForm().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(ctx.Locale, errors))

	title := &form.Text{
		Name:        "title",
		Value:       b.Title,
		Label:       ctx.T("builder.title"),
		Cols:        9,
		Placeholder: ctx.T("builder.details.title.placeholder"),
		Error:       localize.ValidationErrorAt(ctx.Locale, errors, "/details/title"),
		Required:    true,
	}

	url := &form.Text{
		Name:  "url",
		Value: b.URL,
		Label: ctx.T("builder.url"),
		Cols:  3,
		Error: localize.ValidationErrorAt(ctx.Locale, errors, "/details/url"),
	}

	detailsForm.AddSection(title, url)

	publisher := &form.Text{
		Name:        "publisher",
		Value:       b.Publisher,
		Label:       ctx.T("builder.publisher"),
		Cols:        9,
		Placeholder: ctx.T("builder.details.publisher.placeholder"),
		Error:       localize.ValidationErrorAt(ctx.Locale, errors, "/details/publisher"),
		Required:    true,
		// Tooltip
	}

	year := &form.Text{
		Name:        "year",
		Value:       b.Year,
		Label:       ctx.T("builder.year"),
		Cols:        3,
		Placeholder: ctx.T("builder.year.placeholder"),
		Error:       localize.ValidationErrorAt(ctx.Locale, errors, "/details/year"),
		Required:    true,
	}

	detailsForm.AddSection(publisher, year)

	format := &form.TextRepeat{
		Name:            "format",
		Values:          b.Format,
		Label:           ctx.T("builder.format"),
		Cols:            9,
		Error:           localize.ValidationErrorAt(ctx.Locale, errors, "/details/format"),
		Required:        true,
		AutocompleteURL: "media_type_choose",
	}

	keyword := &form.TextRepeat{
		Name:   "keyword",
		Values: b.Keyword,
		Label:  ctx.T("builder.details.keyword"),
		Cols:   9,
		Error:  localize.ValidationErrorAt(ctx.Locale, errors, "/details/keyword"),
	}

	detailsForm.AddSection(format, keyword)

	license := &form.Select{
		Name:    "license",
		Value:   b.License,
		Label:   ctx.T("builder.license"),
		Options: localize.LicenseSelectOptions(ctx.Locale),
		Cols:    3,
		Error:   localize.ValidationErrorAt(ctx.Locale, errors, "/details/license"),
	}

	otherLicense := &form.Text{
		Name:        "other_license",
		Value:       b.OtherLicense,
		Label:       ctx.T("builder.other_license"),
		Cols:        9,
		Placeholder: ctx.T("builder.details.other_license.placeholder"),
		Error:       localize.ValidationErrorAt(ctx.Locale, errors, "/details/other_license"),
		Required:    true,
	}

	accessLevel := &form.Select{
		Template:    "dataset/access_level",
		Name:        "access_level",
		Label:       ctx.T("builder.access_level"),
		Value:       b.AccessLevel,
		Options:     localize.AccessLevelSelectOptions(ctx.Locale),
		Cols:        3,
		Error:       localize.ValidationErrorAt(ctx.Locale, errors, "/details/access_level"),
		Required:    true,
		EmptyChoice: true,
		// Tooltip
		Vars: struct{ ID string }{ID: b.ID},
	}

	disabled := false
	if b.AccessLevel != "info:eu-repo/semantics/embargoedAccess" {
		disabled = true
	}

	embargo := &form.Date{
		Name:     "embargo",
		Value:    b.Embargo,
		Label:    ctx.T("builder.embargo"),
		Cols:     3,
		Error:    localize.ValidationErrorAt(ctx.Locale, errors, "/details/embargo"),
		Disabled: disabled,
	}

	embargoTo := &form.Select{
		Name:        "embargo_to",
		Label:       ctx.T("builder.embargo_to"),
		Value:       b.EmbargoTo,
		Options:     localize.AccessLevelSelectOptions(ctx.Locale),
		Cols:        3,
		Error:       localize.ValidationErrorAt(ctx.Locale, errors, "/details/embargo_to"),
		EmptyChoice: true,
		Disabled:    disabled,
	}

	detailsForm.AddSection(license, otherLicense, accessLevel, embargo, embargoTo)

	return detailsForm

	// @todo Kept for future reference

	// return &form.Form{
	// 	Theme:  "default",
	// 	Errors: localize.ValidationErrors(ctx.Locale, errors),
	// 	Fields: []form.Field{
	// 		&form.FieldGroup{
	// 			Fields: []form.Field{
	// 				&form.Text{
	// 					Name:        "title",
	// 					Value:       b.Title,
	// 					Label:       ctx.T("builder.details.title"),
	// 					Cols:        9,
	// 					Placeholder: ctx.T("builder.details.title.placeholder"),
	// 					Error:       localize.ValidationErrorAt(ctx.Locale, errors, "/details/title"),
	// 					Required:    true,
	// 				},
	// 				// DOI
	// 				&form.Text{
	// 					Name:  "URL",
	// 					Value: b.URL,
	// 					Label: ctx.T("builder.details.url"),
	// 					Cols:  3,
	// 					Error: localize.ValidationErrorAt(ctx.Locale, errors, "/details/url"),
	// 				},
	// 			},
	// 		},
	// 		&form.FieldGroup{
	// 			Fields: []form.Field{
	// 				&form.Text{
	// 					Name:        "publisher",
	// 					Value:       b.Publisher,
	// 					Label:       ctx.T("builder.details.publisher"),
	// 					Cols:        9,
	// 					Placeholder: ctx.T("builder.details.publisher.placeholder"),
	// 					Error:       localize.ValidationErrorAt(ctx.Locale, errors, "/details/publisher"),
	// 					Required:    true,
	// 					// Tooltip
	// 				},
	// 				&form.Text{
	// 					Name:        "year",
	// 					Value:       b.Year,
	// 					Label:       ctx.T("builder.details.year"),
	// 					Cols:        3,
	// 					Placeholder: ctx.T("builder.details.year.placeholder"),
	// 					Error:       localize.ValidationErrorAt(ctx.Locale, errors, "/details/year"),
	// 					Required:    true,
	// 				},
	// 			},
	// 		},
	// 		&form.FieldGroup{
	// 			Fields: []form.Field{
	// 				&form.TextRepeat{
	// 					Name:            "format",
	// 					Values:          b.Format,
	// 					Label:           ctx.T("builder.details.format"),
	// 					Cols:            9,
	// 					Error:           localize.ValidationErrorAt(ctx.Locale, errors, "/details/format"),
	// 					Required:        true,
	// 					AutocompleteURL: "media_type_choose",
	// 				},
	// 			},
	// 		},
	// 		// &form.Select{
	// 		// 	Name:    "lang",
	// 		// 	Value:   b.Lang,
	// 		// 	Label:   ctx.T("builder.abstract.lang"),
	// 		// 	Options: localize.LanguageSelectOptions(ctx.Locale),
	// 		// 	Cols:    12,
	// 		// 	Error:   localize.ValidationErrorAt(ctx.Locale, errors, fmt.Sprintf("/abstract/%d/lang", b.Position)),
	// 		// },
	// 	},
	// }
}

func detailFields(ctx Context) []*fields.Fields {
	return []*fields.Fields{
		{
			Theme: "default",
			Fields: []fields.Field{
				&fields.Text{
					Label:    ctx.T("builder.title"),
					Value:    ctx.Dataset.Title,
					Required: true,
				},
				&fields.Text{
					Label:         ctx.T("builder.doi"),
					Value:         ctx.Dataset.DOI,
					Required:      true,
					ValueTemplate: "format/doi",
				},
				&fields.Text{
					Label:         ctx.T("builder.url"),
					Value:         ctx.Dataset.URL,
					ValueTemplate: "format/link",
				},
			},
		},
		{
			Theme: "default",
			Fields: []fields.Field{
				&fields.Text{
					Label:    ctx.T("builder.publisher"),
					Value:    ctx.Dataset.Publisher,
					Required: true,
				},
				&fields.Text{
					Label:    ctx.T("builder.year"),
					Value:    ctx.Dataset.Year,
					Required: true,
				},
			},
		},
		{
			Theme: "default",
			Fields: []fields.Field{
				&fields.Text{
					Label:    ctx.T("builder.format"),
					Values:   ctx.Dataset.Format,
					List:     true,
					Required: true,
				},
				&fields.Text{
					Label:         ctx.T("builder.keyword"),
					Values:        ctx.Dataset.Keyword,
					ValueTemplate: "format/badge",
				},
			},
		},
		{
			Theme: "default",
			Fields: []fields.Field{
				&fields.Text{
					Label:    ctx.T("builder.license"),
					Value:    ctx.TS("cc_licenses", ctx.Dataset.License),
					Required: true,
				},
				&fields.Text{
					Label: ctx.T("builder.other_license"),
					Value: ctx.Dataset.OtherLicense,
				},
				&fields.Text{
					Label:    ctx.T("builder.access_level"),
					Value:    ctx.TS("access_levels", ctx.Dataset.AccessLevel),
					Required: true,
				},
				&fields.Text{
					Label: ctx.T("builder.embargo"),
					Value: ctx.Dataset.Embargo,
				},
				&fields.Text{
					Label: ctx.T("builder.embargo_to"),
					Value: ctx.TS("access_levels", ctx.Dataset.EmbargoTo),
				},
			},
		},
	}
}
