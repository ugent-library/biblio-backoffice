package form

import (
	"html/template"
	"io"
	"path"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/render"
)

type Form struct {
	Theme    string
	Errors   []string
	Fields   []Field
	Sections []Section
}

func NewForm() *Form {
	return &Form{}
}

func (f *Form) WithTheme(theme string) *Form {
	f.Theme = theme

	return f
}

func (f *Form) WithErrors(errors []string) *Form {
	f.Errors = errors

	return f
}

func (f *Form) AddSection(fields ...Field) *Form {
	section := Section{
		Fields: fields,
		Form:   f,
	}
	f.Sections = append(f.Sections, section)

	return f
}

func (f *Form) RenderErrors() (template.HTML, error) {
	var buf strings.Builder

	if len(f.Errors) > 0 {
		if err := render.Templates().ExecuteTemplate(&buf, "form/errors", f); err != nil {
			return "", err
		}
	}

	return template.HTML(buf.String()), nil
}

type Section struct {
	Fields []Field
	Form   *Form
}

func (s *Section) Render() (template.HTML, error) {
	var buf strings.Builder

	for _, field := range s.Fields {
		if err := field.Render(s.Form, &buf); err != nil {
			return "", err
		}
	}

	return template.HTML(buf.String()), nil
}

type Field interface {
	Render(*Form, io.Writer) error
}

type Text struct {
	Name            string
	Value           string
	Label           string
	Tooltip         string
	Placeholder     string
	Required        bool
	Cols            int
	AutocompleteURL string
	Error           string
	ID              string
	Disabled        bool
}

func (f *Text) Render(form *Form, w io.Writer) error {
	return render.Templates().ExecuteTemplate(w, path.Join("form", form.Theme, "text"), f)
}

type TextRepeat struct {
	Name            string
	Values          []string
	Label           string
	Tooltip         string
	Placeholder     string
	Required        bool
	Cols            int
	AutocompleteURL string
	Error           string
	ID              string
	Disabled        bool
}

func (f *TextRepeat) Render(form *Form, w io.Writer) error {
	return render.Templates().ExecuteTemplate(w, path.Join("form", form.Theme, "text_repeat"), f)
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

func (f *TextArea) Render(form *Form, w io.Writer) error {
	return render.Templates().ExecuteTemplate(w, path.Join("form", form.Theme, "text_area"), f)
}

type Select struct {
	Template    string
	Cols        int
	Disabled    bool
	Error       string
	Label       string
	Name        string
	EmptyChoice bool
	Options     []SelectOption
	Placeholder string
	Required    bool
	Tooltip     string
	Value       string
	Vars        any
}

type SelectOption struct {
	Label string
	Value string
}

func (f *Select) Render(form *Form, w io.Writer) error {
	template := path.Join("form", form.Theme, "select")
	if f.Template != "" {
		template = path.Join("form", form.Theme, f.Template)
	}

	return render.Templates().ExecuteTemplate(w, template, f)
}

type SelectRepeat struct {
	Name            string
	values          []string
	Label           string
	Tooltip         string
	Placeholder     string
	Required        bool
	Checked         bool
	Options         []SelectOption
	EmptyChoice     bool
	Cols            int
	Rows            int
	AutocompleteURL string
	Error           string
	errorPointer    string
	ID              string
	Disabled        bool
}

func (f *SelectRepeat) Render(form *Form, w io.Writer) error {
	return render.Templates().ExecuteTemplate(w, path.Join("form", form.Theme, "select_repeat"), f)
}

type RadioButtonGroup struct {
	Name            string
	values          []string
	Label           string
	Tooltip         string
	Placeholder     string
	Required        bool
	Checked         bool
	Options         []SelectOption
	EmptyChoice     bool
	Cols            int
	Rows            int
	AutocompleteURL string
	Error           string
	ID              string
	Disabled        bool
}

func (f *RadioButtonGroup) Render(form *Form, w io.Writer) error {
	return render.Templates().ExecuteTemplate(w, path.Join("form", form.Theme, "radio_button_group"), f)
}

type Date struct {
	Name            string
	Value           string
	Label           string
	Tooltip         string
	Placeholder     string
	Required        bool
	Cols            int
	Rows            int
	Min             int
	Max             int
	AutocompleteURL string
	Error           string
	ID              string
	Disabled        bool
}

func (f *Date) Render(form *Form, w io.Writer) error {
	return render.Templates().ExecuteTemplate(w, path.Join("form", form.Theme, "date"), f)
}

type Checkbox struct {
	Name            string
	values          []string
	Label           string
	Tooltip         string
	Placeholder     string
	Required        bool
	Checked         bool
	Options         []SelectOption
	EmptyChoice     bool
	Cols            int
	Rows            int
	AutocompleteURL string
	Error           string
	ID              string
	Disabled        bool
}

func (f *Checkbox) Render(form *Form, w io.Writer) error {
	return render.Templates().ExecuteTemplate(w, path.Join("form", form.Theme, "checkbox"), f)
}
