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
	engine *engine.Engine
	render *render.Render
}

func NewPublicationDatasets(e *engine.Engine, r *render.Render) *PublicationDatasets {
	return &PublicationDatasets{
		engine: e,
		render: r,
	}
}

func (c *PublicationDatasets) Choose(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	pubDatasets, err := c.engine.GetPublicationDatasets(id)
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

	hits, err := c.engine.UserDatasets(context.User(r.Context()).ID, searchArgs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.render.HTML(w, 200,
		"publication/datasets/_modal",
		struct {
			Publication *models.Publication
			Hits        *models.DatasetHits
		}{
			pub,
			hits,
		},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDatasets) ActiveSearch(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pubDatasets, err := c.engine.GetPublicationDatasets(id)
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

	hits, err := c.engine.UserDatasets(context.User(r.Context()).ID, searchArgs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.render.HTML(w, 200,
		"publication/datasets/_modal_hits",
		struct {
			Publication *models.Publication
			Hits        *models.DatasetHits
		}{
			pub,
			hits,
		},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDatasets) Add(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	datasetID := mux.Vars(r)["dataset_id"]

	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	_, err = c.engine.GetPublication(datasetID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = c.engine.AddPublicationDataset(id, datasetID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	datasets, err := c.engine.GetPublicationDatasets(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pub.Dataset = datasets

	c.render.HTML(w, 200,
		"publication/datasets/_content",
		views.NewPublicationData(r, c.render, pub),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDatasets) ConfirmRemove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	datasetID := mux.Vars(r)["dataset_id"]

	c.render.HTML(w, 200,
		"publication/datasets/_modal_confirm_removal",
		struct {
			PublicationID string
			DatasetID     string
		}{
			id,
			datasetID,
		},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationDatasets) Remove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	datasetID := mux.Vars(r)["dataset_id"]

	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = c.engine.RemovePublicationDataset(id, datasetID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	datasets, err := c.engine.GetPublicationDatasets(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pub.Dataset = datasets

	c.render.HTML(w, 200,
		"publication/datasets/_content",
		views.NewPublicationData(r, c.render, pub),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
