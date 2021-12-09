package controllers

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
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

func (c *Datasets) Publish(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	savedDataset, err := c.Engine.PublishDataset(dataset)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	savedDataset.RelatedPublicationCount = len(datasetPubs)

	c.Render.HTML(w, http.StatusOK, "dataset/show", views.NewData(c.Render, r, struct {
		PageTitle           string
		Dataset             *models.Dataset
		DatasetPublications []*models.Publication
		Show                *views.ShowBuilder
		SearchArgs          *engine.SearchArgs
	}{
		"Dataset - Biblio",
		savedDataset,
		datasetPubs,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
		searchArgs,
	},
		views.Flash{Type: "success", Message: "Successfully published to Biblio.", DismissAfter: 5 * time.Second},
	))
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
	loc := locale.Get(r.Context())

	dataset, err := c.Engine.ImportUserDatasetByIdentifier(context.GetUser(r.Context()).ID, source, identifier)
	if err != nil {
		log.Println(err)
		messages := make([]string, 0)
		/*
			TODO: when datacite does not provide the necessary information for a record
				  to validate, this list will be too long to show, and there is no
				  way to attach an error to a specific field as we do not have a full form.
		*/
		switch cErr := err.(type) {
		case jsonapi.Errors:
			for _, e := range cErr {
				messages = append(messages, loc.T("dataset.single_import", e.Code))
			}
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		flash := strings.Join(messages, "<br>")
		c.Render.HTML(w, http.StatusOK, "dataset/add", views.NewData(c.Render, r, struct {
			PageTitle string
			Step      int
		}{
			"Add - Datasets - Biblio",
			1,
		},
			views.Flash{Type: "error", Message: flash},
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
		views.NewShowBuilder(c.Render, loc),
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

	c.Render.HTML(w, http.StatusOK, "dataset/add_finish", views.NewData(c.Render, r, struct {
		PageTitle string
		Step      int
		Dataset   *models.Dataset
	}{
		"Add - Datasets - Biblio",
		4,
		savedDataset,
	}))
}

func (c *Datasets) ConfirmDelete(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	searchArgs := engine.NewSearchArgs()
	if err := forms.Decode(searchArgs, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/_confirm_delete", views.NewData(c.Render, r, struct {
		Dataset    *models.Dataset
		SearchArgs *engine.SearchArgs
	}{
		dataset,
		searchArgs,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *Datasets) Delete(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	r.ParseForm()
	searchArgs := engine.NewSearchArgs()
	if err := forms.Decode(searchArgs, r.Form); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := c.Engine.UpdateDataset(dataset); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hits, err := c.Engine.UserDatasets(context.GetUser(r.Context()).ID, searchArgs)
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
		searchArgs,
		hits,
		c.Engine.Vocabularies()["publication_sorts"],
	},
		views.Flash{Type: "success", Message: "Successfully deleted dataset.", DismissAfter: 5 * time.Second},
	),
	)
}
