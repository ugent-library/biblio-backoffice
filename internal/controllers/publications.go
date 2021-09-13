package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/ctx"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/presenters"
	"github.com/ugent-library/go-web/forms"
	"github.com/unrolled/render"
)

type Publications struct {
	engine *engine.Engine
	render *render.Render
}

type PublicationListVars struct {
	SearchArgs *engine.SearchArgs
	Hits       *engine.PublicationHits
}

type PublicationShowVars struct {
	Pub *presenters.Publication
}

type PublicationNewVars struct {
}

func NewPublication(e *engine.Engine, r *render.Render) *Publications {
	return &Publications{engine: e, render: r}
}

func (c *Publications) List(w http.ResponseWriter, r *http.Request) {
	args := engine.NewSearchArgs()
	if err := forms.Decode(args, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hits, err := c.engine.UserPublications(ctx.GetUser(r).ID, args)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.render.HTML(w, http.StatusOK, "publication/list", PublicationListVars{SearchArgs: args, Hits: hits})
}

func (c *Publications) Show(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h := &presenters.Publication{Publication: pub, Render: c.render}
	c.render.HTML(w, http.StatusOK, "publication/show", PublicationShowVars{Pub: h})
}

func (c *Publications) New(w http.ResponseWriter, r *http.Request) {
	c.render.HTML(w, http.StatusOK, "publication/new", PublicationNewVars{})
}
