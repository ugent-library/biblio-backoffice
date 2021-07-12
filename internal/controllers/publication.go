package controllers

import (
	"log"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/unrolled/render"
)

type Publication struct {
	engine *engine.Engine
	render *render.Render
}

type PublicationListVars struct {
	Query engine.Query
	Hits  *engine.PublicationHits
}

func NewPublication(e *engine.Engine, r *render.Render) *Publication {
	return &Publication{engine: e, render: r}
}

func (c *Publication) List(w http.ResponseWriter, r *http.Request) {
	q := engine.Query{}
	if err := form.NewDecoder().Decode(&q, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hits, err := c.engine.UserPublications("F72763F2-F0ED-11E1-A9DE-61C894A0A6B4", q)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.render.HTML(w, http.StatusOK, "publication/list", PublicationListVars{Query: q, Hits: hits})
}
