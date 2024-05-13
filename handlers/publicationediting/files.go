package publicationediting

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	publicationviews "github.com/ugent-library/biblio-backoffice/views/publication"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

type BindFile struct {
	FileID                   string `path:"file_id"`
	AccessLevel              string `query:"access_level" form:"access_level"`
	License                  string `query:"license" form:"license"`
	ContentType              string `query:"content_type" form:"content_type"`
	EmbargoDate              string `query:"embargo_date" form:"embargo_date"`
	AccessLevelAfterEmbargo  string `query:"access_level_after_embargo" form:"access_level_after_embargo"`
	AccessLevelDuringEmbargo string `query:"access_level_during_embargo" form:"access_level_during_embargo"`
	Name                     string `query:"name" form:"name"`
	Size                     int    `query:"size" form:"size"`
	SHA256                   string `query:"sha256" form:"sha256"`
	OtherLicense             string `query:"other_license" form:"other_license"`
	PublicationVersion       string `query:"publication_version" form:"publication_version"`
	Relation                 string `query:"relation" form:"relation"`
	URL                      string `query:"url" form:"url"`
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
	MaxFileSize int
}

func (h *Handler) RefreshFiles(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.View(w, "publication/refresh_files", YieldShowFiles{
		Context:     ctx,
		MaxFileSize: h.MaxFileSize,
	})
}

