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
	engine *engine.Engine
	render *render.Render
}

type PublicationListVars struct {
	views.Data
	SearchArgs       *engine.SearchArgs
	Hits             *models.PublicationHits
	PublicationSorts []string
}

func NewPublications(e *engine.Engine, r *render.Render) *Publications {
	return &Publications{engine: e, render: r}
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

	c.render.HTML(w, http.StatusOK, "publication/list", PublicationListVars{
		Data:             views.NewData(r),
		SearchArgs:       args,
		Hits:             hits,
		PublicationSorts: c.engine.PublicationSorts(),
	})
}

func (c *Publications) Show(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.render.HTML(w, http.StatusOK, "publication/show", views.NewPublicationData(r, c.render, pub))
}

func (c *Publications) New(w http.ResponseWriter, r *http.Request) {
	c.render.HTML(w, http.StatusOK, "publication/new", views.NewData(r))
}
