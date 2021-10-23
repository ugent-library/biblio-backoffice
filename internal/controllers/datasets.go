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

type DatasetListVars struct {
	views.Data
	SearchArgs       *engine.SearchArgs
	Hits             *models.DatasetHits
	PublicationSorts []string
}

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

	hits, err := c.engine.UserDatasets(context.GetUser(r.Context()).ID, args)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.render.HTML(w, http.StatusOK, "dataset/list", DatasetListVars{
		Data:             views.NewData(c.render, r),
		SearchArgs:       args,
		Hits:             hits,
		PublicationSorts: c.engine.Vocabularies()["publication_sorts"],
	})
}

func (c *Datasets) Show(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	dataset, err := c.engine.GetDataset(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.render.HTML(w, http.StatusOK, "dataset/show", views.NewDatasetData(r, c.render, dataset))
}