func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request, ctx Context) {
	// 2GB limit on request body
	r.Body = http.MaxBytesReader(w, r.Body, 2000000000)

	file, handler, err := r.FormFile("file")
	if err != nil {
		h.Logger.Errorf("publication upload file: could not process file", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		views.ShowModal(views.ErrorDialog(ctx.Loc.Get("publication.file_upload_error"))).Render(r.Context(), w)
		return
	}
	defer file.Close()

	// add file to filestore
	checksum, err := h.FileStore.Add(r.Context(), file, "")

	if err != nil {
		h.Logger.Errorf("publication upload file: could not save file", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		views.ShowModal(views.ErrorDialog(ctx.Loc.Get("publication.file_upload_error"))).Render(r.Context(), w)
		return
	}

	// save publication
	// TODO check if file with same checksum is already present
	pubFile := models.PublicationFile{
		Relation:    "main_file",
		AccessLevel: "info:eu-repo/semantics/restrictedAccess",
		Name:        handler.Filename,
		Size:        int(handler.Size),
		ContentType: handler.Header.Get("Content-Type"),
		SHA256:      checksum,
	}
	/*
		automatically generates extra fields:
		id, date_created, date_updated
	*/
	ctx.Publication.AddFile(&pubFile)

	err = h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ShowModal(views.ErrorDialog(ctx.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	// render edit file form
	render.Layout(w, "show_modal", "publication/edit_file", YieldEditFile{
		Context: ctx,
		File:    &pubFile,
		Form:    fileForm(ctx.Loc, ctx.Publication, &pubFile, nil),
	})

}

func EditFile(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	var b BindFile
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	file := p.GetFile(b.FileID)

	if file == nil {
		c.Log.Warnw("publication upload file: could not find file", "fileid", b.FileID, "publication", p.ID, "user", c.User.ID)
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	// Set Relation to default "main_file" if absent in older records.
	if file.Relation == "" {
		file.Relation = "main_file"
	}

	idx := -1
	for i, f := range p.File {
		if f.ID == file.ID {
			idx = i
			break
		}
	}

	views.ShowModal(publicationviews.EditFileDialog(c, p, file, idx, false, nil)).Render(r.Context(), w)
}

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
			Form:     fileForm(ctx.Loc, ctx.Publication, file, nil),
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
		Form:    fileForm(ctx.Loc, ctx.Publication, file, nil),
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
			Form:     fileForm(ctx.Loc, ctx.Publication, file, nil),
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

	validationErrs := ctx.Publication.Validate()
	// check EmbargoDate is in the future at time of submit
	if file.EmbargoDate != "" {
		t, e := time.Parse("2006-01-02", file.EmbargoDate)
		if e == nil && !t.After(time.Now()) {
			validationErrs = okay.Add(validationErrs, okay.NewError(fmt.Sprintf("/file/%d/embargo_date", ctx.Publication.FileIndex(file.ID)), "publication.file.embargo_date.expired"))
		}
	}

	if validationErrs != nil {
		render.Layout(w, "refresh_modal", "publication/edit_file", YieldEditFile{
			Context:  ctx,
			File:     file,
			Form:     fileForm(ctx.Loc, ctx.Publication, file, validationErrs.(*okay.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "publication/edit_file", YieldEditFile{
			Context:  ctx,
			File:     file,
			Form:     fileForm(ctx.Loc, ctx.Publication, file, nil),
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
		Context:     ctx,
		MaxFileSize: h.MaxFileSize,
	})
}

func ConfirmDeleteFile(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	var b BindDeleteFile
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete publication file: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	file := publication.GetFile(b.FileID)

	if b.SnapshotID != publication.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   fmt.Sprintf("Are you sure you want to remove <b>%s</b> from the publication?", file.Name),
		DeleteUrl:  c.PathTo("publication_delete_file", "id", publication.ID, "file_id", file.ID),
		SnapshotID: publication.SnapshotID,
	}).Render(r.Context(), w)
}

func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteFile
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete publication file: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	ctx.Publication.RemoveFile(b.FileID)

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(ctx.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		h.Logger.Errorf("delete publication file: could not save the publication:", "error", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_files", YieldShowFiles{
		Context:     ctx,
		MaxFileSize: h.MaxFileSize,
	})
}

func fileForm(loc *gotext.Locale, publication *models.Publication, file *models.PublicationFile, errors *okay.Errors) *form.Form {
	idx := -1
	for i, f := range publication.File {
		if f.ID == file.ID {
			idx = i
			break
		}
	}

	f := form.New().
		WithTheme("file").
		WithErrors(localize.ValidationErrors(loc, errors))

	if file.Relation == "main_file" {
		f.AddTemplatedSection(
			"sections/document_type",
			struct{}{},
			&form.Select{
				Template: "document_type",
				Name:     "relation",
				Value:    file.Relation,
				Label:    loc.Get("builder.file.relation"),
				Options: localize.VocabularySelectOptions(
					loc,
					"publication_file_relations"),
				Error: localize.ValidationErrorAt(
					loc,
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
				Label:       loc.Get("builder.file.publication_version"),
				EmptyOption: true,
				Options:     localize.VocabularySelectOptions(loc, "publication_versions"),
				Error: localize.ValidationErrorAt(
					loc,
					errors,
					fmt.Sprintf("/file/%d/publication_version", idx)),
			},
		)
	} else {
		f.AddTemplatedSection(
			"sections/document_type",
			struct{}{},
			&form.Select{
				Template: "document_type",
				Name:     "relation",
				Value:    file.Relation,
				Label:    loc.Get("builder.file.relation"),
				Options: localize.VocabularySelectOptions(
					loc,
					"publication_file_relations"),
				Error: localize.ValidationErrorAt(
					loc,
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
			Label:    loc.Get("builder.file.access_level"),
			Options:  localize.VocabularySelectOptions(loc, "publication_file_access_levels"),
			Cols:     9,
			Error: localize.ValidationErrorAt(
				loc,
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
				Options:     localize.VocabularySelectOptions(loc, "publication_file_access_levels_during_embargo"),
				Error: localize.ValidationErrorAt(
					loc,
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
				Options:     localize.VocabularySelectOptions(loc, "publication_file_access_levels_after_embargo"),
				Error: localize.ValidationErrorAt(
					loc,
					errors,
					fmt.Sprintf("/file/%d/access_level_after_embargo", idx)),
			},
		)

		now := time.Now()
		nextDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)

		f.AddSection(
			&form.Date{
				Template: "embargo_end",
				Name:     "embargo_date",
				Value:    file.EmbargoDate,
				Label:    loc.Get("builder.file.embargo_date"),
				Min:      nextDay.Format("2006-01-02"),
				Error: localize.ValidationErrorAt(
					loc,
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
			Label:       loc.Get("builder.file.license"),
			Tooltip:     loc.Get("tooltip.publication.file.license"),
			EmptyOption: true,
			Options:     localize.VocabularySelectOptions(loc, "publication_licenses"),
		},
	)

	return f
}
