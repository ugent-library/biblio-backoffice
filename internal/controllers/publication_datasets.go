package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
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

	pubDatasets, err := c.Engine.GetPublicationDatasets(pub.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pubDatasetIDs := make([]string, len(pubDatasets))
	for i, d := range pubDatasets {
		pubDatasetIDs[i] = d.ID
	}

	searchArgs := engine.NewSearchArgs()
	searchArgs.Filters["exclude"] = pubDatasetIDs

	hits, err := c.Engine.UserDatasets(context.GetUser(r.Context()).ID, searchArgs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/datasets/_modal", views.NewData(c.Render, r, struct {
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

	pubDatasets, err := c.Engine.GetPublicationDatasets(pub.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pubDatasetIDs := make([]string, len(pubDatasets))
	for i, d := range pubDatasets {
		pubDatasetIDs[i] = d.ID
	}

	searchArgs := engine.NewSearchArgs()
	searchArgs.Query = r.Form["search"][0]
	searchArgs.Filters["exclude"] = pubDatasetIDs

	hits, err := c.Engine.UserDatasets(context.GetUser(r.Context()).ID, searchArgs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/datasets/_modal_hits", views.NewData(c.Render, r, struct {
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

	_, err := c.Engine.GetDataset(datasetID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = c.Engine.AddPublicationDataset(pub.ID, datasetID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	datasets, err := c.Engine.GetPublicationDatasets(pub.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/datasets/_show", views.NewData(c.Render, r, struct {
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

	c.Render.HTML(w, http.StatusOK, "publication/datasets/_modal_confirm_removal", views.NewData(c.Render, r, struct {
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

	err := c.Engine.RemovePublicationDataset(pub.ID, datasetID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	datasets, err := c.Engine.GetPublicationDatasets(pub.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/datasets/_show", views.NewData(c.Render, r, struct {
		Publication         *models.Publication
		PublicationDatasets []*models.Dataset
	}{
		pub,
		datasets,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
