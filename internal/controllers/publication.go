package controllers

import (
	"log"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/unrolled/render"
)

type Publication struct {
	engine *engine.Engine
	render *render.Render
}

func NewPublication(e *engine.Engine, r *render.Render) *Publication {
	return &Publication{engine: e, render: r}
}

func (c *Publication) List(w http.ResponseWriter, r *http.Request) {
	hits, err := c.engine.UserPublications("F72763F2-F0ED-11E1-A9DE-61C894A0A6B4")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vars := struct {
		Hits *engine.PublicationHits
	}{
		Hits: hits,
	}

	c.render.HTML(w, http.StatusOK, "publication/list", vars)
}
