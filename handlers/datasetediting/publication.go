package datasetediting

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/bind"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
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
	hits, err := h.searchRelatedPublications(ctx.User, ctx.Dataset, "")
	if err != nil {
		h.Logger.Errorf("add dataset publication: Could find related publications:", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "show_modal", "dataset/add_publication", YieldAddPublication{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) SuggestPublications(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindSuggestPublications{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("suggest dataset publications: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	hits, err := h.searchRelatedPublications(ctx.User, ctx.Dataset, b.Query)
	if err != nil {
		h.Logger.Errorf("add dataset publication: Could find related publications:", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Partial(w, "dataset/suggest_publications", YieldAddPublication{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) CreatePublication(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindPublication{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("create dataset publication: could not bind request arguments", "errors", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	// TODO reduce calls to repository
	p, err := h.Repo.GetPublication(b.PublicationID)
	if err != nil {
		h.Logger.Errorw("create dataset publication: could not get the publication", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// TODO handle validation errors
	// TODO pass If-Match
	err = h.Repo.AddPublicationDataset(p, ctx.Dataset, ctx.User)

	// TODO handle conflict

	if err != nil {
		h.Logger.Errorw("create dataset publication: could not add the publication", "error", err, "dataset", ctx.Dataset.ID, "publication", p.ID)
		render.InternalServerError(w, r, err)
		return
	}

	relatedPublications, err := h.Repo.GetVisibleDatasetPublications(ctx.User, ctx.Dataset)
	if err != nil {
		h.Logger.Errorw("create dataset publication: could not get dataset publications", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_publications", YieldPublications{
		Context:             ctx,
		RelatedPublications: relatedPublications,
	})
}

func (h *Handler) ConfirmDeletePublication(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeletePublication{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("confirm delete dataset publication: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	if b.SnapshotID != ctx.Dataset.SnapshotID {
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Locale.T("dataset.conflict_error_reload"),
		})
		return
	}

	render.Layout(w, "show_modal", "dataset/confirm_delete_publication", YieldDeletePublication{
		Context:       ctx,
		PublicationID: b.PublicationID,
	})
}

func (h *Handler) DeletePublication(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeletePublication{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete dataset publication: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	// TODO reduce calls to repository
	p, err := h.Repo.GetPublication(b.PublicationID)
	if err != nil {
		h.Logger.Errorw("delete dataset publication: could not get the publication", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// TODO handle validation errors
	// TODO pass If-Match
	// TODO handle conflict
	err = h.Repo.RemovePublicationDataset(p, ctx.Dataset, ctx.User)

	if err != nil {
		h.Logger.Errorw("delete dataset publication: could not delete the publication", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// Refresh the dataset since it still caries the old snapshotid
	ctx.Dataset, err = h.Repo.GetDataset(ctx.Dataset.ID)
	if err != nil {
		h.Logger.Errorw("delete dataset publication: could not get the dataset", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	relatedPublications, err := h.Repo.GetVisibleDatasetPublications(ctx.User, ctx.Dataset)
	if err != nil {
		h.Logger.Errorw("create dataset publication: could not get dataset publications", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_publications", YieldPublications{
		Context:             ctx,
		RelatedPublications: relatedPublications,
	})
}

func (h *Handler) searchRelatedPublications(user *models.Person, d *models.Dataset, q string) (*models.PublicationHits, error) {
	args := models.NewSearchArgs().WithQuery(q)

	// add exclusion filter if necessary
	if len(d.RelatedPublication) > 0 {
		datasetPubIDs := make([]string, len(d.RelatedPublication))
		for i, d := range d.RelatedPublication {
			datasetPubIDs[i] = d.ID
		}
		args.Filters["!id"] = datasetPubIDs
	}

	searchService := h.PublicationSearchIndex.WithScope("status", "public")

	return searchService.Search(args)
}
