package datasetediting

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views"
	datasetviews "github.com/ugent-library/biblio-backoffice/views/dataset"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
)

type BindSuggestPublications struct {
	Query string `query:"q"`
}
type BindPublication struct {
	PublicationID string `form:"publication_id"`
}
type BindDeletePublication struct {
	PublicationID string `path:"publication_id"`
	SnapshotID    string `path:"snapshot_id"`
}

func AddPublication(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	hits, err := searchRelatedPublications(c, dataset, "")
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.ShowModal(datasetviews.AddPublicationDialog(c, dataset, hits)).Render(r.Context(), w)
}

func SuggestPublications(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindSuggestPublications{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	hits, err := searchRelatedPublications(c, dataset, b.Query)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	datasetviews.SuggestPublications(c, dataset, hits).Render(r.Context(), w)
}

func CreatePublication(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindPublication{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	// TODO reduce calls to repository
	p, err := c.Repo.GetPublication(b.PublicationID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	// TODO handle validation errors
	// TODO pass If-Match
	// TODO handle conflict
	err = c.Repo.AddPublicationDataset(p, dataset, c.User)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	// Refresh the ctx.Dataset: it still carries the old snapshotID
	dataset, err = c.Repo.GetDataset(dataset.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	relatedPublications, err := c.Repo.GetVisibleDatasetPublications(c.User, dataset)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.CloseModalAndReplace(datasetviews.PublicationsBodySelector, datasetviews.PublicationsBody(c, dataset, relatedPublications)).Render(r.Context(), w)
}

func ConfirmDeletePublication(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindDeletePublication{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	if b.SnapshotID != dataset.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to remove this publication from the dataset?",
		DeleteUrl:  c.PathTo("dataset_delete_publication", "id", dataset.ID, "publication_id", b.PublicationID),
		SnapshotID: dataset.SnapshotID,
	}).Render(r.Context(), w)
}

func DeletePublication(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindDeletePublication{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	// TODO reduce calls to repository
	p, err := c.Repo.GetPublication(b.PublicationID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	// TODO handle validation errors
	// TODO pass If-Match
	// TODO handle conflict
	err = c.Repo.RemovePublicationDataset(p, dataset, c.User)

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	// Refresh the dataset since it still caries the old snapshotid
	dataset, err = c.Repo.GetDataset(dataset.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	relatedPublications, err := c.Repo.GetVisibleDatasetPublications(c.User, dataset)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.CloseModalAndReplace(datasetviews.PublicationsBodySelector, datasetviews.PublicationsBody(c, dataset, relatedPublications)).Render(r.Context(), w)
}

func searchRelatedPublications(c *ctx.Ctx, d *models.Dataset, q string) (*models.PublicationHits, error) {
	args := models.NewSearchArgs().WithQuery(q)

	// add exclusion filter if necessary
	if len(d.RelatedPublication) > 0 {
		datasetPubIDs := make([]string, len(d.RelatedPublication))
		for i, d := range d.RelatedPublication {
			datasetPubIDs[i] = d.ID
		}
		args.Filters["!id"] = datasetPubIDs
	}

	searchService := c.PublicationSearchIndex.WithScope("status", "public")

	return searchService.Search(args)
}
