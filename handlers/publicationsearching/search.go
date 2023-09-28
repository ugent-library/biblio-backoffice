package publicationsearching

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/bind"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"

	"github.com/ugent-library/biblio-backoffice/vocabularies"
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
	CurrentScope string
	IsFirstUse   bool
	Hits         *models.PublicationHits
	ActionItems  []*ActionItem
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request, ctx Context) {
	if ctx.UserRole == "curator" {
		h.CurationSearch(w, r, ctx)
		return
	}

	ctx.SearchArgs.WithFacets(vocabularies.Map["publication_facets"]...)
	if ctx.SearchArgs.FilterFor("scope") == "" {
		ctx.SearchArgs.WithFilter("scope", "all")
	}

	searcher := h.PublicationSearchIndex.WithScope("status", "private", "public", "returned")
	args := ctx.SearchArgs.Clone()
	var currentScope string

	switch args.FilterFor("scope") {
	case "created":
		searcher = searcher.WithScope("creator_id", ctx.User.ID)
		currentScope = "created"
	case "contributed":
		searcher = searcher.WithScope("author_id", ctx.User.ID)
		currentScope = "contributed"
	case "all":
		searcher = searcher.WithScope("creator_id|author_id", ctx.User.ID)
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
		ActionItems:  h.getSearchActions(ctx),
	})
}

/*
globalSearch(searcher)

	returns total number of search hits
	for scoped searcher, regardless of choosen filters
	Used to determine wether user has any records
*/
func globalSearch(searcher backends.PublicationIndex) (*models.PublicationHits, error) {
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

	ctx.SearchArgs.WithFacets(vocabularies.Map["publication_curation_facets"]...)

	searcher := h.PublicationSearchIndex.WithScope("status", "private", "public", "returned")
	hits, err := searcher.Search(ctx.SearchArgs)
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
			h.Logger.Errorw("curation publication search: could not execute global search", "errors", globalHitsErr, "user", ctx.User.ID)
			render.InternalServerError(w, r, globalHitsErr)
			return
		}
		isFirstUse = globalHits.Total == 0
	}

	render.Layout(w, "layouts/default", "publication/pages/search", YieldSearch{
		Context:      ctx,
		PageTitle:    "Overview - Publications - Biblio",
		ActiveNav:    "publications",
		Hits:         hits,
		IsFirstUse:   isFirstUse,
		CurrentScope: "all", //only here to translate first use
		ActionItems:  h.getCurationSearchActions(ctx),
	})
}

func (h *Handler) getCurationSearchActions(ctx Context) []*ActionItem {
	actionItems := make([]*ActionItem, 0)
	// if oa := h.getOrcidAction(ctx); oa != nil {
	// 	actionItems = append(actionItems, oa)
	// }
	u := h.PathFor("export_publications", "format", "xlsx")
	q, _ := bind.EncodeQuery(ctx.SearchArgs)
	u.RawQuery = q.Encode()
	actionItems = append(actionItems, &ActionItem{
		Label:    ctx.Locale.T("export_to.xlsx"),
		URL:      u,
		Template: "actions/export",
	})
	return actionItems
}

func (h *Handler) getSearchActions(ctx Context) []*ActionItem {
	// actionItems := make([]*ActionItem, 0)
	// if oa := h.getOrcidAction(ctx); oa != nil {
	// 	actionItems = append(actionItems, oa)
	// }
	// return actionItems
	return nil
}

// func (h *Handler) getOrcidAction(ctx Context) *ActionItem {
// 	if ctx.User.ORCID == "" || ctx.User.ORCIDToken == "" {
// 		return nil
// 	}
// 	return &ActionItem{
// 		Label:    "Send my publications to ORCID",
// 		URL:      h.PathFor("publication_orcid_add_all"),
// 		Template: "actions/publication_orcid_add_all",
// 	}
// }
