package datasetviewing

import (
	"net/http"

	"slices"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/displays"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/display"
	datasetviews "github.com/ugent-library/biblio-backoffice/views/dataset"
	"github.com/ugent-library/httperror"
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

func (h *Handler) Show(w http.ResponseWriter, r *http.Request, legacyCtx Context) {
	activeSubNav := r.URL.Query().Get("show")
	if !slices.Contains(subNavs, activeSubNav) {
		activeSubNav = "description"
	}

	render.Layout(w, "layouts/default", "dataset/pages/show", YieldShow{
		Context:      legacyCtx,
		PageTitle:    legacyCtx.Loc.Get("dataset.page.show.title"),
		SubNavs:      subNavs,
		ActiveNav:    "datasets",
		ActiveSubNav: activeSubNav,
	})
}

func (h *Handler) ShowDescription(w http.ResponseWriter, r *http.Request, legacyCtx Context) {
	render.View(w, "dataset/show_description", YieldShowDescription{
		Context:        legacyCtx,
		SubNavs:        subNavs,
		ActiveSubNav:   "description",
		DisplayDetails: displays.DatasetDetails(legacyCtx.User, legacyCtx.Loc, legacyCtx.Dataset),
	})
}

func (h *Handler) ShowContributors(w http.ResponseWriter, r *http.Request, legacyCtx Context) {
	render.View(w, "dataset/show_contributors", YieldShowContributors{
		Context:      legacyCtx,
		SubNavs:      subNavs,
		ActiveSubNav: "contributors",
	})
}

func ShowPublications(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	relatedPublications, err := c.Repo.GetVisibleDatasetPublications(c.User, dataset)
	if err != nil {
		c.Log.Errorw("show dataset publications: could not get publications", "errors", err, "dataset", dataset.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	datasetviews.ShowPublications(c, dataset, relatedPublications).Render(r.Context(), w)
}

func ShowActivity(w http.ResponseWriter, r *http.Request) {
	datasetviews.ShowActivity(ctx.Get(r), ctx.GetDataset(r)).Render(r.Context(), w)
}
