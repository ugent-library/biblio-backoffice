package datasetviewing

import (
	"net/http"

	"slices"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views"
	datasetviews "github.com/ugent-library/biblio-backoffice/views/dataset"
)

var (
	subNavs = []string{"description", "contributors", "publications"}
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

func BiblioMessages(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	datasetviews.Messages(c, datasetviews.MessagesArgs{
		Dataset: d,
	}).Render(r.Context(), w)
}

func RecentActivity(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	var (
		snapshots []*models.Dataset
		acts      []views.Activity
	)

	// First take the (max) 21 most recent snapshots
	err := c.Repo.DatasetHistory(d.ID, func(snapshot *models.Dataset) bool {
		snapshots = append(snapshots, snapshot)
		return len(snapshots) <= 21
	})
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	// Convert the 20 most recent snapshots to activities
	for i := 0; i < len(snapshots); i++ {
		var prevSnapshot *models.Dataset
		if len(snapshots) > i+1 {
			prevSnapshot = snapshots[i+1]
		}
		acts = append(acts, handlers.GetDatasetActivity(c, snapshots[i], prevSnapshot))

		// Ignore the 21st snapshot, only used for comparison with the 20th
		if len(acts) >= 20 {
			break
		}
	}

	datasetviews.RecentActivity(c, acts, d).Render(r.Context(), w)
}
