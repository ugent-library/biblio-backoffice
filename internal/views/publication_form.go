package views

import (
	"html/template"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
)

type Pair struct {
	Label string
	Value string
}

type PublicationForm struct {
	Data
	render      *render.Render
	Publication *models.Publication
	FormErrors  []models.FormError
}

func NewPublicationForm(r *http.Request, render *render.Render, p *models.Publication, fe []models.FormError) PublicationForm {
	return PublicationForm{Data: NewData(r), render: render, Publication: p, FormErrors: fe}
}

func (f PublicationForm) RenderFormText(value, name, pointer, label string, tooltip string, required bool, cols int) (template.HTML, error) {

	var formError models.FormError
	hasError := false

	if f.FormErrors != nil {
		for _, err := range f.FormErrors {
			if err.Source["pointer"] == pointer {
				formError = err
				hasError = true
			}
		}
	}

	return (&TextInput{
		Name:     name,
    Label:    label,
    Value:    value,
    Tooltip:  tooltip,
    Required: required,
    Cols:     cols,
    HasError: hasError,
    Error:    formError,
	}).Render(f.render)
}

// TODO: We'll need dedicated functions for Department, Project, etc. because
// Department, Project take specific types ([]PublicationDepartment, []PublicationProject)
func (f PublicationForm) RenderFormTextMultiple(values []string, name, pointer, label string, tooltip string, required bool, cols int) (template.HTML, error) {

	var formError models.FormError
	hasError := false

	if f.FormErrors != nil {
		for _, err := range f.FormErrors {
			if err.Source["pointer"] == pointer {
				formError = err
				hasError = true
			}
		}
	}

	return (&MultiTextInput{
		Name:     name,
		Label:    label,
		Values:   values,
		Tooltip:  tooltip,
		Required: required,
		Cols:     cols,
		HasError: hasError,
		Error:    formError,
	}).Render(f.render)

}

func (f PublicationForm) RenderFormList(name, pointer, label string, selectedTerm string, taxonomy string, tooltip string, required bool, cols int) (template.HTML, error) {
	var formError models.FormError
	hasError := false

	if f.FormErrors != nil {
		for _, err := range f.FormErrors {
			if err.Source["pointer"] == pointer {
				formError = err
				hasError = true
			}
		}
	}

  mapLabelEmpty := map[string]string{
    "classification": "",
		"articleType": "",
  }

	// TODO: should come from a different API, use "taxonomy" to fetch list of values from taxonomy
	vocabulary := make(map[string][]Pair)
	vocabulary["type"] = []Pair{
		Pair{Value: "journal_article", Label: "Journal Article"},
		Pair{Value: "book", Label: "Book"},
		Pair{Value: "book_chapter", Label: "Book Chapter"},
		Pair{Value: "book_editor", Label: "Book editor"},
		Pair{Value: "issue_editor", Label: "Issue editor"},
		Pair{Value: "conference", Label: "Conference"},
		Pair{Value: "dissertation", Label: "Dissertation"},
		Pair{Value: "miscellaneous", Label: "Miscellaneous"},
		Pair{Value: "report", Label: "Report"},
		Pair{Value: "preprint", Label:"Preprint"},
	}

	// TODO: should come from a different API, use "taxonomy" to fetch list of values from taxonomy
	// TODO: This might be hierarchical select with a dep on type: only show "journal article" if type of publication = "Journal Article"
	vocabulary["classification"] = []Pair{
		Pair{Value: "U", Label: "U"},
		Pair{Value: "A1", Label: "A1"},
		Pair{Value: "A2", Label: "A2"},
		Pair{Value: "A3", Label: "A3"},
		Pair{Value: "A4", Label: "A4"},
		Pair{Value: "V", Label: "V"},
	}

	vocabulary["articleType"] = []Pair{
		Pair{Value:	"original", Label: "Original"},
		Pair{Value: "review", Label: "Review"},
		Pair{Value:	"letter_note", Label: "Letter note"},
		Pair{Value: "proceedingsPaper", Label: "Proceedings Paper"},
	}

	// Generate list of dropdown values, set selectedTerm in dropdown to "selected"
	//empty option?
  var emptyOption *SelectOption = nil
  if labelEmpty, ok := mapLabelEmpty[taxonomy]; ok {
    emptyOption = &SelectOption{ Label: labelEmpty }
  }
	var terms []*SelectOption
	if emptyOption != nil {
		terms = append(terms, emptyOption)
	}
	for _, pair := range vocabulary[taxonomy] {
		terms = append(terms, &SelectOption{
			Value: pair.Value,
			Label: pair.Label,
			Selected: pair.Value == selectedTerm,
		})
	}

	return (&Select{
		Name:     name,
		Label:    label,
		Values:   terms,
		Tooltip:  tooltip,
		Required: required,
		Cols:     cols,
		HasError: hasError,
		Error:    formError,
	}).Render(f.render)

}

