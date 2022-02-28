package controllers

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
	"github.com/ugent-library/go-locale/locale"
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

	c.Engine.ServePublicationFile(fileURL, w, r)
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

	c.Engine.ServePublicationFile(thumbnailURL, w, r)
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

	pubFile := &models.PublicationFile{
		Filename:    handler.Filename,
		FileSize:    int(handler.Size),
		ContentType: filetype,
	}

	if err = c.Engine.AddPublicationFile(id, pubFile, file); err != nil {
		flash := views.Flash{Type: "error", Message: "There was a problem adding your file"}
		if apiErrors, ok := err.(jsonapi.Errors); ok && apiErrors[0].Code == "api.create_publication_file.file_already_present" {
			flash = views.Flash{Type: "warning", Message: "A file with the same name is already attached to this publication"}
		}
		c.Render.HTML(w, http.StatusCreated, "publication/files/_list",
			c.ViewData(r, struct {
				Publication *models.Publication
			}{
				context.GetPublication(r.Context()),
			},
				flash,
			),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)
		return
	}

	pub, _ := c.Engine.GetPublication(id)

	c.Render.HTML(w, http.StatusCreated, "publication/files/_upload_edit",
		c.ViewData(r, struct {
			Publication  *models.Publication
			File         *models.PublicationFile
			FileIndex    int
			Form         *views.FormBuilder
			Vocabularies map[string][]string
		}{
			pub,
			pub.File[len(pub.File)-1],
			len(pub.File) - 1,
			views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
			c.Engine.Vocabularies(),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationFiles) Edit(w http.ResponseWriter, r *http.Request) {
	fileID := mux.Vars(r)["file_id"]

	pub := context.GetPublication(r.Context())

	var file *models.PublicationFile
	fileIndex := 0
	for i, f := range pub.File {
		if f.ID == fileID {
			file = f
			fileIndex = i
		}
	}

	if file == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/files/_edit", c.ViewData(r, struct {
		Publication  *models.Publication
		File         *models.PublicationFile
		FileIndex    int
		Form         *views.FormBuilder
		Vocabularies map[string][]string
	}{
		pub,
		file,
		fileIndex,
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
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
	fileIndex := 0
	for i, f := range pub.File {
		if f.ID == fileID {
			file = f
			fileIndex = i
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

	// TODO handle checkbox boolean values elegantly
	if r.FormValue("no_license") != "true" {
		file.NoLicense = false
	}

	// embargo sanity check
	if file.AccessLevel == "open_access" || file.EmbargoTo == file.AccessLevel {
		file.EmbargoTo = ""
		file.Embargo = ""
	}

	log.Printf("%+v", r.Form)
	log.Printf("%+v", file)

	err = c.Engine.UpdatePublicationFile(pub.ID, file)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK, "publication/files/_edit", c.ViewData(r, struct {
			Publication  *models.Publication
			File         *models.PublicationFile
			FileIndex    int
			Form         *views.FormBuilder
			Vocabularies map[string][]string
		}{
			pub,
			file,
			fileIndex,
			views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), formErrors),
			c.Engine.Vocabularies(),
		},
			views.Flash{Type: "error", Message: "There are some problems with your input"},
		),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//reload publication
	pub, err = c.Engine.GetPublication(pub.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "publication/files/_update", c.ViewData(r, struct {
		Publication *models.Publication
	}{
		pub,
	},
		views.Flash{Type: "success", Message: "File metadata updated succesfully", DismissAfter: 5 * time.Second},
	),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationFiles) ConfirmRemove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fileID := mux.Vars(r)["file_id"]

	c.Render.HTML(w, http.StatusOK, "publication/files/_modal_confirm_removal", c.ViewData(r, struct {
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

	c.Render.HTML(w, http.StatusCreated, "publication/files/_list", c.ViewData(r, struct {
		Publication *models.Publication
	}{
		pub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
