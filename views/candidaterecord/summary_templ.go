// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package candidaterecordviews

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	contributorviews "github.com/ugent-library/biblio-backoffice/views/contributor"
	publicationsummaryviews "github.com/ugent-library/biblio-backoffice/views/publication/summary"
	"github.com/ugent-library/biblio-backoffice/views/util"
	"math"
	"time"
)

type SummaryOpts struct {
	Thumbnail   string
	Badge       templ.Component
	Info        templ.Component
	ShowDetails bool
}

func daysUntilDisappearanceDate(date time.Time) string {
	disappearanceDate := util.DateOnly(date).AddDate(0, 0, 90)
	today := util.DateOnly(time.Now())
	daysUntil := math.Max(0, disappearanceDate.Sub(today).Hours()/24)

	return fmt.Sprintf("%d", int(daysUntil))
}

func Summary(c *ctx.Ctx, rec *models.CandidateRecord) templ.Component {
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
		switch rec.Status {
		case "new":
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
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"btn btn-link btn-link-muted\" hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var3 string
				templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("confirm_reject_candidate_record", "id", rec.ID, "redirect-url", c.PathTo("candidate_records").String()).String())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 38, Col: 136}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><div class=\"btn-text\">Reject</div></button> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if c.UserRole == "curator" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"btn btn-link btn-link-muted\" hx-get=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var4 string
					templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("candidate_records_preview", "id", rec.ID, "redirect-url", c.PathTo("candidate_records").String()).String())
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 46, Col: 131}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><div class=\"btn-text\">Preview</div></button> <button class=\"btn btn-outline-primary\" hx-put=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var5 string
					templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("import_candidate_record", "id", rec.ID).String())
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 53, Col: 73}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-swap=\"none\"><div class=\"btn-text\">Import <span class=\"d-lg-none d-xl-inline-block\">&amp; complete</span></div></button>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"btn btn-link btn-link-muted\" hx-put=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var6 string
					templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("import_candidate_record", "id", rec.ID).String())
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 64, Col: 73}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-swap=\"none\"><div class=\"btn-text\">Import <span class=\"d-lg-none d-xl-inline-block\">&amp; complete</span></div></button> <button class=\"btn btn-outline-primary\" hx-get=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var7 string
					templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("candidate_records_preview", "id", rec.ID, "redirect-url", c.PathTo("candidate_records").String()).String())
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 74, Col: 131}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modals\"><div class=\"btn-text\">Preview</div></button>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				return templ_7745c5c3_Err
			})
			templ_7745c5c3_Err = publicationSummary(c, rec.Publication, SummaryOpts{
				Badge:       summaryBadge("badge-default", fmt.Sprintf("Biblio suggestion via %s", rec.SourceName)),
				ShowDetails: true,
				Thumbnail:   c.AssetPath("/images/plato-logo.svg")}).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case "imported":
			templ_7745c5c3_Err = publicationSummary(c, rec.Publication, SummaryOpts{
				Badge: summaryBadge("badge-warning-light", fmt.Sprintf("Imported %s suggestion", rec.SourceName)),
				Info:  importedInfo(c, rec)}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case "rejected":
			templ_7745c5c3_Err = publicationSummary(c, rec.Publication, SummaryOpts{
				Badge: summaryBadge("badge-danger-light", fmt.Sprintf("Rejected %s suggestion", rec.SourceName)),
				Info:  rejectedInfo(c, rec)}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		return templ_7745c5c3_Err
	})
}

