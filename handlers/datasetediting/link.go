package datasetediting

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/validation"
	"github.com/ugent-library/bind"
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
type YieldDeleteLink struct {
	Context
	LinkID string
}

func (h *Handler) AddLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	form := linkForm(ctx.Loc, ctx.Dataset, &models.DatasetLink{}, nil)
	render.Layout(w, "show_modal", "dataset/add_link", YieldAddLink{
		Context: ctx,
		Form:    form,
	})
}

func (h *Handler) CreateLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLink{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("add dataset link: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	datasetLink := models.DatasetLink{
		URL:         b.URL,
		Relation:    b.Relation,
		Description: b.Description,
	}
	ctx.Dataset.AddLink(&datasetLink)

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "dataset/add_link", YieldAddLink{
			Context:  ctx,
			Form:     linkForm(ctx.Loc, ctx.Dataset, &datasetLink, validationErrs.(validation.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "dataset/add_link", YieldAddLink{
			Context:  ctx,
			Form:     linkForm(ctx.Loc, ctx.Dataset, &datasetLink, nil),
			Conflict: true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("add dataset link: Could not save the dataset:", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_links", YieldLinks{
		Context: ctx,
	})
}

func (h *Handler) EditLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLink{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("edit dataset link: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	// TODO catch non-existing item in UI
	link := ctx.Dataset.GetLink(b.LinkID)
	if link == nil {
		h.Logger.Warnw("edit dataset link: could not get link", "link", b.LinkID, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.BadRequest(
			w,
			r,
			fmt.Errorf("no link found for %s in dataset %s", b.LinkID, ctx.Dataset.ID),
		)
		return
	}

	render.Layout(w, "show_modal", "dataset/edit_link", YieldEditLink{
		Context:  ctx,
		LinkID:   b.LinkID,
		Form:     linkForm(ctx.Loc, ctx.Dataset, link, nil),
		Conflict: false,
	})
}

func (h *Handler) UpdateLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLink{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update dataset link: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	link := ctx.Dataset.GetLink(b.LinkID)
	if link == nil {
		h.Logger.Warnw("update dataset link: could not get link", "link", b.LinkID, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("dataset.conflict_error_reload"),
		})
		return
	}

	link.URL = b.URL
	link.Description = b.Description
	link.Relation = b.Relation

	ctx.Dataset.SetLink(link)

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "dataset/edit_link", YieldEditLink{
			Context:  ctx,
			LinkID:   b.LinkID,
			Form:     linkForm(ctx.Loc, ctx.Dataset, link, validationErrs.(validation.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "dataset/edit_link", YieldEditLink{
			Context:  ctx,
			LinkID:   b.LinkID,
			Form:     linkForm(ctx.Loc, ctx.Dataset, link, nil),
			Conflict: true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("update dataset link: Could not save the dataset:", "errors", err, "identifier", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_links", YieldLinks{
		Context: ctx,
	})
}

func (h *Handler) ConfirmDeleteLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteLink
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Errorw("confirm delete dataset link: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	// TODO catch non-existing item in UI
	if b.SnapshotID != ctx.Dataset.SnapshotID {
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("dataset.conflict_error_reload"),
		})
		return
	}

	render.Layout(w, "show_modal", "dataset/confirm_delete_link", YieldDeleteLink{
		Context: ctx,
		LinkID:  b.LinkID,
	})
}

func (h *Handler) DeleteLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteLink
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete dataset link: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	ctx.Dataset.RemoveLink(b.LinkID)

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("dataset.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("delete dataset link: Could not save the dataset:", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_links", YieldLinks{
		Context: ctx,
	})
}

func linkForm(loc *gotext.Locale, dataset *models.Dataset, link *models.DatasetLink, errors validation.Errors) *form.Form {
	idx := -1
	for i, l := range dataset.Link {
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
				Options: localize.VocabularySelectOptions(loc, "dataset_link_relations"),
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
