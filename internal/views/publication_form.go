package views

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
)

type PublicationForm struct {
	Data
	render      *render.Render
	Publication *models.Publication
}

type textFormData struct {
	Key      string
	Text     string
	Label    string
	Required bool
	Tooltip  string
	Cols     int
}

type textMultipleFormData struct {
	Key      string
	Text     []string
	Label    string
	Required bool
	Tooltip  string
	Cols     int
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
}

type listMultipleFormData struct {
	Key        string
	Values     map[int][]*listFormValues
	Vocabulary map[string]string
	Label      string
	Required   bool
	Tooltip    string
	Cols       int
}

func NewPublicationForm(r *http.Request, render *render.Render, p *models.Publication) PublicationForm {
	return PublicationForm{Data: NewData(r), render: render, Publication: p}
}

func (f PublicationForm) RenderFormText(text, key, label string, tooltip string, required bool, cols int) template.HTML {
	return RenderPartial(f.render, "form/_text", &textFormData{
		Key:      key,
		Label:    label,
		Text:     text,
		Tooltip:  tooltip,
		Required: required,
		Cols:     cols,
	})
}

// TODO: We'll need dedicated functions for Department, Project, etc. because
// Department, Project take specific types ([]PublicationDepartment, []PublicationProject)
func (f PublicationForm) RenderFormTextMultiple(text interface{}, key, label string, tooltip string, required bool, cols int) template.HTML {

	// TODO: remove me
	values := []string{"foo", "bar"}

	return RenderPartial(f.render, "form/_text_multiple", &textMultipleFormData{
		Key:      key,
		Label:    label,
		Text:     values,
		Tooltip:  tooltip,
		Required: required,
		Cols:     cols,
	})
}

func (f PublicationForm) RenderFormList(key, label string, selectedTerm string, taxonomy string, tooltip string, required bool, cols int) template.HTML {

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
	})
}

// TODO: We'll need dedicated functions for fields that take specific types ([]PublicationDepartment, []PublicationProject)
//    type assertion in one single function would become too complex quickly.
func (f PublicationForm) RenderFormListMultiple(selectedTerms interface{}, key, label string, taxonomy string, tooltip string, required bool, cols int) template.HTML {

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

		fmt.Println(terms)

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
	})
}
