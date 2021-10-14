package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/models"
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

	hits, err := c.engine.UserDatasets(context.User(r.Context()).ID, engine.NewSearchArgs())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.render.HTML(w, 200,
		"publication/datasets/_modal",
		struct {
			Publication *models.Publication
			Hits        *models.PublicationHits
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

	q := r.Form["search"][0]
	hits, err := c.engine.UserDatasets(context.User(r.Context()).ID, engine.NewSearchArgs().WithQuery(q))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.render.HTML(w, 200,
		"publication/datasets/_modal_hits",
		struct {
			Publication *models.Publication
			Hits        *models.PublicationHits
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

	err = c.engine.AddRelatedPublication(id, datasetID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	related, err := c.engine.GetRelatedPublications(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.render.HTML(w, 200,
		"publication/datasets/_list",
		struct {
			Publication         *models.Publication
			RelatedPublications []*models.RelatedPublication
		}{
			Publication:         pub,
			RelatedPublications: related,
		},
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
