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

type BindAbstract struct {
	AbstractID string `path:"abstract_id"`
	Text       string `form:"text"`
	Lang       string `form:"lang"`
}

type BindDeleteAbstract struct {
	AbstractID string `path:"abstract_id"`
	SnapshotID string `path:"snapshot_id"`
}

func AddAbstract(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	views.ShowModal(publicationviews.EditAbstractDialog(c, p, &models.Text{}, -1, false, nil, true)).Render(r.Context(), w)
}

func CreateAbstract(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindAbstract{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("create publication abstract: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	abstract := &models.Text{
		Lang: b.Lang,
		Text: b.Text,
	}

	p.AddAbstract(abstract)

	idx := -1
	for i, a := range p.Abstract {
		if a.ID == abstract.ID {
			idx = i
		}
	}

	if validationErrs := p.Validate(); validationErrs != nil {
		views.ReplaceModal(publicationviews.EditAbstractDialog(c, p, abstract, idx, false, validationErrs.(*okay.Errors), true)).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(publicationviews.EditAbstractDialog(c, p, abstract, idx, true, nil, true)).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("create publication abstract: could not save the publication:", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(publicationviews.AbstractsBodySelector, publicationviews.AbstractsBody(c, p)).Render(r.Context(), w)
}

func EditAbstract(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindAbstract{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("edit publication abstract: could not bind request arguments", "error", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	abstract := p.GetAbstract(b.AbstractID)

	if abstract == nil {
		c.Log.Warnf("edit publication abstract: Could not fetch the abstract:", "publication", p.ID, "abstract", b.AbstractID, "user", c.User.ID)
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	idx := -1
	for i, a := range p.Abstract {
		if a.ID == abstract.ID {
			idx = i
		}
	}

	views.ShowModal(publicationviews.EditAbstractDialog(c, p, abstract, idx, false, nil, false)).Render(r.Context(), w)
}

func UpdateAbstract(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindAbstract{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("update publication abstract: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	abstract := p.GetAbstract(b.AbstractID)

	if abstract == nil {
		c.Log.Warnw("update publication abstract: could not get abstract", "abstract", b.AbstractID, "publication", p.ID, "user", c.User.ID)
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	abstract.Text = b.Text
	abstract.Lang = b.Lang

	p.SetAbstract(abstract)

	idx := -1
	for i, a := range p.Abstract {
		if a.ID == abstract.ID {
			idx = i
		}
	}

	if validationErrs := p.Validate(); validationErrs != nil {
		views.ReplaceModal(publicationviews.EditAbstractDialog(c, p, abstract, idx, false, validationErrs.(*okay.Errors), false)).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(publicationviews.EditAbstractDialog(c, p, abstract, idx, true, nil, false)).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("update publication abstract: could not save the publication:", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(publicationviews.AbstractsBodySelector, publicationviews.AbstractsBody(c, p)).Render(r.Context(), w)
}

func ConfirmDeleteAbstract(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	var b BindDeleteAbstract
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete publication abstract: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	if b.SnapshotID != p.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to remove this abstract?",
		DeleteUrl:  c.PathTo("publication_delete_abstract", "id", p.ID, "abstract_id", b.AbstractID),
		SnapshotID: p.SnapshotID,
	}).Render(r.Context(), w)
}

func DeleteAbstract(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	var b BindDeleteAbstract
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("delete publication abstract: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	p.RemoveAbstract(b.AbstractID)

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("create publication abstract: could not save the publication:", "error", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(publicationviews.AbstractsBodySelector, publicationviews.AbstractsBody(c, p)).Render(r.Context(), w)
}
