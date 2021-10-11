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
func (f PublicationForm) RenderFormTextMultiple(values []string, key, pointer, label string, tooltip string, required bool, cols int) (template.HTML, error) {

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
		"U":  "U",
		"A1": "A1",
		"A2": "A2",
		"A3": "A3",
		"A4": "A4",
		"V":  "V",
	}

	vocabulary["articleType"] = map[string]string{
		"original":         "Original",
		"review":           "Review",
		"letter_note":      "Letter note",
		"proceedingsPaper": "Proceedings Paper",
	}

	vocabulary["miscellaneousType"] = map[string]string{
		"artReview":         "Art review",
		"artisticWork":      "Artistic Work",
		"bibliography":      "Bibliography",
		"biography":         "Biography",
		"blogPost":          "Blogpost",
		"bookReview":        "Book review",
		"correction":        "Correction",
		"dictionaryEntry":   "Dictionary entry",
		"editorialMaterial": "Editorial material",
		"encyclopediaEntry": "Encyclopedia entry",
		"exhibitionReview":  "Exhibition review",
		"filmReview":        "Film review",
		"lectureSpeech":     "Lecture speech",
		"lemma":             "Lemma",
		"magazinePiece":     "Magazine piece",
		"manual":            "Manual",
		"musicEdition":      "Music edition",
		"musicReview":       "Music review",
		"newsArticle":       "News article",
		"newspaperPiece":    "Newspaper piece",
		"other":             "Other",
		"preprint":          "Preprint",
		"productReview":     "Product review",
		"report":            "Report",
		"technicalStandard": "Technical standard",
		"textEdition":       "Text edition",
		"textTranslation":   "Text translation",
		"theatreReview":     "Theatre review",
		"workingPaper":      "Working paper",
	}

	vocabulary["publicationStatus"] = map[string]string{
		"unpublished": "unpublished",
		"accepted":    "accepted",
		"published":   "published",
	}

	vocabulary["conferenceType"] = map[string]string{
		"proceedingsPaper": "proceedingsPaper",
		"abstract":         "abstract",
		"poster":           "poster",
		"other":            "other",
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
func (f PublicationForm) RenderFormListMultiple(selectedTerms []string, key, pointer, label string, taxonomy string, tooltip string, required bool, cols int) (template.HTML, error) {
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

	// // TODO: should come from a different API, use "taxonomy" to fetch list of values from taxonomy
	vocabulary := make(map[string]map[string]string)
	vocabulary["languages"] = map[string]string{
		"eng": "English",
		"dut": "Dutch",
		"ger": "German",
		"fre": "French",
	}

	values := make(map[int][]*listFormValues)
	for lkey, lterm := range selectedTerms {
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

func (f PublicationForm) RenderFormCheckbox(checked bool, name, pointer, label string, tooltip string, required bool, cols int) (template.HTML, error) {

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

	return RenderPartial(f.render, "form/_checkbox", &CheckboxInput{
		Name:     name,
		Label:    label,
		Checked:  checked,
		Tooltip:  tooltip,
		Required: required,
		Cols:     cols,
		HasError: hasError,
		Error:    formError,
	})

}
