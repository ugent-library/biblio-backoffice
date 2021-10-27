package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
)

type DatasetListVars struct {
	SearchArgs       *engine.SearchArgs
	Hits             *models.DatasetHits
	PublicationSorts []string
}

type Datasets struct {
	Context
}

func NewDatasets(c Context) *Datasets {
	return &Datasets{c}
}

func (c *Datasets) List(w http.ResponseWriter, r *http.Request) {
	args := engine.NewSearchArgs()
	if err := forms.Decode(args, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hits, err := c.Engine.UserDatasets(context.GetUser(r.Context()).ID, args)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/list",
		views.NewData(c.Render, r, DatasetListVars{
			SearchArgs:       args,
			Hits:             hits,
			PublicationSorts: c.Engine.Vocabularies()["publication_sorts"],
		}),
	)
}

func (c *Datasets) Show(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	dataset, err := c.Engine.GetDataset(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/show",
		views.NewData(c.Render, r, struct {
			Dataset *models.Dataset
			Show    *views.ShowBuilder
		}{
			dataset,
			views.NewShowBuilder(c.Render, locale.Get(r.Context())),
		}),
	)
}
