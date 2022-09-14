package publicationsearching

import (
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
)

var (
	userScopes = []string{"all", "contributed", "created"}
)

type YieldSearch struct {
	Context
	PageTitle    string
	ActiveNav    string
	Scopes       []string
	CurrentScope string
	IsFirstUse   bool
	Hits         *models.PublicationHits
}

type YieldHit struct {
	Context
	Publication *models.Publication
}

func (y YieldSearch) YieldHit(d *models.Publication) YieldHit {
	return YieldHit{y.Context, d}
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request, ctx Context) {
	ctx.SearchArgs.WithFacets(vocabularies.Map["publication_facets"]...)
	if ctx.SearchArgs.FilterFor("scope") == "" {
		ctx.SearchArgs.WithFilter("scope", "all")
	}

	searcher := h.PublicationSearchService.WithScope("status", "private", "public")
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
		h.Logger.Warnw("publication search: could not create searcher with passed filters", "errors", errorUnkownScope, "user", ctx.User.ID)
		render.BadRequest(w, r, errorUnkownScope)
		return
	}
	delete(args.Filters, "scope")

	hits, err := searcher.Search(args)
	if err != nil {
		h.Logger.Errorw("publication search: could not execute search", "errors", err, "user", ctx.User.ID)
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

	render.Layout(w, "layouts/default", "publication/pages/search", YieldSearch{
		Context:      ctx,
		PageTitle:    "Overview - Publications - Biblio",
		ActiveNav:    "publications",
		Scopes:       userScopes,
		Hits:         hits,
		IsFirstUse:   isFirstUse,
		CurrentScope: currentScope,
	})
}

/*
	globalSearch(searcher)
		returns total number of search hits
		for scoped searcher, regardless of choosen filters
		Used to determine wether user has any records
*/
func globalSearch(searcher backends.PublicationSearchService) (*models.PublicationHits, error) {
	globalArgs := models.NewSearchArgs()
	globalArgs.Query = ""
	globalArgs.Facets = nil
	globalArgs.Filters = map[string][]string{}
	globalArgs.PageSize = 0
	globalArgs.Page = 1
	return searcher.Search(globalArgs)
}

func (h *Handler) CurationSearch(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCuratePublications() {
		render.Forbidden(w, r)
		return
	}

	ctx.SearchArgs.WithFacets(vocabularies.Map["publication_curation_facets"]...)

	searcher := h.PublicationSearchService.WithScope("status", "private", "public")
	hits, err := searcher.Search(ctx.SearchArgs)
	if err != nil {
		h.Logger.Errorw("publication search: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "layouts/default", "publication/pages/search", YieldSearch{
		Context:   ctx,
		PageTitle: "Overview - Publications - Biblio",
		ActiveNav: "publications",
		Hits:      hits,
	})
}
