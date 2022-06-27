package form

import (
	"html/template"
	"io"
	"path"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/render"
)

type Errors []string

func (e Errors) Render() (template.HTML, error) {
	var buf strings.Builder

	if len(e) > 0 {
		if err := render.Templates().ExecuteTemplate(&buf, "form/errors", e); err != nil {
			return "", err
		}
	}

	return template.HTML(buf.String()), nil
}

type Form struct {
	Theme    string
	Errors   Errors
	Sections []Section
}

func New() *Form {
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
	f.Sections = append(f.Sections, Section{
		Fields: fields,
		Form:   f,
	})

	return f
}

type Section struct {
	Fields []Field
	Form   *Form
}

func (s *Section) Render() (template.HTML, error) {
	var buf strings.Builder

	for _, field := range s.Fields {
		if err := field.Render(s.Form.Theme, &buf); err != nil {
			return "", err
		}
	}

	return template.HTML(buf.String()), nil
}

type Field interface {
	Render(string, io.Writer) error
}

type Text struct {
	AutocompleteURL string
	Cols            int
	Disabled        bool
	Error           string
	ID              string
	Label           string
	Name            string
	Placeholder     string
	Readonly        bool
	Required        bool
	Template        string
	Tooltip         string
	Value           string
}

func (f *Text) Render(theme string, w io.Writer) error {
	return render.Templates().ExecuteTemplate(w, path.Join("form", theme, "text"), f)
}

type TextRepeat struct {
	AutocompleteURL string
	Cols            int
	Disabled        bool
	Error           string
	ID              string
	Label           string
	Name            string
	Placeholder     string
	Readonly        bool
	Required        bool
	Tooltip         string
	Values          []string
}

func (f *TextRepeat) Render(theme string, w io.Writer) error {
	return render.Templates().ExecuteTemplate(w, path.Join("form", theme, "text_repeat"), f)
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

func (f *TextArea) Render(theme string, w io.Writer) error {
	return render.Templates().ExecuteTemplate(w, path.Join("form", theme, "text_area"), f)
}

type Select struct {
	Template    string
	Cols        int
	Disabled    bool
	Error       string
	Label       string
	Name        string
	EmptyOption bool
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

func (f *Select) Render(theme string, w io.Writer) error {
	t := "select"
	if f.Template != "" {
		t = f.Template
	}

	tmpl := path.Join("form", theme, t)
	return render.Templates().ExecuteTemplate(w, tmpl, f)
}

type SelectRepeat struct {
	Name            string
	values          []string
	Label           string
	Tooltip         string
	Placeholder     string
	Required        bool
	Options         []SelectOption
	EmptyOption     bool
	Cols            int
	Rows            int
	AutocompleteURL string
	Error           string
	ID              string
	Disabled        bool
}

func (f *SelectRepeat) Render(theme string, w io.Writer) error {
	return render.Templates().ExecuteTemplate(w, path.Join("form", theme, "select_repeat"), f)
}

type RadioButtonGroup struct {
	Name            string
	values          []string
	Label           string
	Tooltip         string
	Placeholder     string
	Required        bool
	Options         []SelectOption
	EmptyOption     bool
	Cols            int
	Rows            int
	AutocompleteURL string
	Error           string
	ID              string
	Disabled        bool
}

func (f *RadioButtonGroup) Render(theme string, w io.Writer) error {
	return render.Templates().ExecuteTemplate(w, path.Join("form", theme, "radio_button_group"), f)
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

func (f *Date) Render(theme string, w io.Writer) error {
	return render.Templates().ExecuteTemplate(w, path.Join("form", theme, "date"), f)
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
	EmptyOption     bool
	Cols            int
	Rows            int
	AutocompleteURL string
	Error           string
	ID              string
	Disabled        bool
}

func (f *Checkbox) Render(theme string, w io.Writer) error {
	return render.Templates().ExecuteTemplate(w, path.Join("form", theme, "checkbox"), f)
}

var HiddenFieldTemplate = template.Must(template.New("").Parse(`<input type="hidden" name="{{.Name}}" value="{{.Value}}">`))

type Hidden struct {
	Name  string
	Value string
}

func (f *Hidden) Render(theme string, w io.Writer) error {
	return HiddenFieldTemplate.Execute(w, f)
}
