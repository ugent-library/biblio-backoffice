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
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views/form"
	"github.com/ugent-library/okay"
)

type EditAbstractDialogArgs struct {
	IsNew    bool
	Dataset  *models.Dataset
	Abstract models.Text
	Index    int
	Conflict bool
	Errors   *okay.Errors
}

func EditAbstractDialog(c *ctx.Ctx, args EditAbstractDialogArgs) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"modal-dialog modal-dialog-centered modal-lg modal-dialog-scrollable\" role=\"document\"><div class=\"modal-content\"><div class=\"modal-header\"><h2 class=\"modal-title\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if args.IsNew {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("Add abstract")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("Edit abstract")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h2></div><div class=\"modal-body\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if args.Conflict {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"alert alert-danger mb-0\" role=\"alert\"><i class=\"if if--error if-error-circle-fill\"></i> The dataset you are editing has been changed by someone else. Please copy your edits, then close this form.</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = form.Errors(localize.ValidationErrors(c.Loc, args.Errors)).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = form.TextArea(form.TextAreaArgs{
			FieldArgs: form.FieldArgs{
				Label: c.Loc.Get("builder.abstract.text"),
				Name:  "text",
				Error: localize.ValidationErrorAt(c.Loc, args.Errors, fmt.Sprintf("/abstract/%d/text", args.Index)),
				Theme: form.ThemeVertical,
			},
			Value: args.Abstract.Text,
			Rows:  6,
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = form.Select(form.SelectArgs{
			FieldArgs: form.FieldArgs{
				Label: c.Loc.Get("builder.abstract.lang"),
				Name:  "lang",
				Error: localize.ValidationErrorAt(c.Loc, args.Errors, fmt.Sprintf("/abstract/%d/lang", args.Index)),
				Theme: form.ThemeVertical,
			},
			Value:   args.Abstract.Lang,
			Options: localize.LanguageSelectOptions(),
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"modal-footer\"><div class=\"bc-toolbar\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if args.Conflict {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-left\"><button class=\"btn btn-primary modal-close\">Close</button></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-left\"><button class=\"btn btn-link modal-close\">Cancel</button></div><div class=\"bc-toolbar-right\"><button type=\"button\" name=\"create\" class=\"btn btn-primary\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if args.IsNew {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" hx-post=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(c.PathTo("dataset_create_abstract", "id", args.Dataset.ID).String()))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" hx-put=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(c.PathTo("dataset_update_abstract", "id", args.Dataset.ID, "abstract_id", args.Abstract.ID).String()))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" hx-headers=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf(`{"If-Match": "%s"}`, args.Dataset.SnapshotID)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-include=\".modal-body\" hx-swap=\"none\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if args.IsNew {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("Add abstract")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("Update abstract")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</button></div>")
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
