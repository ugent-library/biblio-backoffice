package controllers

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/unrolled/render"
)

type PublicationFiles struct {
	engine *engine.Engine
	render *render.Render
	router *mux.Router
}

func NewPublicationFiles(e *engine.Engine, r *render.Render, router *mux.Router) *PublicationFiles {
	return &PublicationFiles{
		engine: e,
		render: r,
		router: router,
	}
}

func (c *PublicationFiles) Download(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fileID := mux.Vars(r)["file_id"]
	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var fileURL string
	for _, file := range pub.File {
		if file.ID == fileID {
			fileURL = file.URL
			break
		}
	}

	if fileURL == "" {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// quick and dirty reverse proxy
	baseURL, _ := url.Parse(c.engine.Config.URL)
	url, _ := url.Parse(fileURL)
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

func (c *PublicationFiles) Thumbnail(w http.ResponseWriter, r *http.Request) {
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

func (c *PublicationFiles) Upload(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// 2GB limit
	if err := r.ParseMultipartForm(2000000000); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// detect content type
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filetype := http.DetectContentType(buff)
	// rewind
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pubFile := models.PublicationFile{
		Filename:    handler.Filename,
		FileSize:    int(handler.Size),
		ContentType: filetype,
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.engine.AddPublicationFile(id, pubFile, fileContents); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pub, _ := c.engine.GetPublication(id)

	c.render.HTML(w, http.StatusCreated, "publication/show",
		struct {
			views.Data
			Publication  *models.Publication
			Show         *views.ShowBuilder
			Vocabularies map[string][]string
		}{
			views.NewData(c.render, r),
			pub,
			views.NewShowBuilder(c.render, locale.Get(r.Context())),
			c.engine.Vocabularies(),
		},
	)
}
