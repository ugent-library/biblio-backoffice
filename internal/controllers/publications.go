package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-web/forms"
	"github.com/unrolled/render"
)

type Publications struct {
	engine   *engine.Engine
	render   *render.Render
	listView views.Renderer
	newView  views.Renderer
}

type PublicationListVars struct {
	SearchArgs *engine.SearchArgs
	Hits       *models.PublicationHits
}

type PublicationShowVars struct {
	Pub *models.Publication
}

type PublicationNewVars struct {
}

func NewPublication(e *engine.Engine, r *render.Render) *Publications {
	return &Publications{
		engine:   e,
		render:   r,
		listView: views.NewView(r, "publication/list"),
		newView:  views.NewView(r, "publication/new"),
	}
}

func (c *Publications) List(w http.ResponseWriter, r *http.Request) {
	args := engine.NewSearchArgs()
	if err := forms.Decode(args, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hits, err := c.engine.UserPublications(context.User(r.Context()).ID, args)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.listView.Render(w, r, http.StatusOK, PublicationListVars{SearchArgs: args, Hits: hits})
}

func (c *Publications) Show(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.render.HTML(w, http.StatusOK, "publication/show", PublicationShowVars{
		Pub: pub,
	})
}

func (c *Publications) New(w http.ResponseWriter, r *http.Request) {
	c.newView.Render(w, r, http.StatusOK, PublicationNewVars{})
}
