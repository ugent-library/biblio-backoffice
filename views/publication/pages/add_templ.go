// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/views"
	"github.com/ugent-library/biblio-backoffice/views/publication"
)

func addPublicationBreadcrumbs(c *ctx.Ctx) []views.Breadcrumb {
	return []views.Breadcrumb{
		{LabelID: "publications", URL: c.PathTo("publications")},
		{LabelID: "publication_add"},
	}
}

func Add(c *ctx.Ctx, step int) templ.Component {
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
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
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
			templ_7745c5c3_Err = publication.AddMultipleSidebar(c, step).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" <div class=\"w-100 u-scroll-wrapper\"><div class=\"bc-navbar bc-navbar--large bc-navbar-bordered bc-navbar--white bc-navbar--bordered-bottom\"><div class=\"bc-toolbar\"><div class=\"bc-toolbar-left\"><div class=\"bc-toolbar-item\"><div class=\"d-flex flex-column\"><span class=\"text-muted\">Step ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(step))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/pages/add.templ`, Line: 29, Col: 56}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span><h4 class=\"bc-toolbar-title\">Start: add publication(s)</h4></div></div></div><div class=\"bc-toolbar-right\"><div class=\"bc-toolbar-item\"><a class=\"btn btn-tertiary\" href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 templ.SafeURL = templ.URL(c.PathTo("publications").String())
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var4)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">Cancel</a></div></div></div></div><div class=\"p-6 u-scroll-wrapper__body\"><div class=\"row\"><div class=\"col-xl-9 mb-6\"><div class=\"card mb-4\"><div class=\"card-body\"><div class=\"bc-toolbar h-auto\"><div class=\"bc-toolbar-left\"><div><div class=\"d-flex align-items-center flex-wrap\"><h3 class=\"me-3\">Import from Web of Science</h3><span class=\"badge badge-default\">Recommended for records in WoS</span></div><p class=\"text-muted\">Import one or more publications. This option saves you the most time.</p></div></div><div class=\"bc-toolbar-right\"><a class=\"btn btn-primary ms-6\" href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 templ.SafeURL = templ.URL(c.PathTo("publication_add", "method", "wos").String())
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var5)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><i class=\"if if-add\"></i> <span class=\"btn-text\">Add</span></a></div></div></div></div><div class=\"card mb-4\"><div class=\"card-body\"><div class=\"bc-toolbar h-auto\"><div class=\"bc-toolbar-left\"><div><div class=\"d-flex align-items-center flex-wrap\"><h3>Import your publication via an identifier</h3></div><p class=\"text-muted\">Use DOI, PubMed ID or arXiv ID. A good second option.</p></div></div><div class=\"bc-toolbar-right\"><a class=\"btn btn-primary ms-6\" href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var6 templ.SafeURL = templ.URL(c.PathTo("publication_add", "method", "identifier").String())
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var6)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><i class=\"if if-add\"></i> <span class=\"btn-text\">Add</span></a></div></div></div></div><div class=\"card mb-4\"><div class=\"card-body\"><div class=\"bc-toolbar h-auto\"><div class=\"bc-toolbar-left\"><div><div class=\"d-flex align-items-center flex-wrap\"><h3>Enter a publication manually</h3></div><p class=\"text-muted\">Create a publication record from scratch using a template. Recommended for publications such as dissertations.</p></div></div><div class=\"bc-toolbar-right\"><a class=\"btn btn-primary ms-6\" href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var7 templ.SafeURL = templ.URL(c.PathTo("publication_add", "method", "manual").String())
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var7)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><i class=\"if if-add\"></i> <span class=\"btn-text\">Add</span></a></div></div></div></div><div class=\"card mb-4\"><div class=\"card-body\"><div class=\"bc-toolbar h-auto\"><div class=\"bc-toolbar-left\"><div><div class=\"d-flex align-items-center flex-wrap\"><h3>Import via BibTeX file</h3></div><p class=\"text-muted\">Import multiple publications via library files. Use this options if there is no Web of Science import or identifier import available.</p></div></div><div class=\"bc-toolbar-right\"><a class=\"btn btn-primary ms-6\" href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var8 templ.SafeURL = templ.URL(c.PathTo("publication_add", "method", "bibtex").String())
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var8)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><i class=\"if if-add\"></i> <span class=\"btn-text\">Add</span></a></div></div></div></div></div><div class=\"col-xl-3\"><div class=\"card bg-light\"><div class=\"card-header bg-transparent\"><div class=\"card-title\">Tips for adding your publications</div></div><div class=\"card-body\"><div class=\"c-content\"><ol><li class=\"mb-4\">Follow a <a class=\"c-link\" href=\"https://onderzoektips.ugent.be/en/tips/00002065/\" target=\"_blank\">step by step guide</a> about deposit and registration of publications.</li><li>Read general <a class=\"c-link\" href=\"https://onderzoektips.ugent.be/en/tips/00002064/\" target=\"_blank\">documentation</a> about deposit and registration of publications.</li></ol></div></div></div></div></div></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = views.PageLayout(c, views.PageLayoutArgs{
			Title:       "Add - Publications - Biblio",
			Breadcrumbs: addPublicationBreadcrumbs(c),
		}).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
