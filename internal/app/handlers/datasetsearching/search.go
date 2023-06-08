package datasetsearching

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/bind"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/render"
	"github.com/ugent-library/biblio-backoffice/internal/vocabularies"
)

var (
	userScopes = []string{"all", "contributed", "created"}
)

type ActionItem struct {
	Template string
	URL      *url.URL
	Label    string
}

type YieldSearch struct {
	Context
	PageTitle    string
	ActiveNav    string
	Scopes       []string
	Hits         *models.DatasetHits
	IsFirstUse   bool
	CurrentScope string
	ActionItems  []*ActionItem
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request, ctx Context) {
	if ctx.UserRole == "curator" {
		h.CurationSearch(w, r, ctx)
		return
	}

	ctx.SearchArgs.WithFacets(vocabularies.Map["dataset_facets"]...)
	if ctx.SearchArgs.FilterFor("scope") == "" {
		ctx.SearchArgs.WithFilter("scope", "all")
	}

	searcher := h.SearchService.NewDatasetIndex().WithScope("status", "private", "public", "returned")
	args := ctx.SearchArgs.Clone()
	var currentScope string

	switch args.FilterFor("scope") {
	case "created":
		searcher = searcher.WithScope("creator.id", ctx.User.ID)
		currentScope = "created"
	case "contributed":
		searcher = searcher.WithScope("author.id", ctx.User.ID)
		currentScope = "contributed"
	case "all":
		searcher = searcher.WithScope("creator.id|author.id", ctx.User.ID)
		currentScope = "all"
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

	/*
		first use search
		when more applicable, always execute this
		now only when no results are found
	*/
	var isFirstUse bool = false
	if hits.Total == 0 {
		globalHits, globalHitsErr := globalSearch(searcher)
		if globalHitsErr != nil {
			h.Logger.Errorw("publication search: could not execute global search", "errors", globalHitsErr, "user", ctx.User.ID)
			render.InternalServerError(w, r, globalHitsErr)
			return
		}
		isFirstUse = globalHits.Total == 0
	}

	render.Layout(w, "layouts/default", "dataset/pages/search", YieldSearch{
		Context:      ctx,
		PageTitle:    "Overview - Datasets - Biblio",
		ActiveNav:    "datasets",
		Scopes:       userScopes,
		Hits:         hits,
		IsFirstUse:   isFirstUse,
		CurrentScope: currentScope,
		ActionItems:  h.getDatasetActions(ctx),
	})
}

/*
globalSearch(searcher)

	returns total number of search hits
	for scoped searcher, regardless of choosen filters
	Used to determine wether user has any records
*/
func globalSearch(searcher backends.DatasetIndex) (*models.DatasetHits, error) {
	globalArgs := models.NewSearchArgs()
	globalArgs.Query = ""
	globalArgs.Facets = nil
	globalArgs.Filters = map[string][]string{}
	globalArgs.PageSize = 0
	globalArgs.Page = 1
	return searcher.Search(globalArgs)
}

func (h *Handler) CurationSearch(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCurate() {
		render.Forbidden(w, r)
		return
	}

	ctx.SearchArgs.WithFacets(vocabularies.Map["dataset_curation_facets"]...)

	searcher := h.SearchService.NewDatasetIndex().WithScope("status", "private", "public", "returned")
	hits, err := searcher.Search(ctx.SearchArgs)
	if err != nil {
		h.Logger.Errorw("dataset search: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	/*
		first use search
		when more applicable, always execute this
		now only when no results are found
	*/
	var isFirstUse bool = false
	if hits.Total == 0 {
		globalHits, globalHitsErr := globalSearch(searcher)
		if globalHitsErr != nil {
			h.Logger.Errorw("publication search: could not execute global search", "errors", globalHitsErr, "user", ctx.User.ID)
			render.InternalServerError(w, r, globalHitsErr)
			return
		}
		isFirstUse = globalHits.Total == 0
	}

	render.Layout(w, "layouts/default", "dataset/pages/search", YieldSearch{
		Context:      ctx,
		PageTitle:    "Overview - Datasets - Biblio",
		ActiveNav:    "datasets",
		Hits:         hits,
		IsFirstUse:   isFirstUse,
		CurrentScope: "all", //only here to translate first use
		ActionItems:  h.getCurationDatasetActions(ctx),
	})
}

func (h *Handler) getDatasetActions(ctx Context) []*ActionItem {
	return []*ActionItem{}
}

func (h *Handler) getCurationDatasetActions(ctx Context) []*ActionItem {
	actionItems := make([]*ActionItem, 0)
	u := h.PathFor("export_datasets", "format", "xlsx")
	q, _ := bind.EncodeQuery(ctx.SearchArgs)
	u.RawQuery = q.Encode()
	actionItems = append(actionItems, &ActionItem{
		Label:    ctx.Locale.T("export_to.xlsx"),
		URL:      u,
		Template: "actions/export",
	})
	return actionItems
}
