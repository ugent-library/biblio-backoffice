package controllers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

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

type PublicationListVars struct {
	SearchArgs       *engine.SearchArgs
	Hits             *models.PublicationHits
	PublicationSorts []string
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

	c.Render.HTML(w, http.StatusOK, "publication/list",
		views.NewData(c.Render, r, PublicationListVars{
			SearchArgs:       args,
			Hits:             hits,
			PublicationSorts: c.Engine.Vocabularies()["publication_sorts"],
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

	c.Render.HTML(w, http.StatusOK, "publication/show",
		views.NewData(c.Render, r, struct {
			Publication         *models.Publication
			PublicationDatasets []*models.Dataset
			Show                *views.ShowBuilder
			Vocabularies        map[string][]string
			SearchArgs          *engine.SearchArgs
		}{
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
		Step int
	}{
		1,
	}))
}

func (c *Publications) AddSingle(w http.ResponseWriter, r *http.Request) {
	c.Render.HTML(w, http.StatusOK, "publication/add_single", views.NewData(c.Render, r, struct {
		Step int
	}{
		2,
	}))
}

func (c *Publications) AddSingleImport(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var publication *models.Publication

	if identifier := r.FormValue("identifier"); identifier != "" {
		var identifier_type string = r.FormValue("identifier_type")
		if identifier_type == "" {
			log.Println("missing form field identifier_type")
			http.Error(w, "missing form field identifier_type", http.StatusInternalServerError)
			return
		}
		publications, err := c.Engine.ImportUserPublications(context.GetUser(r.Context()).ID, identifier, identifier_type)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO flash messages
		if len(publications) == 0 {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		publication = publications[0]
	} else {
		pt := r.FormValue("publication_type")
		p, err := c.Engine.CreatePublication(pt)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		publication = p
	}

	c.Render.HTML(w, http.StatusOK, "publication/add_single_files", views.NewData(c.Render, r, struct {
		Step        int
		Publication *models.Publication
		Show        *views.ShowBuilder
	}{
		3,
		publication,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
	}))
}

func (c *Publications) AddSingleDescription(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	c.Render.HTML(w, http.StatusOK, "publication/add_single_description", views.NewData(c.Render, r, struct {
		Step        int
		Publication *models.Publication
		Show        *views.ShowBuilder
	}{
		4,
		pub,
		views.NewShowBuilder(c.Render, locale.Get(r.Context())),
	}))
}

func (c *Publications) AddSingleConfirm(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	c.Render.HTML(w, http.StatusOK, "publication/add_single_confirm", views.NewData(c.Render, r, struct {
		Step        int
		Publication *models.Publication
	}{
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
		Step        int
		Publication *models.Publication
	}{
		6,
		savedPub,
	}))
}

func (c *Publications) AddMultiple(w http.ResponseWriter, r *http.Request) {
	c.Render.HTML(w, http.StatusOK, "publication/add_multiple", views.NewData(c.Render, r, struct {
		Step int
	}{
		2,
	}))
}
