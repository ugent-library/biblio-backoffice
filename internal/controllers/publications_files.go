package controllers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/unrolled/render"
)

type PublicationsFiles struct {
	engine *engine.Engine
	render *render.Render
}

func NewPublicationsFiles(e *engine.Engine, r *render.Render) *PublicationsFiles {
	return &PublicationsFiles{
		engine: e,
		render: r,
	}
}

func (c *PublicationsFiles) Thumbnail(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fileID := mux.Vars(r)["file_id"]
	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var thumbnailURL string
	for _, file := range pub.File {
		if file.ID == fileID {
			thumbnailURL = file.ThumbnailURL
			break
		}
	}

	if thumbnailURL == "" {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// quick and dirty reverse proxy
	baseURL, _ := url.Parse(c.engine.Config.URL)
	url, _ := url.Parse(thumbnailURL)
	proxy := httputil.NewSingleHostReverseProxy(baseURL)
	// update the headers to allow for SSL redirection
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.URL.Path = strings.Replace(url.Path, baseURL.Path, "", 1)
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host
	r.SetBasicAuth(c.engine.Config.Username, c.engine.Config.Password)
	proxy.ServeHTTP(w, r)
}
