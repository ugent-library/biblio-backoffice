package render

// derived from https://github.com/biz/templates
// TODO use buffer pool

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
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
	views            map[string]string
	partials         map[string]string
}

func New() *Renderer {
	r := &Renderer{
		contentType:      "text/html; charset=utf-8",
		partialsTemplate: template.New(""),
		viewTemplates:    map[string]*template.Template{},
		views:            map[string]string{},
		partials:         map[string]string{},
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

func (r *Renderer) AddView(name string, content string) *Renderer {
	r.views[name] = content
	return r
}

func (r *Renderer) AddPartial(name string, content string) *Renderer {
	r.partials[name] = content
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

func (r *Renderer) Parse() (*Renderer, error) {
	// reset templates
	r.partialsTemplate = template.New("")
	r.viewTemplates = map[string]*template.Template{}

	// read template files
	for _, dir := range r.dirs {
		if err := r.parseDir(dir); err != nil {
			return r, err
		}
	}

	// create template that contains every partial
	for _, funcs := range r.funcMaps {
		r.partialsTemplate.Funcs(funcs)
	}
	for name, content := range r.partials {
		t, err := r.partialsTemplate.New(name).Parse(content)
		if err != nil {
			return r, err
		}
		r.partialsTemplate = t
	}

	// clone partials template to create view templates
	for name, content := range r.views {
		t, err := r.partialsTemplate.Clone()
		if err != nil {
			return r, err
		}
		t, err = t.Parse(content)
		if err != nil {
			return r, err
		}
		r.viewTemplates[name] = t
	}

	return r, nil
}

func (r *Renderer) parseDir(rootDir string) error {
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
			r.partials[name] = string(content)
		} else {
			r.views[name] = string(content)
		}

		return nil
	})
}

func (r *Renderer) MustRenderView(w http.ResponseWriter, view string, data any) {
	if err := r.RenderView(w, view, data); err != nil {
		panic(err)
	}
}

func (r *Renderer) RenderView(w http.ResponseWriter, view string, data any) error {
	r.SetContentType(w)
	return r.ExecuteView(w, view, data)
}

func (r *Renderer) ExecuteView(w io.Writer, view string, data any) error {
	tmpl, ok := r.viewTemplates[view]
	if !ok {
		return fmt.Errorf("render: view '%v' not found", view)
	}

	if err := tmpl.Execute(w, data); err != nil {
		return errors.Wrapf(err, "render: Execute error, view '%v'", view)
	}

	return nil
}

func (r *Renderer) MustRenderLayout(w http.ResponseWriter, partial, view string, data any) {
	if err := r.RenderLayout(w, partial, view, data); err != nil {
		panic(err)
	}
}

func (r *Renderer) RenderLayout(w http.ResponseWriter, partial, view string, data any) error {
	r.SetContentType(w)
	return r.ExecuteLayout(w, partial, view, data)
}

func (r *Renderer) ExecuteLayout(w io.Writer, partial, view string, data any) error {
	tmpl, ok := r.viewTemplates[view]
	if !ok {
		return fmt.Errorf("render: view '%v' not found", view)
	}

	if err := tmpl.ExecuteTemplate(w, partial, data); err != nil {
		return errors.Wrapf(err, "render: ExecuteTemplate error, partial '%v' view '%v'", partial, view)
	}

	return nil
}

func (r *Renderer) MustRenderPartial(w http.ResponseWriter, partial string, data any) {
	if err := r.RenderPartial(w, partial, data); err != nil {
		panic(err)
	}
}

func (r *Renderer) RenderPartial(w http.ResponseWriter, partial string, data any) error {
	r.SetContentType(w)
	return r.ExecutePartial(w, partial, data)
}

func (r *Renderer) ExecutePartial(w io.Writer, partial string, data any) error {
	if err := r.partialsTemplate.ExecuteTemplate(w, partial, data); err != nil {
		return errors.Wrapf(err, "render: ExecuteTemplate error, partial '%v'", partial)
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

func AddView(name, content string) *Renderer {
	return defaultRenderer.AddView(name, content)
}

func AddPartial(name, content string) *Renderer {
	return defaultRenderer.AddPartial(name, content)
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

func MustRenderView(w http.ResponseWriter, view string, data any) {
	defaultRenderer.MustRenderView(w, view, data)
}

func RenderView(w http.ResponseWriter, view string, data any) error {
	return defaultRenderer.RenderView(w, view, data)
}

func ExecuteView(w io.Writer, view string, data any) error {
	return defaultRenderer.ExecuteView(w, view, data)
}

func MustRenderLayout(w http.ResponseWriter, partial, view string, data any) {
	defaultRenderer.MustRenderLayout(w, partial, view, data)
}

func RenderLayout(w http.ResponseWriter, partial, view string, data any) error {
	return defaultRenderer.RenderLayout(w, partial, view, data)
}

func ExecuteLayout(w io.Writer, partial, view string, data any) error {
	return defaultRenderer.ExecuteLayout(w, partial, view, data)
}

func MustRenderPartial(w http.ResponseWriter, partial string, data any) {
	defaultRenderer.MustRenderPartial(w, partial, data)
}

func RenderPartial(w http.ResponseWriter, partial string, data any) error {
	return defaultRenderer.RenderPartial(w, partial, data)
}

func ExecutePartial(w io.Writer, partial string, data any) error {
	return defaultRenderer.ExecutePartial(w, partial, data)
}

func SetContentType(w http.ResponseWriter) {
	defaultRenderer.SetContentType(w)
}
