package display

import (
	"html/template"
	"io"
	"path"
	"strings"

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
	Render(*Display, io.Writer) error
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
	var buf strings.Builder

	for _, field := range s.Fields {
		if err := field.Render(s.Display, &buf); err != nil {
			return "", err
		}
	}

	return template.HTML(buf.String()), nil
}

type Text struct {
	Label         string
	List          bool
	Required      bool
	Tooltip       string
	Value         string
	Values        []string
	ValueTemplate string
}

type yieldHTML struct {
	Label    string
	List     bool
	Required bool
	Tooltip  string
	Values   []template.HTML
}

func (f *Text) Render(d *Display, w io.Writer) error {
	if f.Value != "" {
		f.Values = []string{f.Value}
	}

	if f.ValueTemplate != "" {
		return f.renderWithValueTemplate(d, w)
	}

	tmpl := path.Join("display", d.Theme, "text")
	return render.Templates().ExecuteTemplate(w, tmpl, f)
}

func (f *Text) renderWithValueTemplate(d *Display, w io.Writer) error {
	y := yieldHTML{
		Label:    f.Label,
		List:     f.List,
		Required: f.Required,
		Tooltip:  f.Tooltip,
	}

	for _, v := range f.Values {
		var buf strings.Builder
		if err := render.Templates().ExecuteTemplate(&buf, f.ValueTemplate, v); err != nil {
			return err
		}
		y.Values = append(y.Values, template.HTML(buf.String()))
	}

	tmpl := path.Join("display", d.Theme, "text")
	return render.Templates().ExecuteTemplate(w, tmpl, y)
}
