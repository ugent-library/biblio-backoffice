package publicationediting

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/cshum/imagor/imagorpath"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/flash"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/ulid"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
)

type BindFile struct {
	AccessLevel        string `form:"access_level"`
	CCLicense          string `form:"cc_license"`
	Description        string `form:"description"`
	Embargo            string `form:"embargo"`
	EmbargoTo          string `form:"embargo_to"`
	NoLicense          bool   `form:"no_license"`
	OtherLicense       string `form:"other_license"`
	PublicationVersion string `form:"publication_version"`
	Relation           string `form:"relation"`
	Title              string `form:"title"`

	// for extraction from URL
	ID string `path:"file_id"`
	// for error extraction after validation
	Index int
}

type YieldEditFile struct {
	Context
	File      *models.PublicationFile
	FileIndex int
	Form      *form.Form
}
type YieldShowFiles struct {
	Context
}

type YieldDeleteFile struct {
	Context
	File *models.PublicationFile
}

func (h *Handler) DownloadFile(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindFile{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	file := ctx.Publication.GetFile(b.ID)

	if file == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename*=UTF-8''%s", url.PathEscape(file.Filename)),
	)
	http.ServeFile(w, r, h.FileStore.FilePath(file.SHA256))
}

func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request, ctx Context) {

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

	pub := ctx.Publication

	// add file to filestore
	checksum, err := h.FileStore.Add(file)

	if err != nil {
		// TODO: render flash
		log.Print(err.Error())
		ctx.Flash = append(ctx.Flash, flash.Flash{
			Type:         "error",
			Body:         "There was a problem adding your file",
			DismissAfter: 5 * time.Second,
		})
		render.Render(w, "publication/refresh_files", YieldShowFiles{
			Context: ctx,
		})
		return
	}

	// save publication
	// TODO check if file with same checksum is already present
	savedPub := pub.Clone()
	now := time.Now()
	pubFile := &models.PublicationFile{
		ID:          ulid.MustGenerate(),
		AccessLevel: "local",
		Filename:    handler.Filename,
		FileSize:    int(handler.Size),
		ContentType: filetype,
		SHA256:      checksum,
		DateCreated: &now,
		DateUpdated: &now,
	}
	savedPub.File = append(savedPub.File, pubFile)

	// add thumbnail(s)
	h.addThumbnailURLs(savedPub)

	if e := h.Repository.SavePublication(savedPub); e != nil {
		// TODO: render flash
		log.Print(e.Error())
		ctx.Flash = append(ctx.Flash, flash.Flash{
			Type:         "error",
			Body:         "There was a problem adding your file",
			DismissAfter: 5 * time.Second,
		})
		render.Render(w, "publication/refresh_files", YieldShowFiles{
			Context: ctx,
		})
		return
	}

	// update publication in context
	ctx.Publication = savedPub

	// populate bind with stored publication file
	bindFile := BindFile{}
	publicationFileToBind(pubFile, &bindFile)
	bindFile.Index = len(savedPub.File) - 1

	// render edit file form
	render.Render(w, "publication/refresh_edit_file", YieldEditFile{
		Context:   ctx,
		File:      pubFile,
		FileIndex: bindFile.Index,
		Form:      fileForm(ctx, bindFile, nil),
	})

}

