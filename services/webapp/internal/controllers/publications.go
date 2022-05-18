package controllers

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/tasks"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-orcid/orcid"
	"github.com/unrolled/render"
	"golang.org/x/text/language"
)

type PublicationAddSingleVars struct {
	PageTitle            string
	Step                 int
	Source               string
	Identifier           string
	DuplicatePublication *models.Publication
}

type Publications struct {
	Base
}

func NewPublications(c Base) *Publications {
	return &Publications{c}
}

func (c *Publications) List(w http.ResponseWriter, r *http.Request) {
	searchArgs := models.NewSearchArgs()
	if err := DecodeForm(searchArgs, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hits, err := c.userPublications(context.GetUser(r.Context()).ID, searchArgs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("publications").URLPath()

	c.Render.HTML(w, http.StatusOK, "publication/list", c.ViewData(r, struct {
		PageTitle  string
		SearchURL  *url.URL
		SearchArgs *models.SearchArgs
		Hits       *models.PublicationHits
	}{
		"Overview - Publications - Biblio",
		searchURL,
		searchArgs,
		hits,
	}),
	)
}

func (c *Publications) Show(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	datasets, err := c.Services.Store.GetPublicationDatasets(pub)
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

	c.Render.HTML(w, http.StatusOK, "publication/show", c.ViewData(r, struct {
		PageTitle           string
		Publication         *models.Publication
		PublicationDatasets []*models.Dataset
		Show                *views.ShowBuilder
		SearchArgs          *models.SearchArgs
		ErrorsTitle         string
		Errors              validation.Errors
	}{
		"Publication - Biblio",
		pub,
		datasets,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		searchArgs,
		"",
		nil,
	}),
	)
}

func (c *Publications) Thumbnail(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	if pub.ThumbnailURL() == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// TODO implement
}

func (c *Publications) Summary(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	c.Render.HTML(w, http.StatusOK, "publication/_summary", c.ViewData(r, struct {
		Publication *models.Publication
	}{
		pub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// Step 1: Start: Present selection form with creation modi (WOS, identifier, manual, Bibtex)

func (c *Publications) Add(w http.ResponseWriter, r *http.Request) {
	c.Render.HTML(w, http.StatusOK, "publication/add", c.ViewData(r, PublicationAddSingleVars{
		PageTitle: "Add - Publications - Biblio",
		Step:      1,
	}))
}

// Step 2: Add publication(s): Present WOS, Identifier, Bibtex or Manual input form based on choice from Step 1

func (c *Publications) AddSelectMethod(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	method := r.FormValue("add-method-select")

	template := ""

	switch method {
	case "wos":
		template = "publication/add_wos"
	case "identifier":
		template = "publication/add_identifier"
	case "manual":
		template = "publication/add_manual"
	case "bibtex":
		template = "publication/add_bibtex"
	default:
		flash := views.Flash{Type: "error"}
		flash.Message = "You didn't specify how you would like to add a publication."

		c.Render.HTML(w, http.StatusOK, "publication/add", c.ViewData(r, PublicationAddSingleVars{
			PageTitle: "Add - Publications - Biblio",
			Step:      1,
		},
			flash,
		))
		return
	}

	c.Render.HTML(w, http.StatusOK, template, c.ViewData(r, PublicationAddSingleVars{
		PageTitle: "Add - Publications - Biblio",
		Step:      2,
	}))
}

// Step 3: Complete description: Via DOI identifier
//   * Process creation of publication via identifier (DOI,...)
//   * Display detail edit form for single publication

func (c *Publications) AddSingleImportConfirm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	source := r.FormValue("source")
	identifier := r.FormValue("identifier")

	// check for duplicates
	if source == "crossref" && identifier != "" {
		args := models.NewSearchArgs().WithFilter("doi", identifier).WithFilter("status", "public")
		if existing, _ := c.Services.PublicationSearchService.SearchPublications(args); existing.Total > 0 {
			c.Render.HTML(w, http.StatusOK, "publication/add_identifier", c.ViewData(r, PublicationAddSingleVars{
				PageTitle:            "Add - Publications - Biblio",
				Step:                 1,
				Source:               source,
				Identifier:           identifier,
				DuplicatePublication: existing.Hits[0],
			}))
			return
		}
	}

	c.AddSingleImport(w, r)
}

func (c *Publications) AddSingleImport(w http.ResponseWriter, r *http.Request) {
	userID := context.GetUser(r.Context()).ID

	r.ParseForm()

	var pub *models.Publication
	loc := locale.Get(r.Context())

	if identifier := r.FormValue("identifier"); identifier != "" {
		var source string = r.FormValue("source")

		p, err := c.importUserPublicationByIdentifier(userID, source, identifier)

		if err != nil {
			log.Println(err)
			flash := views.Flash{Type: "error"}
			flash.Message = loc.T("publication.single_import", "import_by_id.import_failed")

			c.Render.HTML(w, http.StatusOK, "publication/add_identifier", c.ViewData(r, PublicationAddSingleVars{
				PageTitle: "Add - Publications - Biblio",
				Step:      2,
			},
				flash,
			))
			return
		}

		pub = p
	} else {
		pubType := r.FormValue("publication_type")

		if pubType == "" {
			flash := views.Flash{Type: "error"}
			flash.Message = loc.T("publication.single_import", "import_by_id.string.minLength")

			c.Render.HTML(w, http.StatusOK, "publication/add_identifier", c.ViewData(r, PublicationAddSingleVars{
				PageTitle: "Add - Publications - Biblio",
				Step:      2,
			},
				flash,
			))
			return
		}

		p := &models.Publication{
			ID:             uuid.NewString(),
			Type:           pubType,
			Status:         "private",
			Classification: "U",
			CreatorID:      userID,
			UserID:         userID,
		}
		if err := c.Services.Store.UpdatePublication(p); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pub = p
	}

	datasets, err := c.Services.Store.GetPublicationDatasets(pub)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/add_single_description", c.ViewData(r, struct {
		PageTitle           string
		Step                int
		Publication         *models.Publication
		PublicationDatasets []*models.Dataset
		Show                *views.ShowBuilder
	}{
		"Add - Publications - Biblio",
		3,
		pub,
		datasets,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
	}))
}

// Step 3: Complete description: (DOI, Manual flow) Display detail page of the record for editing

func (c *Publications) AddSingleDescription(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	datasets, err := c.Services.Store.GetPublicationDatasets(pub)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/add_single_description", c.ViewData(r, struct {
		PageTitle           string
		Step                int
		Publication         *models.Publication
		PublicationDatasets []*models.Dataset
		Show                *views.ShowBuilder
	}{
		"Add - Publications - Biblio",
		3,
		pub,
		datasets,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
	}))
}

