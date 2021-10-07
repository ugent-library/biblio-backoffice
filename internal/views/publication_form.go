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

type textFormData struct {
	Name     string
	Value    string
	Label    string
	Required bool
	Tooltip  string
	Cols     int
	HasError bool
	Error    models.FormError
}

type textMultipleFormData struct {
	Name     string
	Values   []string
	Label    string
	Required bool
	Tooltip  string
	Cols     int
	HasError bool
	Error    models.FormError
}

type listFormValues struct {
	Value    string
	Label    string
	Selected bool
}

type listFormData struct {
	Name     string
	Values   []*listFormValues
	Label    string
	Required bool
	Tooltip  string
	Cols     int
	HasError bool
	Error    models.FormError
}

type listMultipleFormData struct {
	Name       string
	Values     [][]*listFormValues
	Vocabulary []*listFormValues
	Label      string
	Required   bool
	Tooltip    string
	Cols       int
	HasError   bool
	Error      models.FormError
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

	return RenderPartial(f.render, "form/_text", &textFormData{
		Name:     name,
		Label:    label,
		Value:    value,
		Tooltip:  tooltip,
		Required: required,
		Cols:     cols,
		HasError: hasError,
		Error:    formError,
	})
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

	return RenderPartial(f.render, "form/_text_multiple", &textMultipleFormData{
		Name:     name,
		Label:    label,
		Values:   values,
		Tooltip:  tooltip,
		Required: required,
		Cols:     cols,
		HasError: hasError,
		Error:    formError,
	})

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
  var emptyOption *listFormValues = nil
  if labelEmpty, ok := mapLabelEmpty[taxonomy]; ok {
    emptyOption = &listFormValues{ Label: labelEmpty }
  }
	var terms []*listFormValues
	if emptyOption != nil {
		terms = append(terms, emptyOption)
	}
	for _, pair := range vocabulary[taxonomy] {
		terms = append(terms, &listFormValues{
			Value: pair.Value,
			Label: pair.Label,
			Selected: pair.Value == selectedTerm,
		})
	}

	return RenderPartial(f.render, "form/_list", &listFormData{
		Name:     name,
		Label:    label,
		Values:   terms,
		Tooltip:  tooltip,
		Required: required,
		Cols:     cols,
		HasError: hasError,
		Error:    formError,
	})
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
	var emptyOption *listFormValues = nil
	if labelEmpty, ok := mapLabelEmpty[taxonomy]; ok {
		emptyOption = &listFormValues{ Label: labelEmpty }
	}

	//list of selects
	values := [][]*listFormValues{}
	for _, lterm := range selectedTerms {
		var terms []*listFormValues
		if emptyOption != nil {
			terms = append(terms, emptyOption)
		}
		for _, vpair := range vocabulary[taxonomy] {
			terms = append(terms, &listFormValues{
				Value: vpair.Value,
				Label: vpair.Label,
				Selected: vpair.Value == lterm,
			})
		}

		values = append(values, terms)
	}

	//new select
	var selectableOptions []*listFormValues
	if emptyOption != nil {
  	selectableOptions = append(selectableOptions, emptyOption)
	}
	for _, vpair := range vocabulary[taxonomy] {
    selectableOptions = append(selectableOptions, &listFormValues{
      Value: vpair.Value,
      Label: vpair.Label,
    })
  }

	return RenderPartial(f.render, "form/_list_multiple", &listMultipleFormData{
		Name:       name,
		Label:      label,
		Values:     values,
		Vocabulary: selectableOptions,
		Tooltip:    tooltip,
		Required:   required,
		Cols:       cols,
		HasError:   hasError,
		Error:      formError,
	})
}
