package views

import (
	"html/template"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type AbstractForm struct {
	Data
	render     *render.Render
	ID         string
	Abstract   *models.Text
	Key        string
	FormErrors []jsonapi.Error
}

func NewAbstractForm(r *http.Request, render *render.Render, id string, abstract *models.Text, key string, fe []jsonapi.Error) AbstractForm {
	return AbstractForm{Data: NewData(r), render: render, ID: id, Abstract: abstract, Key: key, FormErrors: fe}
}

func (f AbstractForm) RenderFormTextArea(text, key, pointer, label string, tooltip string, placeholder string, required bool, cols int, rows int) (template.HTML, error) {

	var formError jsonapi.Error
	hasError := false

	if f.FormErrors != nil {
		for _, err := range f.FormErrors {
			if err.Source.Pointer == pointer {
				formError = err
				hasError = true
			}
		}
	}

	return RenderPartial(f.render, "form/_text_area", &textAreaFormData{
		Key:         key,
		Label:       label,
		Text:        text,
		Tooltip:     tooltip,
		Placeholder: placeholder,
		Required:    required,
		Cols:        cols,
		Rows:        rows,
		HasError:    hasError,
		Error:       formError,
	})
}

func (f AbstractForm) RenderFormList(key, pointer, label string, selectedTerm string, taxonomy string, tooltip string, required bool, cols int) (template.HTML, error) {
	var formError jsonapi.Error
	hasError := false

	if f.FormErrors != nil {
		for _, err := range f.FormErrors {
			if err.Source.Pointer == pointer {
				formError = err
				hasError = true
			}
		}
	}

	// TODO: should come from a different API, use "taxonomy" to fetch list of values from taxonomy
	vocabulary := make(map[string]map[string]string)
	vocabulary["languages"] = map[string]string{
		"english": "english",
		"dutch":   "dutch",
		"french":  "french",
	}

	// Generate list of dropdown values, set selectedTerm in dropdown to "selected"
	// TODO: if we get a map back, we'll need to explicitly sort (numerical, alphabetically) since maps are hashmaps
	var terms []*listFormValues
	for key, term := range vocabulary[taxonomy] {
		selected := false
		if key == selectedTerm {
			selected = true
		}
		terms = append(terms, &listFormValues{
			key,
			term,
			selected,
		})
	}

	return RenderPartial(f.render, "form/_list", &listFormData{
		Key:      key,
		Label:    label,
		Values:   terms,
		Tooltip:  tooltip,
		Required: required,
		Cols:     cols,
		HasError: hasError,
		Error:    formError,
	})
}
