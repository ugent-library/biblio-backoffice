// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package publication

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"fmt"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/friendly"
)

func Files(c *ctx.Ctx, p *models.Publication, redirectURL string) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div hx-swap-oob=\"innerHTML:#show-nav\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = SubNav(c, p, redirectURL).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div hx-swap-oob=\"innerHTML:#show-sidebar\"></div><div class=\"alert alert-success mb-6\"><i class=\"if if-check-circle\"></i><div><h3>Full texts are stored and made available in line with UGent’s <a href=\"https://www.ugent.be/intranet/en/research/impact/schol-publishing/policy-ugent.htm\" target=\"_blank\">scholarly publishing policy</a>.</h3><p>Other documents are handled according to the access levels and licences you indicate.</p></div></div><div class=\"card mb-6\"><div class=\"card-header\"><div class=\"bc-toolbar\"><div class=\"bc-toolbar-left\"><div class=\"bc-toolbar-title\">Full text documents</div></div></div></div><div id=\"files-body\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = FilesBody(c, p).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func FilesBody(c *ctx.Ctx, p *models.Publication) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"card-body p-0\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if c.User.CanEditPublication(p) {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form class=\"p-6\" hx-post=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(c.PathTo("publication_upload_file", "id", p.ID).String()))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-encoding=\"multipart/form-data\" hx-headers=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf(`{"If-Match": "%s"}`, p.SnapshotID)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\" hx-trigger=\"change\"><div class=\"c-file-upload file-upload-start\"><input class=\"upload-progress\" type=\"file\" name=\"file\" data-max-size=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprint(c.MaxFileSize)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" data-max-size-error=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("Maximum file size is %s", friendly.Bytes(int64(c.MaxFileSize)))))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><div class=\"c-file-upload__content\"><p>Drag and drop or</p><button class=\"btn btn-outline-primary\">Upload file</button><p class=\"small pt-3 mb-0\">Maximum file size: ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(friendly.Bytes(int64(c.MaxFileSize)))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/files.templ`, Line: 57, Col: 90}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p></div></div><div class=\"c-file-upload c-file-upload--disabled file-upload-busy d-none\"><div class=\"c-file-upload__content\"><p class=\"mt-5\">Uploading your file.<br><span>Hold on, do not refresh the page.</span></p><div class=\"progress w-75\"><div class=\"progress-bar progress-bar-striped progress-bar-animated\" role=\"progressbar\" style=\"width: 0%\" aria-valuenow=\"0\" aria-valuemin=\"0\" aria-valuemax=\"100\"></div></div><p class=\"mt-4 text-muted\"><span class=\"progress-bar-percent\">0</span>%</p></div></div><small class=\"form-text text-muted my-3\"><a href=\"https://onderzoektips.ugent.be/en/tips/00002066\" target=\"_blank\">Which document format or version should I use?</a></small></form><hr>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if len(p.File) > 0 {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<ul class=\"list-group list-group-flush\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, f := range p.File {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<li class=\"list-group-item\"><div class=\"list-group-item-inner\"><div class=\"list-group-item-main u-min-w-0\"><div class=\"c-thumbnail-and-text align-items-start d-block d-lg-flex\"><a href=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var4 templ.SafeURL = templ.URL(c.PathTo("publication_download_file", "id", p.ID, "file_id", f.ID).String())
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var4)))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><div class=\"c-thumbnail c-thumbnail-5-4 c-thumbnail-small c-thumbnail-xl-large mb-6 mb-xl-0 flex-shrink-0 d-none d-lg-block\"><div class=\"c-thumbnail-inner\"><i class=\"if if-article\"></i></div></div></a><div class=\"c-thumbnail-text u-min-w-0\"><div class=\"bc-toolbar bc-toolbar--auto\"><div class=\"bc-toolbar-left flex-wrap\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if f.AccessLevel == "info:eu-repo/semantics/openAccess" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"c-subline text-nowrap me-3 pe-3 my-2 border-end\"><i class=\"if if-download if--small if--muted\"></i> <span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var5 string
					templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_file_access_levels." + f.AccessLevel))
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/files.templ`, Line: 104, Col: 82}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else if f.AccessLevel == "info:eu-repo/semantics/restrictedAccess" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"c-subline text-nowrap me-3 pe-3 my-2 border-end\"><i class=\"if if-ghent-university if--small if--muted\"></i> <span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var6 string
					templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_file_access_levels." + f.AccessLevel))
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/files.templ`, Line: 109, Col: 82}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else if f.AccessLevel == "info:eu-repo/semantics/closedAccess" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"c-subline text-nowrap me-3 pe-3 my-2 border-end\"><i class=\"if if-eye-off if--small if--muted\"></i> <span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var7 string
					templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_file_access_levels." + f.AccessLevel))
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/files.templ`, Line: 114, Col: 82}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else if f.AccessLevel == "info:eu-repo/semantics/embargoedAccess" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"c-subline text-nowrap me-3 pe-3 my-2 border-end\"><i class=\"if if-time if--small\"></i> <span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var8 string
					templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_file_access_levels." + f.AccessLevel))
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/files.templ`, Line: 119, Col: 82}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div><div class=\"c-subline text-nowrap me-3 pe-3 my-2 border-end\">")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					if f.AccessLevelDuringEmbargo == "info:eu-repo/semantics/closedAccess" {
						_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-eye-off if--small if--primary\"></i>")
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
					} else {
						_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-ghent-university if--small if--primary\"></i>")
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var9 string
					templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_file_access_levels_during_embargo." + f.AccessLevelDuringEmbargo))
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/files.templ`, Line: 127, Col: 110}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div><div class=\"c-subline text-nowrap me-3 pe-3 my-2 border-end\">")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					if f.AccessLevelAfterEmbargo == "info:eu-repo/semantics/openAccess" {
						_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-download if--small if--muted\"></i>")
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
					} else {
						_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-ghent-university if--small if--muted\"></i>")
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var10 string
					templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_file_access_levels_after_embargo." + f.AccessLevelAfterEmbargo))
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/files.templ`, Line: 135, Col: 108}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" from ")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var11 string
					templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(f.EmbargoDate)
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/files.templ`, Line: 135, Col: 131}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"c-subline text-nowrap me-3 pe-3 my-2\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if f.License != "" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var12 string
					templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_licenses." + f.License))
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/files.templ`, Line: 140, Col: 68}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var13 string
					templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(f.OtherLicense)
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/files.templ`, Line: 142, Col: 36}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div><div class=\"bc-toolbar-right\"><div class=\"bc-toolbar-item ms-auto ms-lg-0\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if c.User.CanEditPublication(p) {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"c-button-toolbar\"><button class=\"btn btn-icon-only\" type=\"button\" hx-get=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(c.PathTo("publication_edit_file", "id", p.ID, "file_id", f.ID).String()))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-headers=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf(`{"If-Match": "%s"}`, p.SnapshotID)))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-swap=\"innerHTML\" hx-target=\"#modals\"><i class=\"if if-edit\"></i></button> <button class=\"btn btn-icon-only\" type=\"button\" hx-get=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(c.PathTo("publication_confirm_delete_file", "id", p.ID, "snapshot_id", p.SnapshotID, "file_id", f.ID).String()))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\" hx-trigger=\"click\"><i class=\"if if-delete\"></i></button></div>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div></div><h4><a href=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var14 templ.SafeURL = templ.URL(c.PathTo("publication_download_file", "id", p.ID, "file_id", f.ID).String())
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var14)))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><span class=\"list-group-item-title\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var15 string
				templ_7745c5c3_Var15, templ_7745c5c3_Err = templ.JoinStringErrs(f.Name)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/files.templ`, Line: 177, Col: 21}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var15))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></a></h4>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if f.Relation != "" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var16 string
					templ_7745c5c3_Var16, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_file_relations." + f.Relation))
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/files.templ`, Line: 182, Col: 72}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var16))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				if f.PublicationVersion != "" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"ms-3 ps-3 border-start\">")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var17 string
					templ_7745c5c3_Var17, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_versions." + f.PublicationVersion))
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/files.templ`, Line: 185, Col: 107}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var17))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bc-toolbar-left flex-wrap\"><div class=\"bc-toolbar-item\"><span class=\"c-body-small text-muted text-truncate my-2\">Uploaded ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var18 string
				templ_7745c5c3_Var18, templ_7745c5c3_Err = templ.JoinStringErrs(f.DateCreated.In(c.Timezone).Format("2006-01-02 at 15:04"))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `publication/files.templ`, Line: 189, Col: 138}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var18))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div></div></div></div></div></div></li>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</ul>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"m-6\"><span class=\"text-muted\">No files</span></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
