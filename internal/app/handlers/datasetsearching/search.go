package datasetsearching

import (
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
)

var (
	userScopes = []string{"all", "contributed", "created"}
)

type YieldSearch struct {
	Context
	PageTitle string
	ActiveNav string
	Scopes    []string
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
	if ctx.SearchArgs.FilterFor("scope") == "" {
		ctx.SearchArgs.WithFilter("scope", "all")
	}

	searcher := h.DatasetSearchService.WithScope("status", "private", "public")
	args := ctx.SearchArgs.Clone()

	switch args.FilterFor("scope") {
	case "created":
		searcher = searcher.WithScope("creator_id", ctx.User.ID)
	case "contributed":
		searcher = searcher.WithScope("author.id", ctx.User.ID)
	case "all":
		searcher = searcher.WithScope("creator_id|author.id", ctx.User.ID)
	default:
		render.BadRequest(w, r, fmt.Errorf("unknown scope %s", args.FilterFor("scope")))
		return
	}
	delete(args.Filters, "scope")

	hits, err := searcher.IncludeFacets(true).Search(args)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.MustRenderLayout(w, "layouts/default", "dataset/pages/search", YieldSearch{
		Context:   ctx,
		PageTitle: "Overview - Datasets - Biblio",
		ActiveNav: "datasets",
		Scopes:    userScopes,
		Hits:      hits,
	})
}

func (h *Handler) CurationSearch(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCurateDatasets() {
		render.Forbidden(w, r)
		return
	}

	searcher := h.DatasetSearchService.WithScope("status", "private", "public")
	hits, err := searcher.IncludeFacets(true).Search(ctx.SearchArgs)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.MustRenderLayout(w, "layouts/default", "dataset/search_page", YieldSearch{
		Context:   ctx,
		PageTitle: "Overview - Datasets - Biblio",
		ActiveNav: "datasets",
		Hits:      hits,
	})
}
