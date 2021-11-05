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

	c.Render.HTML(w, http.StatusOK, "dataset/list", views.NewData(c.Render, r, DatasetListVars{
		SearchArgs:       args,
		Hits:             hits,
		PublicationSorts: c.Engine.Vocabularies()["publication_sorts"],
	}))
}

func (c *Datasets) Show(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	dataset, err := c.Engine.GetDataset(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	datasetPubs, err := c.Engine.GetDatasetPublications(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchArgs := engine.NewSearchArgs()
	if err := forms.Decode(searchArgs, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/show", views.NewData(c.Render, r, struct {
		Dataset             *models.Dataset
		DatasetPublications []*models.Publication
		Show                *views.ShowBuilder
		SearchArgs          *engine.SearchArgs
	}{
		dataset,
		datasetPubs,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
		searchArgs,
	}))
}

func (c *Datasets) Add(w http.ResponseWriter, r *http.Request) {
	c.Render.HTML(w, http.StatusOK, "dataset/add", views.NewData(c.Render, r, struct {
		Step int
	}{
		1,
	}))
}

func (c *Datasets) AddImport(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	identifier := r.FormValue("identifier")

	datasets, err := c.Engine.ImportUserDatasets(context.GetUser(r.Context()).ID, identifier)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO flash messages
	if len(datasets) == 0 {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/add_description", views.NewData(c.Render, r, struct {
		Step    int
		Dataset *models.Dataset
		Show    *views.ShowBuilder
	}{
		2,
		datasets[0],
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
	}))
}

func (c *Datasets) AddDescription(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	dataset, err := c.Engine.GetDataset(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/add_description", views.NewData(c.Render, r, struct {
		Step    int
		Dataset *models.Dataset
		Show    *views.ShowBuilder
	}{
		2,
		dataset,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
	}))
}

func (c *Datasets) AddConfirm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	dataset, err := c.Engine.GetDataset(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/add_confirm", views.NewData(c.Render, r, struct {
		Step    int
		Dataset *models.Dataset
	}{
		3,
		dataset,
	}))
}

func (c *Datasets) AddPublish(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	dataset, err := c.Engine.GetDataset(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	savedDataset, err := c.Engine.PublishDataset(dataset)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/add_publish", views.NewData(c.Render, r, struct {
		Step    int
		Dataset *models.Dataset
	}{
		4,
		savedDataset,
	}))
}
