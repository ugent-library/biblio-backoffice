// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package publication

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views/form"
	"github.com/ugent-library/okay"
	"time"
)

func EditFileDialog(c *ctx.Ctx, p *models.Publication, f *models.PublicationFile, idx int, conflict bool, errors *okay.Errors, setAutofocus bool) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"modal-dialog modal-dialog-centered modal-lg modal-dialog-scrollable\" role=\"document\"><div class=\"modal-content\"><div class=\"modal-header h-auto py-5\"><h2 class=\"modal-title\">Document details for file ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(f.Name)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 17, Col: 62}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h2></div><div class=\"modal-body file-attributes\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if conflict {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"alert alert-danger mb-5\" role=\"alert\"><i class=\"if if--error if-error-circle-fill\"></i> The publication you are editing has been changed by someone else. Please copy your edits, then close this form.</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = form.Errors(localize.ValidationErrors(c.Loc, errors)).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form><h3 class=\"mb-3\">Document type</h3>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var3 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			var templ_7745c5c3_Var4 = []any{"form-select", "form-control", templ.KV("is-invalid", errors != nil && errors.Get(fmt.Sprintf("/file/%d/relation", idx)) != nil)}
			templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var4...)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<select class=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(templ.CSSClasses(templ_7745c5c3_Var4).String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 1, Col: 0}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" name=\"relation\" id=\"relation\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if setAutofocus {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" autofocus")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" hx-get=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var6 string
			templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("publication_edit_file_refresh_form", "id", p.ID, "file_id", f.ID).String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 40, Col: 100}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-swap=\"none\" hx-include=\".file-attributes\" hx-indicator=\".modal-dialog .spinner-border\" hx-trigger=\"change delay:50ms\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, o := range localize.VocabularySelectOptions(c.Loc, "publication_file_relations") {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var7 string
				templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(o.Value)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 47, Col: 31}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if o.Value == f.Relation {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" selected")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var8 string
				templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(o.Label)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 47, Col: 79}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
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
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = form.Field(form.FieldArgs{
			Label: c.Loc.Get("builder.file.relation"),
			Name:  "relation",
			Error: localize.ValidationErrorAt(c.Loc, errors, fmt.Sprintf("/file/%d/relation", idx)),
			Theme: form.ThemeVertical,
		}, "relation").Render(templ.WithChildren(ctx, templ_7745c5c3_Var3), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if f.Relation == "main_file" {
			templ_7745c5c3_Err = form.Select(form.SelectArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.file.publication_version"),
					Name:  "publication_version",
					Error: localize.ValidationErrorAt(c.Loc, errors, fmt.Sprintf("/file/%d/publication_version", idx)),
					Help:  c.Loc.Get("builder.file.publication_version.help"),
					Theme: form.ThemeVertical,
				},
				Value:       f.PublicationVersion,
				EmptyOption: true,
				Options:     localize.VocabularySelectOptions(c.Loc, "publication_versions"),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"c-divider\"></div><h3 class=\"mb-3\">Who can access this document?</h3>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if f.Relation == "main_file" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"alert alert-success mt-3 mb-5\"><i class=\"if if-check-circle\"></i><div class=\"alert-content\"><p><a href=\"https://www.ugent.be/intranet/en/research/impact/schol-publishing/oa-belgian-law.htm\" target=\"_blank\">Belgian open access law</a> gives you the right to make the author accepted manuscript (AAM) of scientific journal articles publicly available after embargo.<br><small>For articles published in 2023 or later, <a href=\"https://www.ugent.be/intranet/en/research/impact/schol-publishing/policy-ugent.htm#OpenAccess(OA)\" target=\"_blank\">UGent policy</a> assumes you want to make use of this right, unless you opt out by sending us an email at <a href=\"mailto:biblio@ugent.be\">biblio@ugent.be</a>.</small></p></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"mb-6 radio-card-group\"><label class=\"col-form-label\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var9 string
		templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("builder.file.access_level"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 81, Col: 47}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" <a class=\"text-muted px-2\" href=\"#\" data-bs-container=\"body\" data-bs-toggle=\"popover-custom\" data-bs-placement=\"right\" data-popover-content=\"#access-level-tooltip\"><i class=\"if if-info-circle if--small\"></i> <span class=\"visually-hidden\">More info</span></a><div class=\"u-hidden\" id=\"access-level-tooltip\"><div class=\"popover-body p-0\">Full text files are set to <strong>UGent access</strong> by default. However, you are strongly encouraged to make publications available in open access where possible.<br><a class=\"link-primary\" target=\"_blank\" href=\"https://onderzoektips.ugent.be/en/tips/00002074/\"><span class=\"text-decoration-underline\">More info</span> <i class=\"if if--small if-external-link\"></i></a></div></div></label> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, o := range localize.VocabularySelectOptions(c.Loc, "publication_file_access_levels") {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<label")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if o.Value == f.AccessLevel {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" class=\"c-radio-card c-radio-card--selected\" aria-selected=\"true\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" class=\"c-radio-card\" aria-selected=\"false\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" hx-get=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var10 string
			templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("publication_edit_file_refresh_form", "id", p.ID, "file_id", f.ID).String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 109, Col: 101}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-swap=\"none\" hx-include=\".file-attributes\" hx-indicator=\".modal-dialog .spinner-border\" hx-trigger=\"click delay:50ms\"><div class=\"c-radio-card__radio\"><div class=\"form-check\"><input class=\"form-check-input\" id=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var11 string
			templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("access-level-%s", o.Value))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 117, Col: 86}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" type=\"radio\" name=\"access_level\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var12 string
			templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(o.Value)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 117, Col: 137}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if o.Value == f.AccessLevel {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" checked")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("> <label class=\"form-check-label\" for=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var13 string
			templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("access-level-%s", o.Value))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 118, Col: 87}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></label></div></div><div class=\"c-radio-card__content d-flex align-content-center\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			switch o.Value {
			case "info:eu-repo/semantics/openAccess":
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-download text-success align-self-center\"></i><div class=\"ms-5\"><p class=\"mb-1 me-3\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var14 string
				templ_7745c5c3_Var14, templ_7745c5c3_Err = templ.JoinStringErrs(o.Label)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 126, Col: 42}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var14))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p><div><span class=\"badge rounded-pill badge-success-light\"><i class=\"if if-info-circle\"></i> <span class=\"badge-text\">Recommended if legally possible</span></span></div><span class=\"text-muted c-body-small\">Your file will be immediately available to anyone. Select \"Local access – UGent only\" if you are unsure.</span></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			case "info:eu-repo/semantics/embargoedAccess":
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-time text-muted align-self-center\"></i><div class=\"ms-5\"><p class=\"mb-1 me-3\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var15 string
				templ_7745c5c3_Var15, templ_7745c5c3_Err = templ.JoinStringErrs(o.Label)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 138, Col: 42}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var15))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p><span class=\"text-muted c-body-small\">Use to switch access levels after a specified embargo period. UGent selects this by default to open up the author accepted manuscript (AAM) of journal articles published since 2023.</span></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			case "info:eu-repo/semantics/restrictedAccess":
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-ghent-university text-primary align-self-center\"></i><div class=\"ms-5\"><p class=\"mb-1 me-3\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var16 string
				templ_7745c5c3_Var16, templ_7745c5c3_Err = templ.JoinStringErrs(o.Label)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 144, Col: 42}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var16))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p><span class=\"text-muted c-body-small\">Your file will be available to users within the UGent network only. Minimum expected by UGent policy.</span></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			case "info:eu-repo/semantics/closedAccess":
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-eye-off text-muted align-self-center\"></i><div class=\"ms-5\"><p class=\"mb-1 me-3\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var17 string
				templ_7745c5c3_Var17, templ_7745c5c3_Err = templ.JoinStringErrs(o.Label)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 150, Col: 42}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var17))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p><span class=\"text-muted c-body-small\">Only the authors and Biblio staff can access your file. Others will see metadata only. Use by exception.</span></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></label>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if f.AccessLevel == "info:eu-repo/semantics/embargoedAccess" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<h3 class=\"mb-3\">Set your embargo</h3><p class=\"mb-5\">More information about <a href=\"https://onderzoektips.ugent.be/en/tips/00002097\" target=\"_blank\">embargoes</a>.</p><div class=\"row\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.Select(form.SelectArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.file.access_level_during_embargo"),
					Name:  "access_level_during_embargo",
					Error: localize.ValidationErrorAt(c.Loc, errors, fmt.Sprintf("/file/%d/access_level_during_embargo", idx)),
					Theme: form.ThemeVertical,
				},
				Value:       f.AccessLevelDuringEmbargo,
				EmptyOption: true,
				Options:     localize.VocabularySelectOptions(c.Loc, "publication_file_access_levels_during_embargo"),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.Select(form.SelectArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.file.access_level_after_embargo"),
					Name:  "access_level_after_embargo",
					Error: localize.ValidationErrorAt(c.Loc, errors, fmt.Sprintf("/file/%d/access_level_after_embargo", idx)),
					Help:  c.Loc.Get("builder.file.access_level_after_embargo.help"),
					Theme: form.ThemeVertical,
				},
				Value:       f.AccessLevelAfterEmbargo,
				EmptyOption: true,
				Options:     localize.VocabularySelectOptions(c.Loc, "publication_file_access_levels_after_embargo"),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"row\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.Date(form.DateArgs{
				FieldArgs: form.FieldArgs{
					Label: c.Loc.Get("builder.file.embargo_date"),
					Name:  "embargo_date",
					Error: localize.ValidationErrorAt(c.Loc, errors, fmt.Sprintf("/file/%d/embargo_date", idx)),
					Theme: form.ThemeVertical,
					Cols:  6,
					Help:  c.Loc.Get("builder.file.embargo_date.help"),
				},
				Value: f.EmbargoDate,
				Min:   nextDay().Format("2006-01-02"),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = form.Select(form.SelectArgs{
			FieldArgs: form.FieldArgs{
				Label:   c.Loc.Get("builder.file.license"),
				Name:    "license",
				Error:   localize.ValidationErrorAt(c.Loc, errors, fmt.Sprintf("/file/%d/license", idx)),
				Tooltip: c.Loc.Get("tooltip.publication.file.license"),
				Theme:   form.ThemeVertical,
				Cols:    6,
			},
			Value:       f.License,
			EmptyOption: true,
			Options:     localize.VocabularySelectOptions(c.Loc, "publication_licenses"),
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</form></div><div class=\"modal-footer\"><div class=\"spinner-border\"><span class=\"visually-hidden\"></span></div><div class=\"bc-toolbar\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if conflict {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-left\"><button class=\"btn btn-primary modal-close\">Close</button></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-left\"><button class=\"btn btn-link modal-close\" hx-get=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var18 string
			templ_7745c5c3_Var18, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("publication_refresh_files", "id", p.ID).String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 229, Col: 75}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var18))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-swap=\"none\">Cancel</button></div><div class=\"bc-toolbar-right\"><button type=\"button\" name=\"create\" class=\"btn btn-primary\" hx-put=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var19 string
			templ_7745c5c3_Var19, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("publication_update_file", "id", p.ID, "file_id", f.ID).String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 238, Col: 90}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var19))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-headers=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var20 string
			templ_7745c5c3_Var20, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf(`{"If-Match": "%s"}`, p.SnapshotID))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_file.templ`, Line: 239, Col: 68}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var20))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-include=\".file-attributes\" hx-swap=\"none\">Save</button></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func nextDay() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)
}
