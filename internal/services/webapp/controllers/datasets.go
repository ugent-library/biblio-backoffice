package controllers

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/context"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/views"
	"github.com/unrolled/render"
)

type Datasets struct {
	Base
	store                backends.Repository
	datasetSearchService backends.DatasetSearchService
	datasetSources       map[string]backends.DatasetGetter
}

func NewDatasets(base Base, store backends.Repository, datasetSearchService backends.DatasetSearchService,
	datasetSources map[string]backends.DatasetGetter) *Datasets {
	return &Datasets{
		Base:                 base,
		store:                store,
		datasetSearchService: datasetSearchService,
		datasetSources:       datasetSources,
	}
}

func (c *Datasets) ConfirmDelete(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	searchArgs := models.NewSearchArgs()
	if err := DecodeQuery(searchArgs, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/_confirm_delete", c.ViewData(r, struct {
		Dataset    *models.Dataset
		SearchArgs *models.SearchArgs
	}{
		dataset,
		searchArgs,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *Datasets) Delete(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	r.ParseForm()
	searchArgs := models.NewSearchArgs()
	if err := DecodeQuery(searchArgs, r.Form); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataset.Status = "deleted"
	if err := c.store.SaveDataset(dataset); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hits, err := c.userDatasets(context.GetUser(r.Context()).ID, searchArgs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("datasets").URLPath()

	c.Render.HTML(w, http.StatusOK, "dataset/list", c.ViewData(r, struct {
		PageTitle  string
		SearchURL  *url.URL
		SearchArgs *models.SearchArgs
		Hits       *models.DatasetHits
	}{
		"Overview - Datasets - Biblio",
		searchURL,
		searchArgs,
		hits,
	},
		views.Flash{Type: "success", Message: "Successfully deleted dataset.", DismissAfter: 5 * time.Second},
	),
	)
}

func (c *Datasets) userDatasets(userID string, args *models.SearchArgs) (*models.DatasetHits, error) {
	searcher := c.datasetSearchService.WithScope("status", "private", "public")

	switch args.FilterFor("scope") {
	case "created":
		searcher = searcher.WithScope("creator_id", userID)
	case "contributed":
		searcher = searcher.WithScope("author.id", userID)
	default:
		searcher = searcher.WithScope("creator_id|author.id", userID)
	}
	delete(args.Filters, "scope")

	searcher = searcher.IncludeFacets(true)

	return searcher.Search(args)
}
