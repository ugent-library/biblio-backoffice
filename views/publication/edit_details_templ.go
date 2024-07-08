// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package publication

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"fmt"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views/display"
	"github.com/ugent-library/biblio-backoffice/views/form"
	"github.com/ugent-library/okay"
)

func EditDetailsDialog(c *ctx.Ctx, p *models.Publication, conflict bool, errors *okay.Errors) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"modal-dialog modal-dialog-centered modal-fullscreen modal-dialog-scrollable\" role=\"document\"><div class=\"modal-content\"><div class=\"modal-header\"><h2 class=\"modal-title\">Edit publication details</h2></div><div class=\"modal-body\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if conflict {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"alert alert-danger mb-0\" role=\"alert\"><i class=\"if if--error if-error-circle-fill\"></i> The publication you are editing has been changed by someone else. Please copy your edits, then close this form.</div>")
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
		if c.User.CanChangeType(p) {
			templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
				templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
				if !templ_7745c5c3_IsBuffer {
					templ_7745c5c3_Buffer = templ.GetBuffer()
					defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
				}
				var templ_7745c5c3_Var3 = []any{"form-select", "form-control", templ.KV("is-invalid", errors != nil && errors.Get("/type") != nil)}
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
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" name=\"type\" id=\"type\" autofocus aria-details=\"type-help\" hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(c.PathTo("publication_confirm_update_type", "id", p.ID).String()))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-headers=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf(`{"If-Match": "%s"}`, p.SnapshotID)))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				for _, o := range localize.VocabularySelectOptions(c.Loc, "publication_types") {
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
					if o.Value == p.Type {
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
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_details.templ`, Line: 47, Col: 77}
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
				Label: c.Loc.Get("builder.type"),
				Name:  "type",
				Cols:  3,
				Error: localize.ValidationErrorAt(c.Loc, errors, "/type"),
				Help:  c.Loc.Get("builder.type.help"),
			}, "type").Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = display.Field(display.FieldArgs{
				Label:   c.Loc.Get("builder.type"),
				Value:   c.Loc.Get("publication_types." + p.Type),
				Tooltip: c.Loc.Get("tooltip.publication.type"),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesJournalArticleType() {
			templ_7745c5c3_Err = form.Select(form.SelectArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.journal_article_type"),
					Name:  "journal_article_type",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/journal_article_type"),
				},
				Value:       p.JournalArticleType,
				EmptyOption: true,
				Options:     localize.VocabularySelectOptions(c.Loc, "journal_article_types"),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesConferenceType() {
			templ_7745c5c3_Err = form.Select(form.SelectArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.conference_type"),
					Name:  "conference_type",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/conference_type"),
				},
				Value:       p.ConferenceType,
				EmptyOption: true,
				Options:     localize.VocabularySelectOptions(c.Loc, "conference_types"),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesMiscellaneousType() {
			templ_7745c5c3_Err = form.Select(form.SelectArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.miscellaneous_type"),
					Name:  "miscellaneous_type",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/miscellaneous_type"),
				},
				Value:       p.MiscellaneousType,
				EmptyOption: true,
				Options:     localize.VocabularySelectOptions(c.Loc, "miscellaneous_types"),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesDOI() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.doi"),
					Name:  "doi",
					Cols:  9,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/doi"),
					Help:  c.Loc.Get("builder.doi.help"),
				},
				Value: p.DOI,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if c.User.CanCurate() {
			templ_7745c5c3_Err = form.Select(form.SelectArgs{
				FieldArgs: form.FieldArgs{
					Label:    c.Loc.Get("builder.classification"),
					Name:     "classification",
					Cols:     3,
					Error:    localize.ValidationErrorAt(c.Loc, errors, "/classification"),
					Required: true,
				},
				Value:   p.Classification,
				Options: localize.ClassificationSelectOptions(c.Loc, p.ClassificationChoices()),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = display.Field(display.FieldArgs{
				Label:   c.Loc.Get("builder.classification"),
				Value:   c.Loc.Get("publication_classifications." + p.Classification),
				Tooltip: c.Loc.Get("tooltip.publication.classification"),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if c.User.CanCurate() {
			templ_7745c5c3_Err = form.Checkbox(form.CheckboxArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.legacy"),
					Name:  "legacy",
					Cols:  9,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/legacy"),
				},
				Value:   "true",
				Checked: p.Legacy,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li><li class=\"list-group-item\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if p.UsesTitle() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label:    c.Loc.Get("builder.title"),
					Name:     "title",
					Cols:     9,
					Error:    localize.ValidationErrorAt(c.Loc, errors, "/title"),
					Required: true,
				},
				Value: p.Title,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesAlternativeTitle() {
			templ_7745c5c3_Err = form.TextRepeat(form.TextRepeatArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.alternative_title"),
					Name:  "alternative_title",
					Cols:  9,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/alternative_title"),
				},
				Values: p.AlternativeTitle,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesPublication() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label:    c.Loc.Get(fmt.Sprintf("builder.%s.publication", p.Type)),
					Name:     "publication",
					Cols:     9,
					Error:    localize.ValidationErrorAt(c.Loc, errors, "/publication"),
					Required: p.ShowPublicationAsRequired(),
				},
				Value: p.Publication,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesPublicationAbbreviation() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get(fmt.Sprintf("builder.%s.publication_abbreviation", p.Type)),
					Name:  "publication_abbreviation",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/publication_abbreviation"),
				},
				Value: p.PublicationAbbreviation,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li><li class=\"list-group-item\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if p.UsesLanguage() {
			templ_7745c5c3_Err = form.SelectRepeat(form.SelectRepeatArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.language"),
					Name:  "language",
					Cols:  9,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/language"),
				},
				Values:      p.Language,
				EmptyOption: true,
				Options:     localize.LanguageSelectOptions(),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesPublicationStatus() {
			templ_7745c5c3_Err = form.Select(form.SelectArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.publication_status"),
					Name:  "publication_status",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/publication_status"),
				},
				Value:       p.PublicationStatus,
				EmptyOption: true,
				Options:     localize.VocabularySelectOptions(c.Loc, "publication_publishing_statuses"),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = form.Checkbox(form.CheckboxArgs{
			FieldArgs: form.FieldArgs{
				Label: c.Loc.Get("builder.extern"),
				Name:  "extern",
				Cols:  9,
				Error: localize.ValidationErrorAt(c.Loc, errors, "/extern"),
			},
			Value:   "true",
			Checked: p.Extern,
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if p.UsesYear() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label:    c.Loc.Get("builder.year"),
					Name:     "year",
					Cols:     3,
					Error:    localize.ValidationErrorAt(c.Loc, errors, "/year"),
					Required: true,
					Help:     c.Loc.Get("builder.year.help"),
				},
				Value: p.Year,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesPublisher() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.place_of_publication"),
					Name:  "place_of_publication",
					Cols:  9,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/place_of_publication"),
				},
				Value: p.PlaceOfPublication,
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
					Label: c.Loc.Get("builder.publisher"),
					Name:  "publisher",
					Cols:  9,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/publisher"),
				},
				Value: p.Publisher,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li><li class=\"list-group-item\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if p.UsesSeriesTitle() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.series_title"),
					Name:  "series_title",
					Cols:  9,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/series_title"),
				},
				Value: p.SeriesTitle,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesVolume() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.volume"),
					Name:  "volume",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/volume"),
				},
				Value: p.Volume,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesIssue() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.issue"),
					Name:  "issue",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/issue"),
				},
				Value: p.Issue,
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
					Label: c.Loc.Get("builder.issue_title"),
					Name:  "issue_title",
					Cols:  9,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/issue_title"),
				},
				Value: p.IssueTitle,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesEdition() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.edition"),
					Name:  "edition",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/edition"),
				},
				Value: p.Edition,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesPage() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.page_first"),
					Name:  "page_first",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/page_first"),
				},
				Value: p.PageFirst,
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
					Label: c.Loc.Get("builder.page_last"),
					Name:  "page_last",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/page_last"),
				},
				Value: p.PageLast,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesPageCount() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.page_count"),
					Name:  "page_count",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/page_count"),
					Help:  c.Loc.Get("builder.page_count.help"),
				},
				Value: p.PageCount,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesArticleNumber() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.article_number"),
					Name:  "article_number",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/article_number"),
				},
				Value: p.ArticleNumber,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesReportNumber() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.report_number"),
					Name:  "report_number",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/report_number"),
				},
				Value: p.ReportNumber,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if p.UsesDefense() {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<li class=\"list-group-item\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label:    c.Loc.Get("builder.defense_date"),
					Name:     "defense_date",
					Cols:     3,
					Error:    localize.ValidationErrorAt(c.Loc, errors, "/defense_date"),
					Required: p.ShowDefenseAsRequired(),
					Help:     c.Loc.Get("builder.defense_date.help"),
				},
				Value: p.DefenseDate,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label:    c.Loc.Get("builder.defense_place"),
					Name:     "defense_place",
					Cols:     3,
					Error:    localize.ValidationErrorAt(c.Loc, errors, "/defense_place"),
					Required: p.ShowDefenseAsRequired(),
				},
				Value: p.DefensePlace,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesConfirmations() {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<li class=\"list-group-item\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.RadioGroup(form.RadioGroupArgs{
				FieldArgs: form.FieldArgs{
					Label:    c.Loc.Get("builder.has_confidential_data"),
					Name:     "has_confidential_data",
					Cols:     9,
					Error:    localize.ValidationErrorAt(c.Loc, errors, "/has_confidential_data"),
					Required: true,
				},
				Value:   p.HasConfidentialData,
				Options: localize.VocabularySelectOptions(c.Loc, "confirmations"),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.RadioGroup(form.RadioGroupArgs{
				FieldArgs: form.FieldArgs{
					Label:    c.Loc.Get("builder.has_patent_application"),
					Name:     "has_patent_application",
					Cols:     9,
					Error:    localize.ValidationErrorAt(c.Loc, errors, "/has_patent_application"),
					Required: true,
				},
				Value:   p.HasPatentApplication,
				Options: localize.VocabularySelectOptions(c.Loc, "confirmations"),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.RadioGroup(form.RadioGroupArgs{
				FieldArgs: form.FieldArgs{
					Label:    c.Loc.Get("builder.has_publications_planned"),
					Name:     "has_publications_planned",
					Cols:     9,
					Error:    localize.ValidationErrorAt(c.Loc, errors, "/has_publications_planned"),
					Required: true,
				},
				Value:   p.HasPublicationsPlanned,
				Options: localize.VocabularySelectOptions(c.Loc, "confirmations"),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.RadioGroup(form.RadioGroupArgs{
				FieldArgs: form.FieldArgs{
					Label:    c.Loc.Get("builder.has_published_material"),
					Name:     "has_published_material",
					Cols:     9,
					Error:    localize.ValidationErrorAt(c.Loc, errors, "/has_published_material"),
					Required: true,
				},
				Value:   p.HasPublishedMaterial,
				Options: localize.VocabularySelectOptions(c.Loc, "confirmations"),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<li class=\"list-group-item\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if p.UsesWOS() {
			if c.User.CanCurate() {
				templ_7745c5c3_Err = form.Text(form.TextArgs{
					FieldArgs: form.FieldArgs{
						Label:   c.Loc.Get("builder.wos_type"),
						Name:    "wos_type",
						Cols:    3,
						Error:   localize.ValidationErrorAt(c.Loc, errors, "/wos_type"),
						Tooltip: c.Loc.Get("tooltip.publication.wos_type"),
					},
					Value: p.WOSType,
				}).Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				templ_7745c5c3_Err = display.Field(display.FieldArgs{
					Label:   c.Loc.Get("builder.wos_type"),
					Value:   p.WOSType,
					Tooltip: c.Loc.Get("tooltip.publication.wos_type"),
				}).Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.wos_id"),
					Name:  "wos_id",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/wos_id"),
					Help:  c.Loc.Get("builder.wos_id.help")},
				Value: p.WOSID,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesISSN() {
			templ_7745c5c3_Err = form.TextRepeat(form.TextRepeatArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.issn"),
					Name:  "issn",
					Cols:  3,
					Help:  c.Loc.Get("builder.issn.help"),
					Error: localize.ValidationErrorAt(c.Loc, errors, "/issn"),
				},
				Values: p.ISSN,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.TextRepeat(form.TextRepeatArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.eissn"),
					Name:  "eissn",
					Cols:  3,
					Help:  c.Loc.Get("builder.eissn.help"),
					Error: localize.ValidationErrorAt(c.Loc, errors, "/eissn"),
				},
				Values: p.EISSN,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesISBN() {
			templ_7745c5c3_Err = form.TextRepeat(form.TextRepeatArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.isbn"),
					Name:  "isbn",
					Cols:  3,
					Help:  c.Loc.Get("builder.isbn.help"),
					Error: localize.ValidationErrorAt(c.Loc, errors, "/isbn"),
				},
				Values: p.ISBN,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.TextRepeat(form.TextRepeatArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.eisbn"),
					Name:  "eisbn",
					Cols:  3,
					Help:  c.Loc.Get("builder.eisbn.help"),
					Error: localize.ValidationErrorAt(c.Loc, errors, "/eisbn"),
				},
				Values: p.EISBN,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesPubMedID() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.pubmed_id"),
					Name:  "pubmed_id",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/pubmed_id"),
					Help:  c.Loc.Get("builder.pubmed_id.help")},
				Value: p.PubMedID,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesArxivID() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.arxiv_id"),
					Name:  "arxiv_id",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/arxiv_id"),
					Help:  c.Loc.Get("builder.arxiv_id.help")},
				Value: p.ArxivID,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if p.UsesESCIID() {
			templ_7745c5c3_Err = form.Text(form.TextArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.esci_id"),
					Name:  "esci_id",
					Cols:  3,
					Error: localize.ValidationErrorAt(c.Loc, errors, "/esci_id"),
					Help:  c.Loc.Get("builder.esci_id.help")},
				Value: p.ESCIID,
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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(c.PathTo("publication_update_details", "id", p.ID).String()))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-headers=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf(`{"If-Match": "%s"}`, p.SnapshotID)))
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
