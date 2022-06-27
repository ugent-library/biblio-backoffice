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
	ActiveSubNav          string
	SearchArgs            *models.SearchArgs
	DisplayDetails        *display.Display
	DisplayConference     *display.Display
	DisplayAdditionalInfo *display.Display
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
		Context:               ctx,
		ActiveSubNav:          "description",
		SearchArgs:            searchArgs,
		DisplayDetails:        displayDetails(ctx),
		DisplayConference:     displays.DisplayConference(ctx.Locale, ctx.Publication),
		DisplayAdditionalInfo: displays.DisplayAdditionalInfo(ctx.Locale, ctx.Publication),
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
		return displays.DisplayTypeBookChapter(ctx.Locale, ctx.Publication)
	case "book_editor":
		return displays.DisplayTypeBookEditor(ctx.Locale, ctx.Publication)
	case "book":
		return displays.DisplayTypeBook(ctx.Locale, ctx.Publication)
	case "conference":
		return displays.DisplayTypeConference(ctx.Locale, ctx.Publication)
	case "dissertation":
		return displays.DisplayTypeDissertation(ctx.Locale, ctx.Publication)
	case "issue_editor":
		return displays.DisplayTypeIssueEditor(ctx.Locale, ctx.Publication)
	case "journal_article":
		return displays.DisplayTypeJournalArticle(ctx.Locale, ctx.Publication)
	case "miscellaneous":
		return displays.DisplayTypeMiscellaneous(ctx.Locale, ctx.Publication)
	default:
		return display.New()
	}

}
