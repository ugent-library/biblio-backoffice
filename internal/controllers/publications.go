package controllers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
	"github.com/unrolled/render"
)

type Publications struct {
	Context
}

func NewPublications(c Context) *Publications {
	return &Publications{c}
}

func (c *Publications) List(w http.ResponseWriter, r *http.Request) {
	args := engine.NewSearchArgs()
	if err := forms.Decode(args, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hits, err := c.Engine.UserPublications(context.GetUser(r.Context()).ID, args)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("publications").URLPath()

	c.Render.HTML(w, http.StatusOK, "publication/list", views.NewData(c.Render, r, struct {
		PageTitle        string
		SearchURL        *url.URL
		SearchArgs       *engine.SearchArgs
		Hits             *models.PublicationHits
		PublicationSorts []string
	}{
		"Overview - Publications - Biblio",
		searchURL,
		args,
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

	searchArgs := engine.NewSearchArgs()
	if err := forms.Decode(searchArgs, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/show", views.NewData(c.Render, r, struct {
		PageTitle           string
		Publication         *models.Publication
		PublicationDatasets []*models.Dataset
		Show                *views.ShowBuilder
		Vocabularies        map[string][]string
		SearchArgs          *engine.SearchArgs
	}{
		"Publication - Biblio",
		pub,
		datasets,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
		c.Engine.Vocabularies(),
		searchArgs,
	}),
	)
}

func (c *Publications) Thumbnail(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	if pub.ThumbnailURL() == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// quick and dirty reverse proxy
	baseURL, _ := url.Parse(c.Engine.Config.URL)
	url, _ := url.Parse(pub.ThumbnailURL())
	proxy := httputil.NewSingleHostReverseProxy(baseURL)
	// update the headers to allow for SSL redirection
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.URL.Path = strings.Replace(url.Path, baseURL.Path, "", 1)
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Header.Del("Cookie")
	r.Host = url.Host
	r.SetBasicAuth(c.Engine.Config.Username, c.Engine.Config.Password)
	proxy.ServeHTTP(w, r)
}

func (c *Publications) Summary(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	c.Render.HTML(w, http.StatusOK, "publication/_summary", views.NewData(c.Render, r, struct {
		Publication *models.Publication
	}{
		pub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *Publications) Add(w http.ResponseWriter, r *http.Request) {
	c.Render.HTML(w, http.StatusOK, "publication/add", views.NewData(c.Render, r, struct {
		PageTitle string
		Step      int
	}{
		"Add - Publications - Biblio",
		1,
	}))
}

func (c *Publications) AddSingle(w http.ResponseWriter, r *http.Request) {
	c.Render.HTML(w, http.StatusOK, "publication/add_single", views.NewData(c.Render, r, struct {
		PageTitle string
		Step      int
	}{
		"Add - Publications - Biblio",
		2,
	}))
}

func (c *Publications) AddSingleImport(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var pub *models.Publication

	if identifier := r.FormValue("identifier"); identifier != "" {
		var source string = r.FormValue("source")

		p, err := c.Engine.ImportUserPublicationByIdentifier(context.GetUser(r.Context()).ID, source, identifier)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pub = p
	} else {
		pt := r.FormValue("publication_type")
		p, err := c.Engine.CreatePublication(pt)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pub = p
	}

	c.Render.HTML(w, http.StatusOK, "publication/add_single_files", views.NewData(c.Render, r, struct {
		PageTitle   string
		Step        int
		Publication *models.Publication
		Show        *views.ShowBuilder
	}{
		"Add - Publications - Biblio",
		3,
		pub,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
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

	c.Render.HTML(w, http.StatusOK, "publication/add_single_description", views.NewData(c.Render, r, struct {
		PageTitle           string
		Step                int
		Publication         *models.Publication
		PublicationDatasets []*models.Dataset
		Show                *views.ShowBuilder
	}{
		"Add - Publications - Biblio",
		4,
		pub,
		datasets,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
	}))
}

func (c *Publications) AddSingleConfirm(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	c.Render.HTML(w, http.StatusOK, "publication/add_single_confirm", views.NewData(c.Render, r, struct {
		PageTitle   string
		Step        int
		Publication *models.Publication
	}{
		"Add - Publications - Biblio",
		5,
		pub,
	}))
}

func (c *Publications) AddSinglePublish(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	savedPub, err := c.Engine.PublishPublication(pub)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/add_single_publish", views.NewData(c.Render, r, struct {
		PageTitle   string
		Step        int
		Publication *models.Publication
	}{
		"Add - Publications - Biblio",
		6,
		savedPub,
	}))
}

func (c *Publications) AddMultiple(w http.ResponseWriter, r *http.Request) {
	c.Render.HTML(w, http.StatusOK, "publication/add_multiple", views.NewData(c.Render, r, struct {
		PageTitle string
		Step      int
	}{
		"Add - Publications - Biblio",
		2,
	}))
}

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	args := engine.NewSearchArgs()

	hits, err := c.Engine.UserPublications(userID, args.Clone().WithFilter("batch_id", batchID))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("publication_add_multiple_description").URLPath("batch_id", batchID)

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_description", views.NewData(c.Render, r, struct {
		PageTitle        string
		Step             int
		SearchURL        *url.URL
		SearchArgs       *engine.SearchArgs
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

	args := engine.NewSearchArgs()
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

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_description", views.NewData(c.Render, r, struct {
		PageTitle        string
		Step             int
		SearchURL        *url.URL
		SearchArgs       *engine.SearchArgs
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

	searchArgs := engine.NewSearchArgs()
	if err := forms.Decode(searchArgs, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_show", views.NewData(c.Render, r, struct {
		PageTitle           string
		Step                int
		Publication         *models.Publication
		PublicationDatasets []*models.Dataset
		Show                *views.ShowBuilder
		Vocabularies        map[string][]string
		SearchArgs          *engine.SearchArgs
		BatchID             string
	}{
		"Add - Publications - Biblio",
		3,
		pub,
		datasets,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
		c.Engine.Vocabularies(),
		searchArgs,
		batchID,
	}),
	)
}

func (c *Publications) AddMultipleConfirm(w http.ResponseWriter, r *http.Request) {
	userID := context.GetUser(r.Context()).ID
	batchID := mux.Vars(r)["batch_id"]

	args := engine.NewSearchArgs()

	hits, err := c.Engine.UserPublications(userID, args.Clone().WithFilter("batch_id", batchID))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchURL, _ := c.Router.Get("publication_add_multiple_confirm").URLPath("batch_id", batchID)

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_confirm", views.NewData(c.Render, r, struct {
		PageTitle        string
		Step             int
		SearchURL        *url.URL
		SearchArgs       *engine.SearchArgs
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

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_confirm_show", views.NewData(c.Render, r, struct {
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
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
		c.Engine.Vocabularies(),
		batchID,
	}),
	)
}

func (c *Publications) AddMultiplePublish(w http.ResponseWriter, r *http.Request) {
	userID := context.GetUser(r.Context()).ID
	batchID := mux.Vars(r)["batch_id"]

	batchFilter := engine.NewSearchArgs().WithFilter("batch_id", batchID)

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

	c.Render.HTML(w, http.StatusOK, "publication/add_multiple_publish", views.NewData(c.Render, r, struct {
		PageTitle        string
		Step             int
		SearchURL        *url.URL
		SearchArgs       *engine.SearchArgs
		Hits             *models.PublicationHits
		PublicationSorts []string
		BatchID          string
	}{
		"Add - Publications - Biblio",
		5,
		searchURL,
		engine.NewSearchArgs(),
		hits,
		c.Engine.Vocabularies()["publication_sorts"],
		batchID,
	}),
	)
}
