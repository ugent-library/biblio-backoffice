package render

import (
	"html/template"
	"io"
	"strings"
)

type Form struct {
	Errors []string
	Fields []FormField
}

type FormField interface {
	Render(io.Writer) error
}

func (f *Form) RenderHTML() (template.HTML, error) {
	var buf strings.Builder

	if len(f.Errors) > 0 {
		if err := Templates().ExecuteTemplate(&buf, "form/errors", f); err != nil {
			return "", err
		}
	}

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
	return Templates().ExecuteTemplate(w, "form/text_area", f)
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
	return Templates().ExecuteTemplate(w, "form/select", f)
}
