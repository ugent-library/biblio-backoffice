package controllers

import (
	"log"
	"net/http"
	"net/url"

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

	searchURL, _ := c.Router.Get("datasets").URLPath()

	c.Render.HTML(w, http.StatusOK, "dataset/list", views.NewData(c.Render, r, struct {
		PageTitle        string
		SearchURL        *url.URL
		SearchArgs       *engine.SearchArgs
		Hits             *models.DatasetHits
		PublicationSorts []string
	}{
		"Overview - Datasets - Biblio",
		searchURL,
		args,
		hits,
		c.Engine.Vocabularies()["publication_sorts"],
	}))
}

func (c *Datasets) Show(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	datasetPubs, err := c.Engine.GetDatasetPublications(dataset.ID)
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

	dataset.RelatedPublicationCount = len(datasetPubs)

	c.Render.HTML(w, http.StatusOK, "dataset/show", views.NewData(c.Render, r, struct {
		PageTitle           string
		Dataset             *models.Dataset
		DatasetPublications []*models.Publication
		Show                *views.ShowBuilder
		SearchArgs          *engine.SearchArgs
	}{
		"Dataset - Biblio",
		dataset,
		datasetPubs,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
		searchArgs,
	}))
}

func (c *Datasets) Add(w http.ResponseWriter, r *http.Request) {
	c.Render.HTML(w, http.StatusOK, "dataset/add", views.NewData(c.Render, r, struct {
		PageTitle string
		Step      int
	}{
		"Add - Datasets - Biblio",
		1,
	}))
}

func (c *Datasets) AddImport(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	source := r.FormValue("source")
	identifier := r.FormValue("identifier")

	dataset, err := c.Engine.ImportUserDatasetByIdentifier(context.GetUser(r.Context()).ID, source, identifier)
	if err != nil {
		log.Println(err)
		c.Render.HTML(w, http.StatusOK, "dataset/add", views.NewData(c.Render, r, struct {
			Step int
		}{
			1,
		},
			views.Flash{Type: "error", Message: "Sorry, something went wrong. Could not import the dataset."},
		))
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/add_description", views.NewData(c.Render, r, struct {
		PageTitle           string
		Step                int
		Dataset             *models.Dataset
		DatasetPublications []*models.Publication
		Show                *views.ShowBuilder
	}{
		"Add - Datasets - Biblio",
		2,
		dataset,
		nil,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
	}))
}

func (c *Datasets) AddDescription(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	c.Render.HTML(w, http.StatusOK, "dataset/add_description", views.NewData(c.Render, r, struct {
		PageTitle           string
		Step                int
		Dataset             *models.Dataset
		DatasetPublications []*models.Publication
		Show                *views.ShowBuilder
	}{
		"Add - Datasets - Biblio",
		2,
		dataset,
		nil,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
	}))
}

func (c *Datasets) AddConfirm(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	c.Render.HTML(w, http.StatusOK, "dataset/add_confirm", views.NewData(c.Render, r, struct {
		PageTitle string
		Step      int
		Dataset   *models.Dataset
	}{
		"Add - Datasets - Biblio",
		3,
		dataset,
	}))
}

func (c *Datasets) AddPublish(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	savedDataset, err := c.Engine.PublishDataset(dataset)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/add_publish", views.NewData(c.Render, r, struct {
		PageTitle string
		Step      int
		Dataset   *models.Dataset
	}{
		"Add - Datasets - Biblio",
		4,
		savedDataset,
	}))
}
