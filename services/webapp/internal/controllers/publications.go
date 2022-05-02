package controllers

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-orcid/orcid"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type PublicationAddSingleVars struct {
	PageTitle            string
	Step                 int
	Source               string
	Identifier           string
	DuplicatePublication *models.Publication
}

type Publications struct {
	Context
}

func NewPublications(c Context) *Publications {
	return &Publications{c}
}

func (c *Publications) List(w http.ResponseWriter, r *http.Request) {
	searchArgs := models.NewSearchArgs()
	if err := forms.Decode(searchArgs, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hits, err := c.Engine.UserPublications(context.GetUser(r.Context()).ID, searchArgs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("publications").URLPath()

	c.Render.HTML(w, http.StatusOK, "publication/list", c.ViewData(r, struct {
		PageTitle        string
		SearchURL        *url.URL
		SearchArgs       *models.SearchArgs
		Hits             *models.PublicationHits
		PublicationSorts []string
	}{
		"Overview - Publications - Biblio",
		searchURL,
		searchArgs,
		hits,
		c.Engine.Vocabularies()["publication_sorts"],
	}),
	)
}

func (c *Publications) Show(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	datasets, err := c.Engine.GetPublicationDatasets(pub.ID)
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

	pub.RelatedDatasetCount = len(datasets)

	c.Render.HTML(w, http.StatusOK, "publication/show", c.ViewData(r, struct {
		PageTitle           string
		Publication         *models.Publication
		PublicationDatasets []*models.Dataset
		Show                *views.ShowBuilder
		Vocabularies        map[string][]string
		SearchArgs          *models.SearchArgs
		ErrorsTitle         string
		Errors              jsonapi.Errors
	}{
		"Publication - Biblio",
		pub,
		datasets,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		c.Engine.Vocabularies(),
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

	c.Engine.ServePublicationFile(pub.ThumbnailURL(), w, r)
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

// Step 1: Start

func (c *Publications) Add(w http.ResponseWriter, r *http.Request) {
	c.Render.HTML(w, http.StatusOK, "publication/add", c.ViewData(r, PublicationAddSingleVars{
		PageTitle: "Add - Publications - Biblio",
		Step:      1,
	}))
}

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

// func (c *Publications) AddSingle(w http.ResponseWriter, r *http.Request) {
// 	c.Render.HTML(w, http.StatusOK, "publication/add_single", c.ViewData(r, PublicationAddSingleVars{
// 		PageTitle: "Add - Publications - Biblio",
// 		Step:      1,
// 	}))
// }

// func (c *Publications) AddSingleStart(w http.ResponseWriter, r *http.Request) {
// 	c.Render.HTML(w, http.StatusOK, "publication/add_single_start", c.ViewData(r, PublicationAddSingleVars{
// 		PageTitle: "Add - Publications - Biblio",
// 		Step:      2,
// 	}))
// }

func (c *Publications) AddSingleImportConfirm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	source := r.FormValue("source")
	identifier := r.FormValue("identifier")

	// check for duplicates
	if source == "crossref" && identifier != "" {
		if existing, _ := c.Engine.Publications(models.NewSearchArgs().WithFilter("doi", identifier).WithFilter("status", "public")); existing.Total > 0 {
			c.Render.HTML(w, http.StatusOK, "publication/add", c.ViewData(r, PublicationAddSingleVars{
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

		p, err := c.Engine.ImportUserPublicationByIdentifier(userID, source, identifier)

		if err != nil {
			flash := views.Flash{Type: "error"}

			if e, ok := err.(jsonapi.Errors); ok {
				flash.Message = loc.T("publication.single_import", e[0].Code)
			} else {
				log.Println(e)
				flash.Message = loc.T("publication.single_import", "import_by_id.import_failed")
			}

			c.Render.HTML(w, http.StatusOK, "publication/add", c.ViewData(r, PublicationAddSingleVars{
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

			c.Render.HTML(w, http.StatusOK, "publication/add", c.ViewData(r, PublicationAddSingleVars{
				PageTitle: "Add - Publications - Biblio",
				Step:      2,
			},
				flash,
			))
			return
		}

		p, err := c.Engine.CreateUserPublication(userID, pubType)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pub = p
	}

	datasets, err := c.Engine.GetPublicationDatasets(pub.ID)
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

func (c *Publications) AddSingleDescription(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	datasets, err := c.Engine.GetPublicationDatasets(pub.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pub.RelatedDatasetCount = len(datasets)

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

func (c *Publications) AddSinglePublish(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	savedPub, err := c.Engine.PublishPublication(pub)
	if err != nil {

		/*
			TODO: return to /add-single/confirm with flash in session instead of rendering this in the wrong path
			We only use one error, as publishing can only fail on attribute title
		*/
		if e, ok := err.(jsonapi.Errors); ok {

			c.Render.HTML(w, http.StatusOK, "publication/add_single_confirm", c.ViewData(r, struct {
				PageTitle   string
				Step        int
				Publication *models.Publication
			}{
				"Add - Publications - Biblio",
				4,
				pub,
			},
				views.Flash{Type: "error", Message: e[0].Title},
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
		6,
		savedPub,
	}))
}

// func (c *Publications) AddMultiple(w http.ResponseWriter, r *http.Request) {
// 	c.Render.HTML(w, http.StatusOK, "publication/add_multiple", c.ViewData(r, struct {
// 		PageTitle string
// 		Step      int
// 	}{
// 		"Add - Publications - Biblio",
// 		1,
// 	}))
// }

// func (c *Publications) AddMultipleStart(w http.ResponseWriter, r *http.Request) {
// 	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_start", c.ViewData(r, struct {
// 		PageTitle string
// 		Step      int
// 	}{
// 		"Add - Publications - Biblio",
// 		2,
// 	}))
// }

// Step 2: Add publications (WOS, BibTex)

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

	batchID, err := c.Engine.ImportUserPublications(userID, source, file)
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

	hits, err := c.Engine.UserPublications(userID, args.Clone().WithFilter("batch_id", batchID))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("publication_add_multiple_description").URLPath("batch_id", batchID)

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_description", c.ViewData(r, struct {
		PageTitle        string
		Step             int
		SearchURL        *url.URL
		SearchArgs       *models.SearchArgs
		Hits             *models.PublicationHits
		PublicationSorts []string
		BatchID          string
	}{
		"Add - Publications - Biblio",
		3,
		searchURL,
		args,
		hits,
		c.Engine.Vocabularies()["publication_sorts"],
		batchID,
	}),
	)
}

func (c *Publications) AddMultipleDescription(w http.ResponseWriter, r *http.Request) {
	userID := context.GetUser(r.Context()).ID
	batchID := mux.Vars(r)["batch_id"]

	args := models.NewSearchArgs()
	if err := forms.Decode(args, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hits, err := c.Engine.UserPublications(userID, args.Clone().WithFilter("batch_id", batchID))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("publication_add_multiple_description").URLPath("batch_id", batchID)

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_description", c.ViewData(r, struct {
		PageTitle        string
		Step             int
		SearchURL        *url.URL
		SearchArgs       *models.SearchArgs
		Hits             *models.PublicationHits
		PublicationSorts []string
		BatchID          string
	}{
		"Add - Publications - Biblio",
		3,
		searchURL,
		args,
		hits,
		c.Engine.Vocabularies()["publication_sorts"],
		batchID,
	}),
	)
}

func (c *Publications) AddMultipleShow(w http.ResponseWriter, r *http.Request) {
	batchID := mux.Vars(r)["batch_id"]
	pub := context.GetPublication(r.Context())

	datasets, err := c.Engine.GetPublicationDatasets(pub.ID)
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

	pub.RelatedDatasetCount = len(datasets)

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_show", c.ViewData(r, struct {
		PageTitle           string
		Step                int
		Publication         *models.Publication
		PublicationDatasets []*models.Dataset
		Show                *views.ShowBuilder
		Vocabularies        map[string][]string
		SearchArgs          *models.SearchArgs
		BatchID             string
	}{
		"Add - Publications - Biblio",
		3,
		pub,
		datasets,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		c.Engine.Vocabularies(),
		searchArgs,
		batchID,
	}),
	)
}

func (c *Publications) AddMultipleConfirm(w http.ResponseWriter, r *http.Request) {
	userID := context.GetUser(r.Context()).ID
	batchID := mux.Vars(r)["batch_id"]

	args := models.NewSearchArgs()

	hits, err := c.Engine.UserPublications(userID, args.Clone().WithFilter("batch_id", batchID))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("publication_add_multiple_confirm").URLPath("batch_id", batchID)

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_confirm", c.ViewData(r, struct {
		PageTitle        string
		Step             int
		SearchURL        *url.URL
		SearchArgs       *models.SearchArgs
		Hits             *models.PublicationHits
		PublicationSorts []string
		BatchID          string
	}{
		"Add - Publications - Biblio",
		4,
		searchURL,
		args,
		hits,
		c.Engine.Vocabularies()["publication_sorts"],
		batchID,
	}),
	)
}

func (c *Publications) AddMultipleConfirmShow(w http.ResponseWriter, r *http.Request) {
	batchID := mux.Vars(r)["batch_id"]
	pub := context.GetPublication(r.Context())

	datasets, err := c.Engine.GetPublicationDatasets(pub.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pub.RelatedDatasetCount = len(datasets)

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_confirm_show", c.ViewData(r, struct {
		PageTitle           string
		Step                int
		Publication         *models.Publication
		PublicationDatasets []*models.Dataset
		Show                *views.ShowBuilder
		Vocabularies        map[string][]string
		BatchID             string
	}{
		"Add - Publications - Biblio",
		4,
		pub,
		datasets,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		c.Engine.Vocabularies(),
		batchID,
	}),
	)
}

func (c *Publications) AddMultiplePublish(w http.ResponseWriter, r *http.Request) {
	userID := context.GetUser(r.Context()).ID
	batchID := mux.Vars(r)["batch_id"]

	batchFilter := models.NewSearchArgs().WithFilter("batch_id", batchID)

	if err := c.Engine.BatchPublishPublications(userID, batchFilter); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hits, err := c.Engine.UserPublications(userID, batchFilter)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("publication_add_multiple_publish").URLPath("batch_id", batchID)

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_finish", c.ViewData(r, struct {
		PageTitle        string
		Step             int
		SearchURL        *url.URL
		SearchArgs       *models.SearchArgs
		Hits             *models.PublicationHits
		PublicationSorts []string
		BatchID          string
	}{
		"Add - Publications - Biblio",
		5,
		searchURL,
		models.NewSearchArgs(),
		hits,
		c.Engine.Vocabularies()["publication_sorts"],
		batchID,
	}),
	)
}

