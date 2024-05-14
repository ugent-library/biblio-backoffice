// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/views"
	datasetviews "github.com/ugent-library/biblio-backoffice/views/dataset"
	datasetsummaryviews "github.com/ugent-library/biblio-backoffice/views/dataset/summary"
)

type AddIdentifierArgs struct {
	Dataset          *models.Dataset
	Source           string
	Identifier       string
	DuplicateDataset bool
	Errors           form.Errors
}

func AddIdentifier(c *ctx.Ctx, args AddIdentifierArgs) templ.Component {
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
			templ_7745c5c3_Err = datasetviews.AddSidebar(c, 1).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" <form class=\"w-100\" action=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 templ.SafeURL = templ.URL(c.PathTo("dataset_confirm_import").String())
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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"u-scroll-wrapper\"><div class=\"bc-navbar bc-navbar--large bc-navbar-bordered bc-navbar--white bc-navbar--bordered-bottom\"><div class=\"bc-toolbar\"><div class=\"bc-toolbar-left\"><div class=\"bc-toolbar-item\"><div class=\"d-flex flex-column\"><span class=\"text-muted\">Step 1</span><h4 class=\"bc-toolbar-title\">Add dataset</h4></div></div></div><div class=\"bc-toolbar-right\"><div class=\"bc-toolbar-item\"><button type=\"submit\" class=\"btn btn-primary\"><div class=\"btn-text\">Add dataset</div><i class=\"if if-arrow-right\"></i></button></div></div></div></div><div class=\"p-6 u-scroll-wrapper__body\"><div class=\"card mb-6\"><div class=\"card-header\"><div class=\"bc-toolbar\"><div class=\"bc-toolbar-left\"><div class=\"bc-toolbar-item\"><h5 class=\"h6\">Enter the DOI from an external repository to import the metadata of a dataset</h5></div></div></div></div><div class=\"card-body\"><div class=\"row\"><div class=\"col-6\"><div class=\"input-group\"><div class=\"input-group-prepend\"><span class=\"input-group-text\">DOI</span></div><input type=\"hidden\" name=\"source\" value=\"datacite\"> <input class=\"form-control\" type=\"text\" name=\"identifier\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(args.Identifier))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" placeholder=\"e.g. 10.7484/INSPIREHEP.DATA.RF5P.6M3K\"></div></div></div><p class=\"form-text text-muted small mt-3\">If you do not have a DOI (Digital Object Identifier), select the manual registration method in the previous step.</p></div></div><div class=\"card mb-6\"><div class=\"card-header\"><div class=\"bc-toolbar\"><div class=\"bc-toolbar-left\"><div class=\"bc-toolbar-item\"><h5 class=\"h6\">Tips for depositing and registering your data</h5></div></div></div></div><div class=\"card-body\"><ol><li class=\"mb-2\" id=\"notInExternalRepo\"><a href=\"https://onderzoektips.ugent.be/en/tips/00002071/\" target=\"_blank\">Share your data in a repository</a> <em>before</em> registering it in Biblio.<br><span class=\"text-muted\">This step will provide you with an identifier.</span></li><li class=\"mb-2\">Get more information about <a href=\"https://onderzoektips.ugent.be/en/tips/00002054/\" target=\"_blank\">dataset registration in Biblio</a>.</li><li class=\"mb-2\"><a href=\"https://onderzoektips.ugent.be/en/tips/00002055/\" target=\"_blank\">Follow a simple illustrated guide to register your dataset in Biblio</a>.</li></ol></div></div></div></div></form>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if args.DuplicateDataset {
				templ_7745c5c3_Err = addDuplicate(c, args).Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if len(args.Errors) > 0 {
				templ_7745c5c3_Err = views.FormErrorsDialog(c, "Unable to import this dataset due to the following errors", args.Errors).Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = views.PageLayout(c, "Add - Datasets - Biblio", nil).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func addDuplicate(c *ctx.Ctx, args AddIdentifierArgs) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var4 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var4 == nil {
			templ_7745c5c3_Var4 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"modal-backdrop\" class=\"d-block modal-backdrop fade show\"></div><div class=\"d-block modal show\" tabindex=\"-1\" aria-modal=\"true\" role=\"dialog\"><div class=\"modal-dialog modal-dialog-centered modal-lg modal-dialog-scrollable\" role=\"document\"><div class=\"modal-content\"><div class=\"modal-header\"><h2 class=\"modal-title\">Are you sure you want to import this dataset?</h2></div><div class=\"modal-body\"><p>Biblio contains another dataset with the same DOI:</p><ul class=\"list-group mt-6\"><li class=\"list-group-item\"><div class=\"d-flex w-100\"><div class=\"w-100\"><div class=\"d-flex align-items-start\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = datasetsummaryviews.Summary(c, datasetsummaryviews.SummaryArgs{
			Dataset: args.Dataset,
			URL:     c.PathTo("dataset", "id", args.Dataset.ID),
			Actions: datasetsummaryviews.DefaultActions(c, args.Dataset),
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div></div></li></ul></div><div class=\"modal-footer\"><div class=\"bc-toolbar\"><div class=\"bc-toolbar-left\"><button class=\"btn btn-link modal-close\">Cancel</button></div><div class=\"bc-toolbar-right\"><form action=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 templ.SafeURL = templ.URL(c.PathTo("dataset_add_import").String())
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var5)))
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<input type=\"hidden\" name=\"source\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(args.Source))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> <input type=\"hidden\" name=\"identifier\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(args.Identifier))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> <button type=\"submit\" class=\"btn btn-danger\">Import Anyway</button></form></div></div></div></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
