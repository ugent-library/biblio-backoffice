package controllers

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type PublicationFiles struct {
	Context
}

func NewPublicationFiles(c Context) *PublicationFiles {
	return &PublicationFiles{c}
}

func (c *PublicationFiles) Download(w http.ResponseWriter, r *http.Request) {
	fileID := mux.Vars(r)["file_id"]

	pub := context.GetPublication(r.Context())

	var fileURL string
	for _, file := range pub.File {
		if file.ID == fileID {
			fileURL = file.URL
			break
		}
	}

	if fileURL == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// quick and dirty reverse proxy
	baseURL, _ := url.Parse(c.Engine.Config.URL)
	url, _ := url.Parse(fileURL)
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

func (c *PublicationFiles) Thumbnail(w http.ResponseWriter, r *http.Request) {
	fileID := mux.Vars(r)["file_id"]

	pub := context.GetPublication(r.Context())

	var thumbnailURL string
	for _, file := range pub.File {
		if file.ID == fileID {
			thumbnailURL = file.ThumbnailURL
			break
		}
	}

	if thumbnailURL == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// quick and dirty reverse proxy
	baseURL, _ := url.Parse(c.Engine.Config.URL)
	url, _ := url.Parse(thumbnailURL)
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

func (c *PublicationFiles) Upload(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// 2GB limit on request body
	r.Body = http.MaxBytesReader(w, r.Body, 2000000000)

	// buffer limit of 32MB
	if err := r.ParseMultipartForm(32 << 20); err != nil {
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

	if err := c.Engine.AddPublicationFile(id, pubFile, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pub, _ := c.Engine.GetPublication(id)

	c.Render.HTML(w, http.StatusCreated, "publication/files/_show",
		views.NewData(c.Render, r, struct {
			Publication *models.Publication
		}{
			pub,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationFiles) Edit(w http.ResponseWriter, r *http.Request) {
	fileID := mux.Vars(r)["file_id"]

	pub := context.GetPublication(r.Context())

	var file *models.PublicationFile
	for _, f := range pub.File {
		if f.ID == fileID {
			file = f
		}
	}

	if file == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/files/_edit", views.NewData(c.Render, r, struct {
		Publication  *models.Publication
		File         *models.PublicationFile
		Vocabularies map[string][]string
	}{
		pub,
		file,
		c.Engine.Vocabularies(),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// TODO avoid getting publication multiple times
func (c *PublicationFiles) Update(w http.ResponseWriter, r *http.Request) {
	fileID := mux.Vars(r)["file_id"]

	pub := context.GetPublication(r.Context())

	var file *models.PublicationFile
	for _, f := range pub.File {
		if f.ID == fileID {
			file = f
		}
	}

	if file == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := forms.Decode(file, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO save file metadata
	err = c.Engine.UpdatePublicationFile(pub.ID, file)

	// TODO show errors
	if _, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK, "publication/files/_edit", views.NewData(c.Render, r, struct {
			Publication  *models.Publication
			File         *models.PublicationFile
			Vocabularies map[string][]string
		}{
			pub,
			file,
			c.Engine.Vocabularies(),
		}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pub, err = c.Engine.GetPublication(pub.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/files/_show", views.NewData(c.Render, r, struct {
		Publication *models.Publication
	}{
		pub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationFiles) ConfirmRemove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fileID := mux.Vars(r)["file_id"]

	c.Render.HTML(w, http.StatusOK, "publication/files/_modal_confirm_removal", views.NewData(c.Render, r, struct {
		PublicationID string
		FileID        string
	}{
		id,
		fileID,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationFiles) Remove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fileID := mux.Vars(r)["file_id"]

	if err := c.Engine.RemovePublicationFile(id, fileID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pub, _ := c.Engine.GetPublication(id)

	c.Render.HTML(w, http.StatusCreated, "publication/files/_show", views.NewData(c.Render, r, struct {
		Publication *models.Publication
	}{
		pub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
