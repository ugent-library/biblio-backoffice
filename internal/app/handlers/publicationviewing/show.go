package publicationviewing

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/displays"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type YieldShow struct {
	Context
	PageTitle    string
	ActiveNav    string
	ActiveSubNav string
	SearchArgs   *models.SearchArgs
}

type YieldShowDescription struct {
	Context
	ActiveSubNav   string
	SearchArgs     *models.SearchArgs
	DisplayDetails *display.Display
}

type YieldShowContributors struct {
	Context
	ActiveSubNav string
	SearchArgs   *models.SearchArgs
}

type YieldShowDatasets struct {
	Context
	ActiveSubNav    string
	SearchArgs      *models.SearchArgs
	RelatedDatasets []*models.Dataset
}

var allowedSubNavs = []string{
	"description",
	"files",
	"contributors",
	"datasets",
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request, ctx Context) {
	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	// TODO bind and validate
	activeSubNav := r.URL.Query().Get("show")
	if !validation.InArray(allowedSubNavs, activeSubNav) {
		activeSubNav = "description"
	}

	render.Wrap(w, "layouts/default", "publication/show_page", YieldShow{
		Context:      ctx,
		PageTitle:    ctx.T("publication.page.show.title"),
		ActiveNav:    "publications",
		ActiveSubNav: activeSubNav,
		SearchArgs:   searchArgs,
	})
}

func (h *Handler) ShowDescription(w http.ResponseWriter, r *http.Request, ctx Context) {
	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.Render(w, "publication/show_description", YieldShowDescription{
		Context:        ctx,
		ActiveSubNav:   "description",
		SearchArgs:     searchArgs,
		DisplayDetails: displayDetails(ctx),
	})
}

func (h *Handler) ShowFiles(w http.ResponseWriter, r *http.Request, ctx Context) {
	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.Render(w, "publication/show_files", YieldShowContributors{
		Context:      ctx,
		ActiveSubNav: "files",
		SearchArgs:   searchArgs,
	})
}

func (h *Handler) ShowContributors(w http.ResponseWriter, r *http.Request, ctx Context) {
	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.Render(w, "publication/show_contributors", YieldShowContributors{
		Context:      ctx,
		ActiveSubNav: "contributors",
		SearchArgs:   searchArgs,
	})
}

func (h *Handler) ShowDatasets(w http.ResponseWriter, r *http.Request, ctx Context) {
	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	relatedDatasets, err := h.Repository.GetPublicationDatasets(ctx.Publication)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "publication/show_datasets", YieldShowDatasets{
		Context:         ctx,
		ActiveSubNav:    "datasets",
		SearchArgs:      searchArgs,
		RelatedDatasets: relatedDatasets,
	})
}

func displayDetails(ctx Context) *display.Display {

	switch ctx.Publication.Type {
	case "book_chapter":
		return displays.DisplayBookChapter(ctx.Locale, ctx.Publication)
	case "book_editor":
		return displays.DisplayBookEditor(ctx.Locale, ctx.Publication)
	case "book":
		return displays.DisplayBook(ctx.Locale, ctx.Publication)
	case "conference":
		return displays.DisplayConference(ctx.Locale, ctx.Publication)
	case "dissertation":
		return displays.DisplayDissertation(ctx.Locale, ctx.Publication)
	case "issue_editor":
		return displays.DisplayIssueEditor(ctx.Locale, ctx.Publication)
	case "journal_article":
		return displays.DisplayJournalArticle(ctx.Locale, ctx.Publication)
	case "miscellaneous":
		return displays.DisplayMiscellaneous(ctx.Locale, ctx.Publication)
	default:
		return display.New()
	}

}
