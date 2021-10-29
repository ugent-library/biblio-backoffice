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
	id := mux.Vars(r)["id"]
	pub, err := c.Engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	datasets, err := c.Engine.GetPublicationDatasets(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pub.Dataset = datasets

	c.Render.HTML(w, http.StatusOK, "publication/show",
		views.NewData(c.Render, r, struct {
			Publication  *models.Publication
			Show         *views.ShowBuilder
			Vocabularies map[string][]string
		}{
			pub,
			views.NewShowBuilder(c.Render, locale.Get(r.Context())),
			c.Engine.Vocabularies(),
		}),
	)
}

func (c *Publications) Thumbnail(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	pub, err := c.Engine.GetPublication(id)
	if err != nil || pub.ThumbnailURL() == "" {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
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
	r.Host = url.Host
	r.SetBasicAuth(c.Engine.Config.Username, c.Engine.Config.Password)
	proxy.ServeHTTP(w, r)
}
func (c *Publications) Summary(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := c.Engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.Render.HTML(w, http.StatusOK,
		"publication/_summary",
		pub,
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
