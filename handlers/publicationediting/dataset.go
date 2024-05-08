package publicationediting

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
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

type YieldAddDataset struct {
	Context
	Hits *models.DatasetHits
}
type YieldDatasets struct {
	Context
	RelatedDatasets []*models.Dataset
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

func (h *Handler) CreateDataset(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDataset{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("create publication dataset: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	// TODO reduce calls to repository
	d, err := h.Repo.GetDataset(b.DatasetID)
	if err != nil {
		h.Logger.Errorw("create publication dataset: could not get dataset", "errors", err, "publication", ctx.Publication.ID, "dataset", b.DatasetID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// TODO handle validation errors
	// TODO pass If-Match
	// TODO handle conflict
	err = h.Repo.AddPublicationDataset(ctx.Publication, d, ctx.User)
	if err != nil {
		h.Logger.Errorw("create publication dataset: could not add dataset to publication", "errors", err, "publication", ctx.Publication.ID, "dataset", b.DatasetID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// Refresh the ctx.Publication: it still carries the old snapshotID
	ctx.Publication, err = h.Repo.GetPublication(ctx.Publication.ID)
	if err != nil {
		h.Logger.Errorw("create publication dataset: could not get publication", "errors", err, "publication", ctx.Publication.ID, "dataset", b.DatasetID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	relatedDatasets, err := h.Repo.GetVisiblePublicationDatasets(ctx.User, ctx.Publication)
	if err != nil {
		h.Logger.Errorw("create publication dataset: could not get related datasets", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_datasets", YieldDatasets{
		Context:         ctx,
		RelatedDatasets: relatedDatasets,
	})
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

func (h *Handler) DeleteDataset(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeleteDataset{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete publication dataset: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	if !ctx.Publication.HasRelatedDataset(b.DatasetID) {
		views.ReplaceModal(views.ErrorDialog(ctx.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	// TODO reduce calls to repository
	d, err := h.Repo.GetDataset(b.DatasetID)
	if err != nil {
		h.Logger.Errorw("delete publication dataset: could not get dataset", "errors", err, "publication", ctx.Publication.ID, "dataset", b.DatasetID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// TODO handle validation errors
	// TODO pass If-Match
	// TODO handle conflict
	err = h.Repo.RemovePublicationDataset(ctx.Publication, d, ctx.User)

	if err != nil {
		h.Logger.Errorw("delete publication dataset: could not remove dataset", "errors", err, "publication", ctx.Publication.ID, "dataset", b.DatasetID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// Refresh the ctx.Publication: it still carries the old snapshotID
	ctx.Publication, err = h.Repo.GetPublication(ctx.Publication.ID)
	if err != nil {
		h.Logger.Errorw("delete publication dataset: could not get publication", "errors", err, "publication", ctx.Publication.ID, "dataset", b.DatasetID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	relatedDatasets, err := h.Repo.GetVisiblePublicationDatasets(ctx.User, ctx.Publication)
	if err != nil {
		h.Logger.Errorw("create publication dataset: could not get related datasets", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_datasets", YieldDatasets{
		Context:         ctx,
		RelatedDatasets: relatedDatasets,
	})
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
