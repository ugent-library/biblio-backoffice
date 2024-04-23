// Code generated by templ@v0.2.334 DO NOT EDIT.

package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
)

func AddImpersonationSuggest(c *ctx.Ctx, firstName string, lastName string, hits []*models.Person) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_1 := templ.GetChildren(ctx)
		if var_1 == nil {
			var_1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if len(hits) > 0 {
			_, err = templBuffer.WriteString("<h3 class=\"mt-6\">")
			if err != nil {
				return err
			}
			var_2 := `Search results`
			_, err = templBuffer.WriteString(var_2)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h3> <ul class=\"list-group\">")
			if err != nil {
				return err
			}
			for _, hit := range hits {
				_, err = templBuffer.WriteString("<li class=\"list-group-item\"><div class=\"list-group-item-inner\"><div class=\"list-group-item-main\">")
				if err != nil {
					return err
				}
				err = contributorShowSummary(c, models.ContributorFromPerson(hit)).Render(ctx, templBuffer)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div><div class=\"c-button-toolbar\"><form action=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(c.PathTo("create_impersonation").String()))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\" method=\"POST\">")
				if err != nil {
					return err
				}
				err = csrfField(c).Render(ctx, templBuffer)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("<input type=\"hidden\" name=\"id\" value=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(hit.ID))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"><button type=\"submit\" class=\"btn btn-primary\">")
				if err != nil {
					return err
				}
				var_3 := `Change user`
				_, err = templBuffer.WriteString(var_3)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</button></form></div></div></li>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</ul>")
			if err != nil {
				return err
			}
		} else if firstName != "" || lastName != "" {
			_, err = templBuffer.WriteString("<div class=\"c-blank-slate c-blank-slate-muted c-blank-slate-large\"><div class=\"bc-avatar bc-avatar--small\"><i class=\"if if-info-circle\"></i></div><p>")
			if err != nil {
				return err
			}
			var_4 := `No users found.`
			_, err = templBuffer.WriteString(var_4)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</p></div>")
			if err != nil {
				return err
			}
		} else {
			_, err = templBuffer.WriteString("<div class=\"c-blank-slate c-blank-slate-muted c-blank-slate-large\"><div class=\"bc-avatar bc-avatar--small\"><i class=\"if if-info-circle\"></i></div><p>")
			if err != nil {
				return err
			}
			var_5 := `Type a first or last name above.`
			_, err = templBuffer.WriteString(var_5)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</p></div>")
			if err != nil {
				return err
			}
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
