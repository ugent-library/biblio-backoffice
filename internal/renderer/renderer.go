package renderer

import (
	"strings"

	"github.com/ugent-library/biblio-backend/internal/fields"
)

// Render an entire tree

func DoRender(element interface{}) string {
	var output []string

	// Process a single field
	field, ok := element.(fields.Field)
	if ok {
		output = append(output, field.Render())
		return strings.Join(output, "")
	}

	// Process a fieldSet containing fields
	fieldset, ok := element.(fields.FieldSet)
	if ok {
		var children []string
		for _, field := range fieldset.Items {
			children = append(children, DoRender(field))
		}

		output = append(output, fieldset.Render(children))
		return strings.Join(output, "")
	}

	// Start from root
	tree, ok := element.(map[string]fields.FieldSet)
	if ok {
		for _, fieldset := range tree {
			output = append(output, DoRender(fieldset))
		}
		return strings.Join(output, "")
	}

	return strings.Join(output, "")
}

// Render a single element

func render() {
}
