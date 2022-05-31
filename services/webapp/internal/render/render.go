package render

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
)

var (
	TemplateDir = "templates/"
	TemplateExt = ".gohtml"
	FuncMaps    = []template.FuncMap{}
)

type View struct {
	Template *template.Template
}

func NewView(files ...string) View {
	addTemplateDirExt(files)
	tmpl := template.New("")
	for _, funcs := range FuncMaps {
		tmpl.Funcs(funcs)
	}
	tmpl = template.Must(tmpl.ParseFiles(files...))
	return View{Template: tmpl}
}

func (p View) Render(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")

	var buf bytes.Buffer
	if err := p.Template.ExecuteTemplate(&buf, name, data); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func addTemplateDirExt(files []string) {
	for i, f := range files {
		files[i] = path.Join(TemplateDir, f+TemplateExt)
	}
}
