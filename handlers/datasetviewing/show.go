package datasetviewing

import (
	"net/http"

	"slices"

	"github.com/ugent-library/biblio-backoffice/ctx"
	datasetviews "github.com/ugent-library/biblio-backoffice/views/dataset"
)

var (
	subNavs = []string{"description", "contributors", "publications", "activity"}
)

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

func ShowContributors(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	redirectURL := r.URL.Query().Get("redirect-url")
	if redirectURL == "" {
		redirectURL = c.PathTo("datasets").String()
	}
	datasetviews.Contributors(ctx.Get(r), ctx.GetDataset(r), redirectURL).Render(r.Context(), w)
}

func ShowPublications(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	relatedPublications, err := c.Repo.GetDatasetPublications(dataset)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	datasetviews.Publications(c, dataset, relatedPublications).Render(r.Context(), w)
}

func ShowActivity(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	redirectURL := r.URL.Query().Get("redirect-url")
	if redirectURL == "" {
		redirectURL = c.PathTo("datasets").String()
	}
	datasetviews.Activity(ctx.Get(r), ctx.GetDataset(r), redirectURL).Render(r.Context(), w)
}
