package datasetviewing

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/displays"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

var (
	subNavs = []string{"description", "contributors", "publications"}
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
	SubNavs        []string
	ActiveSubNav   string
	DisplayDetails *display.Display
}

type YieldShowContributors struct {
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

type YieldShowPublications struct {
	Context
	SubNavs             []string
	ActiveSubNav        string
	RelatedPublications []*models.Publication
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request, ctx Context) {
	activeSubNav := r.URL.Query().Get("show")
	if !validation.InArray(subNavs, activeSubNav) {
		activeSubNav = "description"
	}

	render.Layout(w, "layouts/default", "dataset/pages/show", YieldShow{
		Context:      ctx,
		PageTitle:    ctx.T("dataset.page.show.title"),
		SubNavs:      subNavs,
		ActiveNav:    "datasets",
		ActiveSubNav: activeSubNav,
	})
}

func (h *Handler) ShowDescription(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.View(w, "dataset/show_description", YieldShowDescription{
		Context:        ctx,
		SubNavs:        subNavs,
		ActiveSubNav:   "description",
		DisplayDetails: displays.DatasetDetails(ctx.Locale, ctx.Dataset),
	})
}

func (h *Handler) ShowContributors(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.View(w, "dataset/show_contributors", YieldShowContributors{
		Context:      ctx,
		SubNavs:      subNavs,
		ActiveSubNav: "contributors",
	})
}

func (h *Handler) ShowPublications(w http.ResponseWriter, r *http.Request, ctx Context) {
	relatedPublications, err := h.Repository.GetDatasetPublications(ctx.Dataset)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/show_publications", YieldShowPublications{
		Context:             ctx,
		SubNavs:             subNavs,
		ActiveSubNav:        "publications",
		RelatedPublications: relatedPublications,
	})
}
