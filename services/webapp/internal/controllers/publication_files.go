package controllers

import (
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/unrolled/render"
)

type PublicationFiles struct {
	Context
}

func NewPublicationFiles(c Context) *PublicationFiles {
	return &PublicationFiles{c}
}

func (c *PublicationFiles) Download(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())
	fileID := mux.Vars(r)["file_id"]
	file := pub.GetFile(fileID)

	if file == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, c.Engine.FileStore.FilePath(file.SHA256))
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

	c.Engine.ServePublicationThumbnail(thumbnailURL, w, r)
}

func (c *PublicationFiles) Upload(w http.ResponseWriter, r *http.Request) {
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

	pub := context.GetPublication(r.Context())

	checksum, err := c.Engine.FileStore.Add(file)

	if err != nil {
		flash := views.Flash{Type: "error", Message: "There was a problem adding your file"}
		c.Render.HTML(w, http.StatusCreated, "publication/files/_list",
			c.ViewData(r, struct {
				Publication *models.Publication
			}{
				pub,
			},
				flash,
			),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)
		return
	}

	// TODO check if file with same checksum is already present
	savedPub := pub.Clone()
	now := time.Now()
	pubFile := &models.PublicationFile{
		ID:          uuid.New().String(),
		AccessLevel: "local",
		Filename:    handler.Filename,
		FileSize:    int(handler.Size),
		ContentType: filetype,
		SHA256:      checksum,
		DateCreated: &now,
		DateUpdated: &now,
	}
	savedPub.File = append(savedPub.File, pubFile)
	err = c.Engine.UpdatePublication(savedPub)

	if err != nil {
		flash := views.Flash{Type: "error", Message: "There was a problem adding your file"}
		c.Render.HTML(w, http.StatusCreated, "publication/files/_list",
			c.ViewData(r, struct {
				Publication *models.Publication
			}{
				pub,
			},
				flash,
			),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)
		return
	}

	c.Render.HTML(w, http.StatusCreated, "publication/files/_upload_edit",
		c.ViewData(r, struct {
			Publication *models.Publication
			File        *models.PublicationFile
			FileIndex   int
			Form        *views.FormBuilder
		}{
			savedPub,
			savedPub.File[len(savedPub.File)-1],
			len(savedPub.File) - 1,
			views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
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
		Publication *models.Publication
		File        *models.PublicationFile
		FileIndex   int
		Form        *views.FormBuilder
	}{
		pub,
		file,
		fileIndex,
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationFiles) License(w http.ResponseWriter, r *http.Request) {
	fileID := mux.Vars(r)["file_id"]
	triggerEl := r.Header.Get("HX-Trigger-name")

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

	if err := DecodeForm(file, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Conditionally voiding licensing fields based on a triggering element
	// and dependent element inputs

	// Clear embargo fields when Access Level is set to "open access"
	if triggerEl == "" {
		if file.AccessLevel == "open_access" {
			file.Embargo = ""
			file.EmbargoTo = ""
		}
	}

	// Clear CC License when Embargo To is set to "local" or "empty"
	if triggerEl == "embargo_to" {
		if (file.EmbargoTo == "" || file.EmbargoTo == "local") && file.CCLicense != "" {
			file.CCLicense = ""
		}
	}

	// Clear No License (Copyright) when CC License is set to a non-empty value
	if triggerEl == "cc_license" {
		if file.CCLicense != "" && file.NoLicense == true {
			file.NoLicense = false
		}
	}

	// Clear CC License when No License (Copyright) is checked
	if triggerEl == "no_license" {
		if file.CCLicense != "" && file.NoLicense == true {
			file.CCLicense = ""
		}
	}

	c.Render.HTML(w, http.StatusOK, "publication/files/_edit", c.ViewData(r, struct {
		Publication *models.Publication
		File        *models.PublicationFile
		FileIndex   int
		Form        *views.FormBuilder
	}{
		pub,
		file,
		fileIndex,
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
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

	if err := DecodeForm(file, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Embargo sanity / paranoia check
	// Ensure embargo fields are definitely empty when Acces Level is set to "Open Access"
	if file.AccessLevel == "open_access" || file.EmbargoTo == file.AccessLevel {
		file.EmbargoTo = ""
		file.Embargo = ""
	}

	err = c.Engine.UpdatePublication(pub)

	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {
		c.Render.HTML(w, http.StatusOK, "publication/files/_edit", c.ViewData(r, struct {
			Publication *models.Publication
			File        *models.PublicationFile
			FileIndex   int
			Form        *views.FormBuilder
		}{
			pub,
			file,
			fileIndex,
			views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), validationErrors),
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
	pub, err = c.Engine.Store.GetPublication(pub.ID)
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
	pub := context.GetPublication(r.Context())
	fileID := mux.Vars(r)["file_id"]
	newFile := []*models.PublicationFile{}
	for _, f := range pub.File {
		if f.ID != fileID {
			newFile = append(newFile, f)
		}
	}
	pub.File = newFile

	err := c.Engine.UpdatePublication(pub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusCreated, "publication/files/_list", c.ViewData(r, struct {
		Publication *models.Publication
	}{
		pub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
