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
)

type EditMessageDialogArgs struct {
	Publication *models.Publication
	Errors      *okay.Errors
	Conflict    bool
}

func EditMessageDialog(c *ctx.Ctx, args EditMessageDialogArgs) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"modal-dialog modal-dialog-centered modal-xl modal-dialog-scrollable\" role=\"document\"><div class=\"modal-content\"><div class=\"modal-header\"><h2 class=\"modal-title\">Edit messages from and for Biblio team</h2></div><div class=\"modal-body\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if args.Conflict {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"alert alert-danger mb-0\" role=\"alert\"><i class=\"if if--error if-error-circle-fill\"></i> The publication you are editing has been changed by someone else. Please copy your edits, then close this form.</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = form.Errors(localize.ValidationErrors(c.Loc, args.Errors)).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<ul class=\"list-group list-group-flush\" data-panel-state=\"edit\"><li class=\"list-group-item\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = form.TextArea(form.TextAreaArgs{
			FieldArgs: form.FieldArgs{
				Label:     c.Loc.Get("builder.message"),
				Name:      "message",
				Error:     localize.ValidationErrorAt(c.Loc, args.Errors, "/message"),
				Theme:     form.ThemeVertical,
				AutoFocus: true,
			},
			Value: args.Publication.Message,
			Rows:  4,
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li></ul></div><div class=\"modal-footer\"><div class=\"bc-toolbar\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if args.Conflict {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-left\"><button class=\"btn btn-primary modal-close\">Close</button></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-left\"><button class=\"btn btn-link modal-close\">Cancel</button></div><div class=\"bc-toolbar-right\"><button type=\"button\" name=\"create\" class=\"btn btn-primary\" hx-put=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("publication_update_message", "id", args.Publication.ID).String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_message.templ`, Line: 63, Col: 91}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-headers=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf(`{"If-Match": "%s"}`, args.Publication.SnapshotID))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/edit_message.templ`, Line: 64, Col: 83}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
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
		return templ_7745c5c3_Err
	})
}
