package publicationediting

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/bind"
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
type YieldDeleteDataset struct {
	Context
	DatasetID string
}

func (h *Handler) AddDataset(w http.ResponseWriter, r *http.Request, ctx Context) {
	hits, err := h.searchRelatedDatasets(ctx.User, ctx.Publication, "")
	if err != nil {
		h.Logger.Errorw("add dataset publication: could not execute search", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "show_modal", "publication/add_dataset", YieldAddDataset{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) SuggestDatasets(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindSuggestDatasets{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("suggest publication datasets: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	hits, err := h.searchRelatedDatasets(ctx.User, ctx.Publication, b.Query)
	if err != nil {
		h.Logger.Errorw("add dataset publication: could not execute search", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Partial(w, "publication/suggest_datasets", YieldAddDataset{
		Context: ctx,
		Hits:    hits,
	})
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
		h.Logger.Errorw("create dataset publication: could not get dataset", "errors", err, "publication", ctx.Publication.ID, "dataset", b.DatasetID, "user", ctx.User.ID)
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
		h.Logger.Errorw("create dataset publication: could not get publication", "errors", err, "publication", ctx.Publication.ID, "dataset", b.DatasetID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	relatedDatasets, err := h.Repo.GetVisiblePublicationDatasets(ctx.User, ctx.Publication)
	if err != nil {
		h.Logger.Errorw("create dataset publication: could not get related datasets", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_datasets", YieldDatasets{
		Context:         ctx,
		RelatedDatasets: relatedDatasets,
	})
}

func (h *Handler) ConfirmDeleteDataset(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeleteDataset{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("confirm delete publication dataset: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	if b.SnapshotID != ctx.Publication.SnapshotID {
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("publication.conflict_error_reload"),
		})
		return
	}

	render.Layout(w, "show_modal", "publication/confirm_delete_dataset", YieldDeleteDataset{
		Context:   ctx,
		DatasetID: b.DatasetID,
	})
}

func (h *Handler) DeleteDataset(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeleteDataset{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete publication dataset: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	if !ctx.Publication.HasRelatedDataset(b.DatasetID) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("publication.conflict_error_reload"),
		})
		return
	}

	// TODO reduce calls to repository
	d, err := h.Repo.GetDataset(b.DatasetID)
	if err != nil {
		h.Logger.Errorw("delete dataset publication: could not get dataset", "errors", err, "publication", ctx.Publication.ID, "dataset", b.DatasetID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// TODO handle validation errors
	// TODO pass If-Match
	// TODO handle conflict
	err = h.Repo.RemovePublicationDataset(ctx.Publication, d, ctx.User)

	if err != nil {
		h.Logger.Errorw("delete dataset publication: could not remove dataset", "errors", err, "publication", ctx.Publication.ID, "dataset", b.DatasetID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// Refresh the ctx.Publication: it still carries the old snapshotID
	ctx.Publication, err = h.Repo.GetPublication(ctx.Publication.ID)
	if err != nil {
		h.Logger.Errorw("delete dataset publication: could not get publication", "errors", err, "publication", ctx.Publication.ID, "dataset", b.DatasetID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	relatedDatasets, err := h.Repo.GetVisiblePublicationDatasets(ctx.User, ctx.Publication)
	if err != nil {
		h.Logger.Errorw("create dataset publication: could not get related datasets", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_datasets", YieldDatasets{
		Context:         ctx,
		RelatedDatasets: relatedDatasets,
	})
}

func (h *Handler) searchRelatedDatasets(user *models.Person, p *models.Publication, q string) (*models.DatasetHits, error) {
	args := models.NewSearchArgs().WithQuery(q)

	// add exclusion filter if necessary
	if len(p.RelatedDataset) > 0 {
		pubDatasetIDs := make([]string, len(p.RelatedDataset))
		for i, d := range p.RelatedDataset {
			pubDatasetIDs[i] = d.ID
		}
		args.Filters["!id"] = pubDatasetIDs
	}

	searchService := h.DatasetSearchIndex.WithScope("status", "public")

	return searchService.Search(args)
}
