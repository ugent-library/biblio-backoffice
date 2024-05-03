package datasetsearching

import (
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	datasetviews "github.com/ugent-library/biblio-backoffice/views/dataset"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
)

var (
	userScopes = []string{"all", "contributed", "created"}
)

func Search(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if c.UserRole == "curator" {
		CurationSearch(w, r)
		return
	}

	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		c.Log.Warnw("dataset search: could not bind search arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}
	searchArgs.Cleanup()

	searchArgs.WithFacetLines(vocabularies.Facets["dataset"])
	if searchArgs.FilterFor("scope") == "" {
		searchArgs.WithFilter("scope", "all")
	}

	searcher := c.DatasetSearchIndex.WithScope("status", "private", "public", "returned")
	args := searchArgs.Clone()
	var currentScope string

	switch args.FilterFor("scope") {
	case "created":
		searcher = searcher.WithScope("creator_id", c.User.ID)
		currentScope = "created"
	case "contributed":
		searcher = searcher.WithScope("author_id", c.User.ID)
		currentScope = "contributed"
	case "all":
		searcher = searcher.WithScope("creator_id|author_id", c.User.ID)
		currentScope = "all"
	default:
		errorUnkownScope := fmt.Errorf("unknown scope: %s", args.FilterFor("scope"))
		c.Log.Warnw("dataset search: could not create searcher with passed filters", "errors", errorUnkownScope, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}
	delete(args.Filters, "scope")

	hits, err := searcher.Search(args)
	if err != nil {
		c.Log.Errorw("dataset search: could not execute search", "errors", err, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
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
			c.Log.Errorw("publication search: could not execute global search", "errors", globalHitsErr, "user", c.User.ID)
			c.HandleError(w, r, httperror.InternalServerError)
			return
		}
		isFirstUse = globalHits.Total == 0
	}

	// you are on the wrong page: cap page to last available page
	if hits.Total > 0 && len(hits.Hits) == 0 {
		query := c.CurrentURL.Query()
		query.Set("page", fmt.Sprintf("%d", hits.TotalPages()))
		c.CurrentURL.RawQuery = query.Encode()
		http.Redirect(w, r, c.CurrentURL.String(), http.StatusTemporaryRedirect)
		return
	}

	datasetviews.Search(c, &datasetviews.SearchArgs{
		Scopes:       userScopes,
		Hits:         hits,
		IsFirstUse:   isFirstUse,
		CurrentScope: currentScope,
		SearchArgs:   searchArgs,
	}).Render(r.Context(), w)
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

func CurationSearch(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if !c.User.CanCurate() {
		c.HandleError(w, r, httperror.Forbidden)
		return
	}

	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		c.Log.Warnw("dataset search: could not bind search arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}
	searchArgs.Cleanup()

	searchArgs.WithFacetLines(vocabularies.Facets["dataset_curation"])

	searcher := c.DatasetSearchIndex.WithScope("status", "private", "public", "returned")
	hits, err := searcher.Search(searchArgs)
	if err != nil {
		c.Log.Errorw("dataset search: could not execute search", "errors", err, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
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
			c.Log.Errorw("publication search: could not execute global search", "errors", globalHitsErr, "user", c.User.ID)
			c.HandleError(w, r, httperror.InternalServerError)
			return
		}
		isFirstUse = globalHits.Total == 0
	}

	// you are on the wrong page: cap page to last available page
	if hits.Total > 0 && len(hits.Hits) == 0 {
		query := c.CurrentURL.Query()
		query.Set("page", fmt.Sprintf("%d", hits.TotalPages()))
		c.CurrentURL.RawQuery = query.Encode()
		http.Redirect(w, r, c.CurrentURL.String(), http.StatusTemporaryRedirect)
		return
	}

	datasetviews.Search(c, &datasetviews.SearchArgs{
		Scopes:       nil,
		Hits:         hits,
		IsFirstUse:   isFirstUse,
		CurrentScope: "all", //only here to translate first use
		SearchArgs:   searchArgs,
	}).Render(r.Context(), w)
}