func (c *Publications) Publish(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	savedPub, err := c.Engine.PublishPublication(pub)

	flashes := make([]views.Flash, 0)
	var publicationErrors jsonapi.Errors
	var publicationErrorsTitle string

	if err != nil {

		savedPub = pub

		if e, ok := err.(jsonapi.Errors); ok {

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

	pubDatasets, err := c.Engine.GetPublicationDatasets(pub.ID)
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

	savedPub.RelatedDatasetCount = len(pubDatasets)

	c.Render.HTML(w, http.StatusOK, "publication/show", c.ViewData(r, struct {
		PageTitle           string
		Publication         *models.Publication
		PublicationDatasets []*models.Dataset
		Show                *views.ShowBuilder
		SearchArgs          *models.SearchArgs
		ErrorsTitle         string
		Errors              jsonapi.Errors
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
	if err := forms.Decode(searchArgs, r.URL.Query()); err != nil {
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
	if err := forms.Decode(searchArgs, r.Form); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := c.Engine.DeletePublication(pub.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hits, err := c.Engine.UserPublications(context.GetUser(r.Context()).ID, searchArgs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("publications").URLPath()

	c.Render.HTML(w, http.StatusOK, "publication/list", c.ViewData(r, struct {
		PageTitle        string
		SearchURL        *url.URL
		SearchArgs       *models.SearchArgs
		Hits             *models.PublicationHits
		PublicationSorts []string
	}{
		"Overview - Publications - Biblio",
		searchURL,
		searchArgs,
		hits,
		c.Engine.Vocabularies()["publication_sorts"],
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

	pub, err := c.Engine.AddPublicationToORCID(user.ORCID, user.ORCIDToken, pub)
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
	// TODO handle error
	id, err := c.Engine.AddPublicationsToORCID(
		context.GetUser(r.Context()).ID,
		models.NewSearchArgs().WithFilter("status", "public"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "task/_status", c.ViewData(r, struct {
		TaskID    string
		TaskState models.TaskState
	}{
		id,
		models.TaskState{Status: models.Waiting},
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