// Step 3: Complete Description: overview of records to be added (single record, DOI / Manual flow)

func (c *Publications) AddSingleConfirm(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	c.Render.HTML(w, http.StatusOK, "publication/add_single_confirm", c.ViewData(r, struct {
		PageTitle   string
		Step        int
		Publication *models.Publication
	}{
		"Add - Publications - Biblio",
		3,
		pub,
	}))
}

// Step 4: Finish publishing to biblio: (Single record, DOI / Manual flow)

func (c *Publications) AddSinglePublish(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	savedPub := pub.Clone()
	savedPub.Status = "public"
	err := c.Services.Store.UpdatePublication(savedPub)
	if err != nil {

		/*
			TODO: return to /add-single/confirm with flash in session instead of rendering this in the wrong path
			TODO: replace hardcoded error by validation errors report
		*/
		if _, ok := err.(validation.Errors); ok {

			c.Render.HTML(w, http.StatusOK, "publication/add_single_confirm", c.ViewData(r, struct {
				PageTitle   string
				Step        int
				Publication *models.Publication
			}{
				"Add - Publications - Biblio",
				4,
				pub,
			},
				views.Flash{Type: "error", Message: "Title is required"},
			))
			return
		}

		// unexpected error
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/add_single_finish", c.ViewData(r, struct {
		PageTitle   string
		Step        int
		Publication *models.Publication
	}{
		"Add - Publications - Biblio",
		5,
		savedPub,
	}))
}

// Step 3: Complete description (WOS, Bibtex import file)
//  * Process upload from a WOS / BibTex import file
//  * Display a list of imported records

