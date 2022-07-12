package publicationviewing

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/displays"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

var (
	subNavs = []string{"description", "files", "contributors", "datasets"}
)

type YieldShow struct {
	Context
	PageTitle    string
	SubNavs      []string
	ActiveNav    string
	ActiveSubNav string
}

type YieldShowDescription struct {
	Context
	SubNavs               []string
	ActiveSubNav          string
	DisplayDetails        *display.Display
	DisplayConference     *display.Display
	DisplayAdditionalInfo *display.Display
}

type YieldShowContributors struct {
	Context
	SubNavs      []string
	ActiveSubNav string
}

type YieldShowFiles struct {
	Context
	SubNavs      []string
	ActiveSubNav string
}

type YieldShowContributorsRole struct {
	YieldShowContributors
	Role string
}

func (y YieldShowContributors) YieldRole(role string) YieldShowContributorsRole {
	return YieldShowContributorsRole{y, role}
}

type YieldShowDatasets struct {
	Context
	SubNavs         []string
	ActiveSubNav    string
	RelatedDatasets []*models.Dataset
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request, ctx Context) {
	activeSubNav := r.URL.Query().Get("show")
	if !validation.InArray(subNavs, activeSubNav) {
		activeSubNav = "description"
	}

	render.Layout(w, "layouts/default", "publication/pages/show", YieldShow{
		Context:      ctx,
		PageTitle:    ctx.Locale.T("publication.page.show.title"),
		SubNavs:      subNavs,
		ActiveNav:    "publications",
		ActiveSubNav: activeSubNav,
	})
}

func (h *Handler) ShowDescription(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.View(w, "publication/show_description", YieldShowDescription{
		Context:               ctx,
		SubNavs:               subNavs,
		ActiveSubNav:          "description",
		DisplayDetails:        displays.PublicationDetails(ctx.Locale, ctx.Publication),
		DisplayConference:     displays.PublicationConference(ctx.Locale, ctx.Publication.Conference),
		DisplayAdditionalInfo: displays.PublicationAdditionalInfo(ctx.Locale, ctx.Publication),
	})
}

func (h *Handler) ShowFiles(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.View(w, "publication/show_files", YieldShowFiles{
		Context:      ctx,
		SubNavs:      subNavs,
		ActiveSubNav: "files",
	})
}

func (h *Handler) ShowContributors(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.View(w, "publication/show_contributors", YieldShowContributors{
		Context:      ctx,
		SubNavs:      subNavs,
		ActiveSubNav: "contributors",
	})
}

func (h *Handler) ShowDatasets(w http.ResponseWriter, r *http.Request, ctx Context) {
	relatedDatasets, err := h.Repository.GetPublicationDatasets(ctx.Publication)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/show_datasets", YieldShowDatasets{
		Context:         ctx,
		SubNavs:         subNavs,
		ActiveSubNav:    "datasets",
		RelatedDatasets: relatedDatasets,
	})
}
