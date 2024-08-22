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
		c.HandleError(w, r, err)
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

	// view publications of proxy
	personID := c.User.ID
	if proxiedPersonID := args.FilterFor("person"); proxiedPersonID != "" {
		personID = proxiedPersonID
	}

	switch args.FilterFor("scope") {
	case "created":
		searcher = searcher.WithScope("creator_id", personID)
		currentScope = "created"
	case "contributed":
		searcher = searcher.WithScope("author_id", personID)
		currentScope = "contributed"
	case "all":
		searcher = searcher.WithScope("creator_id|author_id", personID)
		currentScope = "all"
	default:
		c.HandleError(w, r, httperror.BadRequest.Wrap(fmt.Errorf("unknown scope: %s", args.FilterFor("scope"))))
		return
	}

	delete(args.Filters, "person")
	delete(args.Filters, "scope")

	hits, err := searcher.Search(args)
	if err != nil {
		c.HandleError(w, r, err)
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
			c.HandleError(w, r, globalHitsErr)
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

	if !c.Repo.CanCurate(c.User) {
		c.HandleError(w, r, httperror.Forbidden)
		return
	}

	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		c.HandleError(w, r, err)
		return
	}
	searchArgs.Cleanup()

	searchArgs.WithFacetLines(vocabularies.Facets["dataset_curation"])

	searcher := c.DatasetSearchIndex.WithScope("status", "private", "public", "returned")
	hits, err := searcher.Search(searchArgs)
	if err != nil {
		c.HandleError(w, r, err)
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
			c.HandleError(w, r, globalHitsErr)
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
