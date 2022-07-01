package publicationediting

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
)

type BindSuggestDatasets struct {
	Query string `query:"q"`
}
type BindDataset struct {
	DatasetID string `form:"dataset_id"`
}
type BindDeleteDataset struct {
	DatasetID string `path:"dataset_id"`
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
	hits, err := h.searchRelatedDatasets(ctx.User.ID, ctx.Publication, "")
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "publication/add_dataset", YieldAddDataset{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) SuggestDatasets(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindSuggestDatasets{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	hits, err := h.searchRelatedDatasets(ctx.User.ID, ctx.Publication, b.Query)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "publication/suggest_datasets", YieldAddDataset{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) CreateDataset(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDataset{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	// TODO reduce calls to repository
	d, err := h.Repository.GetDataset(b.DatasetID)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	// TODO handle validation errors
	// TODO pass If-Match
	err = h.Repository.AddPublicationDataset(ctx.Publication, d)

	// TODO handle conflict

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	relatedDatasets, err := h.Repository.GetPublicationDatasets(ctx.Publication)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "publication/refresh_datasets", YieldDatasets{
		Context:         ctx,
		RelatedDatasets: relatedDatasets,
	})
}

func (h *Handler) ConfirmDeleteDataset(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeleteDataset{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.Render(w, "publication/confirm_delete_dataset", YieldDeleteDataset{
		Context:   ctx,
		DatasetID: b.DatasetID,
	})
}

func (h *Handler) DeleteDataset(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeleteDataset{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	// TODO reduce calls to repository
	d, err := h.Repository.GetDataset(b.DatasetID)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	// TODO handle validation errors
	// TODO pass If-Match
	err = h.Repository.RemovePublicationDataset(ctx.Publication, d)

	// TODO handle conflict

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	relatedDatasets, err := h.Repository.GetPublicationDatasets(ctx.Publication)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "publication/refresh_datasets", YieldDatasets{
		Context:         ctx,
		RelatedDatasets: relatedDatasets,
	})
}

func (h *Handler) searchRelatedDatasets(userID string, p *models.Publication, q string) (*models.DatasetHits, error) {
	args := models.NewSearchArgs().WithQuery(q)

	// add exclusion filter if necessary
	if len(p.RelatedDataset) > 0 {
		pubDatasetIDs := make([]string, len(p.RelatedDataset))
		for i, d := range p.RelatedDataset {
			pubDatasetIDs[i] = d.ID
		}
		args.Filters["!id"] = pubDatasetIDs
	}

	return h.DatasetSearchService.
		WithScope("status", "private", "public").
		WithScope("creator_id|author.id", userID).
		Search(args)
}
