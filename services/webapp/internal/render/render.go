package render

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strings"
)

var (
	TemplateDir = "templates/"
	TemplateExt = ".gohtml"
	FuncMaps    = []template.FuncMap{}

	formTemplates = template.Must(template.New("").ParseGlob("templates/form/*.gohtml"))
)

type Partial struct {
	Name     string
	Template *template.Template
}

func NewPartial(name string, files ...string) Partial {
	addTemplateDirExt(files)
	tmpl := template.New("")
	for _, funcs := range FuncMaps {
		tmpl.Funcs(funcs)
	}
	tmpl = template.Must(tmpl.ParseFiles(files...))
	return Partial{Name: name, Template: tmpl}
}

func (p Partial) Render(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html")

	var buf bytes.Buffer
	if err := p.Template.ExecuteTemplate(&buf, p.Name, data); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

type FormField interface {
	Render(io.Writer) error
}

type Form struct {
	Fields []FormField
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

func addTemplateDirExt(files []string) {
	for i, f := range files {
		files[i] = path.Join(TemplateDir, f+TemplateExt)
	}
}
