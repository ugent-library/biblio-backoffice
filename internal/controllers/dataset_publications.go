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

type DatasetPublications struct {
	Context
}

func NewDatasetPublications(c Context) *DatasetPublications {
	return &DatasetPublications{c}
}

func (c *DatasetPublications) Choose(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	datasetPubs, err := c.Engine.GetDatasetPublications(dataset.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	datasetPubIDs := make([]string, len(datasetPubs))
	for i, d := range datasetPubs {
		datasetPubIDs[i] = d.ID
	}

	searchArgs := engine.NewSearchArgs()
	searchArgs.Filters["exclude"] = datasetPubIDs

	hits, err := c.Engine.UserPublications(context.GetUser(r.Context()).ID, searchArgs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/publications/_modal", struct {
		Dataset *models.Dataset
		Hits    *models.PublicationHits
	}{
		dataset,
		hits,
	},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetPublications) ActiveSearch(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	datasetPubs, err := c.Engine.GetDatasetPublications(dataset.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	datasetPubIDs := make([]string, len(datasetPubs))
	for i, d := range datasetPubs {
		datasetPubIDs[i] = d.ID
	}

	searchArgs := engine.NewSearchArgs()
	searchArgs.Query = r.Form["search"][0]
	searchArgs.Filters["exclude"] = datasetPubIDs

	hits, err := c.Engine.UserPublications(context.GetUser(r.Context()).ID, searchArgs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/publications/_modal_hits", struct {
		Dataset *models.Dataset
		Hits    *models.PublicationHits
	}{
		dataset,
		hits,
	},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetPublications) Add(w http.ResponseWriter, r *http.Request) {
	pubID := mux.Vars(r)["publication_id"]

	dataset := context.GetDataset(r.Context())

	_, err := c.Engine.GetPublication(pubID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = c.Engine.AddPublicationDataset(pubID, dataset.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	publications, err := c.Engine.GetDatasetPublications(dataset.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/publications/_show", views.NewData(c.Render, r, struct {
		Dataset             *models.Dataset
		DatasetPublications []*models.Publication
	}{
		dataset,
		publications,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetPublications) ConfirmRemove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	pubID := mux.Vars(r)["publication_id"]

	c.Render.HTML(w, http.StatusOK, "dataset/publications/_modal_confirm_removal", struct {
		DatasetID     string
		PublicationID string
	}{
		id,
		pubID,
	},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetPublications) Remove(w http.ResponseWriter, r *http.Request) {
	pubID := mux.Vars(r)["publication_id"]

	dataset := context.GetDataset(r.Context())

	err := c.Engine.RemovePublicationDataset(pubID, dataset.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	publications, err := c.Engine.GetDatasetPublications(dataset.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/publications/_show", views.NewData(c.Render, r, struct {
		Dataset             *models.Dataset
		DatasetPublications []*models.Publication
	}{
		dataset,
		publications,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
