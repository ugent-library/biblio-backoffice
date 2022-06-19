package form

import (
	"html/template"
	"path"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/render"
)

type Form struct {
	Theme  string
	Errors []string
	Fields []Field
}

type Field interface {
	FieldTemplate() string
}

func (f *Form) RenderHTML() (template.HTML, error) {
	var buf strings.Builder

	if len(f.Errors) > 0 {
		if err := render.Templates().ExecuteTemplate(&buf, "form/errors", f); err != nil {
			return "", err
		}
	}

	for _, field := range f.Fields {
		tmpl := path.Join("form", f.Theme, field.FieldTemplate())
		if err := render.Templates().ExecuteTemplate(&buf, tmpl, field); err != nil {
			return "", err
		}
	}

	return template.HTML(buf.String()), nil
}

type TextArea struct {
	Cols        int
	Error       string
	Label       string
	Name        string
	Placeholder string
	Required    bool
	Rows        int
	Tooltip     string
	Value       string
}

func (f *TextArea) FieldTemplate() string {
	return "text_area"
}

type Select struct {
	Cols        int
	Disabled    bool
	Error       string
	Label       string
	Name        string
	Options     []SelectOption
	Placeholder string
	Required    bool
	Tooltip     string
	Value       string
}

type SelectOption struct {
	Label string
	Value string
}

func (f *Select) FieldTemplate() string {
	return "select"
}
