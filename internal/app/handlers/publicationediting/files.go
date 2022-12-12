package publicationediting

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type BindFile struct {
	FileID                   string `path:"file_id"`
	AccessLevel              string `form:"access_level"`
	License                  string `form:"license"`
	ContentType              string `form:"content_type"`
	EmbargoDate              string `form:"embargo_date"`
	AccessLevelAfterEmbargo  string `form:"access_level_after_embargo"`
	AccessLevelDuringEmbargo string `form:"access_level_during_embargo"`
	Name                     string `form:"name"`
	Size                     int    `form:"size"`
	SHA256                   string `form:"sha256"`
	OtherLicense             string `form:"other_license"`
	PublicationVersion       string `form:"publication_version"`
	Relation                 string `form:"relation"`
	URL                      string `form:"url"`
}

type BindDeleteFile struct {
	Context
	FileID     string `path:"file_id"`
	SnapshotID string `path:"snapshot_id"`
	Name       string `form:"name"`
	Conflict   bool
}

type YieldEditFile struct {
	Context
	File     *models.PublicationFile
	Form     *form.Form
	Conflict bool
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

	file, handler, err := r.FormFile("file")
	if err != nil {
		h.Logger.Errorf("publication upload file: could not process file", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Locale.T("publication.file_upload_error"),
		})
		return
	}
	defer file.Close()

	// detect content type
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		h.Logger.Errorf("publication upload file: could not read file", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Locale.T("publication.file_upload_error"),
		})
		return
	}
	filetype := http.DetectContentType(buff)

	// rewind
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		h.Logger.Errorf("publication upload file: could not read file", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Locale.T("publication.file_upload_error"),
		})
		return
	}

	// add file to filestore
	checksum, err := h.FileStore.Add(file)

	if err != nil {
		h.Logger.Errorf("publication upload file: could not save file", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Locale.T("publication.file_upload_error"),
		})
		return
	}

	// save publication
	// TODO check if file with same checksum is already present
	pubFile := models.PublicationFile{
		AccessLevel: "info:eu-repo/semantics/restrictedAccess",
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

	err = h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Locale.T("publication.conflict_error_reload"),
		})
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
	var b BindFile
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	file := ctx.Publication.GetFile(b.FileID)

	if file == nil {
		h.Logger.Warnw("publication upload file: could not find file", "fileid", b.FileID, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Locale.T("publication.conflict_error_reload"),
		})
		return
	}

	render.Layout(w, "show_modal", "publication/edit_file", YieldEditFile{
		Context:  ctx,
		File:     file,
		Form:     fileForm(ctx.Locale, ctx.Publication, file, nil),
		Conflict: false,
	})
}

