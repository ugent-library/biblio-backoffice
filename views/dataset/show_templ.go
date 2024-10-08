// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package dataset

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views"
)

func Show(c *ctx.Ctx, dataset *models.Dataset, redirectURL string) templ.Component {
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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"w-100 u-scroll-wrapper\"><div class=\"bg-white\"><div id=\"summary\"><div class=\"mx-6\"><div class=\"bc-toolbar bc-toolbar-md-responsive flex-column-reverse flex-md-row w-100\"><div class=\"bc-toolbar-left\"><div class=\"d-inline-flex align-items-center flex-wrap\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = views.BadgeStatus(dataset.Status).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if dataset.Locked {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"c-subline text-nowrap me-3 pe-3 border-end\" data-bs-toggle=\"tooltip\" data-bs-title=\"Locked for editing\" data-bs-placement=\"bottom\"><i class=\"if if-lock if--small if--muted\"></i> <span class=\"text-muted c-body-small ms-2\">Locked</span></span> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if dataset.AccessLevel != "" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"c-subline pe-3\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if dataset.AccessLevel == "info:eu-repo/semantics/openAccess" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-download if--small if--success\"></i> <span class=\"c-subline text-truncate ps-2\">")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var3 string
					templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("dataset_access_levels." + dataset.AccessLevel))
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 37, Col: 72}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else if dataset.AccessLevel == "info:eu-repo/semantics/embargoedAccess" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-time if--small if--muted\"></i> <span class=\"c-subline text-truncate ps-2\">")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var4 string
					templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("dataset_access_levels." + dataset.AccessLevel))
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 42, Col: 72}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else if dataset.AccessLevel == "info:eu-repo/semantics/restrictedAccess" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-ghent-university if--small if--primary\"></i> <span class=\"c-subline text-muted ps-2\">")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var5 string
					templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("dataset_access_levels." + dataset.AccessLevel))
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 47, Col: 72}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else if dataset.AccessLevel == "info:eu-repo/semantics/closedAccess" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-eye-off if--small if--muted\"></i> <span class=\"c-subline text-muted ps-2\">")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var6 string
					templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("dataset_access_levels." + dataset.AccessLevel))
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 52, Col: 72}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if dataset.AccessLevel == "info:eu-repo/semantics/embargoedAccess" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"c-subline me-3 px-3 border-end border-start\"><i class=\"if if-eye-off if--small if--muted\"></i> <span class=\"c-subline text-muted ps-2\">Closed access</span></span> <span class=\"c-subline me-3\">")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					if dataset.AccessLevelAfterEmbargo == "info:eu-repo/semantics/openAccess" {
						_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-download if--small if--success\"></i> ")
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
					} else {
						_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-ghent-university if--small if--primary\"></i> ")
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"c-subline text-truncate ps-2\">")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var7 string
					templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("dataset_access_levels." + dataset.AccessLevelAfterEmbargo))
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 71, Col: 84}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" from ")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var8 string
					templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(dataset.EmbargoDate)
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 73, Col: 34}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></span> ")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"c-subline me-3 ps-3 border-start\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if dataset.License != "" {
				if dataset.License == "LicenseNotListed" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("Licensed")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else {
					var templ_7745c5c3_Var9 string
					templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(dataset.License)
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 84, Col: 29}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
			} else {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<em>Add license</em>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div></div><div class=\"bc-toolbar-right justify-content-end\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if c.Repo.CanDeleteDataset(c.User, dataset) {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\"><div class=\"dropdown\"><button class=\"btn btn-outline-secondary btn-icon-only\" type=\"button\" data-bs-toggle=\"dropdown\" aria-haspopup=\"true\" aria-expanded=\"false\"><i class=\"if if-more\"></i></button><div class=\"dropdown-menu\"><button class=\"dropdown-item\" hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var10 string
				templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(views.URL(c.PathTo("dataset_confirm_delete", "id", dataset.ID)).QuerySet("redirect-url", redirectURL).String())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 108, Col: 132}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><i class=\"if if-delete\"></i> <span>Delete</span></button></div></div></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if c.Repo.CanCurate(c.User) && dataset.Locked {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"btn btn-outline-secondary\" hx-post=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var11 string
				templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(views.URL(c.PathTo("dataset_unlock", "id", dataset.ID)).QuerySet("redirect-url", c.CurrentURL.String()).String())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 122, Col: 133}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-swap=\"none\"><i class=\"if if-lock-unlock\"></i> <span class=\"btn-text d-md-none d-lg-inline\">Unlock record</span></button>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else if c.Repo.CanCurate(c.User) {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"btn btn-outline-secondary\" hx-post=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var12 string
				templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(views.URL(c.PathTo("dataset_lock", "id", dataset.ID)).QuerySet("redirect-url", c.CurrentURL.String()).String())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 131, Col: 131}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-swap=\"none\"><i class=\"if if-lock\"></i> <span class=\"btn-text d-md-none d-lg-inline\">Lock record</span></button>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if c.Repo.CanWithdrawDataset(c.User, dataset) {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\"><button class=\"btn btn-outline-danger\" hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var13 string
				templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(views.URL(c.PathTo("dataset_confirm_withdraw", "id", dataset.ID)).QuerySet("redirect-url", redirectURL).String())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 143, Col: 132}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><i class=\"if if-arrow-go-back\"></i> <span class=\"btn-text\">Withdraw</span></button></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if c.Repo.CanPublishDataset(c.User, dataset) && dataset.Status == "returned" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\"><button class=\"btn btn-success\" hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var14 string
				templ_7745c5c3_Var14, templ_7745c5c3_Err = templ.JoinStringErrs(views.URL(c.PathTo("dataset_confirm_republish", "id", dataset.ID)).QuerySet("redirect-url", redirectURL).String())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 155, Col: 133}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var14))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><div class=\"btn-text\">Republish <span class=\"d-md-none d-lg-inline\">to Biblio</span></div></button></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if c.Repo.CanPublishDataset(c.User, dataset) && dataset.Status != "returned" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-item\"><button class=\"btn btn-success\" hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var15 string
				templ_7745c5c3_Var15, templ_7745c5c3_Err = templ.JoinStringErrs(views.URL(c.PathTo("dataset_confirm_publish", "id", dataset.ID)).QuerySet("redirect-url", redirectURL).String())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 166, Col: 131}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var15))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><div class=\"btn-text\">Publish <span class=\"d-md-none d-lg-inline\">to Biblio</span></div></button></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div><h4 class=\"w-100 mb-0 mb-lg-2\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if len(dataset.Author) == 1 && dataset.Author[0] != nil {
				var templ_7745c5c3_Var16 string
				templ_7745c5c3_Var16, templ_7745c5c3_Err = templ.JoinStringErrs(dataset.Author[0].LastName())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 178, Col: 38}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var16))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(", ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var17 string
				templ_7745c5c3_Var17, templ_7745c5c3_Err = templ.JoinStringErrs(dataset.Author[0].FirstName())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 178, Col: 73}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var17))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(", creator. ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if len(dataset.Author) >= 2 {
				var templ_7745c5c3_Var18 string
				templ_7745c5c3_Var18, templ_7745c5c3_Err = templ.JoinStringErrs(dataset.Author[0].LastName())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 181, Col: 38}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var18))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(", ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var19 string
				templ_7745c5c3_Var19, templ_7745c5c3_Err = templ.JoinStringErrs(dataset.Author[0].FirstName())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 181, Col: 73}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var19))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(", et al., creators. ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if dataset.Title != "" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<em>\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var20 string
				templ_7745c5c3_Var20, templ_7745c5c3_Err = templ.JoinStringErrs(dataset.Title)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 184, Col: 28}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var20))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(".\"</em> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<em>\"Untitled record.\"</em> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if dataset.Publisher != "" {
				var templ_7745c5c3_Var21 string
				templ_7745c5c3_Var21, templ_7745c5c3_Err = templ.JoinStringErrs(dataset.Publisher)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 189, Col: 27}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var21))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(", ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if dataset.Year != "" {
				var templ_7745c5c3_Var22 string
				templ_7745c5c3_Var22, templ_7745c5c3_Err = templ.JoinStringErrs(dataset.Year)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 192, Col: 22}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var22))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(". ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if dataset.IdentifierType() != "" {
				var templ_7745c5c3_Var23 string
				templ_7745c5c3_Var23, templ_7745c5c3_Err = templ.JoinStringErrs(dataset.IdentifierType())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 195, Col: 34}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var23))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(".")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h4><div class=\"bc-toolbar bc-toolbar-lg-responsive flex-md-row\"><div class=\"bc-toolbar-left\"><div class=\"bc-toolbar-item d-block\"><div class=\"d-block pe-3 my-2\"><div class=\"input-group\"><button type=\"button\" class=\"btn btn-outline-secondary btn-sm\" data-clipboard=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var24 string
			templ_7745c5c3_Var24, templ_7745c5c3_Err = templ.JoinStringErrs(dataset.ID)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 203, Col: 101}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var24))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><i class=\"if if-copy text-muted \"></i> <span class=\"btn-text ms-0 me-1\">Biblio ID</span></button> <code class=\"c-code\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var25 string
			templ_7745c5c3_Var25, templ_7745c5c3_Err = templ.JoinStringErrs(dataset.ID)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 207, Col: 44}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var25))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</code></div></div></div></div><div class=\"bc-toolbar-right\"><div class=\"bc-toolbar-item\"><p class=\"c-body-small text-muted d-inline-block me-2 me-md-0\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var26 string
			templ_7745c5c3_Var26, templ_7745c5c3_Err = templ.JoinStringErrs(views.CreatedBy(c, dataset.DateCreated, dataset.Creator))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 215, Col: 68}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var26))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p><p class=\"c-body-small text-muted d-inline-block ms-md-2\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var27 string
			templ_7745c5c3_Var27, templ_7745c5c3_Err = templ.JoinStringErrs(views.UpdatedBy(c, dataset.DateUpdated, dataset.User, dataset.LastUser))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 218, Col: 83}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var27))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p></div></div></div></div></div><div class=\"bc-navbar bc-navbar--auto bc-navbar--white bc-navbar--bordered-bottom\"><div class=\"bc-toolbar bc-toolbar--auto\"><div class=\"bc-toolbar-left\"><div class=\"bc-toolbar-item\" id=\"show-nav\">")
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
			var templ_7745c5c3_Var28 string
			templ_7745c5c3_Var28, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("dataset_"+c.SubNav, "id", dataset.ID).String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `dataset/show.templ`, Line: 241, Col: 70}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var28))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-trigger=\"load delay:100ms\" hx-target=\"#show-content\"></div></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = views.PageLayout(c, views.PageLayoutArgs{
			Title: c.Loc.Get("dataset.page.show.title"),
			Breadcrumbs: []views.Breadcrumb{
				{LabelID: "datasets", URL: c.PathTo("datasets")},
				{LabelID: "dataset"},
			},
		}).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
