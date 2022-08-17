package publicationediting

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type BindLink struct {
	LinkID      string `path:"link_id"`
	URL         string `form:"url"`
	Relation    string `form:"relation"`
	Description string `form:"description"`
}

type BindDeleteLink struct {
	LinkID string `path:"link_id"`
}

type YieldLinks struct {
	Context
}
type YieldAddLink struct {
	Context
	Form *form.Form
}
type YieldEditLink struct {
	Context
	LinkID string
	Form   *form.Form
}
type YieldDeleteLink struct {
	Context
	LinkID string
}

func (h *Handler) AddLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	form := linkForm(ctx.Locale, ctx.Publication, &models.PublicationLink{}, nil)
	render.Layout(w, "show_modal", "publication/add_link", YieldAddLink{
		Context: ctx,
		Form:    form,
	})
}

func (h *Handler) CreateLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLink{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
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
			Context: ctx,
			Form:    linkForm(ctx.Locale, ctx.Publication, &publicationLink, validationErrs.(validation.Errors)),
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
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
		render.BadRequest(w, r, err)
		return
	}

	link := ctx.Publication.GetLink(b.LinkID)
	if link == nil {
		render.BadRequest(
			w,
			r,
			fmt.Errorf("no link found for %s in publication %s", b.LinkID, ctx.Publication.ID),
		)
		return
	}

	render.Layout(w, "show_modal", "publication/edit_link", YieldEditLink{
		Context: ctx,
		LinkID:  b.LinkID,
		Form:    linkForm(ctx.Locale, ctx.Publication, link, nil),
	})
}

func (h *Handler) UpdateLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLink{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	/*
		TODO: throw a conflict error when
		trying to update a non existing id?
	*/
	link := ctx.Publication.GetLink(b.LinkID)
	if link == nil {
		render.BadRequest(
			w,
			r,
			fmt.Errorf("no link found for %s in publication %s", b.LinkID, ctx.Publication.ID),
		)
		return
	}
	link.URL = b.URL
	link.Description = b.Description
	link.Relation = b.Relation

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		form := linkForm(ctx.Locale, ctx.Publication, link, validationErrs.(validation.Errors))

		render.Layout(w, "refresh_modal", "publication/edit_link", YieldEditLink{
			Context: ctx,
			LinkID:  b.LinkID,
			Form:    form,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_links", YieldLinks{
		Context: ctx,
	})
}

func (h *Handler) ConfirmDeleteLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteLink
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.Layout(w, "show_modal", "publication/confirm_delete_link", YieldDeleteLink{
		Context: ctx,
		LinkID:  b.LinkID,
	})
}

func (h *Handler) DeleteLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteLink
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	/*
		Note: link possibly already removed:
		conflict resolving will solve this
	*/
	ctx.Publication.RemoveLink(b.LinkID)

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_links", YieldLinks{
		Context: ctx,
	})
}

func linkForm(l *locale.Locale, publication *models.Publication, link *models.PublicationLink, errors validation.Errors) *form.Form {
	idx := -1
	for i, l := range publication.Link {
		if l.ID == link.ID {
			idx = i
			break
		}
	}
	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(l, errors)).
		AddSection(
			&form.Text{
				Name:  "url",
				Value: link.URL,
				Label: l.T("builder.link.url"),
				Cols:  12,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					fmt.Sprintf("/link/%d/url", idx),
				),
			},
			&form.Select{
				Name:    "relation",
				Value:   link.Relation,
				Label:   l.T("builder.link.relation"),
				Options: localize.VocabularySelectOptions(l, "publication_link_relations"),
				Cols:    12,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					fmt.Sprintf("/link/%d/relation", idx),
				),
			},
			&form.Text{
				Name:  "description",
				Value: link.Description,
				Label: l.T("builder.link.description"),
				Cols:  12,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					fmt.Sprintf("/link/%d/description", idx),
				),
			},
		)
}