// TODO add more rules
func (h *Handler) RefreshEditFileForm(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindFile
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("edit publication file license: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	file := ctx.Publication.GetFile(b.FileID)

	if file == nil {
		file := &models.PublicationFile{
			AccessLevel:              b.AccessLevel,
			License:                  b.License,
			EmbargoDate:              b.EmbargoDate,
			AccessLevelAfterEmbargo:  b.AccessLevelAfterEmbargo,
			AccessLevelDuringEmbargo: b.AccessLevelDuringEmbargo,
			OtherLicense:             b.OtherLicense,
			PublicationVersion:       b.PublicationVersion,
			Relation:                 b.Relation,
		}
		render.Layout(w, "refresh_modal", "publication/edit_file", YieldEditFile{
			Context:  ctx,
			File:     file,
			Form:     fileForm(ctx.Locale, ctx.Publication, file, nil),
			Conflict: true,
		})
		return
	}

	// clear embargo fields when access level is set to anything else
	if b.AccessLevel != "info:eu-repo/semantics/embargoedAccess" {
		b.EmbargoDate = ""
		b.AccessLevelAfterEmbargo = ""
		b.AccessLevelDuringEmbargo = ""
	}

	// TODO apply other license && access level related rules, if any

	// Copy everything
	file.AccessLevel = b.AccessLevel
	file.License = b.License
	file.EmbargoDate = b.EmbargoDate
	file.AccessLevelAfterEmbargo = b.AccessLevelAfterEmbargo
	file.AccessLevelDuringEmbargo = b.AccessLevelDuringEmbargo
	file.OtherLicense = b.OtherLicense
	file.PublicationVersion = b.PublicationVersion
	file.Relation = b.Relation

	ctx.Publication.SetFile(file)

	render.Layout(w, "refresh_modal", "publication/edit_file", YieldEditFile{
		Context: ctx,
		File:    file,
		Form:    fileForm(ctx.Locale, ctx.Publication, file, nil),
	})
}

func (h *Handler) UpdateFile(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindFile{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("update publication file: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	file := ctx.Publication.GetFile(b.FileID)

	// TODO catch non-existing item in UI (show message)
	if file == nil {
		file := &models.PublicationFile{
			AccessLevel:              b.AccessLevel,
			License:                  b.License,
			EmbargoDate:              b.EmbargoDate,
			AccessLevelAfterEmbargo:  b.AccessLevelAfterEmbargo,
			AccessLevelDuringEmbargo: b.AccessLevelDuringEmbargo,
			OtherLicense:             b.OtherLicense,
			PublicationVersion:       b.PublicationVersion,
			Relation:                 b.Relation,
		}
		render.Layout(w, "refresh_modal", "publication/edit_file", YieldEditFile{
			Context:  ctx,
			File:     file,
			Form:     fileForm(ctx.Locale, ctx.Publication, file, nil),
			Conflict: false,
		})
		return
	}

	file.AccessLevel = b.AccessLevel
	file.License = b.License
	file.EmbargoDate = b.EmbargoDate
	file.AccessLevelAfterEmbargo = b.AccessLevelAfterEmbargo
	file.AccessLevelDuringEmbargo = b.AccessLevelDuringEmbargo
	file.OtherLicense = b.OtherLicense
	file.PublicationVersion = b.PublicationVersion
	file.Relation = b.Relation

	ctx.Publication.SetFile(file)

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "publication/edit_file", YieldEditFile{
			Context:  ctx,
			File:     file,
			Form:     fileForm(ctx.Locale, ctx.Publication, file, validationErrs.(validation.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "publication/edit_file", YieldEditFile{
			Context:  ctx,
			File:     file,
			Form:     fileForm(ctx.Locale, ctx.Publication, file, nil),
			Conflict: true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("update publication file: could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_files", YieldShowFiles{
		Context: ctx,
	})
}

func (h *Handler) ConfirmDeleteFile(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteFile
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("confirm delete publication file: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	file := ctx.Publication.GetFile(b.FileID)

	if b.SnapshotID != ctx.Publication.SnapshotID {
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Locale.T("publication.conflict_error_reload"),
		})
		return
	}

	render.Layout(w, "show_modal", "publication/confirm_delete_file", YieldDeleteFile{
		Context: ctx,
		File:    file,
	})
}

func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteFile
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete publication file: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	ctx.Publication.RemoveFile(b.FileID)

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Locale.T("publication.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("delete publication file: could not save the publication:", "error", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

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
		WithTheme("file").
		WithErrors(localize.ValidationErrors(l, errors))

	if file.Relation == "main_file" {
		f.AddTemplatedSection(
			"sections/document_type",
			struct{}{},
			&form.Select{
				Template:    "document_type",
				Name:        "relation",
				Value:       file.Relation,
				Label:       l.T("builder.file.relation"),
				EmptyOption: true,
				Options: localize.VocabularySelectOptions(
					l,
					"publication_file_relations"),
				Error: localize.ValidationErrorAt(
					l,
					errors,
					fmt.Sprintf("/file/%d/relation", idx)),
				Vars: struct {
					ID     string
					FileID string
				}{
					ID:     publication.ID,
					FileID: file.ID,
				},
			},
			&form.Select{
				Template:    "publication_version",
				Name:        "publication_version",
				Value:       file.PublicationVersion,
				Label:       l.T("builder.file.publication_version"),
				EmptyOption: true,
				Options:     localize.VocabularySelectOptions(l, "publication_versions"),
				Error: localize.ValidationErrorAt(
					l,
					errors,
					fmt.Sprintf("/file/%d/publication_version", idx)),
			},
		)
	} else {
		f.AddTemplatedSection(
			"sections/document_type",
			struct{}{},
			&form.Select{
				Template:    "document_type",
				Name:        "relation",
				Value:       file.Relation,
				Label:       l.T("builder.file.relation"),
				EmptyOption: true,
				Options: localize.VocabularySelectOptions(
					l,
					"publication_file_relations"),
				Error: localize.ValidationErrorAt(
					l,
					errors,
					fmt.Sprintf("/file/%d/relation", idx)),
				Vars: struct {
					ID     string
					FileID string
				}{
					ID:     publication.ID,
					FileID: file.ID,
				},
			})
	}

	// Calculate access level for embargo for the "Access Level" field
	// The "Embargo" field carries the end date of the embargo. Having a value implicitly
	// signals that the file has an "embargo access". The "embargo" option in the "Access
	// level" field needs to be active in the display. We need this because "embargo"
	// isn't a status value we store on the level of the data.
	// accessLevel := file.AccessLevel
	// if file.EmbargoDate != "" {
	// 	accessLevel = "embargo"
	// }

	f.AddTemplatedSection(
		"sections/access_level",
		struct {
			Relation string
		}{
			Relation: file.Relation,
		},
		&form.RadioButtonGroup{
			Template: "file_access_level",
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

	if file.AccessLevel == "info:eu-repo/semantics/embargoedAccess" {
		f.AddTemplatedSection(
			"sections/embargo",
			struct{}{},
			&form.Select{
				Template: "embargo_during",
				Name:     "access_level_during_embargo",
				Value:    file.AccessLevelDuringEmbargo,
				// TODO html in l.T is transformed into html entities
				// Label:       l.T("builder.file.embargo_during"),
				EmptyOption: true,
				Options:     localize.VocabularySelectOptions(l, "publication_file_access_levels_during_embargo"),
				Error: localize.ValidationErrorAt(
					l,
					errors,
					fmt.Sprintf("/file/%d/access_level_during_embargo", idx)),
			},
			&form.Select{
				Template: "embargo_after",
				Name:     "access_level_after_embargo",
				Value:    file.AccessLevelAfterEmbargo,
				// TODO html in l.T is transformed into html entities
				// Label:       l.T("builder.file.embargo_after"),
				EmptyOption: true,
				Options:     localize.VocabularySelectOptions(l, "publication_file_access_levels_after_embargo"),
				Error: localize.ValidationErrorAt(
					l,
					errors,
					fmt.Sprintf("/file/%d/access_level_after_embargo", idx)),
			},
		)

		f.AddSection(
			&form.Date{
				Template: "embargo_end",
				Name:     "embargo_date",
				Value:    file.EmbargoDate,
				Label:    l.T("builder.file.embargo_date"),
				Min:      time.Now().Format("2006-01-02"),
				Error: localize.ValidationErrorAt(
					l,
					errors,
					fmt.Sprintf("/file/%d/embargo_date", idx)),
			},
		)
	}

	f.AddTemplatedSection(
		"sections/license",
		struct{}{},
		&form.Select{
			Template:    "license",
			Name:        "license",
			Value:       file.License,
			Label:       l.T("builder.file.license"),
			Tooltip:     l.T("tooltip.publication.file.license"),
			EmptyOption: true,
			Options:     localize.VocabularySelectOptions(l, "publication_licenses"),
		},
	)

	return f
}
