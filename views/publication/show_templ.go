// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package publication

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views"
)

func Show(c *ctx.Ctx, p *models.Publication, redirectURL string) templ.Component {
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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"w-100 u-scroll-wrapper\"><div class=\"bg-white\" data-collapsed=\"true\"><div class=\"bc-navbar bc-navbar--large bc-navbar--white\"><div class=\"bc-toolbar\"><div class=\"bc-toolbar-left\"><div class=\"bc-toolbar-item\"><a class=\"btn btn-link btn-link-muted\" href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 templ.SafeURL = templ.URL(redirectURL)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var3)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><i class=\"if if-arrow-left\"></i><div class=\"btn-text\">Publications overview</div></a></div></div><div class=\"bc-toolbar-right\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if c.User.CanWithdrawPublication(p) {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\"><button class=\"btn btn-outline-danger\" hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var4 string
				templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("publication_confirm_withdraw", "id", p.ID, "redirect-url", c.CurrentURL.String()).String())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/show.templ`, Line: 28, Col: 119}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><i class=\"if if-arrow-go-back\"></i> <span class=\"btn-text\">Withdraw</span></button></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if c.User.CanPublishPublication(p) && p.Status == "returned" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\"><button class=\"btn btn-success\" hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var5 string
				templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("publication_confirm_republish", "id", p.ID, "redirect-url", c.CurrentURL.String()).String())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/show.templ`, Line: 40, Col: 120}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><div class=\"btn-text\">Republish to Biblio</div></button></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if c.User.CanPublishPublication(p) && p.Status != "returned" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\"><button class=\"btn btn-success\" hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var6 string
				templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("publication_confirm_publish", "id", p.ID, "redirect-url", c.CurrentURL.String()).String())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/show.templ`, Line: 51, Col: 118}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
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
			if c.User.CanCurate() && p.Locked {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"btn btn-outline-secondary\" hx-post=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var7 string
				templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("publication_unlock", "id", p.ID, "redirect-url", c.CurrentURL.String()).String())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/show.templ`, Line: 62, Col: 110}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
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
				var templ_7745c5c3_Var8 string
				templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("publication_lock", "id", p.ID, "redirect-url", c.CurrentURL.String()).String())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/show.templ`, Line: 71, Col: 108}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
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
			if c.User.CanDeletePublication(p) {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\"><div class=\"dropdown dropleft\"><button class=\"btn btn-outline-primary btn-icon-only\" type=\"button\" data-bs-toggle=\"dropdown\" aria-haspopup=\"true\" aria-expanded=\"false\"><i class=\"if if-more\"></i></button><div class=\"dropdown-menu\"><a class=\"dropdown-item\" href=\"#\" hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var9 string
				templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("publication_confirm_delete", "id", p.ID, "redirect-url", redirectURL).String())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/show.templ`, Line: 95, Col: 109}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><i class=\"if if-delete\"></i> <span>Delete</span></a></div></div></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div></div><div class=\"c-divider mt-0 mx-6 mb-4\"></div><div id=\"summary\"><div class=\"mx-6\"><div class=\"c-thumbnail-text u-min-w-0\"><div class=\"bc-toolbar bc-toolbar--auto bc-toolbar--responsive mb-3\"><div class=\"bc-toolbar-left\"><div class=\"bc-toolbar-item\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = views.BadgeStatus(p.Status).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if p.PublicationStatus != "" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\"><p class=\"c-subline text-nowrap pe-5 border-end\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var10 string
				templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_publishing_statuses." + p.PublicationStatus))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/show.templ`, Line: 119, Col: 129}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if p.Locked {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\" data-bs-toggle=\"tooltip\" data-bs-title=\"Locked for editing\" data-bs-placement=\"bottom\"><span class=\"c-subline text-nowrap pe-5 border-end\"><i class=\"if if-lock if--small if--muted\"></i> <span class=\"text-muted c-body-small ms-2\">Locked</span></span></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\"><span class=\"c-subline text-nowrap\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var11 string
			templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_types." + p.Type))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/show.templ`, Line: 132, Col: 53}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if p.Classification != "" {
				var templ_7745c5c3_Var12 string
				templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(p.Classification)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/show.templ`, Line: 134, Col: 30}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div></div><div class=\"bc-toolbar-right\"><div class=\"bc-toolbar-item me-auto me-lg-0\"><div class=\"d-flex align-items-center flex-wrap justify-content-end\"><span class=\"c-subline text-truncate text-nowrap\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var13 string
			templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(views.CreatedBy(c, p.DateCreated, p.Creator))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/show.templ`, Line: 143, Col: 58}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span> <span class=\"c-subline text-truncate text-nowrap ps-5\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var14 string
			templ_7745c5c3_Var14, templ_7745c5c3_Err = templ.JoinStringErrs(views.UpdatedBy(c, p.DateUpdated, p.User, p.LastUser))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/show.templ`, Line: 146, Col: 67}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var14))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div></div></div></div><h4 class=\"list-group-item-title\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if p.Title != "" {
				var templ_7745c5c3_Var15 string
				templ_7745c5c3_Var15, templ_7745c5c3_Err = templ.JoinStringErrs(p.Title)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/show.templ`, Line: 154, Col: 18}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var15))
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
			templ_7745c5c3_Err = SubNav(c, p, redirectURL).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div><div class=\"bc-toolbar-right\"></div></div></div></div><div class=\"d-flex flex-grow-1 flex-shrink-1 overflow-hidden\"><div id=\"show-sidebar\"></div><div class=\"u-scroll-wrapper__body p-6\" id=\"show-content\" hx-get=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var16 string
			templ_7745c5c3_Var16, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("publication_"+c.SubNav, "id", p.ID, "redirect-url", redirectURL).String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/show.templ`, Line: 179, Col: 97}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var16))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-trigger=\"load delay:100ms\" hx-target=\"#show-content\"></div></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = views.PageLayout(c, c.Loc.Get("publication.page.show.title"), nil).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
