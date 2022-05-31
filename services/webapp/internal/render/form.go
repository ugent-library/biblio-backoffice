package render

import (
	"html/template"
	"io"
	"strings"
)

var (
	formTemplates = template.Must(template.New("").ParseGlob("templates/form/*.gohtml"))
)

type Form struct {
	Fields []FormField
}

type FormField interface {
	Render(io.Writer) error
}

func (f *Form) RenderHTML() (template.HTML, error) {
	var buf strings.Builder
	for _, field := range f.Fields {
		if err := field.Render(&buf); err != nil {
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

func (f *TextArea) Render(w io.Writer) error {
	return formTemplates.ExecuteTemplate(w, "text-area-field", f)
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

func (f *Select) Render(w io.Writer) error {
	return formTemplates.ExecuteTemplate(w, "select-field", f)
}
