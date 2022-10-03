package datasetediting

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type BindAbstract struct {
	AbstractID string `path:"abstract_id"`
	Text       string `form:"text"`
	Lang       string `form:"lang"`
}

type BindDeleteAbstract struct {
	AbstractID string `path:"abstract_id"`
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
type YieldDeleteAbstract struct {
	Context
	AbstractID string
	Conflict   bool
}

func (h *Handler) AddAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	form := abstractForm(ctx.Locale, ctx.Dataset, &models.Text{}, nil)

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
			Form:     abstractForm(ctx.Locale, ctx.Dataset, &abstract, validationErrs.(validation.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "dataset/add_abstract", YieldAddAbstract{
			Context:  ctx,
			Form:     abstractForm(ctx.Locale, ctx.Dataset, &abstract, nil),
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

func (h *Handler) EditAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAbstract{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("edit dataset abstract: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	abstract := ctx.Dataset.GetAbstract(b.AbstractID)

	// TODO catch non-existing item in UI
	if abstract == nil {
		h.Logger.Warnf("edit dataset abstract: Could not fetch the abstract:", "dataset", ctx.Dataset.ID, "abstract", b.AbstractID, "user", ctx.User.ID)
		render.NotFoundError(
			w,
			r,
			fmt.Errorf("no abstract found for %s in dataset %s", b.AbstractID, ctx.Dataset.ID),
		)
		return
	}

	render.Layout(w, "show_modal", "dataset/edit_abstract", YieldEditAbstract{
		Context:    ctx,
		AbstractID: b.AbstractID,
		Form:       abstractForm(ctx.Locale, ctx.Dataset, abstract, nil),
		Conflict:   false,
	})
}

func (h *Handler) UpdateAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAbstract{}
	if err := bind.Request(r, &b); err != nil {
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
			Form:       abstractForm(ctx.Locale, ctx.Dataset, abstract, nil),
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
			Form:       abstractForm(ctx.Locale, ctx.Dataset, abstract, validationErrs.(validation.Errors)),
			Conflict:   false,
		})
		return
	}

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "dataset/edit_abstract", YieldEditAbstract{
			Context:    ctx,
			AbstractID: b.AbstractID,
			Form:       abstractForm(ctx.Locale, ctx.Dataset, abstract, nil),
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

func (h *Handler) ConfirmDeleteAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteAbstract
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("confirm delete dataset: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	// TODO catch non-existing item in UI

	render.Layout(w, "show_modal", "dataset/confirm_delete_abstract", YieldDeleteAbstract{
		Context:    ctx,
		AbstractID: b.AbstractID,
		Conflict:   false,
	})
}

func (h *Handler) DeleteAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteAbstract
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete datase abstract: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	ctx.Dataset.RemoveAbstract(b.AbstractID)

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Locale.T("dataset.conflict_error"),
		})
		return
	}

	if err != nil {
		h.Logger.Warnf("delete dataset abstract: Could not save the dataset:", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_abstracts", YieldAbstracts{
		Context: ctx,
	})
}

func abstractForm(l *locale.Locale, dataset *models.Dataset, abstract *models.Text, errors validation.Errors) *form.Form {

	idx := -1
	for i, a := range dataset.Abstract {
		if a.ID == abstract.ID {
			idx = i
			break
		}
	}

	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(l, errors)).
		AddSection(
			&form.TextArea{
				Name:  "text",
				Value: abstract.Text,
				Label: l.T("builder.abstract.text"),
				Cols:  12,
				Rows:  6,
				Error: localize.ValidationErrorAt(l, errors, fmt.Sprintf("/abstract/%d/text", idx)),
			},
			&form.Select{
				Name:    "lang",
				Value:   abstract.Lang,
				Label:   l.T("builder.abstract.lang"),
				Options: localize.LanguageSelectOptions(l),
				Cols:    12,
				Error:   localize.ValidationErrorAt(l, errors, fmt.Sprintf("/abstract/%d/lang", idx)),
			},
		)
}