func (c *Publications) AddMultipleImport(w http.ResponseWriter, r *http.Request) {
	// 2GB limit on request body
	r.Body = http.MaxBytesReader(w, r.Body, 2000000000)

	// buffer limit of 32MB
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	source := r.FormValue("source")

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	userID := context.GetUser(r.Context()).ID

	batchID, err := c.importUserPublications(userID, source, file)
	if err != nil {
		log.Println(err)
		c.Render.HTML(w, http.StatusOK, "publication/add", c.ViewData(r, PublicationAddSingleVars{
			PageTitle: "Add - Publications - Biblio",
			Step:      1,
		},
			views.Flash{Type: "error", Message: "Sorry, something went wrong. Could not import the publications."},
		))
		return
	}

	args := models.NewSearchArgs()

	hits, err := c.userPublications(userID, args.Clone().WithFilter("batch_id", batchID))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("publication_add_multiple_description").URLPath("batch_id", batchID)

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_description", c.ViewData(r, struct {
		PageTitle  string
		Step       int
		SearchURL  *url.URL
		SearchArgs *models.SearchArgs
		Hits       *models.PublicationHits
		BatchID    string
	}{
		"Add - Publications - Biblio",
		3,
		searchURL,
		args,
		hits,
		batchID,
	}),
	)
}

// Step 3: Complete description (WOS, Bibtex import file)
//  * Just display the overview of records which were imported for the current batch id
//    This is used when returning to this step via the sidebar navigation

func (c *Publications) AddMultipleDescription(w http.ResponseWriter, r *http.Request) {
	userID := context.GetUser(r.Context()).ID
	batchID := mux.Vars(r)["batch_id"]

	args := models.NewSearchArgs()
	if err := DecodeForm(args, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hits, err := c.userPublications(userID, args.Clone().WithFilter("batch_id", batchID))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("publication_add_multiple_description").URLPath("batch_id", batchID)

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_description", c.ViewData(r, struct {
		PageTitle  string
		Step       int
		SearchURL  *url.URL
		SearchArgs *models.SearchArgs
		Hits       *models.PublicationHits
		BatchID    string
	}{
		"Add - Publications - Biblio",
		3,
		searchURL,
		args,
		hits,
		batchID,
	}),
	)
}

// Step 3: Complete description (WOS, Bibtex import file)
//  * Show the detail / edit page for a single record from a WOS / BibTex batch
//    Used for providing the BatchID & allowing returning back to the "Add publication flow"

func (c *Publications) AddMultipleShow(w http.ResponseWriter, r *http.Request) {
	batchID := mux.Vars(r)["batch_id"]
	pub := context.GetPublication(r.Context())

	datasets, err := c.Services.Store.GetPublicationDatasets(pub)
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

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_show", c.ViewData(r, struct {
		PageTitle           string
		Step                int
		Publication         *models.Publication
		PublicationDatasets []*models.Dataset
		Show                *views.ShowBuilder
		SearchArgs          *models.SearchArgs
		BatchID             string
	}{
		"Add - Publications - Biblio",
		3,
		pub,
		datasets,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		searchArgs,
		batchID,
	}),
	)
}

// Step 4: Publish to Biblio (Multiple record, WOS / BibTex)
// * Present a list of records ready to be published to Biblio / BibTex

func (c *Publications) AddMultipleConfirm(w http.ResponseWriter, r *http.Request) {
	userID := context.GetUser(r.Context()).ID
	batchID := mux.Vars(r)["batch_id"]

	args := models.NewSearchArgs()

	hits, err := c.userPublications(userID, args.Clone().WithFilter("batch_id", batchID))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("publication_add_multiple_confirm").URLPath("batch_id", batchID)

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_confirm", c.ViewData(r, struct {
		PageTitle  string
		Step       int
		SearchURL  *url.URL
		SearchArgs *models.SearchArgs
		Hits       *models.PublicationHits
		BatchID    string
	}{
		"Add - Publications - Biblio",
		4,
		searchURL,
		args,
		hits,
		batchID,
	}),
	)
}

// Step 4: Publish to Biblio (Multiple record, WOS / BibTex)
//  * Show the detail / edit page for a single record from a WOS / BibTex batch
//    Used for providing the BatchID & allowing returning back to the "Add publication flow"

func (c *Publications) AddMultipleConfirmShow(w http.ResponseWriter, r *http.Request) {
	batchID := mux.Vars(r)["batch_id"]
	pub := context.GetPublication(r.Context())

	datasets, err := c.Services.Store.GetPublicationDatasets(pub)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_confirm_show", c.ViewData(r, struct {
		PageTitle           string
		Step                int
		Publication         *models.Publication
		PublicationDatasets []*models.Dataset
		Show                *views.ShowBuilder
		BatchID             string
	}{
		"Add - Publications - Biblio",
		4,
		pub,
		datasets,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		batchID,
	}),
	)
}

