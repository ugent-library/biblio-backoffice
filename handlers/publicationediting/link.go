package publicationediting

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	views "github.com/ugent-library/biblio-backoffice/views/publication"
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

type YieldLinks struct {
	Context
}
type YieldAddLink struct {
	Context
	Form     *form.Form
	Conflict bool
}
type YieldEditLink struct {
	Context
	LinkID   string
	Form     *form.Form
	Conflict bool
}

func (h *Handler) AddLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	form := linkForm(ctx.Loc, ctx.Publication, &models.PublicationLink{}, nil)
	render.Layout(w, "show_modal", "publication/add_link", YieldAddLink{
		Context: ctx,
		Form:    form,
	})
}

func (h *Handler) CreateLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLink{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("add publication link: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	publicationLink := models.PublicationLink{
		URL:         b.URL,
		Relation:    b.Relation,
		Description: b.Description,
	}
	ctx.Publication.AddLink(&publicationLink)

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "publication/add_link", YieldAddLink{
			Context:  ctx,
			Form:     linkForm(ctx.Loc, ctx.Publication, &publicationLink, validationErrs.(*okay.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "publication/add_link", YieldAddLink{
			Context:  ctx,
			Form:     linkForm(ctx.Loc, ctx.Publication, &publicationLink, nil),
			Conflict: true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("add publication link: Could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_links", YieldLinks{
		Context: ctx,
	})
}

func (h *Handler) EditLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLink{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("edit publication link: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	// TODO catch non-existing item in UI
	link := ctx.Publication.GetLink(b.LinkID)
	if link == nil {
		h.Logger.Warnw("edit publication link: could not get link", "link", b.LinkID, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.BadRequest(
			w,
			r,
			fmt.Errorf("no link found for %s in publication %s", b.LinkID, ctx.Publication.ID),
		)
		return
	}

	render.Layout(w, "show_modal", "publication/edit_link", YieldEditLink{
		Context:  ctx,
		LinkID:   b.LinkID,
		Form:     linkForm(ctx.Loc, ctx.Publication, link, nil),
		Conflict: false,
	})
}

func (h *Handler) UpdateLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLink{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update publication link: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	link := ctx.Publication.GetLink(b.LinkID)
	if link == nil {
		h.Logger.Warnw("update publication link: could not get link", "link", b.LinkID, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("publication.conflict_error_reload"),
		})
		return
	}

	link.URL = b.URL
	link.Description = b.Description
	link.Relation = b.Relation

	ctx.Publication.SetLink(link)

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "publication/edit_link", YieldEditLink{
			Context:  ctx,
			LinkID:   b.LinkID,
			Form:     linkForm(ctx.Loc, ctx.Publication, link, validationErrs.(*okay.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "publication/edit_link", YieldEditLink{
			Context:  ctx,
			LinkID:   b.LinkID,
			Form:     linkForm(ctx.Loc, ctx.Publication, link, nil),
			Conflict: true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("update publication link: Could not save the publication:", "errors", err, "identifier", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_links", YieldLinks{
		Context: ctx,
	})
}

func ConfirmDeleteLink(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	var b BindDeleteLink
	if err := bind.Request(r, &b); err != nil {
		c.Log.Errorw("confirm delete publication link: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	// TODO catch non-existing item in UI
	if b.SnapshotID != publication.SnapshotID {
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: c.Loc.Get("publication.conflict_error_reload"),
		})
		return
	}

	views.ConfirmDeleteLink(c, publication, b.LinkID).Render(r.Context(), w)
}

func (h *Handler) DeleteLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteLink
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete publication link: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	ctx.Publication.RemoveLink(b.LinkID)

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("publication.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("delete publication link: Could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_links", YieldLinks{
		Context: ctx,
	})
}

func linkForm(loc *gotext.Locale, publication *models.Publication, link *models.PublicationLink, errors *okay.Errors) *form.Form {
	idx := -1
	for i, l := range publication.Link {
		if l.ID == link.ID {
			idx = i
			break
		}
	}
	return form.New().
		WithTheme("cols").
		WithErrors(localize.ValidationErrors(loc, errors)).
		AddSection(
			&form.Text{
				Name:     "url",
				Value:    link.URL,
				Label:    loc.Get("builder.link.url"),
				Required: true,
				Cols:     12,
				Error: localize.ValidationErrorAt(
					loc,
					errors,
					fmt.Sprintf("/link/%d/url", idx),
				),
			},
			&form.Select{
				Name:    "relation",
				Value:   link.Relation,
				Label:   loc.Get("builder.link.relation"),
				Options: localize.VocabularySelectOptions(loc, "publication_link_relations"),
				Cols:    12,
				Error: localize.ValidationErrorAt(
					loc,
					errors,
					fmt.Sprintf("/link/%d/relation", idx),
				),
			},
			&form.Text{
				Name:  "description",
				Value: link.Description,
				Label: loc.Get("builder.link.description"),
				Cols:  12,
				Error: localize.ValidationErrorAt(
					loc,
					errors,
					fmt.Sprintf("/link/%d/description", idx),
				),
			},
		)
}
