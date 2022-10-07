package form

import (
	"bytes"
	"html/template"
	"io"
	"path"

	"github.com/ugent-library/biblio-backend/internal/render"
)

type Errors []string

func (e Errors) Render() (template.HTML, error) {
	b := &bytes.Buffer{}

	if len(e) > 0 {
		if err := render.ExecuteView(b, "form/errors", e); err != nil {
			return "", err
		}
	}

	return template.HTML(b.String()), nil
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
	b := &bytes.Buffer{}

	for _, field := range s.Fields {
		if err := field.Render(s.Form.Theme, b); err != nil {
			return "", err
		}
	}

	return template.HTML(b.String()), nil
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
	Help            string
	Readonly        bool
	Required        bool
	Template        string
	Tooltip         string
	Value           string
	Vars            any
}

func (f *Text) Render(theme string, w io.Writer) error {
	t := "text"
	if f.Template != "" {
		t = f.Template
	}
	return render.ExecuteView(w, path.Join("form", theme, t), f)
}

type TextRepeat struct {
	AutocompleteURL string
	Cols            int
	Disabled        bool
	Error           string
	ID              string
	Label           string
	Name            string
	Help            string
	Readonly        bool
	Required        bool
	Template        string
	Tooltip         string
	Values          []string
	Vars            any
}

func (f *TextRepeat) Render(theme string, w io.Writer) error {
	t := "text_repeat"
	if f.Template != "" {
		t = f.Template
	}
	return render.ExecuteView(w, path.Join("form", theme, t), f)
}

type TextArea struct {
	Cols     int
	Error    string
	Label    string
	Name     string
	Help     string
	Required bool
	Rows     int
	Template string
	Tooltip  string
	Value    string
	Vars     any
}

func (f *TextArea) Render(theme string, w io.Writer) error {
	t := "text_area"
	if f.Template != "" {
		t = f.Template
	}
	return render.ExecuteView(w, path.Join("form", theme, t), f)
}

type Select struct {
	Cols        int
	Disabled    bool
	EmptyOption bool
	Error       string
	Help        string
	Label       string
	Name        string
	Options     []SelectOption
	Required    bool
	Template    string
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
	return render.ExecuteView(w, path.Join("form", theme, t), f)
}

type SelectRepeat struct {
	Cols        int
	EmptyOption bool
	Error       string
	Label       string
	Name        string
	Options     []SelectOption
	Required    bool
	Template    string
	Tooltip     string
	Values      []string
	Vars        any
}

func (f *SelectRepeat) Render(theme string, w io.Writer) error {
	t := "select_repeat"
	if f.Template != "" {
		t = f.Template
	}
	return render.ExecuteView(w, path.Join("form", theme, t), f)
}

type Date struct {
	Cols     int
	Disabled bool
	Error    string
	Label    string
	Max      string
	Min      string
	Name     string
	Required bool
	Template string
	Tooltip  string
	Value    string
	Vars     any
}

func (f *Date) Render(theme string, w io.Writer) error {
	t := "date"
	if f.Template != "" {
		t = f.Template
	}
	return render.ExecuteView(w, path.Join("form", theme, t), f)
}

type Checkbox struct {
	Checked  bool
	Cols     int
	Error    string
	Label    string
	Name     string
	Template string
	Value    string
	Vars     any
}

func (f *Checkbox) Render(theme string, w io.Writer) error {
	t := "checkbox"
	if f.Template != "" {
		t = f.Template
	}
	return render.ExecuteView(w, path.Join("form", theme, t), f)
}

var HiddenFieldTemplate = template.Must(template.New("").Parse(`<input type="hidden" name="{{.Name}}" value="{{.Value}}">`))

type Hidden struct {
	Name  string
	Value string
}

func (f *Hidden) Render(theme string, w io.Writer) error {
	return HiddenFieldTemplate.Execute(w, f)
}

type RadioButtonGroup struct {
	Cols     int
	Error    string
	Label    string
	Name     string
	Options  []SelectOption
	Required bool
	Template string
	Tooltip  string
	Value    string
	Vars     any
}

func (f *RadioButtonGroup) Render(theme string, w io.Writer) error {
	t := "radio_button_group"
	if f.Template != "" {
		t = f.Template
	}
	return render.ExecuteView(w, path.Join("form", theme, t), f)
}
