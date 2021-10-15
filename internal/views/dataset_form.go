package views

import (
	"html/template"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type DatasetForm struct {
	Data
	render     *render.Render
	Dataset    *models.Dataset
	FormErrors []jsonapi.Error
}

func NewDatasetForm(r *http.Request, render *render.Render, d *models.Dataset, fe []jsonapi.Error) DatasetForm {
	return DatasetForm{Data: NewData(r), render: render, Dataset: d, FormErrors: fe}
}

func (f DatasetForm) RenderFormText(text, key, pointer, label string, tooltip string, required bool, cols int) (template.HTML, error) {

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

func (f DatasetForm) RenderFormList(key, pointer, label string, selectedTerm string, taxonomy string, tooltip string, required bool, cols int) (template.HTML, error) {
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

	vocabulary["dataFormat"] = map[string]string{
		"xml":   "XML",
		"pdf":   "PDF",
		"json":  "JSON",
		"txt":   "TXT",
		"docx":  "DOCX",
		"zip":   "ZIP",
		"xlsx":  "XLSX",
		"pptx":  "PPTX",
		"rdf":   "RDF",
		"other": "Other",
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
