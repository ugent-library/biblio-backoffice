package datasetediting

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
)

type BindSuggestPublications struct {
	Query string `query:"q"`
}
type BindPublication struct {
	PublicationID string `form:"publication_id"`
}
type BindDeletePublication struct {
	PublicationID string `path:"publication_id"`
}

type YieldAddPublication struct {
	Context
	Hits *models.PublicationHits
}
type YieldPublications struct {
	Context
	RelatedPublications []*models.Publication
}
type YieldDeletePublication struct {
	Context
	PublicationID string
}

func (h *Handler) AddPublication(w http.ResponseWriter, r *http.Request, ctx Context) {
	hits, err := h.searchRelatedPublications(ctx.User.ID, ctx.Dataset, "")
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "dataset/add_publication", YieldAddPublication{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) SuggestPublications(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindSuggestPublications{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	hits, err := h.searchRelatedPublications(ctx.User.ID, ctx.Dataset, b.Query)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "dataset/suggest_publications", YieldAddPublication{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) CreatePublication(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindPublication{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	// TODO reduce calls to repository
	p, err := h.Repository.GetPublication(b.PublicationID)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	// TODO handle validation errors
	// TODO pass If-Match
	err = h.Repository.AddPublicationDataset(p, ctx.Dataset)

	// TODO handle conflict

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	relatedPublications, err := h.Repository.GetDatasetPublications(ctx.Dataset)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "dataset/refresh_publications", YieldPublications{
		Context:             ctx,
		RelatedPublications: relatedPublications,
	})
}

func (h *Handler) ConfirmDeletePublication(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeletePublication{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.Render(w, "dataset/confirm_delete_publication", YieldDeletePublication{
		Context:       ctx,
		PublicationID: b.PublicationID,
	})
}

func (h *Handler) DeletePublication(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeletePublication{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	// TODO reduce calls to repository
	p, err := h.Repository.GetPublication(b.PublicationID)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	// TODO handle validation errors
	// TODO pass If-Match
	err = h.Repository.RemovePublicationDataset(p, ctx.Dataset)

	// TODO handle conflict

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	relatedPublications, err := h.Repository.GetDatasetPublications(ctx.Dataset)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "dataset/refresh_publications", YieldPublications{
		Context:             ctx,
		RelatedPublications: relatedPublications,
	})
}

func (h *Handler) searchRelatedPublications(userID string, d *models.Dataset, q string) (*models.PublicationHits, error) {
	args := models.NewSearchArgs().WithQuery(q)

	// add exclusion filter if necessary
	if len(d.RelatedPublication) > 0 {
		datasetPubIDs := make([]string, len(d.RelatedPublication))
		for i, d := range d.RelatedPublication {
			datasetPubIDs[i] = d.ID
		}
		args.Filters["!id"] = datasetPubIDs
	}

	return h.PublicationSearchService.
		WithScope("status", "private", "public").
		WithScope("creator_id|author.id", userID).
		Search(args)
}
