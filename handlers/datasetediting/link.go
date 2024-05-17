package datasetediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	datasetviews "github.com/ugent-library/biblio-backoffice/views/dataset"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

type BindLink struct {
	LinkID      string `path:"link_id"`
	URL         string `form:"url"`
	Relation    string `form:"relation"`
	Description string `form:"description"`
}

type BindDeleteLink struct {
	LinkID     string `path:"link_id"`
	SnapshotID string `path:"snapshot_id"`
}

func AddLink(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)
	views.ShowModal(datasetviews.EditLinkDialog(c, d, &models.DatasetLink{}, -1, false, nil, true)).Render(r.Context(), w)
}

func CreateLink(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	b := BindLink{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("add dataset link: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	link := &models.DatasetLink{
		URL:         b.URL,
		Relation:    b.Relation,
		Description: b.Description,
	}
	d.AddLink(link)

	idx := -1
	for i, l := range d.Link {
		if l.ID == link.ID {
			idx = i
		}
	}

	if validationErrs := d.Validate(); validationErrs != nil {
		views.ReplaceModal(datasetviews.EditLinkDialog(c, d, link, idx, false, validationErrs.(*okay.Errors), true)).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), d, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(datasetviews.EditLinkDialog(c, d, link, idx, true, nil, true)).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("add dataset link: Could not save the dataset:", "errors", err, "dataset", d.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace("#links-body", datasetviews.LinksBody(c, d)).Render(r.Context(), w)
}

func EditLink(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	b := BindLink{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("edit dataset link: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	// TODO catch non-existing item in UI
	link := d.GetLink(b.LinkID)
	if link == nil {
		c.Log.Warnw("edit dataset link: could not get link", "link", b.LinkID, "dataset", d.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	idx := -1
	for i, l := range d.Link {
		if l.ID == link.ID {
			idx = i
		}
	}

	views.ShowModal(datasetviews.EditLinkDialog(c, d, link, idx, false, nil, false)).Render(r.Context(), w)
}

func UpdateLink(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	b := BindLink{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("update dataset link: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	link := d.GetLink(b.LinkID)
	if link == nil {
		c.Log.Warnw("update dataset link: could not get link", "link", b.LinkID, "dataset", d.ID, "user", c.User.ID)
		views.ShowModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	link.URL = b.URL
	link.Description = b.Description
	link.Relation = b.Relation

	d.SetLink(link)

	idx := -1
	for i, l := range d.Link {
		if l.ID == link.ID {
			idx = i
		}
	}

	if validationErrs := d.Validate(); validationErrs != nil {
		views.ReplaceModal(datasetviews.EditLinkDialog(c, d, link, idx, false, validationErrs.(*okay.Errors), false)).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), d, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(datasetviews.EditLinkDialog(c, d, link, idx, true, nil, false)).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("update dataset link: Could not save the dataset:", "errors", err, "identifier", d.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace("#links-body", datasetviews.LinksBody(c, d)).Render(r.Context(), w)
}

func ConfirmDeleteLink(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	var b BindDeleteLink
	if err := bind.Request(r, &b); err != nil {
		c.Log.Errorw("confirm delete dataset link: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	// TODO catch non-existing item in UI
	if b.SnapshotID != d.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to remove this link?",
		DeleteUrl:  c.PathTo("dataset_delete_link", "id", d.ID, "link_id", b.LinkID),
		SnapshotID: d.SnapshotID,
	}).Render(r.Context(), w)
}

func DeleteLink(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	var b BindDeleteLink
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("delete dataset link: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	d.RemoveLink(b.LinkID)

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), d, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("delete dataset link: Could not save the dataset:", "errors", err, "dataset", d.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace("#links-body", datasetviews.LinksBody(c, d)).Render(r.Context(), w)
}