func publicationSummary(c *ctx.Ctx, p *models.Publication, opts SummaryOpts) templ.Component {
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
		templ_7745c5c3_Var8 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var8 == nil {
			templ_7745c5c3_Var8 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"w-100\"><div class=\"c-thumbnail-and-text align-items-start w-100\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if opts.Thumbnail != "" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"c-thumbnail c-thumbnail-1-1 c-thumbnail-img c-thumbnail-small c-thumbnail-lg-large\"><div class=\"c-thumbnail-inner\"><img src=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var9 string
			templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(opts.Thumbnail)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 98, Col: 31}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"c-thumbnail-text\"><div class=\"hstack-lg-responsive align-items-start gap-3 w-100\"><div class=\"vstack gap-5\"><div class=\"vstack gap-2\"><div class=\"d-inline-flex align-items-center flex-wrap\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if opts.Badge != nil {
			templ_7745c5c3_Err = opts.Badge.Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"c-subline ps-2 me-3 pe-3 border-end\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if p.Classification != "" {
			var templ_7745c5c3_Var10 string
			templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%s: %s", c.Loc.Get("publication_types."+p.Type), p.Classification))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 112, Col: 93}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			var templ_7745c5c3_Var11 string
			templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_types." + p.Type))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 114, Col: 52}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if mainFile := p.MainFile(); mainFile != nil {
			var templ_7745c5c3_Var12 = []any{"c-subline", "me-3", "pe-3", templ.KV("border-end", mainFile.AccessLevel == "info:eu-repo/semantics/embargoedAccess")}
			templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var12...)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var13 string
			templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(templ.CSSClasses(templ_7745c5c3_Var12).String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 1, Col: 0}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if mainFile.AccessLevel == "info:eu-repo/semantics/openAccess" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-download if--small if--success\"></i> <span class=\"c-subline text-truncate\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var14 string
				templ_7745c5c3_Var14, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_file_access_levels." + mainFile.AccessLevel))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 121, Col: 118}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var14))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else if mainFile.AccessLevel == "info:eu-repo/semantics/embargoedAccess" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-time if--small if--warning\"></i> <span class=\"c-subline text-muted\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var15 string
				templ_7745c5c3_Var15, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_file_access_levels." + mainFile.AccessLevel))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 124, Col: 115}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var15))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else if mainFile.AccessLevel == "info:eu-repo/semantics/restrictedAccess" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-ghent-university if--small if--primary\"></i> <span class=\"c-subline text-muted\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var16 string
				templ_7745c5c3_Var16, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_file_access_levels." + mainFile.AccessLevel))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 127, Col: 115}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var16))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else if mainFile.AccessLevel == "info:eu-repo/semantics/closedAccess" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-forbid if--small if--danger\"></i> <span class=\"c-subline text-muted\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var17 string
				templ_7745c5c3_Var17, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_file_access_levels." + mainFile.AccessLevel))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 130, Col: 115}
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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span> ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if mainFile.AccessLevel == "info:eu-repo/semantics/embargoedAccess" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"c-subline me-3 pe-3 border-end\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if mainFile.AccessLevelDuringEmbargo == "info:eu-repo/semantics/closedAccess" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-eye-off if--small if--muted\"></i> ")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<i class=\"if if-ghent-university if--small if--primary\"></i> ")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"c-subline text-truncate\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var18 string
				templ_7745c5c3_Var18, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_file_access_levels_during_embargo." + mainFile.AccessLevelDuringEmbargo))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 140, Col: 146}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var18))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></span> <span class=\"c-subline me-3 pe-3\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if mainFile.AccessLevelAfterEmbargo == "info:eu-repo/semantics/openAccess" {
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
				var templ_7745c5c3_Var19 string
				templ_7745c5c3_Var19, templ_7745c5c3_Err = templ.JoinStringErrs(c.Loc.Get("publication_file_access_levels_after_embargo." + mainFile.AccessLevelAfterEmbargo))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 148, Col: 106}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var19))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" from ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var20 string
				templ_7745c5c3_Var20, templ_7745c5c3_Err = templ.JoinStringErrs(mainFile.EmbargoDate)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 148, Col: 136}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var20))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><h4 class=\"mb-0\"><span class=\"list-group-item-title\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if p.Title != "" {
			var templ_7745c5c3_Var21 string
			templ_7745c5c3_Var21, templ_7745c5c3_Err = templ.JoinStringErrs(p.Title)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 156, Col: 19}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var21))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("Untitled record")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></h4>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if opts.ShowDetails {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<ul class=\"c-meta-list c-meta-list-inline\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, summaryPart := range p.SummaryParts() {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<li class=\"c-meta-item\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var22 string
				templ_7745c5c3_Var22, templ_7745c5c3_Err = templ.JoinStringErrs(summaryPart)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 165, Col: 47}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var22))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</ul><div class=\"d-inline-flex align-items-center flex-wrap\"><div class=\"d-inline-flex align-items-center flex-wrap\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = contributorviews.Summary(c, contributorviews.SummaryArgs{
				Role:                    "author",
				Contributors:            p.Author,
				CanViewMoreContributors: c.Repo.CanViewPublication(c.User, p),
				CanEditContributors:     c.Repo.CanEditPublication(c.User, p),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"text-muted fst-italic me-5 my-1\">supervised by</div><div class=\"d-inline-flex align-items-center flex-wrap\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = contributorviews.Summary(c, contributorviews.SummaryArgs{
				Role:                    "supervisor",
				Contributors:            p.Supervisor,
				CurrentUserRoles:        publicationsummaryviews.GetUserContributorRoles(p, c.User),
				CanViewMoreContributors: c.Repo.CanViewPublication(c.User, p),
				CanEditContributors:     c.Repo.CanEditPublication(c.User, p),
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if len(p.RelatedOrganizations) > 0 {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"d-inline-flex align-items-center flex-wrap\"><span class=\"badge rounded-pill badge-light me-4\">Suggested departments</span><ul class=\"c-meta-list c-meta-list-inline\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				for _, o := range p.RelatedOrganizations {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<li class=\"c-meta-item\"><i class=\"if if-building if--small if--muted pe-1\"></i> <span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var23 string
					templ_7745c5c3_Var23, templ_7745c5c3_Err = templ.JoinStringErrs(o.OrganizationID)
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 195, Col: 37}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var23))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></li>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</ul></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div><div class=\"c-button-toolbar flex-row-reverse flex-lg-row\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var8.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if opts.Info != nil {
			templ_7745c5c3_Err = opts.Info.Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func summaryBadge(colorClass string, text string) templ.Component {
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
		templ_7745c5c3_Var24 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var24 == nil {
			templ_7745c5c3_Var24 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var templ_7745c5c3_Var25 = []any{"badge", "rounded-pill", "me-3", "my-2", colorClass}
		templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var25...)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var26 string
		templ_7745c5c3_Var26, templ_7745c5c3_Err = templ.JoinStringErrs(templ.CSSClasses(templ_7745c5c3_Var25).String())
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 1, Col: 0}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var26))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><span class=\"badge-circle\"></span> <span class=\"badge-text\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var27 string
		templ_7745c5c3_Var27, templ_7745c5c3_Err = templ.JoinStringErrs(text)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 219, Col: 33}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var27))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></span>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func importedInfo(c *ctx.Ctx, rec *models.CandidateRecord) templ.Component {
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
		templ_7745c5c3_Var28 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var28 == nil {
			templ_7745c5c3_Var28 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<p class=\"text-muted text-md-end text-nowrap\"><span>Claimed: ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if rec.StatusPerson != nil {
			var templ_7745c5c3_Var29 string
			templ_7745c5c3_Var29, templ_7745c5c3_Err = templ.JoinStringErrs(rec.StatusPerson.FullName)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 228, Col: 31}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var29))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(". ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("[missing]. ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if rec.ImportedID != "" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<a class=\"c-link c-link-muted\" type=\"button\" href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var30 templ.SafeURL = templ.URL(c.PathTo("publication", "id", rec.ImportedID).String())
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var30)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">View</a>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span><br>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if rec.StatusDate != nil {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<small class=\"fst-italic d-inline-block pt-2\">On ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var31 string
			templ_7745c5c3_Var31, templ_7745c5c3_Err = templ.JoinStringErrs(rec.StatusDate.Format("2006-01-02"))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 238, Col: 89}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var31))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(". Reminder disappears in ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var32 string
			templ_7745c5c3_Var32, templ_7745c5c3_Err = templ.JoinStringErrs(daysUntilDisappearanceDate(*rec.StatusDate))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 238, Col: 161}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var32))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" day(s).</small>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func rejectedInfo(c *ctx.Ctx, rec *models.CandidateRecord) templ.Component {
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
		templ_7745c5c3_Var33 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var33 == nil {
			templ_7745c5c3_Var33 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<p class=\"text-muted text-md-end text-nowrap\"><span>Rejected: ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if rec.StatusPerson != nil {
			var templ_7745c5c3_Var34 string
			templ_7745c5c3_Var34, templ_7745c5c3_Err = templ.JoinStringErrs(rec.StatusPerson.FullName)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 248, Col: 31}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var34))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(". ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("[missing]. ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<a class=\"c-link c-link-muted\" type=\"button\" hx-put=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var35 string
		templ_7745c5c3_Var35, templ_7745c5c3_Err = templ.JoinStringErrs(c.PathTo("restore_rejected_candidate_record", "id", rec.ID).String())
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 255, Col: 81}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var35))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var36 string
		templ_7745c5c3_Var36, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("#candidate-record-%s", rec.ID))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 256, Col: 59}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var36))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-swap=\"outerHTML\">Restore</a></span><br><small class=\"fst-italic d-inline-block pt-2\">On ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var37 string
		templ_7745c5c3_Var37, templ_7745c5c3_Err = templ.JoinStringErrs(rec.StatusDate.Format("2006-01-02"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 261, Col: 88}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var37))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(". Reminder disappears in ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var38 string
		templ_7745c5c3_Var38, templ_7745c5c3_Err = templ.JoinStringErrs(daysUntilDisappearanceDate(*rec.StatusDate))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `candidaterecord/summary.templ`, Line: 261, Col: 160}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var38))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" day(s).</small></p>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
