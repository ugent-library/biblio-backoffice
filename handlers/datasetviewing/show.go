package datasetviewing

import (
	"net/http"

	"slices"

	"github.com/ugent-library/biblio-backoffice/displays"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/display"
)

var (
	subNavs = []string{"description", "contributors", "publications", "activity"}
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

type YieldShowPublications struct {
	Context
	SubNavs             []string
	ActiveSubNav        string
	RelatedPublications []*models.Publication
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

	render.Layout(w, "layouts/default", "dataset/pages/show", YieldShow{
		Context:      ctx,
		PageTitle:    ctx.Locale.T("dataset.page.show.title"),
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
		DisplayDetails: displays.DatasetDetails(ctx.User, ctx.Locale, ctx.Dataset),
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
	relatedPublications, err := h.Repo.GetVisibleDatasetPublications(ctx.User, ctx.Dataset)
	if err != nil {
		h.Logger.Errorw("show dataset publications: could not get publications", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
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

func (h *Handler) ShowActivity(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.View(w, "dataset/show_activity", YieldShowActivity{
		Context:      ctx,
		SubNavs:      subNavs,
		ActiveSubNav: "activity",
	})
}
