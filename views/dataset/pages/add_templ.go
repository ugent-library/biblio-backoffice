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
	"github.com/ugent-library/biblio-backoffice/views"
	datasetviews "github.com/ugent-library/biblio-backoffice/views/dataset"
)

func Add(c *ctx.Ctx) templ.Component {
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
			var templ_7745c5c3_Var3 templ.SafeURL = templ.URL(c.PathTo("dataset_add").String())
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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"u-scroll-wrapper\"><div class=\"bc-navbar bc-navbar--large bc-navbar-bordered bc-navbar--white bc-navbar--bordered-bottom\"><div class=\"bc-toolbar\"><div class=\"bc-toolbar-left\"><div class=\"bc-toolbar-item\"><div class=\"d-flex flex-column\"><span class=\"text-muted\">Step 1</span><h4 class=\"bc-toolbar-title\">Add dataset</h4></div></div></div><div class=\"bc-toolbar-right\"><div class=\"bc-toolbar-item\"><button type=\"submit\" class=\"btn btn-primary\"><div class=\"btn-text\">Add dataset</div><i class=\"if if-arrow-right\"></i></button></div></div></div></div><div class=\"p-6 u-scroll-wrapper__body\"><div class=\"card mb-6\"><div class=\"card-header\"><div class=\"bc-toolbar\"><div class=\"bc-toolbar-left\"><h5>Add a dataset</h5></div></div></div><div class=\"card-body radio-card-group\"><label class=\"c-radio-card\" aria-selected=\"false\"><div class=\"c-radio-card__radio\"><div class=\"form-check\"><input class=\"form-check-input\" id=\"add-method-identifier\" type=\"radio\" name=\"method\" value=\"identifier\"> <label class=\"form-check-label\" for=\"add-method-identifier\"></label></div></div><div class=\"c-radio-card__content\"><h6>Register your dataset via a DOI <span class=\"badge rounded-pill badge-success-light ms-3\">Recommended</span></h6><p class=\"text-muted\">Automated retrieval of metadata. For datasets that have a <a href=\"https://onderzoektips.ugent.be/en/tips/00001743/#WhatisaDigitalObjectIdentifier(DOI)?\" target=\"_blank\">DOI (Digital Object Identifier)</a>, created by the external data repository during data deposit.</p></div></label> <label class=\"c-radio-card\"><div class=\"c-radio-card__radio\" aria-selected=\"false\"><div class=\"form-check\"><input class=\"form-check-input\" id=\"add-method-manual\" type=\"radio\" name=\"method\" value=\"manual\"> <label class=\"form-check-label\" for=\"add-method-manual\"></label></div></div><div class=\"c-radio-card__content\"><h6>Register a dataset manually <span class=\"badge rounded-pill badge-default ms-3\">Beta</span></h6><p class=\"text-muted\">Manual input of metadata. Recommended for <a href=\"https://onderzoektips.ugent.be/en/tips/00001743/\" target=\"_blank\">datasets with identifiers</a> such as ENA BioProject, BioStudies, ENA, Ensembl or Handle. The identifiers are created by external data repositories during data deposit.</p></div></label></div></div><div class=\"card mb-6\"><div class=\"card-header\"><div class=\"bc-toolbar\"><div class=\"bc-toolbar-left\"><div class=\"bc-toolbar-item\"><h5 class=\"h6\">Tips for depositing and registering your data</h5></div></div></div></div><div class=\"card-body\"><ol><li class=\"mb-2\"><a href=\"https://onderzoektips.ugent.be/en/tips/00002071/\" target=\"_blank\">Share your data in a repository</a> <em>before</em> registering it in Biblio.<br><span class=\"text-muted\">This step will provide you with an identifier.</span></li><li class=\"mb-2\">Get more information about <a href=\"https://onderzoektips.ugent.be/en/tips/00002054/\" target=\"_blank\">dataset registration in Biblio</a>.</li><li class=\"mb-2\"><a href=\"https://onderzoektips.ugent.be/en/tips/00002055/\" target=\"_blank\">Follow a simple illustrated guide to register your dataset in Biblio</a>.</li></ol></div></div></div></div></form>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
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