// TODO: We'll need dedicated functions for fields that take specific types ([]PublicationDepartment, []PublicationProject)
//    type assertion in one single function would become too complex quickly.
func (f PublicationForm) RenderFormListMultiple(selectedTerms []string, name, pointer, label string, taxonomy string, tooltip string, required bool, cols int) (template.HTML, error) {
	var formError models.FormError
	hasError := false

	if f.FormErrors != nil {
		for _, err := range f.FormErrors {
			if err.Source["pointer"] == pointer {
				formError = err
				hasError = true
			}
		}
	}

	mapLabelEmpty := map[string]string{
		"languages": "",
	}

	// // TODO: should come from a different API, use "taxonomy" to fetch list of values from taxonomy
	vocabulary := make(map[string][]Pair)
	vocabulary["languages"] = []Pair{
		Pair{Value:	"eng", Label: "English"},
		Pair{Value:	"dut", Label: "Dutch"},
		Pair{Value:	"ger", Label: "German"},
		Pair{Value:	"fre", Label:"French"},
	}

	//empty option?
	var emptyOption *SelectOption = nil
	if labelEmpty, ok := mapLabelEmpty[taxonomy]; ok {
		emptyOption = &SelectOption{ Label: labelEmpty }
	}

	//list of selects
	values := [][]*SelectOption{}
	for _, lterm := range selectedTerms {
		var terms []*SelectOption
		if emptyOption != nil {
			terms = append(terms, emptyOption)
		}
		for _, vpair := range vocabulary[taxonomy] {
			terms = append(terms, &SelectOption{
				Value: vpair.Value,
				Label: vpair.Label,
				Selected: vpair.Value == lterm,
			})
		}

		values = append(values, terms)
	}

	//new select
	var selectableOptions []*SelectOption
	if emptyOption != nil {
  	selectableOptions = append(selectableOptions, emptyOption)
	}
	for _, vpair := range vocabulary[taxonomy] {
    selectableOptions = append(selectableOptions, &SelectOption{
      Value: vpair.Value,
      Label: vpair.Label,
    })
  }

	return (&MultiSelect{
		Name:       name,
		Label:      label,
		Values:     values,
		Vocabulary: selectableOptions,
		Tooltip:    tooltip,
		Required:   required,
		Cols:       cols,
		HasError:   hasError,
		Error:      formError,
	}).Render(f.render)

}

func (f PublicationForm) RenderFormCheckbox(checked bool, name, pointer, label string, tooltip string, required bool, cols int) (template.HTML, error) {

	var formError models.FormError
	hasError := false

	if f.FormErrors != nil {
		for _, err := range f.FormErrors {
			if err.Source["pointer"] == pointer {
				formError = err
				hasError = true
			}
		}
	}

	return (&CheckboxInput{
		Name:     name,
		Label:    label,
		Checked:  checked,
		Tooltip:  tooltip,
		Required: required,
		Cols:     cols,
		HasError: hasError,
		Error:    formError,
	}).Render(f.render)
}
