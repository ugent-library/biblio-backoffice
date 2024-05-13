package datasetviewing

import (
	"net/http"

	"slices"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/render"
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

type YieldShowContributors struct {
	Context
	SubNavs      []string
	ActiveSubNav string
}

type YieldShowActivity struct {
	Context
	SubNavs      []string
	ActiveSubNav string
}

func Show(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	activeSubNav := r.URL.Query().Get("show")
	if !slices.Contains(subNavs, activeSubNav) {
		activeSubNav = "description"
	}
	c.SubNav = activeSubNav

	datasetviews.Show(c, ctx.GetDataset(r), r.URL.Query().Get("redirect-url")).Render(r.Context(), w)
}

func ShowDescription(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	redirectURL := r.URL.Query().Get("redirect-url")
	if redirectURL == "" {
		redirectURL = c.PathTo("datasets").String()
	}

	datasetviews.Description(c, ctx.GetDataset(r), redirectURL).Render(r.Context(), w)
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

func (h *Handler) ShowActivity(w http.ResponseWriter, r *http.Request, legacyCtx Context) {
	render.View(w, "dataset/show_activity", YieldShowActivity{
		Context:      legacyCtx,
		SubNavs:      subNavs,
		ActiveSubNav: "activity",
	})
}
