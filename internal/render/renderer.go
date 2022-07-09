package render

// TODO use buffer pool

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Renderer struct {
	contentType      string
	dirs             []string
	exts             []string
	funcMaps         []template.FuncMap
	partialsTemplate *template.Template
	viewTemplates    map[string]*template.Template
}

func New() *Renderer {
	r := &Renderer{
		contentType:      "text/html; charset=utf-8",
		partialsTemplate: template.New(""),
		viewTemplates:    map[string]*template.Template{},
	}

	r.AddFuncs(template.FuncMap{
		"view": func(view string, data any) (template.HTML, error) {
			var b strings.Builder
			if err := r.ExecuteView(&b, view, data); err != nil {
				return "", err
			}
			return template.HTML(b.String()), nil
		},
		"layout": func(partial, view string, data any) (template.HTML, error) {
			var b strings.Builder
			if err := r.ExecuteLayout(&b, partial, view, data); err != nil {
				return "", err
			}
			return template.HTML(b.String()), nil
		},
		"partial": func(partial string, data any) (template.HTML, error) {
			var b strings.Builder
			if err := r.ExecutePartial(&b, partial, data); err != nil {
				return "", err
			}
			return template.HTML(b.String()), nil
		},
	})

	return r
}

func (r *Renderer) ContentType(contentType string) *Renderer {
	r.contentType = contentType
	return r
}

func (r *Renderer) AddExt(ext string) *Renderer {
	r.exts = append(r.exts, ext)
	return r
}

func (r *Renderer) AddDir(dir string) *Renderer {
	r.dirs = append(r.dirs, dir)
	return r
}

func (r *Renderer) AddFuncs(funcMap template.FuncMap) *Renderer {
	r.funcMaps = append(r.funcMaps, funcMap)
	return r
}

func (r *Renderer) MustParse() *Renderer {
	if _, err := r.Parse(); err != nil {
		panic(err)
	}
	return r
}

// TODO we don't need top keep all views in memory during parsing
func (r *Renderer) Parse() (*Renderer, error) {
	var (
		views            = map[string]string{}
		partials         = map[string]string{}
		partialsTemplate = template.New("")
		viewTemplates    = map[string]*template.Template{}
	)

	// read template files
	for _, dir := range r.dirs {
		if err := r.parseDir(dir, views, partials); err != nil {
			return r, err
		}
	}

	// create template that contains every partial
	for name, content := range partials {
		t := partialsTemplate.New(name)
		for _, funcs := range r.funcMaps {
			t.Funcs(funcs)
		}
		if _, err := t.Parse(content); err != nil {
			return r, err
		}
		partialsTemplate = t
	}

	// clone partials template to create view templates
	for name, content := range views {
		t, err := partialsTemplate.Clone()
		if err != nil {
			return r, err
		}
		t = t.New(name)
		for _, funcs := range r.funcMaps {
			t.Funcs(funcs)
		}
		if _, err := t.Parse(content); err != nil {
			return r, err
		}
		viewTemplates[name] = t
	}

	r.partialsTemplate = partialsTemplate
	r.viewTemplates = viewTemplates

	return r, nil
}

func (r *Renderer) parseDir(rootDir string, views, partials map[string]string) error {
	rootDir = filepath.Clean(rootDir)

	return filepath.Walk(rootDir, func(f string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		var ext string
		for _, e := range r.exts {
			if strings.HasSuffix(f, e) {
				ext = e
				break
			}
		}

		if ext == "" {
			return nil
		}

		content, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}

		name := f[len(rootDir)+1 : len(f)-len(ext)]
		dir, file := path.Split(name)

		// if template name starts with an underscore it's a partial
		if strings.HasPrefix(file, "_") {
			name = path.Join(dir, file[1:])
			partials[name] = string(content)
		} else {
			views[name] = string(content)
		}

		return nil
	})
}

func (r *Renderer) View(w http.ResponseWriter, view string, data any) {
	var b bytes.Buffer

	if err := r.ExecuteView(&b, view, data); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	r.SetContentType(w)
	io.Copy(w, &b)
}

func (r *Renderer) ExecuteView(w io.Writer, view string, data any) error {
	tmpl, ok := r.viewTemplates[view]
	if !ok {
		return fmt.Errorf("render: view '%s' not found", view)
	}

	if err := tmpl.ExecuteTemplate(w, view, data); err != nil {
		return errors.Wrapf(err, "render: Execute error, view '%s'", view)
	}

	return nil
}

func (r *Renderer) Layout(w http.ResponseWriter, partial, view string, data any) {
	var b bytes.Buffer

	if err := r.ExecuteLayout(&b, partial, view, data); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	r.SetContentType(w)
	io.Copy(w, &b)
}

func (r *Renderer) ExecuteLayout(w io.Writer, partial, view string, data any) error {
	tmpl, ok := r.viewTemplates[view]
	if !ok {
		return fmt.Errorf("render: view '%s' not found", view)
	}

	if err := tmpl.ExecuteTemplate(w, partial, data); err != nil {
		return errors.Wrapf(err, "render: ExecuteTemplate error, partial '%s' view '%s'", partial, view)
	}

	return nil
}

func (r *Renderer) Partial(w http.ResponseWriter, partial string, data any) {
	var b bytes.Buffer

	if err := r.ExecutePartial(&b, partial, data); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	r.SetContentType(w)
	io.Copy(w, &b)
}

func (r *Renderer) ExecutePartial(w io.Writer, partial string, data any) error {
	if err := r.partialsTemplate.ExecuteTemplate(w, partial, data); err != nil {
		return errors.Wrapf(err, "render: ExecuteTemplate error, partial '%s'", partial)
	}

	return nil
}

func (r *Renderer) SetContentType(w http.ResponseWriter) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", r.contentType)
	}
}

var defaultRenderer = New()

func ContentType(contentType string) *Renderer {
	return defaultRenderer.ContentType(contentType)
}

func AddExt(ext string) *Renderer {
	return defaultRenderer.AddExt(ext)
}

func AddDir(dir string) *Renderer {
	return defaultRenderer.AddDir(dir)
}

func AddFuncs(funcs template.FuncMap) *Renderer {
	return defaultRenderer.AddFuncs(funcs)
}

func MustParse() *Renderer {
	return defaultRenderer.MustParse()
}

func Parse() (*Renderer, error) {
	return defaultRenderer.Parse()
}

func View(w http.ResponseWriter, view string, data any) {
	defaultRenderer.View(w, view, data)
}

func ExecuteView(w io.Writer, view string, data any) error {
	return defaultRenderer.ExecuteView(w, view, data)
}

func Layout(w http.ResponseWriter, partial, view string, data any) {
	defaultRenderer.Layout(w, partial, view, data)
}

func ExecuteLayout(w io.Writer, partial, view string, data any) error {
	return defaultRenderer.ExecuteLayout(w, partial, view, data)
}

func Partial(w http.ResponseWriter, partial string, data any) {
	defaultRenderer.Partial(w, partial, data)
}

func ExecutePartial(w io.Writer, partial string, data any) error {
	return defaultRenderer.ExecutePartial(w, partial, data)
}

func SetContentType(w http.ResponseWriter) {
	defaultRenderer.SetContentType(w)
}
