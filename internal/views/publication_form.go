package views

import (
	"html/template"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type PublicationForm struct {
	Data
	render      *render.Render
	Publication *models.Publication
	FormErrors  []jsonapi.Error
}

type textFormData struct {
	Key      string
	Text     string
	Label    string
	Required bool
	Tooltip  string
	Cols     int
	HasError bool
	Error    jsonapi.Error
}

type textMultipleFormData struct {
	Key      string
	Text     []string
	Label    string
	Required bool
	Tooltip  string
	Cols     int
	HasError bool
	Error    jsonapi.Error
}

type listFormValues struct {
	Key      string
	Value    string
	Selected bool
}

type listFormData struct {
	Key      string
	Values   []*listFormValues
	Label    string
	Required bool
	Tooltip  string
	Cols     int
	HasError bool
	Error    jsonapi.Error
}

type listMultipleFormData struct {
	Key        string
	Values     map[int][]*listFormValues
	Vocabulary map[string]string
	Label      string
	Required   bool
	Tooltip    string
	Cols       int
	HasError   bool
	Error      jsonapi.Error
}

func NewPublicationForm(r *http.Request, render *render.Render, p *models.Publication, fe []jsonapi.Error) PublicationForm {
	return PublicationForm{Data: NewData(r), render: render, Publication: p, FormErrors: fe}
}

func (f PublicationForm) RenderFormText(text, key, pointer, label string, tooltip string, required bool, cols int) (template.HTML, error) {

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

	return RenderPartial(f.render, "form/_text", &textFormData{
		Key:      key,
		Label:    label,
		Text:     text,
		Tooltip:  tooltip,
		Required: required,
		Cols:     cols,
		HasError: hasError,
		Error:    formError,
	})
}

// TODO: We'll need dedicated functions for Department, Project, etc. because
// Department, Project take specific types ([]PublicationDepartment, []PublicationProject)
func (f PublicationForm) RenderFormTextMultiple(text interface{}, key, pointer, label string, tooltip string, required bool, cols int) (template.HTML, error) {
	// TODO: remove me
	values := []string{"foo", "bar"}

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

	return RenderPartial(f.render, "form/_text_multiple", &textMultipleFormData{
		Key:      key,
		Label:    label,
		Text:     values,
		Tooltip:  tooltip,
		Required: required,
		Cols:     cols,
		HasError: hasError,
		Error:    formError,
	})

}

func (f PublicationForm) RenderFormList(key, pointer, label string, selectedTerm string, taxonomy string, tooltip string, required bool, cols int) (template.HTML, error) {
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
	vocabulary["type"] = map[string]string{
		"journal_article": "Journal Article",
		"book":            "Book",
		"book_chapter":    "Book Chapter",
		"book_editor":     "Book editor",
		"issue_editor":    "Issue editor",
		"conference":      "Conference",
		"dissertation":    "Dissertation",
		"miscellaneous":   "Miscellaneous",
		"report":          "Report",
		"preprint":        "Preprint",
	}

	// TODO: should come from a different API, use "taxonomy" to fetch list of values from taxonomy
	// TODO: This might be hierarchical select with a dep on type: only show "journal article" if type of publication = "Journal Article"
	vocabulary["classification"] = map[string]string{
		"journal_article_u":  "Journal Article - U",
		"journal_article_a1": "Journal Article - A1",
		"journal_article_a2": "Journal Article - A2",
		"journal_article_a3": "Journal Article - A3",
		"journal_article_a4": "Journal Article - A4",
		"journal_article_v":  "Journal Article - V",
	}

	vocabulary["articleType"] = map[string]string{
		"original":         "Original",
		"review":           "Review",
		"letter_note":      "Letter note",
		"proceedingsPaper": "Proceedings Paper",
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

// TODO: We'll need dedicated functions for fields that take specific types ([]PublicationDepartment, []PublicationProject)
//    type assertion in one single function would become too complex quickly.
func (f PublicationForm) RenderFormListMultiple(selectedTerms interface{}, key, pointer, label string, taxonomy string, tooltip string, required bool, cols int) (template.HTML, error) {
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

	// TODO: remove me / Fetch me from a struct field
	languages := []string{"eng", "dut"}

	// // TODO: should come from a different API, use "taxonomy" to fetch list of values from taxonomy
	vocabulary := make(map[string]map[string]string)
	vocabulary["languages"] = map[string]string{
		"eng": "English",
		"dut": "Dutch",
		"ger": "German",
		"fre": "French",
	}

	values := make(map[int][]*listFormValues)
	for lkey, lterm := range languages {
		var terms []*listFormValues
		for vkey, vterm := range vocabulary[taxonomy] {
			selected := false
			if vkey == lterm {
				selected = true
			}
			terms = append(terms, &listFormValues{
				vkey,
				vterm,
				selected,
			})
		}

		values[lkey] = terms
	}

	return RenderPartial(f.render, "form/_list_multiple", &listMultipleFormData{
		Key:        key,
		Label:      label,
		Values:     values,
		Vocabulary: vocabulary["languages"],
		Tooltip:    tooltip,
		Required:   required,
		Cols:       cols,
		HasError:   hasError,
		Error:      formError,
	})
}
