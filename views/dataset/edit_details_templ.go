// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package dataset

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"fmt"
	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views/form"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
	"github.com/ugent-library/okay"
	"time"
)

func datasetIdentifierAndType(dataset *models.Dataset) (string, string) {
	var identifierType, identifier string
	for _, key := range vocabularies.Map["dataset_identifier_types"] {
		if val := dataset.Identifiers.Get(key); val != "" {
			identifierType = key
			identifier = val
			break
		}
	}
	return identifier, identifierType
}

func datasetIdentifierOptions(loc *gotext.Locale) []form.Option {
	identifierTypeOptions := make([]form.Option, len(vocabularies.Map["dataset_identifier_types"]))
	for i, v := range vocabularies.Map["dataset_identifier_types"] {
		identifierTypeOptions[i].Label = loc.Get("identifier." + v)
		identifierTypeOptions[i].Value = v
	}
	return identifierTypeOptions
}

func nextDay() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)
}

func EditDetailsDialog(c *ctx.Ctx, dataset *models.Dataset, conflict bool, errors *okay.Errors) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"modal-dialog modal-dialog-centered modal-fullscreen modal-dialog-scrollable\" role=\"document\"><div class=\"modal-content\"><div class=\"modal-header\"><h2 class=\"modal-title\">Edit dataset details</h2></div><div class=\"modal-body\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if conflict {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"alert alert-danger mb-0\" role=\"alert\"><i class=\"if if--error if-error-circle-fill\"></i> The dataset you are editing has been changed by someone else. Please copy your edits, then close this form.</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = form.Errors(localize.ValidationErrors(c.Loc, errors)).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<ul class=\"list-group list-group-flush\" data-panel-state=\"edit\"><li class=\"list-group-item\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = form.Text(form.TextArgs{
			FieldArgs: form.FieldArgs{
				Name:     "title",
				Label:    c.Loc.Get("builder.title"),
				Cols:     9,
				Error:    localize.ValidationErrorAt(c.Loc, errors, "/title"),
				Required: true,
			},
			Value: dataset.Title,
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if identifier, identifierType := datasetIdentifierAndType(dataset); identifierType != "" {
			templ_7745c5c3_Err = form.Select(form.SelectArgs{
				FieldArgs: form.FieldArgs{
					Name:     "identifier_type",
					Label:    c.Loc.Get("builder.identifier_type"),
					Cols:     3,
					Help:     c.Loc.Get("builder.identifier_type.help"),
					Error:    localize.ValidationErrorAt(c.Loc, errors, "/identifier"),
					Required: true,
				},
				Value:       identifierType,
				EmptyOption: true,
				Options:     datasetIdentifierOptions(c.Loc),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Name:     "identifier",
					Label:    c.Loc.Get("builder.identifier"),
					Required: true,
					Cols:     3,
					Help:     c.Loc.Get("builder.identifier.help"),
					Error:    localize.ValidationErrorAt(c.Loc, errors, "/identifier"),
					Tooltip:  c.Loc.Get("tooltip.dataset.identifier"),
				},
				Value: identifier,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li><li class=\"list-group-item\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = form.SelectRepeat(form.SelectRepeatArgs{
			FieldArgs: form.FieldArgs{
				Name:  "language",
				Label: c.Loc.Get("builder.language"),
				Cols:  9,
				Error: localize.ValidationErrorAt(c.Loc, errors, "/language"),
			},
			Options:     localize.LanguageSelectOptions(),
			Values:      dataset.Language,
			EmptyOption: true,
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = form.Text(form.TextArgs{
			FieldArgs: form.FieldArgs{
				Name:     "year",
				Label:    c.Loc.Get("builder.year"),
				Cols:     3,
				Help:     c.Loc.Get("builder.year.help"),
				Error:    localize.ValidationErrorAt(c.Loc, errors, "/year"),
				Required: true,
			},
			Value: dataset.Year,
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = form.Text(form.TextArgs{
			FieldArgs: form.FieldArgs{
				Name:     "publisher",
				Label:    c.Loc.Get("builder.publisher"),
				Cols:     9,
				Error:    localize.ValidationErrorAt(c.Loc, errors, "/publisher"),
				Required: true,
				Tooltip:  c.Loc.Get("tooltip.dataset.publisher"),
			},
			Value: dataset.Publisher,
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li><li class=\"list-group-item\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = form.TextRepeatWithAutocomplete(form.TextRepeatWithAutocompleteArgs{
			FieldArgs: form.FieldArgs{
				Name:     "format",
				Label:    c.Loc.Get("builder.format"),
				Cols:     9,
				Error:    localize.ValidationErrorAt(c.Loc, errors, "/format"),
				Required: true,
				Tooltip:  c.Loc.Get("tooltip.dataset.format"),
			},
			Values:          dataset.Format,
			AutocompleteURL: c.PathTo("suggest_media_types").String(),
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = form.Tags(form.TagsArgs{
			FieldArgs: form.FieldArgs{
				Name:  "keyword",
				Label: c.Loc.Get("builder.keyword"),
				Cols:  9,
				Error: localize.ValidationErrorAt(c.Loc, errors, "/keyword"),
			},
			Values: dataset.Keyword,
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li><li class=\"list-group-item\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			var templ_7745c5c3_Var3 = []any{"form-select", "form-control", templ.KV("is-invalid", errors != nil && errors.Get("/license") != nil)}
			templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var3...)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<select class=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ.CSSClasses(templ_7745c5c3_Var3).String()))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" id=\"license\" name=\"license\" hx-put=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(c.PathTo("dataset_refresh_edit_details", "id", dataset.ID).String()))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-include=\".modal-body\" hx-swap=\"none scroll:.modal-body:bottom\"><option></option> ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, o := range localize.VocabularySelectOptions(c.Loc, "dataset_licenses") {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(o.Value))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if o.Value == dataset.License {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" selected")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var4 string
				templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(o.Label)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/edit_details.templ`, Line: 171, Col: 85}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</option>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</select>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = form.Field(form.FieldArgs{
			Name:     "license",
			Label:    c.Loc.Get("builder.license"),
			Cols:     3,
			Error:    localize.ValidationErrorAt(c.Loc, errors, "/license"),
			Tooltip:  c.Loc.Get("tooltip.dataset.license"),
			Required: true,
		}, "license").Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if dataset.License == "LicenseNotListed" {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Name:     "other_license",
					Label:    c.Loc.Get("builder.other_license"),
					Cols:     9,
					Help:     c.Loc.Get("builder.other_license.help"),
					Error:    localize.ValidationErrorAt(c.Loc, errors, "/other_license"),
					Required: true,
				},
				Value: dataset.OtherLicense,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li><li class=\"list-group-item\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var5 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			var templ_7745c5c3_Var6 = []any{"form-select", "form-control", templ.KV("is-invalid", errors != nil && errors.Get("/access_level") != nil)}
			templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var6...)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<select class=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ.CSSClasses(templ_7745c5c3_Var6).String()))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" id=\"access_level\" name=\"access_level\" hx-put=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(c.PathTo("dataset_refresh_edit_details", "id", dataset.ID).String()))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-include=\".modal-body\" hx-swap=\"none scroll:.modal-body:bottom\"><option></option> ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, o := range localize.VocabularySelectOptions(c.Loc, "dataset_access_levels") {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(o.Value))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if o.Value == dataset.AccessLevel {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" selected")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var7 string
				templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(o.Label)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/edit_details.templ`, Line: 208, Col: 89}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</option>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</select>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = form.Field(form.FieldArgs{
			Name:     "access_level",
			Label:    c.Loc.Get("builder.access_level"),
			Cols:     3,
			Error:    localize.ValidationErrorAt(c.Loc, errors, "/access_level"),
			Tooltip:  c.Loc.Get("tooltip.dataset.access_level"),
			Required: true,
		}, "access_level").Render(templ.WithChildren(ctx, templ_7745c5c3_Var5), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if dataset.AccessLevel == "info:eu-repo/semantics/embargoedAccess" {
			templ_7745c5c3_Err = form.Date(form.DateArgs{
				FieldArgs: form.FieldArgs{
					Name:  "embargo_date",
					Label: c.Loc.Get("builder.embargo_date"),
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/embargo_date"),
				},
				Value: dataset.EmbargoDate,
				Min:   nextDay().Format("2006-01-02"),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.Select(form.SelectArgs{
				FieldArgs: form.FieldArgs{
					Name:  "access_level_after_embargo",
					Label: c.Loc.Get("builder.access_level_after_embargo"),
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/access_level_after_embargo"),
				},
				EmptyOption: true,
				Options:     localize.VocabularySelectOptions(c.Loc, "dataset_access_levels_after_embargo"),
				Value:       dataset.AccessLevelAfterEmbargo,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li></ul></div><div class=\"modal-footer\"><div class=\"bc-toolbar\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if conflict {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-left\"><button class=\"btn btn-primary modal-close\">Close</button></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-left\"><button class=\"btn btn-link modal-close\">Cancel</button></div><div class=\"bc-toolbar-right\"><button type=\"button\" name=\"create\" class=\"btn btn-primary\" hx-put=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(c.PathTo("dataset_update_details", "id", dataset.ID).String()))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-headers=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf(`{"If-Match": "%s"}`, dataset.SnapshotID)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-include=\".modal-body\" hx-swap=\"none\">Save</button></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}