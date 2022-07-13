package display

import (
	"bytes"
	"html/template"
	"io"
	"path"

	"github.com/ugent-library/biblio-backend/internal/render"
)

type Display struct {
	Theme    string
	Sections []Section
}

type Section struct {
	Display *Display
	Fields  []Field
}

type Field interface {
	Render(string, io.Writer) error
}

func New() *Display {
	return &Display{}
}

func (d *Display) WithTheme(theme string) *Display {
	d.Theme = theme

	return d
}

func (d *Display) AddSection(fields ...Field) *Display {
	d.Sections = append(d.Sections, Section{
		Fields:  fields,
		Display: d,
	})

	return d
}

func (s Section) Render() (template.HTML, error) {
	b := &bytes.Buffer{}

	for _, field := range s.Fields {
		if err := field.Render(s.Display.Theme, b); err != nil {
			return "", err
		}
	}

	return template.HTML(b.String()), nil
}

type Text struct {
	Label         string
	Required      bool
	Tooltip       string
	Value         string
	ValueTemplate string
}

type yieldHTML struct {
	Label    string
	List     bool
	Required bool
	Tooltip  string
	Value    template.HTML
}

func (f *Text) Render(theme string, w io.Writer) error {
	if f.ValueTemplate != "" {
		return f.renderWithValueTemplate(theme, w)
	}

	tmpl := path.Join("display", theme, "text")
	return render.ExecuteView(w, tmpl, f)
}

func (f *Text) renderWithValueTemplate(theme string, w io.Writer) error {
	y := yieldHTML{
		Label:    f.Label,
		Required: f.Required,
		Tooltip:  f.Tooltip,
	}

	b := &bytes.Buffer{}
	if err := render.ExecuteView(b, f.ValueTemplate, f.Value); err != nil {
		return err
	}
	y.Value = template.HTML(b.String())

	tmpl := path.Join("display", theme, "text")
	return render.ExecuteView(w, tmpl, y)
}

type List struct {
	Inline        bool
	Label         string
	Required      bool
	Tooltip       string
	Values        []string
	ValueTemplate string
}

type yieldListHTML struct {
	Inline   bool
	Label    string
	Required bool
	Tooltip  string
	Values   []template.HTML
}

func (f *List) Render(theme string, w io.Writer) error {
	if f.ValueTemplate != "" {
		return f.renderWithValueTemplate(theme, w)
	}

	tmpl := path.Join("display", theme, "list")
	return render.ExecuteView(w, tmpl, f)
}

func (f *List) renderWithValueTemplate(theme string, w io.Writer) error {
	y := yieldListHTML{
		Inline:   f.Inline,
		Label:    f.Label,
		Required: f.Required,
		Tooltip:  f.Tooltip,
	}

	for _, v := range f.Values {
		b := &bytes.Buffer{}
		if err := render.ExecuteView(b, f.ValueTemplate, v); err != nil {
			return err
		}
		y.Values = append(y.Values, template.HTML(b.String()))
	}

	tmpl := path.Join("display", theme, "list")
	return render.ExecuteView(w, tmpl, y)
}
