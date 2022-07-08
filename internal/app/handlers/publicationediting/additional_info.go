package publicationediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/displays"
	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
)

type BindAdditionalInfo struct {
	AdditionalInfo string   `form:"additional_info"`
	Keyword        []string `form:"keyword"`
	ResearchField  []string `form:"research_field"`
}

type YieldAdditionalInfo struct {
	Context
	DisplayAdditionalInfo *display.Display
}

type YieldEditAdditionalInfo struct {
	Context
	Form *form.Form
}

func (h *Handler) EditAdditionalInfo(w http.ResponseWriter, r *http.Request, ctx Context) {
	p := ctx.Publication
	b := BindAdditionalInfo{
		AdditionalInfo: p.AdditionalInfo,
		Keyword:        p.Keyword,
		ResearchField:  p.ResearchField,
	}

	render.Layout(w, "show_modal", "publication/edit_additional_info", YieldEditAdditionalInfo{
		Context: ctx,
		Form:    additionalInfoForm(ctx, b, nil),
	})
}

func (h *Handler) UpdateAdditionalInfo(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAdditionalInfo{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	p := ctx.Publication
	p.AdditionalInfo = b.AdditionalInfo
	p.Keyword = b.Keyword
	p.ResearchField = b.ResearchField

	if validationErrs := p.Validate(); validationErrs != nil {
		form := additionalInfoForm(ctx, b, validationErrs.(validation.Errors))

		render.Layout(w, "refresh_modal", "publication/edit_additional_info", YieldEditAdditionalInfo{
			Context: ctx,
			Form:    form,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_additional_info", YieldAdditionalInfo{
		Context:               ctx,
		DisplayAdditionalInfo: displays.PublicationAdditionalInfo(ctx.Locale, p),
	})
}

func additionalInfoForm(ctx Context, b BindAdditionalInfo, errors validation.Errors) *form.Form {
	l := ctx.Locale

	researchFieldOptions := make([]form.SelectOption, len(vocabularies.Map["research_fields"]))
	for i, v := range vocabularies.Map["research_fields"] {
		researchFieldOptions[i].Label = v
		researchFieldOptions[i].Value = v
	}

	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(ctx.Locale, errors)).
		AddSection(
			&form.SelectRepeat{
				Name:        "research_field",
				Options:     researchFieldOptions,
				Values:      b.ResearchField,
				EmptyOption: true,
				Label:       l.T("builder.research_field"),
				Cols:        9,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					"/research_field",
				),
			},
			&form.TextRepeat{
				Name:   "keyword",
				Values: b.Keyword,
				Label:  l.T("builder.keyword"),
				Cols:   9,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					"/keyword",
				),
			},
			&form.TextArea{
				Name:  "additional_info",
				Value: b.AdditionalInfo,
				Label: l.T("builder.additional_info"),
				Cols:  9,
				Rows:  4,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					"/additional_info",
				),
			},
		)
}
