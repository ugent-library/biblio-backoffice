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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<ul class=\"list-group list-group-flush\" data-panel-state=\"edit\"><li class=\"list-group-item\"></li><li class=\"list-group-item\">")
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</ul></div><div class=\"modal-footer\"><div class=\"bc-toolbar\">")
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
