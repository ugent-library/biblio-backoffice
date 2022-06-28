package publicationviewing

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/displays"
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
}

type YieldShowDescription struct {
	Context
	ActiveSubNav          string
	DisplayDetails        *display.Display
	DisplayConference     *display.Display
	DisplayAdditionalInfo *display.Display
}

type YieldShowContributors struct {
	Context
	ActiveSubNav string
}

type YieldShowDatasets struct {
	Context
	ActiveSubNav    string
	RelatedDatasets []*models.Dataset
}

var allowedSubNavs = []string{
	"description",
	"files",
	"contributors",
	"datasets",
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request, ctx Context) {
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
	})
}

func (h *Handler) ShowDescription(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Render(w, "publication/show_description", YieldShowDescription{
		Context:               ctx,
		ActiveSubNav:          "description",
		DisplayDetails:        displays.PublicationDetails(ctx.Locale, ctx.Publication),
		DisplayConference:     displays.PublicationConference(ctx.Locale, ctx.Publication.Conference),
		DisplayAdditionalInfo: displays.PublicationAdditionalInfo(ctx.Locale, ctx.Publication),
	})
}

func (h *Handler) ShowFiles(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Render(w, "publication/show_files", YieldShowContributors{
		Context:      ctx,
		ActiveSubNav: "files",
	})
}

func (h *Handler) ShowContributors(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Render(w, "publication/show_contributors", YieldShowContributors{
		Context:      ctx,
		ActiveSubNav: "contributors",
	})
}

func (h *Handler) ShowDatasets(w http.ResponseWriter, r *http.Request, ctx Context) {
	relatedDatasets, err := h.Repository.GetPublicationDatasets(ctx.Publication)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "publication/show_datasets", YieldShowDatasets{
		Context:         ctx,
		ActiveSubNav:    "datasets",
		RelatedDatasets: relatedDatasets,
	})
}
