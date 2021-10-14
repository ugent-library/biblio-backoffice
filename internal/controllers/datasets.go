package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-web/forms"
	"github.com/unrolled/render"
)

type Datasets struct {
	engine *engine.Engine
	render *render.Render
}

func NewDatasets(e *engine.Engine, r *render.Render) *Datasets {
	return &Datasets{engine: e, render: r}
}

func (c *Datasets) List(w http.ResponseWriter, r *http.Request) {
	args := engine.NewSearchArgs()
	if err := forms.Decode(args, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hits, err := c.engine.UserDatasets(context.User(r.Context()).ID, args)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: change PublicationListVars into DatasetListVars
	c.render.HTML(w, http.StatusOK, "dataset/list", PublicationListVars{
		Data:             views.NewData(r),
		SearchArgs:       args,
		Hits:             hits,
		PublicationSorts: c.engine.PublicationSorts(),
	})
}

func (c *Datasets) Show(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	pub, err := c.engine.GetPublication(id) // TODO constrain to research_data type
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.render.HTML(w, http.StatusOK, "dataset/show", views.NewDatasetData(r, c.render, pub))
}