// Step 4: Publish to Biblio (Multiple record, WOS / BibTex)
//   * Process the records & publish them to Biblio
//   * Present a "Congratulations" landing page.

func (c *Publications) AddMultiplePublish(w http.ResponseWriter, r *http.Request) {
	userID := context.GetUser(r.Context()).ID
	batchID := mux.Vars(r)["batch_id"]

	batchFilter := models.NewSearchArgs().WithFilter("batch_id", batchID)

	if err := c.batchPublishPublications(userID, batchFilter); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hits, err := c.userPublications(userID, batchFilter)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("publication_add_multiple_publish").URLPath("batch_id", batchID)

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_finish", c.ViewData(r, struct {
		PageTitle  string
		Step       int
		SearchURL  *url.URL
		SearchArgs *models.SearchArgs
		Hits       *models.PublicationHits
		BatchID    string
	}{
		"Add - Publications - Biblio",
		5,
		searchURL,
		models.NewSearchArgs(),
		hits,
		batchID,
	}),
	)
}

func (c *Publications) Publish(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	savedPub := pub.Clone()
	err := c.Services.Store.UpdatePublication(savedPub)

	flashes := make([]views.Flash, 0)
	var publicationErrors validation.Errors
	var publicationErrorsTitle string

	if err != nil {

		savedPub = pub

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

	pubDatasets, err := c.Services.Store.GetPublicationDatasets(pub)
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

	c.Render.HTML(w, http.StatusOK, "publication/show", c.ViewData(r, struct {
		PageTitle           string
		Publication         *models.Publication
		PublicationDatasets []*models.Dataset
		Show                *views.ShowBuilder
		SearchArgs          *models.SearchArgs
		ErrorsTitle         string
		Errors              validation.Errors
	}{
		"Publication - Biblio",
		savedPub,
		pubDatasets,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		searchArgs,
		publicationErrorsTitle,
		publicationErrors,
	},
		flashes...,
	))
}

func (c *Publications) ConfirmDelete(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	searchArgs := models.NewSearchArgs()
	if err := DecodeForm(searchArgs, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/_confirm_delete", c.ViewData(r, struct {
		Publication *models.Publication
		SearchArgs  *models.SearchArgs
	}{
		pub,
		searchArgs,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *Publications) Delete(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	r.ParseForm()
	searchArgs := models.NewSearchArgs()
	if err := DecodeForm(searchArgs, r.Form); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pub.Status = "deleted"
	if err := c.Services.Store.UpdatePublication(pub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hits, err := c.userPublications(context.GetUser(r.Context()).ID, searchArgs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("publications").URLPath()

	c.Render.HTML(w, http.StatusOK, "publication/list", c.ViewData(r, struct {
		PageTitle  string
		SearchURL  *url.URL
		SearchArgs *models.SearchArgs
		Hits       *models.PublicationHits
	}{
		"Overview - Publications - Biblio",
		searchURL,
		searchArgs,
		hits,
	},
		views.Flash{Type: "success", Message: "Successfully deleted publication.", DismissAfter: 5 * time.Second},
	),
	)
}

func (c *Publications) ORCIDAdd(w http.ResponseWriter, r *http.Request) {
	user := context.GetUser(r.Context())
	pub := context.GetPublication(r.Context())

	var (
		flash views.Flash
	)

	pub, err := c.addPublicationToORCID(user.ORCID, user.ORCIDToken, pub)
	if err != nil {
		if err == orcid.ErrDuplicate {
			flash = views.Flash{Type: "info", Message: "This publication is already part of your ORCID works."}
		} else {
			flash = views.Flash{Type: "error", Message: "Couldn't add this publication to your ORCID works."}
		}
	} else {
		flash = views.Flash{Type: "success", Message: "Successfully added the publication to your ORCID works.", DismissAfter: 5 * time.Second}
	}

	c.Render.HTML(w, http.StatusOK, "publication/_orcid_status", c.ViewData(r, struct {
		Publication *models.Publication
	}{
		pub,
	},
		flash,
	),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *Publications) ORCIDAddAll(w http.ResponseWriter, r *http.Request) {
	userID := context.GetUser(r.Context()).ID
	// TODO handle error
	id, err := c.addPublicationsToORCID(
		userID,
		models.NewSearchArgs().WithFilter("status", "public").WithFilter("author.id", userID),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "task/_status", c.ViewData(r, struct {
		ID      string
		Status  tasks.Status
		Message string
	}{
		id,
		tasks.Status{},
		"",
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *Publications) userPublications(userID string, args *models.SearchArgs) (*models.PublicationHits, error) {
	if args.FilterInRange("status", "private", "public") {
		args = args.Clone()
	} else {
		args = args.Clone().WithFilter("status", "private", "public")
	}
	switch args.FilterFor("scope") {
	case "created":
		args.WithFilter("creator_id", userID)
	case "contributed":
		args.WithFilter("author.id", userID)
	default:
		args.WithFilter("creator_id|author.id", userID)
	}
	delete(args.Filters, "scope")
	return c.Services.PublicationSearchService.SearchPublications(args)
}

// TODO should be async task
func (c *Publications) importUserPublications(userID, source string, file io.Reader) (string, error) {
	batchID := uuid.New().String()
	decFactory, ok := c.Services.PublicationDecoders[source]
	if !ok {
		return "", errors.New("unknown publication source")
	}
	dec := decFactory(file)

	var indexWG sync.WaitGroup

	// indexing channel
	indexC := make(chan *models.Publication)

	// start bulk indexer
	go func() {
		indexWG.Add(1)
		defer indexWG.Done()
		c.Services.PublicationSearchService.IndexPublications(indexC)
	}()

	var importErr error
	for {
		p := models.Publication{
			ID:             uuid.NewString(),
			BatchID:        batchID,
			Status:         "private",
			Classification: "U",
			CreatorID:      userID,
			UserID:         userID,
		}
		if err := dec.Decode(&p); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			importErr = err
			break
		}
		if err := c.Services.Store.UpdatePublication(&p); err != nil {
			importErr = err
			break
		}

		indexC <- &p
	}

	// close indexing channel when all recs are stored
	close(indexC)
	// wait for indexing to finish
	indexWG.Wait()

	// TODO rollback if error
	if importErr != nil {
		return "", importErr
	}

	return batchID, nil
}

// TODO should be async task
func (c *Publications) importUserPublicationByIdentifier(userID, source, identifier string) (*models.Publication, error) {
	s, ok := c.Services.PublicationSources[source]
	if !ok {
		return nil, errors.New("unknown dataset source")
	}
	p, err := s.GetPublication(identifier)
	if err != nil {
		return nil, err
	}

	p.ID = uuid.NewString()
	p.CreatorID = userID
	p.UserID = userID
	p.Status = "private"
	p.Classification = "U"

	if err := c.Services.Store.UpdatePublication(p); err != nil {
		return nil, err
	}

	return p, nil
}

// TODO should be async task
func (c *Publications) batchPublishPublications(userID string, args *models.SearchArgs) (err error) {
	var hits *models.PublicationHits
	for {
		hits, err = c.userPublications(userID, args)
		for _, pub := range hits.Hits {
			pub.Status = "public"
			if err = c.Services.Store.UpdatePublication(pub); err != nil {
				break
			}
		}
		if !hits.NextPage() {
			break
		}
		args.Page = args.Page + 1
	}
	return
}

func (c *Publications) addPublicationsToORCID(userID string, s *models.SearchArgs) (string, error) {
	user, err := c.Services.GetUser(userID)
	if err != nil {
		return "", err
	}

	taskID := "orcid:" + uuid.NewString()

	c.Services.Tasks.Add(taskID, func(t tasks.Task) error {
		return c.sendPublicationsToORCIDTask(t, userID, user.ORCID, user.ORCIDToken, s)
	})

	return taskID, nil
}

// TODO make workflow
func (c *Publications) addPublicationToORCID(orcidID, orcidToken string, p *models.Publication) (*models.Publication, error) {
	client := orcid.NewMemberClient(orcid.Config{
		Token:   orcidToken,
		Sandbox: c.Services.ORCIDSandbox,
	})

	work := publicationToORCID(p)
	putCode, res, err := client.AddWork(orcidID, work)
	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		log.Printf("orcid error: %s", body)
		return p, err
	}

	p.ORCIDWork = append(p.ORCIDWork, models.PublicationORCIDWork{
		ORCID:   orcidID,
		PutCode: putCode,
	})

	if err := c.Services.Store.UpdatePublication(p); err != nil {
		return nil, err
	}

	return p, nil
}

// TODO move to workflows
func (c *Publications) sendPublicationsToORCIDTask(t tasks.Task, userID, orcidID, orcidToken string, searchArgs *models.SearchArgs) error {
	orcidClient := orcid.NewMemberClient(orcid.Config{
		Token:   orcidToken,
		Sandbox: c.Services.ORCIDSandbox,
	})

	var numDone int

	for {
		hits, _ := c.Services.PublicationSearchService.SearchPublications(searchArgs)

		for _, pub := range hits.Hits {
			numDone++

			var done bool
			for _, ow := range pub.ORCIDWork {
				if ow.ORCID == orcidID { // already sent to orcid
					done = true
					break
				}
			}
			if done {
				continue
			}

			work := publicationToORCID(pub)
			putCode, res, err := orcidClient.AddWork(orcidID, work)
			if res.StatusCode == 409 { // duplicate
				continue
			} else if err != nil {
				body, _ := ioutil.ReadAll(res.Body)
				log.Printf("orcid error: %s", body)
				return err
			}

			pub.ORCIDWork = append(pub.ORCIDWork, models.PublicationORCIDWork{
				ORCID:   orcidID,
				PutCode: putCode,
			})

			if err := c.Services.Store.UpdatePublication(pub); err != nil {
				return err
			}
		}

		t.Progress(numDone, hits.Total)

		if !hits.NextPage() {
			break
		}
		searchArgs.Page = searchArgs.Page + 1
	}

	return nil
}

func publicationToORCID(p *models.Publication) *orcid.Work {
	w := &orcid.Work{
		URL:     orcid.String(fmt.Sprintf("https://biblio.ugent.be/publication/%s", p.ID)),
		Country: orcid.String("BE"),
		ExternalIDs: &orcid.ExternalIDs{
			ExternalID: []orcid.ExternalID{{
				Type:         "handle",
				Relationship: "SELF",
				Value:        fmt.Sprintf("http://hdl.handle.net/1854/LU-%s", p.ID),
			}},
		},
		Title: &orcid.Title{
			Title: orcid.String(p.Title),
		},
		PublicationDate: &orcid.PublicationDate{
			Year: orcid.String(p.Year),
		},
	}

	for _, role := range []string{"author", "editor"} {
		for _, c := range p.Contributors(role) {
			wc := orcid.Contributor{
				CreditName: orcid.String(strings.Join([]string{c.FirstName, c.LastName}, " ")),
				Attributes: &orcid.ContributorAttributes{
					Role: strings.ToUpper(role),
				},
			}
			if c.ORCID != "" {
				wc.ORCID = &orcid.URI{Path: c.ORCID}
			}
			if w.Contributors == nil {
				w.Contributors = &orcid.Contributors{}
			}
			w.Contributors.Contributor = append(w.Contributors.Contributor, wc)
		}
	}

	switch p.Type {
	case "journal_article":
		w.Type = "JOURNAL_ARTICLE"
	case "book":
		w.Type = "BOOK"
	case "book_chapter":
		w.Type = "BOOK_CHAPTER"
	case "book_editor":
		w.Type = "EDITED_BOOK"
	case "dissertation":
		w.Type = "DISSERTATION"
	case "conference":
		switch p.ConferenceType {
		case "meetingAbstract":
			w.Type = "CONFERENCE_ABSTRACT"
		case "poster":
			w.Type = "CONFERENCE_POSTER"
		default:
			w.Type = "CONFERENCE_PAPER"
		}
	case "miscellaneous":
		switch p.MiscellaneousType {
		case "bookReview":
			w.Type = "BOOK_REVIEW"
		case "report":
			w.Type = "REPORT"
		default:
			w.Type = "OTHER"
		}
	default:
		w.Type = "OTHER"
	}

	if len(p.AlternativeTitle) > 0 {
		w.Title.SubTitle = orcid.String(p.AlternativeTitle[0])
	}

	if len(p.Abstract) > 0 {
		w.ShortDescription = p.Abstract[0].Text
	}

	if p.DOI != "" {
		w.ExternalIDs.ExternalID = append(w.ExternalIDs.ExternalID, orcid.ExternalID{
			Type:         "doi",
			Relationship: "SELF",
			Value:        p.DOI,
		})
	}

	if len(p.Language) > 0 {
		if tag, err := language.Parse(p.Language[0]); err == nil {
			w.LanguageCode = tag.String()
		}
	}

	return w
}
