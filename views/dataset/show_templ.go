// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package dataset

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views"
)

func Show(c *ctx.Ctx, dataset *models.Dataset, redirectURL string) templ.Component {
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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"w-100 u-scroll-wrapper\"><div class=\"bg-white\"><div class=\"bc-navbar bc-navbar--large bc-navbar--white\"><div class=\"bc-toolbar\"><div class=\"bc-toolbar-left\"><div class=\"bc-toolbar-item\"><a class=\"btn btn-link btn-link-muted\" href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 templ.SafeURL = templ.URL(redirectURL)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var3)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><i class=\"if if-arrow-left\"></i><div class=\"btn-text\">Datasets overview</div></a></div></div><div class=\"bc-toolbar-right\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if c.User.CanWithdrawDataset(dataset) {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\"><button class=\"btn btn-outline-danger\" hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(views.URL(c.PathTo("dataset_confirm_withdraw", "id", dataset.ID)).QuerySet("redirect-url", redirectURL).String()))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><i class=\"if if-arrow-go-back\"></i> <span class=\"btn-text\">Withdraw</span></button></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if c.User.CanPublishDataset(dataset) && dataset.Status == "returned" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\"><button class=\"btn btn-success\" hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(views.URL(c.PathTo("dataset_confirm_republish", "id", dataset.ID)).QuerySet("redirect-url", redirectURL).String()))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><div class=\"btn-text\">Republish to Biblio</div></button></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if c.User.CanPublishDataset(dataset) && dataset.Status != "returned" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\"><button class=\"btn btn-success\" hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(views.URL(c.PathTo("dataset_confirm_publish", "id", dataset.ID)).QuerySet("redirect-url", redirectURL).String()))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><div class=\"btn-text\">Publish to Biblio</div></button></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if c.User.CanCurate() && dataset.Locked {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"btn btn-outline-secondary\" hx-post=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(views.URL(c.PathTo("dataset_unlock", "id", dataset.ID)).QuerySet("redirect-url", redirectURL).String()))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-swap=\"none\"><i class=\"if if-lock-unlock\"></i> <span class=\"btn-text\">Unlock record</span></button>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else if c.User.CanCurate() {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"btn btn-outline-secondary\" hx-post=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(views.URL(c.PathTo("dataset_lock", "id", dataset.ID)).QuerySet("redirect-url", redirectURL).String()))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-swap=\"none\"><i class=\"if if-lock\"></i> <span class=\"btn-text\">Lock record</span></button>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if c.User.CanDeleteDataset(dataset) {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\"><div class=\"dropdown\"><button class=\"btn btn-outline-primary btn-icon-only\" type=\"button\" data-bs-toggle=\"dropdown\" aria-haspopup=\"true\" aria-expanded=\"false\"><i class=\"if if-more\"></i></button><div class=\"dropdown-menu\"><button class=\"dropdown-item\" hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(views.URL(c.PathTo("dataset_confirm_delete", "id", dataset.ID)).QuerySet("redirect-url", redirectURL).String()))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><i class=\"if if-delete\"></i> <span>Delete</span></button></div></div></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div></div><div class=\"c-divider mt-0 mx-6 mb-4\"></div><div id=\"summary\"><div class=\"mx-6\"><div class=\"c-thumbnail-text u-min-w-0\"><div class=\"bc-toolbar bc-toolbar--auto bc-toolbar--responsive mb-3\"><div class=\"bc-toolbar-left\"><div class=\"bc-toolbar-item\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = views.BadgeStatus(dataset.Status).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if dataset.Locked {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\" data-bs-toggle=\"tooltip\" data-bs-placement=\"bottom\" data-bs-title=\"Locked for editing\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = views.BadgeLocked().Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"bc-toolbar-right\"><div class=\"bc-toolbar-item me-auto me-lg-0\"><div class=\"d-flex align-items-center flex-wrap justify-content-end\"><span class=\"c-subline text-truncate text-nowrap pe-5\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(views.CreatedBy(c, dataset.DateCreated, dataset.Creator))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 125, Col: 70}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span> <span class=\"c-subline text-truncate text-nowrap\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(views.UpdatedBy(c, dataset.DateUpdated, dataset.User, dataset.LastUser))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 128, Col: 85}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div></div></div></div><h4 class=\"list-group-item-title\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if dataset.Title != "" {
				var templ_7745c5c3_Var6 string
				templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(dataset.Title)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 136, Col: 24}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("Untitled record")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h4></div></div></div><div class=\"c-divider mb-2 mx-6 mt-4\"></div><div class=\"bc-navbar bc-navbar--auto bc-navbar--white bc-navbar--bordered-bottom\"><div class=\"bc-toolbar bc-toolbar--auto\"><div class=\"bc-toolbar-left\"><div class=\"bc-toolbar-item\" id=\"show-nav\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = SubNav(c, dataset, redirectURL).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div><div class=\"bc-toolbar-right\"></div></div></div></div><div class=\"d-flex flex-grow-1 flex-shrink-1 overflow-hidden\"><div id=\"show-sidebar\"></div><div class=\"u-scroll-wrapper__body p-6\" id=\"show-content\" hx-get=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(c.PathTo("dataset_"+c.SubNav, "id", dataset.ID).String()))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-trigger=\"load delay:100ms\" hx-target=\"#show-content\"></div></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = views.PageLayout(c, c.Loc.Get("dataset.page.show.title"), nil).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
