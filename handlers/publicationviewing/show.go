package publicationviewing

import (
	"net/http"

	"slices"

	"github.com/ugent-library/biblio-backoffice/displays"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/display"
)

var subNavs = []string{"description", "files", "contributors", "datasets", "activity"}

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
	MaxFileSize  int
}

type YieldShowDatasets struct {
	Context
	SubNavs         []string
	ActiveSubNav    string
	RelatedDatasets []*models.Dataset
}

type YieldShowActivity struct {
	Context
	SubNavs      []string
	ActiveSubNav string
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request, ctx Context) {
	activeSubNav := r.URL.Query().Get("show")
	if !slices.Contains(subNavs, activeSubNav) {
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
		DisplayDetails:        displays.PublicationDetails(ctx.User, ctx.Locale, ctx.Publication),
		DisplayConference:     displays.PublicationConference(ctx.User, ctx.Locale, ctx.Publication),
		DisplayAdditionalInfo: displays.PublicationAdditionalInfo(ctx.User, ctx.Locale, ctx.Publication),
	})
}

func (h *Handler) ShowFiles(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.View(w, "publication/show_files", YieldShowFiles{
		Context:      ctx,
		SubNavs:      subNavs,
		ActiveSubNav: "files",
		MaxFileSize:  h.MaxFileSize,
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
	relatedDatasets, err := h.Repo.GetVisiblePublicationDatasets(ctx.User, ctx.Publication)
	if err != nil {
		h.Logger.Warn("show publication datasets: could not get publication datasets:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
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

func (h *Handler) ShowActivity(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.View(w, "publication/show_activity", YieldShowActivity{
		Context:      ctx,
		SubNavs:      subNavs,
		ActiveSubNav: "activity",
	})
}
