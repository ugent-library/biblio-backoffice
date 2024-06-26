// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"fmt"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/views"
	"github.com/ugent-library/biblio-backoffice/views/publication"
)

func AddManual(c *ctx.Ctx, step int) templ.Component {
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
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			templ_7745c5c3_Err = publication.AddSingleSidebar(c, step).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" <form class=\"w-100 u-scroll-wrapper\" action=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 templ.SafeURL = templ.URL(c.PathTo("publication_add_single_import_confirm").String())
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var3)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" method=\"POST\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = views.CSRFTag(c).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"w-100 u-scroll-wrapper\"><div class=\"bc-navbar bc-navbar--large bc-navbar--white bc-navbar--bordered-bottom\"><div class=\"bc-toolbar\"><div class=\"bc-toolbar-left\"><div class=\"bc-toolbar-item\"><div class=\"d-flex flex-column\"><span class=\"text-muted\">Step ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(step))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/pages/add_manual.templ`, Line: 20, Col: 57}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span><h4 class=\"bc-toolbar-title\">Add publication(s)</h4></div></div></div><div class=\"bc-toolbar-right\"><div class=\"bc-toolbar-item\"><button type=\"submit\" class=\"btn btn-primary\"><div class=\"btn-text\">Add publication(s)</div><i class=\"if if-arrow-right\"></i></button></div></div></div></div><div class=\"p-6 u-scroll-wrapper__body\"><div class=\"card mb-6\"><div class=\"card-header\"><div class=\"collapse-trigger\" data-bs-toggle=\"collapse\" data-bs-target=\"#notInExternalRepo\" aria-expanded=\"true\" aria-controls=\"collapse1\"></div><div class=\"bc-toolbar\"><div class=\"bc-toolbar-left\">Publication does not have an identifier. Enter manually</div></div></div><div class=\"collapse show\" id=\"notInExternalRepo\"><div class=\"card-body radio-card-group\"><h4 class=\"mb-5\">As author</h4><div class=\"row mb-5\"><div class=\"col\"><label class=\"c-radio-card c-radio-card--selected\" aria-selected=\"true\"><div class=\"c-radio-card__radio\"><div class=\"form-check\"><input class=\"form-check-input\" id=\"radio-journal-article\" type=\"radio\" name=\"publication_type\" value=\"journal_article\" checked=\"checked\"> <label class=\"form-check-label\" for=\"radio-journal-article\"></label></div></div><div class=\"c-radio-card__content\">Journal Article</div></label></div><div class=\"col\"><label class=\"c-radio-card\" aria-selected=\"false\"><div class=\"c-radio-card__radio\"><div class=\"form-check\"><input class=\"form-check-input\" id=\"radio-book-chapter\" type=\"radio\" name=\"publication_type\" value=\"book_chapter\"> <label class=\"form-check-label\" for=\"radio-book-chapter\"></label></div></div><div class=\"c-radio-card__content\">Book Chapter</div></label></div></div><div class=\"row mb-5\"><div class=\"col\"><label class=\"c-radio-card\" aria-selected=\"false\"><div class=\"c-radio-card__radio\"><div class=\"form-check\"><input class=\"form-check-input\" id=\"radio-book\" type=\"radio\" name=\"publication_type\" value=\"book\"> <label class=\"form-check-label\" for=\"radio-book\"></label></div></div><div class=\"c-radio-card__content\">Book</div></label></div><div class=\"col\"><label class=\"c-radio-card\" aria-selected=\"false\"><div class=\"c-radio-card__radio\"><div class=\"form-check\"><input class=\"form-check-input\" id=\"radio-conference\" type=\"radio\" name=\"publication_type\" value=\"conference\"> <label class=\"form-check-label\" for=\"radio-conference\"></label></div></div><div class=\"c-radio-card__content\">Conference contribution</div></label></div></div><div class=\"row mb-5\"><div class=\"col\"><label class=\"c-radio-card\" aria-selected=\"false\"><div class=\"c-radio-card__radio\"><div class=\"form-check\"><input class=\"form-check-input\" id=\"radio-dissertation\" type=\"radio\" name=\"publication_type\" value=\"dissertation\"> <label class=\"form-check-label\" for=\"radio-dissertation\"></label></div></div><div class=\"c-radio-card__content\">Dissertation</div></label></div><div class=\"col\"><label class=\"c-radio-card\" aria-selected=\"false\"><div class=\"c-radio-card__radio\"><div class=\"form-check\"><input class=\"form-check-input\" id=\"radio-miscellaneous\" type=\"radio\" name=\"publication_type\" value=\"miscellaneous\"> <label class=\"form-check-label\" for=\"radio-miscellaneous\"></label></div></div><div class=\"c-radio-card__content\">Miscellaneous</div></label></div></div><h4 class=\"mb-5\">As editor</h4><div class=\"row\"><div class=\"col\"><label class=\"c-radio-card\" aria-selected=\"false\"><div class=\"c-radio-card__radio\"><div class=\"form-check\"><input class=\"form-check-input\" id=\"radio-book-editor\" type=\"radio\" name=\"publication_type\" value=\"book_editor\"> <label class=\"form-check-label\" for=\"radio-book-editor\"></label></div></div><div class=\"c-radio-card__content\">Book</div></label></div><div class=\"col\"><label class=\"c-radio-card\" aria-selected=\"false\"><div class=\"c-radio-card__radio\"><div class=\"form-check\"><input class=\"form-check-input\" id=\"radio-issue-editor\" type=\"radio\" name=\"publication_type\" value=\"issue_editor\"> <label class=\"form-check-label\" for=\"radio-issue-editor\"></label></div></div><div class=\"c-radio-card__content\">Issue</div></label></div></div></div></div></div></div></div></form>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = views.PageLayout(c, "Add - Publications - Biblio", nil).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
