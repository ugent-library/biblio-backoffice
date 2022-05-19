package controllers

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/unrolled/render"
)

type DatasetListVars struct {
	SearchArgs *models.SearchArgs
	Hits       *models.DatasetHits
}

type DatasetAddVars struct {
	PageTitle        string
	Step             int
	Source           string
	Identifier       string
	DuplicateDataset *models.Dataset
}

type Datasets struct {
	Base
	store                backends.Store
	datasetSearchService backends.DatasetSearchService
	datasetSources       map[string]backends.DatasetGetter
}

func NewDatasets(base Base, store backends.Store, datasetSearchService backends.DatasetSearchService,
	datasetSources map[string]backends.DatasetGetter) *Datasets {
	return &Datasets{
		Base:                 base,
		store:                store,
		datasetSearchService: datasetSearchService,
		datasetSources:       datasetSources,
	}
}

func (c *Datasets) List(w http.ResponseWriter, r *http.Request) {
	args := models.NewSearchArgs()
	if err := DecodeForm(args, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hits, err := c.userDatasets(context.GetUser(r.Context()).ID, args)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("datasets").URLPath()

	c.Render.HTML(w, http.StatusOK, "dataset/list", c.ViewData(r, struct {
		PageTitle  string
		SearchURL  *url.URL
		SearchArgs *models.SearchArgs
		Hits       *models.DatasetHits
	}{
		"Overview - Datasets - Biblio",
		searchURL,
		args,
		hits,
	}))
}

func (c *Datasets) Show(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	datasetPubs, err := c.store.GetDatasetPublications(dataset)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchArgs := models.NewSearchArgs()
	if err := DecodeForm(searchArgs, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/show", c.ViewData(r, struct {
		PageTitle           string
		Dataset             *models.Dataset
		DatasetPublications []*models.Publication
		Show                *views.ShowBuilder
		SearchArgs          *models.SearchArgs
		ErrorsTitle         string
		Errors              validation.Errors
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
	var publicationErrors validation.Errors
	var publicationErrorsTitle string

	datasetCopy := *dataset
	datasetCopy.Status = "public"
	savedDataset := datasetCopy.Clone()
	err := c.store.UpdateDataset(savedDataset)
	if err != nil {
		savedDataset = dataset

		if e, ok := err.(validation.Errors); ok {
			publicationErrors = e
			publicationErrorsTitle = "Unable to publish record due to following errors"

		} else {

			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}

	} else {

		flashes = append(flashes, views.Flash{Type: "success", Message: "Successfully published to Biblio.", DismissAfter: 5 * time.Second})

	}

	datasetPubs, err := c.store.GetDatasetPublications(dataset)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchArgs := models.NewSearchArgs()
	if err := DecodeForm(searchArgs, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/show", c.ViewData(r, struct {
		PageTitle           string
		Dataset             *models.Dataset
		DatasetPublications []*models.Publication
		Show                *views.ShowBuilder
		SearchArgs          *models.SearchArgs
		ErrorsTitle         string
		Errors              validation.Errors
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
		args := models.NewSearchArgs().WithFilter("doi", identifier).WithFilter("status", "public")
		if existing, _ := c.datasetSearchService.Search(args); existing.Total > 0 {
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

	dataset, err := c.importUserDatasetByIdentifier(context.GetUser(r.Context()).ID, source, identifier)

	if err != nil {
		log.Println(err)
		flash := views.Flash{Type: "error"}
		flash.Message = loc.T("dataset.single_import", "import_by_id.import_failed")

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
	savedDataset := dataset.Clone()
	err := c.store.UpdateDataset(dataset)
	if err != nil {

		/*
		   TODO: return to dataset - add_confirm with flash in session instead of rendering this in the wrong path
		   We only use one error, as publishing can only fail on attribute title
		*/
		if e, ok := err.(validation.Errors); ok {
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
	if err := DecodeForm(searchArgs, r.URL.Query()); err != nil {
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
	if err := DecodeForm(searchArgs, r.Form); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataset.Status = "deleted"
	if err := c.store.UpdateDataset(dataset); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hits, err := c.userDatasets(context.GetUser(r.Context()).ID, searchArgs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("datasets").URLPath()

	c.Render.HTML(w, http.StatusOK, "dataset/list", c.ViewData(r, struct {
		PageTitle  string
		SearchURL  *url.URL
		SearchArgs *models.SearchArgs
		Hits       *models.DatasetHits
	}{
		"Overview - Datasets - Biblio",
		searchURL,
		searchArgs,
		hits,
	},
		views.Flash{Type: "success", Message: "Successfully deleted dataset.", DismissAfter: 5 * time.Second},
	),
	)
}

func (c *Datasets) userDatasets(userID string, args *models.SearchArgs) (*models.DatasetHits, error) {
	searcher := c.datasetSearchService.WithScope("status", "private", "public")

	switch args.FilterFor("scope") {
	case "created":
		searcher = searcher.WithScope("creator_id", userID)
	case "contributed":
		searcher = searcher.WithScope("author.id", userID)
	default:
		searcher = searcher.WithScope("creator_id|author.id", userID)
	}
	delete(args.Filters, "scope")

	return searcher.Search(args)
}

func (c *Datasets) importUserDatasetByIdentifier(userID, source, identifier string) (*models.Dataset, error) {
	s, ok := c.datasetSources[source]
	if !ok {
		return nil, errors.New("unknown dataset source")
	}
	d, err := s.GetDataset(identifier)
	if err != nil {
		return nil, err
	}
	d.Vacuum()
	d.ID = uuid.NewString()
	d.CreatorID = userID
	d.UserID = userID
	d.Status = "private"

	if err = c.store.UpdateDataset(d); err != nil {
		return nil, err
	}

	return d, nil
}
