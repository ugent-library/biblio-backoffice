package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/unrolled/render"
)

type PublicationDatasets struct {
	Context
}

func NewPublicationDatasets(c Context) *PublicationDatasets {
	return &PublicationDatasets{c}
}

func (c *PublicationDatasets) Choose(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	pubDatasets, err := c.Engine.Store.GetPublicationDatasets(pub)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pubDatasetIDs := make([]string, len(pubDatasets))
	for i, d := range pubDatasets {
		pubDatasetIDs[i] = d.ID
	}

	searchArgs := models.NewSearchArgs()
	searchArgs.Filters["!id"] = pubDatasetIDs

	hits, err := c.userDatasets(context.GetUser(r.Context()).ID, searchArgs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/datasets/_modal", c.ViewData(r, struct {
		Publication *models.Publication
		Hits        *models.DatasetHits
	}{
		pub,
		hits,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDatasets) ActiveSearch(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pubDatasets, err := c.Engine.Store.GetPublicationDatasets(pub)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pubDatasetIDs := make([]string, len(pubDatasets))
	for i, d := range pubDatasets {
		pubDatasetIDs[i] = d.ID
	}

	searchArgs := models.NewSearchArgs()
	searchArgs.Query = r.Form["search"][0]
	searchArgs.Filters["exclude"] = pubDatasetIDs

	hits, err := c.userDatasets(context.GetUser(r.Context()).ID, searchArgs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/datasets/_modal_hits", c.ViewData(r, struct {
		Publication *models.Publication
		Hits        *models.DatasetHits
	}{
		pub,
		hits,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDatasets) Add(w http.ResponseWriter, r *http.Request) {
	datasetID := mux.Vars(r)["dataset_id"]

	pub := context.GetPublication(r.Context())

	dataset, err := c.Engine.Store.GetDataset(datasetID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = c.Engine.Store.AddPublicationDataset(pub, dataset); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	datasets, err := c.Engine.Store.GetPublicationDatasets(pub)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/datasets/_show", c.ViewData(r, struct {
		Publication         *models.Publication
		PublicationDatasets []*models.Dataset
	}{
		pub,
		datasets,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDatasets) ConfirmRemove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	datasetID := mux.Vars(r)["dataset_id"]

	c.Render.HTML(w, http.StatusOK, "publication/datasets/_modal_confirm_removal", c.ViewData(r, struct {
		PublicationID string
		DatasetID     string
	}{
		id,
		datasetID,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDatasets) Remove(w http.ResponseWriter, r *http.Request) {
	datasetID := mux.Vars(r)["dataset_id"]

	pub := context.GetPublication(r.Context())

	dataset, err := c.Engine.Store.GetDataset(datasetID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newPub := pub.Clone()
	if err := c.Engine.Store.RemovePublicationDataset(newPub, dataset); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	datasets, err := c.Engine.Store.GetPublicationDatasets(newPub)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/datasets/_show", c.ViewData(r, struct {
		Publication         *models.Publication
		PublicationDatasets []*models.Dataset
	}{
		newPub,
		datasets,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDatasets) userDatasets(userID string, args *models.SearchArgs) (*models.DatasetHits, error) {
	args = args.Clone().WithFilter("status", "private", "public")
	switch args.FilterFor("scope") {
	case "created":
		args.WithFilter("creator_id", userID)
	case "contributed":
		args.WithFilter("author.id", userID)
	default:
		args.WithFilter("creator_id|author.id", userID)
	}
	delete(args.Filters, "scope")
	return c.Engine.DatasetSearchService.SearchDatasets(args)
}
