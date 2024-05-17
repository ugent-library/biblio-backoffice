package publicationediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	publicationviews "github.com/ugent-library/biblio-backoffice/views/publication"
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
	p := ctx.GetPublication(r)

	views.ShowModal(publicationviews.EditLinkDialog(c, p, &models.PublicationLink{}, -1, false, nil, true)).Render(r.Context(), w)
}

func CreateLink(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindLink{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("add publication link: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	link := &models.PublicationLink{
		URL:         b.URL,
		Relation:    b.Relation,
		Description: b.Description,
	}
	p.AddLink(link)

	idx := -1
	for i, a := range p.Link {
		if a.ID == link.ID {
			idx = i
		}
	}

	if validationErrs := p.Validate(); validationErrs != nil {
		views.ReplaceModal(publicationviews.EditLinkDialog(c, p, link, idx, false, validationErrs.(*okay.Errors), true)).Render(r.Context(), w)
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(publicationviews.EditLinkDialog(c, p, link, idx, true, nil, true)).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("add publication link: Could not save the publication:", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(publicationviews.LinksBodySelector, publicationviews.LinksBody(c, p)).Render(r.Context(), w)
}

func EditLink(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindLink{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("edit publication link: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	link := p.GetLink(b.LinkID)
	if link == nil {
		c.Log.Warnw("edit publication link: could not get link", "link", b.LinkID, "publication", p.ID, "user", c.User.ID)
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	idx := -1
	for i, a := range p.Link {
		if a.ID == link.ID {
			idx = i
		}
	}

	views.ShowModal(publicationviews.EditLinkDialog(c, p, link, idx, false, nil, false)).Render(r.Context(), w)
}

func UpdateLink(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindLink{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("update publication link: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	link := p.GetLink(b.LinkID)
	if link == nil {
		c.Log.Warnw("update publication link: could not get link", "link", b.LinkID, "publication", p.ID, "user", c.User.ID)
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	link.URL = b.URL
	link.Description = b.Description
	link.Relation = b.Relation

	p.SetLink(link)

	idx := -1
	for i, a := range p.Link {
		if a.ID == link.ID {
			idx = i
		}
	}

	if validationErrs := p.Validate(); validationErrs != nil {
		views.ReplaceModal(publicationviews.EditLinkDialog(c, p, link, idx, false, validationErrs.(*okay.Errors), false)).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(publicationviews.EditLinkDialog(c, p, link, idx, true, nil, false)).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("update publication link: Could not save the publication:", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(publicationviews.LinksBodySelector, publicationviews.LinksBody(c, p)).Render(r.Context(), w)
}

func ConfirmDeleteLink(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	var b BindDeleteLink
	if err := bind.Request(r, &b); err != nil {
		c.Log.Errorw("confirm delete publication link: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	// TODO catch non-existing item in UI
	if b.SnapshotID != p.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to remove this link?",
		DeleteUrl:  c.PathTo("publication_delete_link", "id", p.ID, "link_id", b.LinkID),
		SnapshotID: p.SnapshotID,
	}).Render(r.Context(), w)
}

func DeleteLink(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	var b BindDeleteLink
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("delete publication link: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	p.RemoveLink(b.LinkID)

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("delete publication link: Could not save the publication:", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(publicationviews.LinksBodySelector, publicationviews.LinksBody(c, p)).Render(r.Context(), w)
}
