package publicationediting

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views"
	publicationviews "github.com/ugent-library/biblio-backoffice/views/publication"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
)

type BindSuggestDatasets struct {
	Query string `query:"q"`
}
type BindDataset struct {
	DatasetID string `form:"dataset_id"`
}
type BindDeleteDataset struct {
	DatasetID  string `path:"dataset_id"`
	SnapshotID string `path:"snapshot_id"`
}

func AddDataset(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	p := ctx.GetPublication(r)

	hits, err := searchRelatedDatasets(c, p, "")
	if err != nil {
		c.Log.Errorw("add publication dataset: could not execute search", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	publicationviews.AddDataset(c, p, hits).Render(r.Context(), w)
}

func SuggestDatasets(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindSuggestDatasets{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("suggest publication datasets: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	p := ctx.GetPublication(r)

	hits, err := searchRelatedDatasets(c, p, b.Query)
	if err != nil {
		c.Log.Errorw("add publication dataset: could not execute search", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	publicationviews.SuggestDatasets(c, p, hits).Render(r.Context(), w)
}

func CreateDataset(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	b := BindDataset{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("create publication dataset: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	// TODO reduce calls to repository
	d, err := c.Repo.GetDataset(b.DatasetID)
	if err != nil {
		c.Log.Errorw("create publication dataset: could not get dataset", "errors", err, "publication", publication.ID, "dataset", b.DatasetID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	// TODO handle validation errors
	// TODO pass If-Match
	// TODO handle conflict
	err = c.Repo.AddPublicationDataset(publication, d, c.User)
	if err != nil {
		c.Log.Errorw("create publication dataset: could not add dataset to publication", "errors", err, "publication", publication.ID, "dataset", b.DatasetID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	// Refresh the ctx.Publication: it still carries the old snapshotID
	publication, err = c.Repo.GetPublication(publication.ID)
	if err != nil {
		c.Log.Errorw("create publication dataset: could not get publication", "errors", err, "publication", publication.ID, "dataset", b.DatasetID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	relatedDatasets, err := c.Repo.GetVisiblePublicationDatasets(c.User, publication)
	if err != nil {
		c.Log.Errorw("create publication dataset: could not get related datasets", "errors", err, "publication", publication.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(publicationviews.DatasetsBodySelector, publicationviews.DatasetsBody(c, publication, relatedDatasets)).Render(r.Context(), w)
}

func ConfirmDeleteDataset(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	b := BindDeleteDataset{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete publication dataset: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	if b.SnapshotID != publication.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to remove this dataset from the publication?",
		DeleteUrl:  c.PathTo("publication_delete_dataset", "id", publication.ID, "dataset_id", b.DatasetID),
		SnapshotID: publication.SnapshotID,
	}).Render(r.Context(), w)
}

func DeleteDataset(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	b := BindDeleteDataset{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("delete publication dataset: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	if !publication.HasRelatedDataset(b.DatasetID) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	// TODO reduce calls to repository
	d, err := c.Repo.GetDataset(b.DatasetID)
	if err != nil {
		c.Log.Errorw("delete publication dataset: could not get dataset", "errors", err, "publication", publication.ID, "dataset", b.DatasetID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	// TODO handle validation errors
	// TODO pass If-Match
	// TODO handle conflict
	err = c.Repo.RemovePublicationDataset(publication, d, c.User)

	if err != nil {
		c.Log.Errorw("delete publication dataset: could not remove dataset", "errors", err, "publication", publication.ID, "dataset", b.DatasetID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	// Refresh the ctx.Publication: it still carries the old snapshotID
	publication, err = c.Repo.GetPublication(publication.ID)
	if err != nil {
		c.Log.Errorw("delete publication dataset: could not get publication", "errors", err, "publication", publication.ID, "dataset", b.DatasetID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	relatedDatasets, err := c.Repo.GetVisiblePublicationDatasets(c.User, publication)
	if err != nil {
		c.Log.Errorw("create publication dataset: could not get related datasets", "errors", err, "publication", publication.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(publicationviews.DatasetsBodySelector, publicationviews.DatasetsBody(c, publication, relatedDatasets)).Render(r.Context(), w)
}

func searchRelatedDatasets(c *ctx.Ctx, p *models.Publication, q string) (*models.DatasetHits, error) {
	args := models.NewSearchArgs().WithQuery(q)

	// add exclusion filter if necessary
	if len(p.RelatedDataset) > 0 {
		pubDatasetIDs := make([]string, len(p.RelatedDataset))
		for i, d := range p.RelatedDataset {
			pubDatasetIDs[i] = d.ID
		}
		args.Filters["!id"] = pubDatasetIDs
	}

	searchService := c.DatasetSearchIndex.WithScope("status", "public")

	return searchService.Search(args)
}