func (h *Handler) EditFile(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindFile{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	var file *models.PublicationFile

	for i, f := range ctx.Publication.File {
		if f.ID == b.ID {
			file = f
			b.Index = i
			break
		}
	}

	if file == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	publicationFileToBind(file, &b)

	render.Render(w, "publication/edit_file", YieldEditFile{
		Context:   ctx,
		File:      file,
		FileIndex: b.Index,
		Form:      fileForm(ctx, b, nil),
	})
}

func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindFile{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	pub := ctx.Publication.Clone()
	newFile := []*models.PublicationFile{}
	for _, f := range pub.File {
		if f.ID != b.ID {
			newFile = append(newFile, f)
		}
	}
	pub.File = newFile

	// add thumbnail urls
	h.addThumbnailURLs(pub)

	// save publication
	if err := h.Repository.SavePublication(pub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update publication in context
	ctx.Publication = pub

	// render
	render.Render(w, "publication/refresh_files", YieldShowFiles{
		Context: ctx,
	})

}

func (h *Handler) ConfirmDeleteFile(w http.ResponseWriter, r *http.Request, ctx Context) {

	b := BindFile{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	file := ctx.Publication.GetFile(b.ID)

	if file == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	render.Render(w, "publication/confirm_delete_file", YieldDeleteFile{
		Context: ctx,
		File:    file,
	})

}

func (h *Handler) SwitchEditFileByLicense(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindFile{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	triggerEl := r.Header.Get("HX-Trigger-name")
	var file *models.PublicationFile
	for i, f := range ctx.Publication.File {
		if f.ID == b.ID {
			b.Index = i
			file = f
			break
		}
	}

	if file == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Conditionally voiding licensing fields based on a triggering element
	// and dependent element inputs

	// Clear embargo fields when Access Level is set to "open access"
	if triggerEl == "" {
		if b.AccessLevel == "open_access" {
			b.Embargo = ""
			b.EmbargoTo = ""
		}
	}

	// Clear CC License when Embargo To is set to "local" or "empty"
	if triggerEl == "embargo_to" {
		if (b.EmbargoTo == "" || b.EmbargoTo == "local") && b.CCLicense != "" {
			b.CCLicense = ""
		}
	}

	// Clear No License (Copyright) when CC License is set to a non-empty value
	if triggerEl == "cc_license" {
		if b.CCLicense != "" && b.NoLicense {
			b.NoLicense = false
		}
	}

	// Clear CC License when No License (Copyright) is checked
	if triggerEl == "no_license" {
		if b.CCLicense != "" && b.NoLicense {
			b.CCLicense = ""
		}
	}

	render.Render(w, "publication/switch_edit_file_by_license", YieldEditFile{
		Context:   ctx,
		File:      file,
		FileIndex: b.Index,
		Form:      fileForm(ctx, b, nil),
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
	for i, f := range pub.File {
		if f.ID == b.ID {
			b.Index = i
			file = f
			break
		}
	}

	if file == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Embargo sanity / paranoia check
	// Ensure embargo fields are definitely empty when Acces Level is set to "Open Access"
	if b.AccessLevel == "open_access" || b.EmbargoTo == b.AccessLevel {
		b.EmbargoTo = ""
		b.Embargo = ""
	}

	// copy attributes from bind to file
	bindToPublicationFile(&b, file)

	// add thumbnails (changes record!)
	h.addThumbnailURLs(pub)

	// save publication
	err := h.Repository.SavePublication(pub)

	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {
		render.Render(w, "publication/refresh_edit_file", YieldEditFile{
			Context:   ctx,
			File:      file,
			FileIndex: b.Index,
			Form:      fileForm(ctx, b, validationErrors),
		})
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// load publication again
	pub, err = h.Repository.GetPublication(pub.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update publication in context
	ctx.Publication = pub

	// TODO: render success
	ctx.Flash = append(ctx.Flash, flash.Flash{
		Type:         "success",
		Body:         "File was updated successfully",
		DismissAfter: 5 * time.Second,
	})
	render.Render(w, "publication/refresh_files", YieldShowFiles{
		Context: ctx,
	})
}

func fileForm(ctx Context, b BindFile, errors validation.Errors) *form.Form {

	l := ctx.Locale
	optsAccessLevels := []form.SelectOption{}
	for _, acl := range vocabularies.Map["publication_file_access_levels"] {
		optsAccessLevels = append(optsAccessLevels, form.SelectOption{
			Label: l.TS("publication_file_access_levels", acl),
			Value: acl,
		})
	}

	fields := []form.Field{
		&form.Text{
			Name:  "title",
			Value: b.Title,
			Label: ctx.T("builder.file.title"),
			Cols:  12,
			Error: localize.ValidationErrorAt(
				ctx.Locale, errors,
				fmt.Sprintf("/file/%d/title", b.Index)),
		},
		&form.Select{
			Name:        "relation",
			Value:       b.Relation,
			Label:       l.T("builder.file.relation"),
			EmptyOption: true,
			Options: localize.VocabularySelectOptions(
				l,
				"publication_file_relations"),
			Cols: 12,
			Error: localize.ValidationErrorAt(
				l,
				errors,
				fmt.Sprintf("/file/%d/relation", b.Index)),
		},
		&form.Select{
			Name:        "publication_version",
			Value:       b.PublicationVersion,
			Label:       l.T("builder.file.publication_version"),
			EmptyOption: true,
			Options:     localize.VocabularySelectOptions(l, "publication_versions"),
			Cols:        12,
			Error: localize.ValidationErrorAt(
				ctx.Locale,
				errors,
				fmt.Sprintf("/file/%d/publication_version", b.Index)),
		},
		&form.RadioButtonGroup{
			Name:    "access_level",
			Value:   b.AccessLevel,
			Label:   l.T("builder.file.access_level"),
			Options: optsAccessLevels,
			Cols:    9,
			Error: localize.ValidationErrorAt(
				ctx.Locale,
				errors,
				fmt.Sprintf("/file/%d/access_level", b.Index)),
		},
	}

	// TODO: make license switch work
	// based on certain conditions add certain fields
	// e.g. do not add embargo and embargo_to if access_level == "open_access"

	fields = append(fields,
		&form.Date{
			Name:  "embargo",
			Value: b.Embargo,
			Label: l.T("builder.file.embargo"),
			Min:   time.Now().Format("2006-01-02"),
			Cols:  9,
			Error: localize.ValidationErrorAt(
				ctx.Locale,
				errors,
				fmt.Sprintf("/file/%d/embargo", b.Index)),
		},
		&form.Select{
			Name:    "embargo_to",
			Value:   b.EmbargoTo,
			Label:   l.T("builder.file.embargo_to"),
			Options: optsAccessLevels,
			Cols:    9,
			Error: localize.ValidationErrorAt(
				ctx.Locale,
				errors,
				fmt.Sprintf("/file/%d/embargo_to", b.Index)),
		})

	fields = append(fields,
		&form.Select{
			Name:        "cc_license",
			Value:       b.CCLicense,
			Label:       l.T("builder.file.cc_license"),
			EmptyOption: true,
			Options:     localize.VocabularySelectOptions(l, "cc_licenses"),
			Cols:        9,
			Error: localize.ValidationErrorAt(
				ctx.Locale,
				errors,
				fmt.Sprintf("/file/%d/cc_license", b.Index)),
		},
		&form.Checkbox{
			Template: "publication/checkbox_no_license",
			Name:     "no_license",
			Value:    "true",
			Label:    l.T("builder.file.no_license"),
			Checked:  b.NoLicense,
			Error: localize.ValidationErrorAt(
				ctx.Locale,
				errors,
				fmt.Sprintf("/file/%d/no_license", b.Index)),
			Vars: struct {
				ID     string
				FileID string
			}{
				ID:     ctx.Publication.ID,
				FileID: b.ID,
			},
		},
		&form.Text{
			Name:  "other_license",
			Value: b.OtherLicense,
			Label: l.T("builder.file.other_license"),
			Cols:  12,
			Error: localize.ValidationErrorAt(
				ctx.Locale,
				errors,
				fmt.Sprintf("/file/%d/other_license", b.Index)),
		},
		&form.TextArea{
			Name:  "description",
			Value: b.Description,
			Label: l.T("builder.file.description"),
			Rows:  4,
			Cols:  12,
			Error: localize.ValidationErrorAt(
				ctx.Locale,
				errors,
				fmt.Sprintf("/file/%d/description", b.Index)),
		})

	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(l, errors)).
		AddSection(fields...)

}

// TODO clean this up
func (h *Handler) addThumbnailURLs(p *models.Publication) {
	var u string
	for _, f := range p.File {
		if f.ContentType == "application/pdf" && f.FileSize <= 25000000 {
			params := imagorpath.Params{
				Image:  h.FileStore.RelativeFilePath(f.SHA256),
				FitIn:  true,
				Width:  156,
				Height: 156,
			}
			p := imagorpath.Generate(params, imagorpath.NewDefaultSigner(viper.GetString("imagor-secret")))
			u = viper.GetString("imagor-url") + "/" + p
			f.ThumbnailURL = u
		}
	}
}

func publicationFileToBind(f *models.PublicationFile, b *BindFile) {
	b.AccessLevel = f.AccessLevel
	b.CCLicense = f.CCLicense
	b.Description = f.Description
	b.Embargo = f.Embargo
	b.EmbargoTo = f.EmbargoTo
	b.NoLicense = f.NoLicense
	b.OtherLicense = f.OtherLicense
	b.PublicationVersion = f.PublicationVersion
	b.Relation = f.Relation
	b.Title = f.Title
	b.ID = f.ID
}

func bindToPublicationFile(b *BindFile, f *models.PublicationFile) {
	f.AccessLevel = b.AccessLevel
	f.CCLicense = b.CCLicense
	f.Description = b.Description
	f.Embargo = b.Embargo
	f.EmbargoTo = b.EmbargoTo
	f.NoLicense = b.NoLicense
	f.OtherLicense = b.OtherLicense
	f.PublicationVersion = b.PublicationVersion
	f.Relation = b.Relation
	f.Title = b.Title
	f.ID = b.ID
}
