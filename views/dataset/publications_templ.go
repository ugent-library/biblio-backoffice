// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package dataset

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	publicationsummaryviews "github.com/ugent-library/biblio-backoffice/views/publication/summary"
)

const PublicationsBodySelector = "#publications-body"

func Publications(c *ctx.Ctx, dataset *models.Dataset, publications []*models.Publication) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div hx-swap-oob=\"innerHTML:#show-nav\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = SubNav(c, dataset, c.CurrentURL.Query().Get("redirect-url")).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div hx-swap-oob=\"innerHTML:#show-sidebar\"></div><div class=\"card mb-6\"><div class=\"card-header\"><div class=\"bc-toolbar\"><div class=\"bc-toolbar-left\"><h1 class=\"bc-toolbar-title\">Related publications</h1><a class=\"text-muted px-2\" href=\"#\" data-bs-container=\"body\" data-bs-toggle=\"popover\" data-bs-placement=\"right\" data-bs-content=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("tooltip.dataset.related_publications"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/publications.templ`, Line: 27, Col: 73}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><i class=\"if if-info-circle if--small\"></i> <span class=\"visually-hidden\">More info</span></a></div><div class=\"bc-toolbar-right\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if c.Repo.CanEditDataset(c.User, dataset) {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"c-button-toolbar\" data-panel-state=\"read\"><button class=\"btn btn-outline-primary\" hx-get=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("dataset_add_publication", "id", dataset.ID).String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/publications.templ`, Line: 38, Col: 79}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><i class=\"if if-add\"></i><div class=\"btn-text\">Add publication</div></button></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div></div><div id=\"publications-body\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = PublicationsBody(c, dataset, publications).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func PublicationsBody(c *ctx.Ctx, dataset *models.Dataset, relatedPublications []*models.Publication) templ.Component {
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
		templ_7745c5c3_Var4 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var4 == nil {
			templ_7745c5c3_Var4 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if len(relatedPublications) > 0 {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"card-body p-0\"><ul class=\"list-group list-group-flush\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, pub := range relatedPublications {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<li class=\"list-group-item\"><div class=\"d-flex w-100\"><div class=\"w-100\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if !c.Repo.CanViewPublication(c.User, pub) && pub.Status == "returned" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"hstack-md-responsive align-items-start gap-3 w-100\"><div class=\"vstack gap-2\"><div class=\"d-inline-flex align-items-center flex-wrap\"><span class=\"badge rounded-pill badge-danger-light me-3 my-2\"><span class=\"badge-circle\"></span> <span class=\"badge-text\">Biblio withdrawn</span></span></div><h4 class=\"mb-0 list-group-item-title\">You cannot see the metadata of linked withdrawn publications you do not own.</h4><p class=\"opacity-50\">The link will not be visible on biblio.ugent.be while withdrawn. When the publication is republished, you will see the metadata again and the link will be visible on biblio.ugent.be.</p></div><div class=\"bc-toolbar-item ms-auto ms-lg-0\"><button class=\"btn btn-tertiary\" type=\"button\" hx-get=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var5 string
					templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("dataset_confirm_delete_publication", "id", dataset.ID, "snapshot_id", dataset.SnapshotID, "publication_id", pub.ID).String())
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/publications.templ`, Line: 83, Col: 155}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><i class=\"if if-delete\"></i> <span class=\"btn-text\">Remove link</span></button></div></div>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else {
					templ_7745c5c3_Err = publicationsummaryviews.Summary(c, publicationsummaryviews.SummaryArgs{
						Publication: pub,
						URL:         c.PathTo("publication", "id", pub.ID),
						Actions:     publicationSummaryActions(c, dataset, pub),
					}).Render(ctx, templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div></li>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</ul></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"card-body\"><span class=\"text-muted\">No publications</span></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		return templ_7745c5c3_Err
	})
}

func publicationSummaryActions(c *ctx.Ctx, dataset *models.Dataset, pub *models.Publication) templ.Component {
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
		templ_7745c5c3_Var6 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var6 == nil {
			templ_7745c5c3_Var6 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if c.Repo.CanViewPublication(c.User, pub) || c.Repo.CanEditDataset(c.User, dataset) {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item ms-auto ms-lg-0\"><div class=\"c-button-toolbar\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if c.Repo.CanViewPublication(c.User, pub) {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<a class=\"btn btn-tertiary\" href=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var7 templ.SafeURL = templ.URL(c.PathTo("publication", "id", pub.ID).String())
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var7)))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><i class=\"if if-eye\"></i> <span class=\"btn-text\">View</span></a> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if c.Repo.CanEditDataset(c.User, dataset) {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"btn btn-tertiary\" type=\"button\" hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var8 string
				templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("dataset_confirm_delete_publication", "id", dataset.ID, "snapshot_id", dataset.SnapshotID, "publication_id", pub.ID).String())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/publications.templ`, Line: 128, Col: 149}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><i class=\"if if-delete\"></i> <span class=\"btn-text\">Remove link</span></button>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		return templ_7745c5c3_Err
	})
}
