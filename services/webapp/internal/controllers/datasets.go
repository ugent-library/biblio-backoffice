package controllers

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type DatasetListVars struct {
	SearchArgs       *models.SearchArgs
	Hits             *models.DatasetHits
	PublicationSorts []string
}

type DatasetAddVars struct {
	PageTitle        string
	Step             int
	Source           string
	Identifier       string
	DuplicateDataset *models.Dataset
}

type Datasets struct {
	Context
}

func NewDatasets(c Context) *Datasets {
	return &Datasets{c}
}

func (c *Datasets) List(w http.ResponseWriter, r *http.Request) {
	args := models.NewSearchArgs()
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

	c.Render.HTML(w, http.StatusOK, "dataset/list", c.ViewData(r, struct {
		PageTitle        string
		SearchURL        *url.URL
		SearchArgs       *models.SearchArgs
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

	datasetPubs, err := c.Engine.GetDatasetPublications(dataset)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchArgs := models.NewSearchArgs()
	if err := forms.Decode(searchArgs, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// dataset.RelatedPublicationCount = len(datasetPubs)

	c.Render.HTML(w, http.StatusOK, "dataset/show", c.ViewData(r, struct {
		PageTitle           string
		Dataset             *models.Dataset
		DatasetPublications []*models.Publication
		Show                *views.ShowBuilder
		SearchArgs          *models.SearchArgs
		ErrorsTitle         string
		Errors              jsonapi.Errors
	}{
		"Dataset - Biblio",
		dataset,
		datasetPubs,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		searchArgs,
		"",
		nil,
	}))
}

func (c *Datasets) Publish(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	flashes := make([]views.Flash, 0)
	var publicationErrors jsonapi.Errors
	var publicationErrorsTitle string

	datasetCopy := *dataset
	datasetCopy.Status = "public"
	savedDataset, err := c.Engine.UpdateDataset(&datasetCopy)
	if err != nil {
		savedDataset = dataset

		if e, ok := err.(models.ValidationErrors); ok {
			formErrors := jsonapi.Errors{jsonapi.Error{
				Detail: e.Error(),
				Title:  e.Error(),
			}}

			publicationErrors = formErrors
			publicationErrorsTitle = "Unable to publish record due to following errors"

		} else {

			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}

	} else {

		flashes = append(flashes, views.Flash{Type: "success", Message: "Successfully published to Biblio.", DismissAfter: 5 * time.Second})

	}

	datasetPubs, err := c.Engine.GetDatasetPublications(dataset)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchArgs := models.NewSearchArgs()
	if err := forms.Decode(searchArgs, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// savedDataset.RelatedPublicationCount = len(datasetPubs)

	c.Render.HTML(w, http.StatusOK, "dataset/show", c.ViewData(r, struct {
		PageTitle           string
		Dataset             *models.Dataset
		DatasetPublications []*models.Publication
		Show                *views.ShowBuilder
		SearchArgs          *models.SearchArgs
		ErrorsTitle         string
		Errors              jsonapi.Errors
	}{
		"Dataset - Biblio",
		savedDataset,
		datasetPubs,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		searchArgs,
		publicationErrorsTitle,
		publicationErrors,
	},
		flashes...,
	))
}

func (c *Datasets) Add(w http.ResponseWriter, r *http.Request) {
	c.Render.HTML(w, http.StatusOK, "dataset/add", c.ViewData(r, DatasetAddVars{
		PageTitle: "Add - Datasets - Biblio",
		Step:      1,
	}))
}

func (c *Datasets) AddImportConfirm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	source := r.FormValue("source")
	identifier := r.FormValue("identifier")

	// check for duplicates
	if source == "datacite" {
		if existing, _ := c.Engine.Datasets(models.NewSearchArgs().WithFilter("doi", identifier).WithFilter("status", "public")); existing.Total > 0 {
			c.Render.HTML(w, http.StatusOK, "dataset/add", c.ViewData(r, DatasetAddVars{
				PageTitle:        "Add - Datasets - Biblio",
				Step:             1,
				Source:           source,
				Identifier:       identifier,
				DuplicateDataset: existing.Hits[0],
			}))
			return
		}
	}

	c.AddImport(w, r)
}

func (c *Datasets) AddImport(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	source := r.FormValue("source")
	identifier := r.FormValue("identifier")
	loc := locale.Get(r.Context())

	dataset, err := c.Engine.ImportUserDatasetByIdentifier(context.GetUser(r.Context()).ID, source, identifier)

	if err != nil {
		flash := views.Flash{Type: "error"}

		if e, ok := err.(jsonapi.Errors); ok {
			flash.Message = loc.T("dataset.single_import", e[0].Code)
		} else {
			log.Println(e)
			flash.Message = loc.T("dataset.single_import", "import_by_id.import_failed")
		}

		c.Render.HTML(w, http.StatusOK, "dataset/add", c.ViewData(r, DatasetAddVars{
			PageTitle:  "Add - Datasets - Biblio",
			Step:       1,
			Source:     source,
			Identifier: identifier,
		},
			flash,
		))
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/add_description", c.ViewData(r, struct {
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
		views.NewShowBuilder(c.RenderPartial, loc),
	}))
}

func (c *Datasets) AddDescription(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	c.Render.HTML(w, http.StatusOK, "dataset/add_description", c.ViewData(r, struct {
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
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
	}))
}

func (c *Datasets) AddConfirm(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	c.Render.HTML(w, http.StatusOK, "dataset/add_confirm", c.ViewData(r, struct {
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

	dataset.Status = "public"
	savedDataset, err := c.Engine.UpdateDataset(dataset)
	if err != nil {

		/*
		   TODO: return to dataset - add_confirm with flash in session instead of rendering this in the wrong path
		   We only use one error, as publishing can only fail on attribute title
		*/
		if e, ok := err.(models.ValidationErrors); ok {
			c.Render.HTML(w, http.StatusOK, "dataset/add_confirm", c.ViewData(r, struct {
				PageTitle string
				Step      int
				Dataset   *models.Dataset
			}{
				"Add - Datasets - Biblio",
				3,
				dataset,
			},
				views.Flash{Type: "error", Message: e.Error()},
			))
			return
		}

		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/add_finish", c.ViewData(r, struct {
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

	searchArgs := models.NewSearchArgs()
	if err := forms.Decode(searchArgs, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/_confirm_delete", c.ViewData(r, struct {
		Dataset    *models.Dataset
		SearchArgs *models.SearchArgs
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
	searchArgs := models.NewSearchArgs()
	if err := forms.Decode(searchArgs, r.Form); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataset.Status = "deleted"
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

	c.Render.HTML(w, http.StatusOK, "dataset/list", c.ViewData(r, struct {
		PageTitle        string
		SearchURL        *url.URL
		SearchArgs       *models.SearchArgs
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
