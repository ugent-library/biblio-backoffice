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
	Form *form.Form
}
type YieldEditAbstract struct {
	Context
	AbstractID string
	Form       *form.Form
}
type YieldDeleteAbstract struct {
	Context
	AbstractID string
}

func (h *Handler) AddAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	form := abstractForm(ctx.Locale, ctx.Publication, &models.Text{}, nil)

	render.Layout(w, "show_modal", "publication/add_abstract", YieldAddAbstract{
		Context: ctx,
		Form:    form,
	})
}

func (h *Handler) CreateAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAbstract{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
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
			Context: ctx,
			Form:    abstractForm(ctx.Locale, ctx.Publication, &abstract, validationErrs.(validation.Errors)),
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

	render.View(w, "publication/refresh_abstracts", YieldAbstracts{
		Context: ctx,
	})
}

func (h *Handler) EditAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAbstract{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	abstract := ctx.Publication.GetAbstract(b.AbstractID)
	if abstract == nil {
		render.BadRequest(
			w,
			r,
			fmt.Errorf("no abstract found for %s in publication %s", b.AbstractID, ctx.Publication.ID),
		)
		return
	}

	render.Layout(w, "show_modal", "publication/edit_abstract", YieldEditAbstract{
		Context:    ctx,
		AbstractID: b.AbstractID,
		Form:       abstractForm(ctx.Locale, ctx.Publication, abstract, nil),
	})
}

func (h *Handler) UpdateAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAbstract{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	/*
		TODO: when abstract already removed,
		this code throws a bad request when
		it should be throwing a conflict error.
		But at this point, one cannot distinguish
		between "id never existed" or "id existed before"
	*/
	abstract := ctx.Publication.GetAbstract(b.AbstractID)
	if abstract == nil {
		render.BadRequest(
			w,
			r,
			fmt.Errorf("no abstract found for %s in publication %s", b.AbstractID, ctx.Publication.ID),
		)
		return
	}
	abstract.Text = b.Text
	abstract.Lang = b.Lang

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		form := abstractForm(ctx.Locale, ctx.Publication, abstract, validationErrs.(validation.Errors))

		render.Layout(w, "refresh_modal", "publication/edit_abstract", YieldEditAbstract{
			Context:    ctx,
			AbstractID: b.AbstractID,
			Form:       form,
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

	render.View(w, "publication/refresh_abstracts", YieldAbstracts{
		Context: ctx,
	})
}

func (h *Handler) ConfirmDeleteAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteAbstract
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.Layout(w, "show_modal", "publication/confirm_delete_abstract", YieldDeleteAbstract{
		Context:    ctx,
		AbstractID: b.AbstractID,
	})
}

func (h *Handler) DeleteAbstract(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteAbstract
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	ctx.Publication.RemoveAbstract(b.AbstractID)

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

	render.View(w, "publication/refresh_abstracts", YieldAbstracts{
		Context: ctx,
	})
}

func abstractForm(l *locale.Locale, publication *models.Publication, abstract *models.Text, errors validation.Errors) *form.Form {
	idx := -1
	for i, a := range publication.Abstract {
		if a.ID == abstract.ID {
			idx = i
		}
	}

	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(l, errors)).
		AddSection(
			&form.TextArea{
				Name:        "text",
				Value:       abstract.Text,
				Label:       l.T("builder.abstract.text"),
				Cols:        12,
				Rows:        6,
				Placeholder: l.T("builder.abstract.text.placeholder"),
				Error:       localize.ValidationErrorAt(l, errors, fmt.Sprintf("/abstract/%d/text", idx)),
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
