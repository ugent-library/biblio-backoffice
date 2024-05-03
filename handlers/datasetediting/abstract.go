package datasetediting

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
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
	form := abstractForm(ctx.Loc, ctx.Dataset, &models.Text{}, nil)

	render.Layout(w, "show_modal", "dataset/add_abstract", YieldAddAbstract{
		Context:  ctx,
		Form:     form,
		Conflict: false,
	})
}

func (h *Handler) CreateAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAbstract{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("create dataset abstract: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	abstract := models.Text{
		Lang: b.Lang,
		Text: b.Text,
	}

	ctx.Dataset.AddAbstract(&abstract)

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "dataset/add_abstract", YieldAddAbstract{
			Context:  ctx,
			Form:     abstractForm(ctx.Loc, ctx.Dataset, &abstract, validationErrs.(*okay.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "dataset/add_abstract", YieldAddAbstract{
			Context:  ctx,
			Form:     abstractForm(ctx.Loc, ctx.Dataset, &abstract, nil),
			Conflict: true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("create dataset abstract: could not save the dataset:", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_abstracts", YieldAbstracts{
		Context: ctx,
	})
}

func EditAbstract(w http.ResponseWriter, r *http.Request, legacyContext Context) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindAbstract{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("edit dataset abstract: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	abstract := dataset.GetAbstract(b.AbstractID)

	// TODO catch non-existing item in UI
	if abstract == nil {
		c.Log.Warnf("edit dataset abstract: Could not fetch the abstract:", "dataset", dataset.ID, "abstract", b.AbstractID, "user", c.User.ID)
		views.ShowModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"), "")).Render(r.Context(), w)
		return
	}

	render.Layout(w, "show_modal", "dataset/edit_abstract", YieldEditAbstract{
		Context:    legacyContext,
		AbstractID: b.AbstractID,
		Form:       abstractForm(c.Loc, dataset, abstract, nil),
		Conflict:   false,
	})
}

func (h *Handler) UpdateAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAbstract{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update dataset abstract: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	// get pointer to abstract and manipulate in place
	abstract := ctx.Dataset.GetAbstract(b.AbstractID)

	if abstract == nil {
		abstract := &models.Text{
			Text: b.Text,
			Lang: b.Lang,
		}
		render.Layout(w, "refresh_modal", "dataset/edit_abstract", YieldEditAbstract{
			Context:    ctx,
			AbstractID: b.AbstractID,
			Form:       abstractForm(ctx.Loc, ctx.Dataset, abstract, nil),
			Conflict:   true,
		})
		return
	}

	abstract.Text = b.Text
	abstract.Lang = b.Lang

	ctx.Dataset.SetAbstract(abstract)

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "dataset/edit_abstract", YieldEditAbstract{
			Context:    ctx,
			AbstractID: b.AbstractID,
			Form:       abstractForm(ctx.Loc, ctx.Dataset, abstract, validationErrs.(*okay.Errors)),
			Conflict:   false,
		})
		return
	}

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "dataset/edit_abstract", YieldEditAbstract{
			Context:    ctx,
			AbstractID: b.AbstractID,
			Form:       abstractForm(ctx.Loc, ctx.Dataset, abstract, nil),
			Conflict:   true,
		})
		return
	}

	if err != nil {
		h.Logger.Warnf("update dataset abstract: Could not save the dataset:", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_abstracts", YieldAbstracts{
		Context: ctx,
	})
}

func ConfirmDeleteAbstract(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	var b BindDeleteAbstract
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete dataset: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	if b.SnapshotID != dataset.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"), "")).Render(r.Context(), w)
		return
	}

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to remove this abstract?",
		DeleteUrl:  c.PathTo("dataset_delete_abstract", "id", dataset.ID, "abstract_id", b.AbstractID),
		SnapshotID: dataset.SnapshotID,
	}).Render(r.Context(), w)
}

func (h *Handler) DeleteAbstract(w http.ResponseWriter, r *http.Request, legacyContext Context) {
	c := ctx.Get(r)

	var b BindDeleteAbstract
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete datase abstract: could not bind request arguments", "errors", err, "request", r, "user", legacyContext.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	legacyContext.Dataset.RemoveAbstract(b.AbstractID)

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), legacyContext.Dataset, legacyContext.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"), "")).Render(r.Context(), w)
		return
	}

	if err != nil {
		h.Logger.Warnf("delete dataset abstract: Could not save the dataset:", "errors", err, "dataset", legacyContext.Dataset.ID, "user", legacyContext.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_abstracts", YieldAbstracts{
		Context: legacyContext,
	})
}

func abstractForm(loc *gotext.Locale, dataset *models.Dataset, abstract *models.Text, errors *okay.Errors) *form.Form {
	idx := -1
	for i, a := range dataset.Abstract {
		if a.ID == abstract.ID {
			idx = i
			break
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
