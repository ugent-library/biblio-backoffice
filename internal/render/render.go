package render

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	TemplateDir = "templates/"
	TemplateExt = ".gohtml"
	FuncMaps    = []template.FuncMap{}
	templates   *template.Template
	parseOnce   sync.Once
)

func Templates() *template.Template {
	parseOnce.Do(func() {
		templates = template.Must(parseTemplates(TemplateDir, TemplateExt, FuncMaps))
	})
	return templates
}

func Render(w http.ResponseWriter, tmpl string, data any) {
	w.Header().Set("Content-Type", "text/html")

	var buf bytes.Buffer

	if err := Templates().ExecuteTemplate(&buf, tmpl, data); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buf)

}

func Wrap(w http.ResponseWriter, wrapperTmpl, tmpl string, data any) {
	w.Header().Set("Content-Type", "text/html")

	var (
		buf        strings.Builder
		wrapperBuf bytes.Buffer
	)

	if err := Templates().ExecuteTemplate(&buf, tmpl, data); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	wrapperData := struct {
		Wrapped     any
		WrappedHTML template.HTML
	}{
		Wrapped:     data,
		WrappedHTML: template.HTML(buf.String()),
	}

	if err := Templates().ExecuteTemplate(&wrapperBuf, wrapperTmpl, wrapperData); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.Copy(w, &wrapperBuf)
}

func parseTemplates(rootDir, ext string, funcMaps []template.FuncMap) (*template.Template, error) {
	cleanRootDir := filepath.Clean(rootDir)
	pathStart := len(cleanRootDir) + 1
	tmpl := template.New("")

	err := filepath.Walk(cleanRootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ext) {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			pathEnd := len(path) - len(ext)
			name := path[pathStart:pathEnd]
			t := tmpl.New(name)
			for _, funcs := range funcMaps {
				t.Funcs(funcs)
			}
			if _, err := t.Parse(string(b)); err != nil {
				return err
			}
		}

		return nil
	})

	return tmpl, err
}
