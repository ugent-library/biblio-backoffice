package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/go-web/forms"
	"github.com/unrolled/render"
)

type Publication struct {
	engine *engine.Engine
	render *render.Render
}

type PublicationListVars struct {
	Query *engine.Query
	Hits  *engine.PublicationHits
}

type PublicationShowVars struct {
	Pub *engine.Publication
}

type PublicationNewVars struct {
}

func NewPublication(e *engine.Engine, r *render.Render) *Publication {
	return &Publication{engine: e, render: r}
}

func (c *Publication) List(w http.ResponseWriter, r *http.Request) {
	q := engine.NewQuery()
	if err := forms.Decode(q, r.URL.Query()); err != nil {
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

func (c *Publication) Show(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.render.HTML(w, http.StatusOK, "publication/show", PublicationShowVars{Pub: pub})
}

func (c *Publication) New(w http.ResponseWriter, r *http.Request) {
	c.render.HTML(w, http.StatusOK, "publication/new", PublicationNewVars{})
}
