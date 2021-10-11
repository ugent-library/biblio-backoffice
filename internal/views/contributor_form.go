package views

import (
	"html/template"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type ContributorForm struct {
	Data
	render      *render.Render
	Publication *models.Publication
	AuthorDelta int
	FormErrors  []jsonapi.Error
}

func NewContributorForm(r *http.Request, render *render.Render, p *models.Publication, ad int, fe []jsonapi.Error) ContributorForm {
	return ContributorForm{Data: NewData(r), render: render, Publication: p, AuthorDelta: ad, FormErrors: fe}
}

func (f ContributorForm) RenderFormText(text, key, pointer, label string, tooltip string, required bool, cols int) (template.HTML, error) {

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

func (f ContributorForm) RenderFormMultiSelectList(key, pointer, label string, selectedTerms []string, taxonomy string, tooltip string, required bool, cols int) (template.HTML, error) {
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
	vocabulary["creditRoles"] = map[string]string{
		"first_author":           "First author",
		"last_author":            "Last author",
		"conceptualization":      "Conceptualization",
		"data_curation":          "Datacuration",
		"formal_analysis":        "Formala nalysis",
		"funding_acquisition":    "Funding acquisition",
		"investigation":          "Investigation",
		"methodology":            "Methodology",
		"project_administration": "Project administration",
		"resources":              "Resources",
		"software":               "Software",
		"supervision":            "Supervision",
		"validation":             "Validation",
		"visualization":          "Visualization",
		"writing_original_draft": "Writing - original draft",
		"writing_review_editing": "Writing - review & editing",
	}

	// Generate list of dropdown values, set selectedTerm in dropdown to "selected"
	// TODO: if we get a map back, we'll need to explicitly sort (numerical, alphabetically) since maps are hashmaps
	var terms []*listFormValues
	for key, term := range vocabulary[taxonomy] {
		selected := false
		for skey := range selectedTerms {
			if key == selectedTerms[skey] {
				selected = true
			}
		}

		terms = append(terms, &listFormValues{
			key,
			term,
			selected,
		})
	}

	return RenderPartial(f.render, "form/_list_multi_select", &listFormData{
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
