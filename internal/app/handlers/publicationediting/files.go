package publicationediting

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type BindFile struct {
	FileID string `path:"file_id"`

	AccessLevel        string `form:"access_level"`
	License            string `form:"license"`
	Description        string `form:"description"`
	Embargo            string `form:"embargo"`
	EmbargoTo          string `form:"embargo_to"`
	OtherLicense       string `form:"other_license"`
	PublicationVersion string `form:"publication_version"`
	Relation           string `form:"relation"`
	Title              string `form:"title"`
}

type YieldEditFile struct {
	Context
	File *models.PublicationFile
	Form *form.Form
}

type YieldShowFiles struct {
	Context
}

type YieldDeleteFile struct {
	Context
	File *models.PublicationFile
}

func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request, ctx Context) {
	// 2GB limit on request body
	r.Body = http.MaxBytesReader(w, r.Body, 2000000000)

	// buffer limit of 32MB
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.Logger.Errorf("publication upload file: exceeded buffer limit:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		h.Logger.Errorf("publication upload file: could not process file", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}
	defer file.Close()

	// detect content type
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		h.Logger.Errorf("publication upload file: could not read file", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}
	filetype := http.DetectContentType(buff)

	// rewind
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		h.Logger.Errorf("publication upload file: could not read file", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// add file to filestore
	checksum, err := h.FileStore.Add(file)

	if err != nil {
		h.Logger.Errorf("publication upload file: could not save file", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// save publication
	// TODO check if file with same checksum is already present
	pubFile := models.PublicationFile{
		AccessLevel: "local",
		Name:        handler.Filename,
		Size:        int(handler.Size),
		ContentType: filetype,
		SHA256:      checksum,
	}
	/*
		automatically generates extra fields:
		id, date_created, date_updated
	*/
	ctx.Publication.AddFile(&pubFile)

	// TODO conflict handling
	if err := h.Repository.SavePublication(ctx.Publication); err != nil {
		h.Logger.Errorf("publication upload file: could not save publication", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// render edit file form
	render.Layout(w, "show_modal", "publication/edit_file", YieldEditFile{
		Context: ctx,
		File:    &pubFile,
		Form:    fileForm(ctx.Locale, ctx.Publication, &pubFile, nil),
	})

}

func (h *Handler) EditFile(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindFile{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	var file *models.PublicationFile

	for _, f := range ctx.Publication.File {
		if f.ID == b.FileID {
			file = f
			break
		}
	}

	if file == nil {
		h.Logger.Warnw("publication upload file: could not find file", "fileid", b.FileID, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	render.Layout(w, "show_modal", "publication/edit_file", YieldEditFile{
		Context: ctx,
		File:    file,
		Form:    fileForm(ctx.Locale, ctx.Publication, file, nil),
	})
}

// TODO add more rules
func (h *Handler) EditFileLicense(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindFile{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	var file *models.PublicationFile
	for _, f := range ctx.Publication.File {
		if f.ID == b.FileID {
			file = f
			break
		}
	}

	if file == nil {
		h.Logger.Warnw("publication upload file: could not find file", "fileid", b.FileID, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// clear embargo fields when access level is set to open access
	if b.AccessLevel == "open_access" {
		b.Embargo = ""
		b.EmbargoTo = ""
	}

	// TODO apply other license && access level related rules

	// copy attributes from bind to file
	bindToPublicationFile(&b, file)

	render.Layout(w, "refresh_modal", "publication/edit_file", YieldEditFile{
		Context: ctx,
		File:    file,
		Form:    fileForm(ctx.Locale, ctx.Publication, file, nil),
	})

}

func (h *Handler) UpdateFile(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindFile{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	pub := ctx.Publication
	var file *models.PublicationFile
	for _, f := range pub.File {
		if f.ID == b.FileID {
			file = f
			break
		}
	}

	if file == nil {
		h.Logger.Warnw("publication edit file license: could not find file", "fileid", b.FileID, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// copy attributes from bind to file
	bindToPublicationFile(&b, file)

	// TODO conflict detection
	// save publication
	err := h.Repository.SavePublication(pub)

	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {
		render.Layout(w, "refresh_modal", "publication/edit_file", YieldEditFile{
			Context: ctx,
			File:    file,
			Form:    fileForm(ctx.Locale, ctx.Publication, file, validationErrors),
		})
		return
	} else if err != nil {
		h.Logger.Errorf("publication edit file license: could not save file", "fileid", b.FileID, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// TODO why?
	// load publication again
	ctx.Publication, err = h.Repository.GetPublication(pub.ID)
	if err != nil {
		h.Logger.Warnw("publication edit file: could not get publication", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_files", YieldShowFiles{
		Context: ctx,
	})
}

func (h *Handler) ConfirmDeleteFile(w http.ResponseWriter, r *http.Request, ctx Context) {

	b := BindFile{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("confirm delete publication file: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	file := ctx.Publication.GetFile(b.FileID)

	if file == nil {
		h.Logger.Warnw("confirm delete publication file: could not find file", "fileid", b.FileID, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	render.Layout(w, "show_modal", "publication/confirm_delete_file", YieldDeleteFile{
		Context: ctx,
		File:    file,
	})

}

func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindFile{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete publication file: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	ctx.Publication.RemoveFile(b.FileID)

	// save publication
	if err := h.Repository.SavePublication(ctx.Publication); err != nil {
		h.Logger.Errorw("delete publication file: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// render
	render.View(w, "publication/refresh_files", YieldShowFiles{
		Context: ctx,
	})

}

func fileForm(l *locale.Locale, publication *models.Publication, file *models.PublicationFile, errors validation.Errors) *form.Form {
	idx := -1
	for i, f := range publication.File {
		if f.ID == file.ID {
			idx = i
			break
		}
	}

	f := form.New().
		WithTheme("cols").
		WithErrors(localize.ValidationErrors(l, errors))

	f.AddSection(
		&form.Select{
			Name:        "relation",
			Value:       file.Relation,
			Label:       l.T("builder.file.relation"),
			EmptyOption: true,
			Options: localize.VocabularySelectOptions(
				l,
				"publication_file_relations"),
			Cols: 12,
			Error: localize.ValidationErrorAt(
				l,
				errors,
				fmt.Sprintf("/file/%d/relation", idx)),
		},
		&form.Select{
			Name:        "publication_version",
			Value:       file.PublicationVersion,
			Label:       l.T("builder.file.publication_version"),
			EmptyOption: true,
			Options:     localize.VocabularySelectOptions(l, "publication_versions"),
			Cols:        12,
			Error: localize.ValidationErrorAt(
				l,
				errors,
				fmt.Sprintf("/file/%d/publication_version", idx)),
		},
	)

	f.AddSection(
		&form.RadioButtonGroup{
			Template: "publication/file_access_level",
			Name:     "access_level",
			Value:    file.AccessLevel,
			Label:    l.T("builder.file.access_level"),
			Options:  localize.VocabularySelectOptions(l, "publication_file_access_levels"),
			Cols:     9,
			Error: localize.ValidationErrorAt(
				l,
				errors,
				fmt.Sprintf("/file/%d/access_level", idx)),
			Vars: struct {
				ID     string
				FileID string
			}{
				ID:     publication.ID,
				FileID: file.ID,
			},
		},
	)

	f.AddSection(
		&form.Date{
			Name:  "embargo",
			Value: file.Embargo,
			Label: l.T("builder.file.embargo"),
			Min:   time.Now().Format("2006-01-02"),
			Cols:  9,
			Error: localize.ValidationErrorAt(
				l,
				errors,
				fmt.Sprintf("/file/%d/embargo", idx)),
		},
		&form.Select{
			Name:        "embargo_to",
			Value:       file.EmbargoTo,
			Label:       l.T("builder.file.embargo_to"),
			EmptyOption: true,
			Options:     localize.VocabularySelectOptions(l, "publication_file_access_levels"),
			Cols:        9,
			Error: localize.ValidationErrorAt(
				l,
				errors,
				fmt.Sprintf("/file/%d/embargo_to", idx)),
		},
	)

	f.AddSection(
		&form.Select{
			Name:        "license",
			Value:       file.License,
			Label:       l.T("builder.file.license"),
			EmptyOption: true,
			Options:     localize.VocabularySelectOptions(l, "licenses"),
			Cols:        9,
			Error: localize.ValidationErrorAt(
				l,
				errors,
				fmt.Sprintf("/file/%d/license", idx)),
		},
	)

	f.AddSection(
		&form.Text{
			Name:  "other_license",
			Value: file.OtherLicense,
			Label: l.T("builder.file.other_license"),
			Help:  l.T("builder.file.other_license.help"),
			Cols:  12,
			Error: localize.ValidationErrorAt(
				l,
				errors,
				fmt.Sprintf("/file/%d/other_license", idx)),
		},
	)

	return f
}

func bindToPublicationFile(b *BindFile, f *models.PublicationFile) {
	f.AccessLevel = b.AccessLevel
	f.License = b.License
	f.Embargo = b.Embargo
	f.EmbargoTo = b.EmbargoTo
	f.OtherLicense = b.OtherLicense
	f.PublicationVersion = b.PublicationVersion
	f.Relation = b.Relation
	f.ID = b.FileID
}
