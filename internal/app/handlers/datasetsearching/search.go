package datasetsearching

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
)

type YieldSearch struct {
	Context
	PageTitle string
	ActiveNav string
	Hits      *models.DatasetHits
}

type YieldHit struct {
	Context
	Dataset *models.Dataset
}

func (y YieldSearch) YieldHit(d *models.Dataset) YieldHit {
	return YieldHit{y.Context, d}
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request, ctx Context) {
	hits, err := h.userDatasets(ctx.User.ID, ctx.SearchArgs)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Wrap(w, "layouts/default", "dataset/search_page", YieldSearch{
		Context:   ctx,
		PageTitle: "Overview - Datasets - Biblio",
		ActiveNav: "datasets",
		Hits:      hits,
	})
}

func (h *Handler) userDatasets(userID string, args *models.SearchArgs) (*models.DatasetHits, error) {
	searcher := h.DatasetSearchService.WithScope("status", "private", "public")

	switch args.FilterFor("scope") {
	case "created":
		searcher = searcher.WithScope("creator_id", userID)
	case "contributed":
		searcher = searcher.WithScope("author.id", userID)
	default:
		searcher = searcher.WithScope("creator_id|author.id", userID)
	}
	delete(args.Filters, "scope")

	return searcher.IncludeFacets(true).Search(args)
}
