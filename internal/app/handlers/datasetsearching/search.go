package datasetsearching

import (
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
)

var (
	userScopes = []string{"all", "contributed", "created"}
)

type YieldSearch struct {
	Context
	PageTitle   string
	ActiveNav   string
	Scopes      []string
	Hits        *models.DatasetHits
	ActionItems []*models.ActionItem
}

type YieldHit struct {
	Context
	Dataset *models.Dataset
}

func (y YieldSearch) YieldHit(d *models.Dataset) YieldHit {
	return YieldHit{y.Context, d}
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request, ctx Context) {
	ctx.SearchArgs.WithFacets(vocabularies.Map["dataset_facets"]...)
	if ctx.SearchArgs.FilterFor("scope") == "" {
		ctx.SearchArgs.WithFilter("scope", "all")
	}

	searcher := h.DatasetSearchService.WithScope("status", "private", "public")
	args := ctx.SearchArgs.Clone()

	switch args.FilterFor("scope") {
	case "created":
		searcher = searcher.WithScope("creator.id", ctx.User.ID)
	case "contributed":
		searcher = searcher.WithScope("author.id", ctx.User.ID)
	case "all":
		searcher = searcher.WithScope("creator.id|author.id", ctx.User.ID)
	default:
		errorUnkownScope := fmt.Errorf("unknown scope: %s", args.FilterFor("scope"))
		h.Logger.Warnw("dataset search: could not create searcher with passed filters", "errors", errorUnkownScope, "user", ctx.User.ID)
		render.BadRequest(w, r, errorUnkownScope)
		return
	}
	delete(args.Filters, "scope")

	hits, err := searcher.Search(args)
	if err != nil {
		h.Logger.Errorw("dataset search: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "layouts/default", "dataset/pages/search", YieldSearch{
		Context:     ctx,
		PageTitle:   "Overview - Datasets - Biblio",
		ActiveNav:   "datasets",
		Scopes:      userScopes,
		Hits:        hits,
		ActionItems: h.getDatasetActions(ctx),
	})
}

func (h *Handler) CurationSearch(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCurateDatasets() {
		render.Forbidden(w, r)
		return
	}

	ctx.SearchArgs.WithFacets(vocabularies.Map["dataset_curation_facets"]...)

	searcher := h.DatasetSearchService.WithScope("status", "private", "public")
	hits, err := searcher.Search(ctx.SearchArgs)
	if err != nil {
		h.Logger.Errorw("dataset search: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "layouts/default", "dataset/pages/search", YieldSearch{
		Context:     ctx,
		PageTitle:   "Overview - Datasets - Biblio",
		ActiveNav:   "datasets",
		Hits:        hits,
		ActionItems: h.getCurationDatasetActions(ctx),
	})
}

func (h *Handler) getDatasetActions(ctx Context) []*models.ActionItem {
	return []*models.ActionItem{}
}

func (h *Handler) getCurationDatasetActions(ctx Context) []*models.ActionItem {
	actionItems := make([]*models.ActionItem, 0)
	u := h.PathFor("export_curation_datasets", "format", "xlsx")
	q, _ := bind.EncodeQuery(ctx.SearchArgs)
	u.RawQuery = q.Encode()
	actionItems = append(actionItems, &models.ActionItem{
		Label:    ctx.Locale.T("export_to.xlsx"),
		URL:      u,
		Template: "actions/export",
	})
	return actionItems
}
