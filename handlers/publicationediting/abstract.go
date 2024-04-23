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

type BindAbstract struct {
	AbstractID string `path:"abstract_id"`
	Text       string `form:"text"`
	Lang       string `form:"lang"`
}

type BindDeleteAbstract struct {
	AbstractID string `path:"abstract_id"`
	SnapshotID string `path:"snapshot_id"`
}

type YieldAbstracts struct {
	Context
}
type YieldAddAbstract struct {
	Context
	Form     *form.Form
	Conflict bool
}
type YieldEditAbstract struct {
	Context
	AbstractID string
	Form       *form.Form
	Conflict   bool
}

func (h *Handler) AddAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "show_modal", "publication/add_abstract", YieldAddAbstract{
		Context: ctx,
		Form:    abstractForm(ctx.Loc, ctx.Publication, &models.Text{}, nil),
	})
}

func (h *Handler) CreateAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAbstract{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("create publication abstract: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	abstract := models.Text{
		Lang: b.Lang,
		Text: b.Text,
	}
	ctx.Publication.AddAbstract(&abstract)

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "publication/add_abstract", YieldAddAbstract{
			Context:  ctx,
			Form:     abstractForm(ctx.Loc, ctx.Publication, &abstract, validationErrs.(*okay.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "publication/add_abstract", YieldAddAbstract{
			Context:  ctx,
			Form:     abstractForm(ctx.Loc, ctx.Publication, &abstract, nil),
			Conflict: true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("create publication abstract: could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_abstracts", YieldAbstracts{
		Context: ctx,
	})
}

func (h *Handler) EditAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAbstract{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("edit publication abstract: could not bind request arguments", "error", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	abstract := ctx.Publication.GetAbstract(b.AbstractID)

	if abstract == nil {
		h.Logger.Warnf("edit publication abstract: Could not fetch the abstract:", "publication", ctx.Publication.ID, "abstract", b.AbstractID, "user", ctx.User.ID)
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("publication.conflict_error_reload"),
		})
		return
	}

	render.Layout(w, "show_modal", "publication/edit_abstract", YieldEditAbstract{
		Context:    ctx,
		AbstractID: b.AbstractID,
		Form:       abstractForm(ctx.Loc, ctx.Publication, abstract, nil),
	})
}

func (h *Handler) UpdateAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAbstract{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update publication abstract: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	abstract := ctx.Publication.GetAbstract(b.AbstractID)

	if abstract == nil {
		abstract := &models.Text{
			Text: b.Text,
			Lang: b.Lang,
		}
		render.Layout(w, "refresh_modal", "publication/edit_abstract", YieldEditAbstract{
			Context:    ctx,
			AbstractID: b.AbstractID,
			Form:       abstractForm(ctx.Loc, ctx.Publication, abstract, nil),
			Conflict:   true,
		})
		return
	}

	abstract.Text = b.Text
	abstract.Lang = b.Lang

	ctx.Publication.SetAbstract(abstract)

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "publication/edit_abstract", YieldEditAbstract{
			Context:    ctx,
			AbstractID: b.AbstractID,
			Form:       abstractForm(ctx.Loc, ctx.Publication, abstract, validationErrs.(*okay.Errors)),
			Conflict:   false,
		})
		return
	}

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "publication/edit_abstract", YieldEditAbstract{
			Context:    ctx,
			AbstractID: b.AbstractID,
			Form:       abstractForm(ctx.Loc, ctx.Publication, abstract, nil),
			Conflict:   true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("update publication abstract: could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_abstracts", YieldAbstracts{
		Context: ctx,
	})
}

func ConfirmDeleteAbstract(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	var b BindDeleteAbstract
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete publication abstract: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	if b.SnapshotID != publication.SnapshotID {
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: c.Loc.Get("publication.conflict_error_reload"),
		})
		return
	}

	views.ConfirmDeleteAbstract(c, publication, b.AbstractID).Render(r.Context(), w)
}

func (h *Handler) DeleteAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteAbstract
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete publication abstract: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	ctx.Publication.RemoveAbstract(b.AbstractID)

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("publication.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("create publication abstract: could not save the publication:", "error", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_abstracts", YieldAbstracts{
		Context: ctx,
	})
}

func abstractForm(loc *gotext.Locale, publication *models.Publication, abstract *models.Text, errors *okay.Errors) *form.Form {
	idx := -1
	for i, a := range publication.Abstract {
		if a.ID == abstract.ID {
			idx = i
		}
	}

	return form.New().
		WithTheme("cols").
		WithErrors(localize.ValidationErrors(loc, errors)).
		AddSection(
			&form.TextArea{
				Name:  "text",
				Value: abstract.Text,
				Label: loc.Get("builder.abstract.text"),
				Cols:  12,
				Rows:  6,
				Error: localize.ValidationErrorAt(loc, errors, fmt.Sprintf("/abstract/%d/text", idx)),
			},
			&form.Select{
				Name:    "lang",
				Value:   abstract.Lang,
				Label:   loc.Get("builder.abstract.lang"),
				Options: localize.LanguageSelectOptions(),
				Cols:    12,
				Error:   localize.ValidationErrorAt(loc, errors, fmt.Sprintf("/abstract/%d/lang", idx)),
			},
		)
}
