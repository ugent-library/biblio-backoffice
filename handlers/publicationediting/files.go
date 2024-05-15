package publicationediting

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
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

func RefreshFiles(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	views.CloseModalAndReplace(publicationviews.FilesBodySelector, publicationviews.FilesBody(c, p)).Render(r.Context(), w)
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	// 2GB limit on request body
	r.Body = http.MaxBytesReader(w, r.Body, 2000000000)

	file, handler, err := r.FormFile("file")
	if err != nil {
		c.Log.Errorf("publication upload file: could not process file", "errors", err, "publication", p.ID, "user", c.User.ID)
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.file_upload_error"))).Render(r.Context(), w)
		return
	}
	defer file.Close()

	// add file to filestore
	checksum, err := c.FileStore.Add(r.Context(), file, "")

	if err != nil {
		c.Log.Errorf("publication upload file: could not save file", "errors", err, "publication", p.ID, "user", c.User.ID)
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.file_upload_error"))).Render(r.Context(), w)
		return
	}

	// save publication
	// TODO check if file with same checksum is already present
	pubFile := &models.PublicationFile{
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
	p.AddFile(pubFile)

	err = c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	views.ShowModal(publicationviews.EditFileDialog(c, p, pubFile, p.FileIndex(pubFile.ID), false, nil)).Render(r.Context(), w)
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

	views.ShowModal(publicationviews.EditFileDialog(c, p, file, p.FileIndex(file.ID), false, nil)).Render(r.Context(), w)
}

func RefreshEditFileForm(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	var b BindFile
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("edit publication file license: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	file := p.GetFile(b.FileID)

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
		views.ReplaceModal(publicationviews.EditFileDialog(c, p, file, -1, true, nil)).Render(r.Context(), w)
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

	p.SetFile(file)

	views.ReplaceModal(publicationviews.EditFileDialog(c, p, file, p.FileIndex(file.ID), false, nil)).Render(r.Context(), w)
}

func UpdateFile(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindFile{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("update publication file: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	file := p.GetFile(b.FileID)

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
		views.ReplaceModal(publicationviews.EditFileDialog(c, p, file, -1, false, nil)).Render(r.Context(), w)
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

	p.SetFile(file)

	validationErrs := p.Validate()
	// check EmbargoDate is in the future at time of submit
	if file.EmbargoDate != "" {
		t, e := time.Parse("2006-01-02", file.EmbargoDate)
		if e == nil && !t.After(time.Now()) {
			c.Log.Infof("%+v", file)
			validationErrs = okay.Add(validationErrs, okay.NewError(fmt.Sprintf("/file/%d/embargo_date", p.FileIndex(file.ID)), "publication.file.embargo_date.expired"))
		}
	}

	if validationErrs != nil {
		views.ReplaceModal(publicationviews.EditFileDialog(c, p, file, p.FileIndex(file.ID), false, validationErrs.(*okay.Errors))).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(publicationviews.EditFileDialog(c, p, file, p.FileIndex(file.ID), true, nil)).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("update publication file: could not save the publication:", "errors", err, "publication", p.ID, "user", c.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	views.CloseModalAndReplace(publicationviews.FilesBodySelector, publicationviews.FilesBody(c, p)).Render(r.Context(), w)
}

func ConfirmDeleteFile(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	var b BindDeleteFile
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete publication file: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	file := p.GetFile(b.FileID)

	if b.SnapshotID != p.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   fmt.Sprintf("Are you sure you want to remove <b>%s</b> from the publication?", file.Name),
		DeleteUrl:  c.PathTo("publication_delete_file", "id", p.ID, "file_id", file.ID),
		SnapshotID: p.SnapshotID,
	}).Render(r.Context(), w)
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	var b BindDeleteFile
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("delete publication file: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	p.RemoveFile(b.FileID)

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("delete publication file: could not save the publication:", "error", err, "publication", p.ID, "user", c.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	views.CloseModalAndReplace(publicationviews.FilesBodySelector, publicationviews.FilesBody(c, p)).Render(r.Context(), w)
}
