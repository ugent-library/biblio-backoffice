package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

var (
	DefaultDir       = "views"
	DefaultExt       = ".gohtml"
	DefaultLayoutExt = ".layout"
)

type Renderer struct {
	dir              string
	ext              string
	layoutExt        string
	contentType      string
	funcMaps         []template.FuncMap
	partialsTemplate *template.Template
	viewTemplates    map[string]*template.Template
	bufPool          sync.Pool
}

func New() *Renderer {
	r := &Renderer{
		dir:              filepath.Clean(DefaultDir),
		ext:              DefaultExt,
		layoutExt:        DefaultLayoutExt,
		contentType:      "text/html; charset=utf-8",
		partialsTemplate: template.New(""),
		bufPool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}

	r.Funcs(template.FuncMap{
		"view": func(view string, data any) (template.HTML, error) {
			b := r.bufPool.Get().(*bytes.Buffer)
			defer func() {
				b.Reset()
				r.bufPool.Put(b)
			}()

			if err := r.ExecuteView(b, view, data); err != nil {
				return "", err
			}
			return template.HTML(b.String()), nil
		},
		"layout": func(partial, view string, data any) (template.HTML, error) {
			b := r.bufPool.Get().(*bytes.Buffer)
			defer func() {
				b.Reset()
				r.bufPool.Put(b)
			}()

			if err := r.ExecuteLayout(b, partial, view, data); err != nil {
				return "", err
			}
			return template.HTML(b.String()), nil
		},
		"partial": func(partial string, data any) (template.HTML, error) {
			b := r.bufPool.Get().(*bytes.Buffer)
			defer func() {
				b.Reset()
				r.bufPool.Put(b)
			}()

			if err := r.ExecutePartial(b, partial, data); err != nil {
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

func (r *Renderer) Ext(ext string) *Renderer {
	r.ext = ext
	return r
}

func (r *Renderer) Dir(dir string) *Renderer {
	r.dir = filepath.Clean(dir)
	return r
}

func (r *Renderer) Funcs(funcMap template.FuncMap) *Renderer {
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
	var (
		layoutsTemplate = template.New("")
		viewTemplates   = map[string]*template.Template{}
	)

	for _, funcs := range r.funcMaps {
		layoutsTemplate.Funcs(funcs)
	}

	// parse layouts
	err := filepath.WalkDir(r.dir, func(f string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(f, r.layoutExt+r.ext) {
			return nil
		}

		content, err := os.ReadFile(f)
		if err != nil {
			return err
		}

		name := f[len(r.dir)+1 : len(f)-len(r.layoutExt+r.ext)]

		t, err := layoutsTemplate.New(name).Parse(string(content))
		if err != nil {
			return err
		}
		layoutsTemplate = t

		return nil
	})

	if err != nil {
		return r, err
	}

	partialsTemplate, err := layoutsTemplate.Clone()
	if err != nil {
		return r, err
	}

	// parse views and partials
	err = filepath.WalkDir(r.dir, func(f string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || strings.HasSuffix(f, r.layoutExt+r.ext) || !strings.HasSuffix(f, r.ext) {
			return nil
		}

		content, err := os.ReadFile(f)
		if err != nil {
			return err
		}

		name := f[len(r.dir)+1 : len(f)-len(r.ext)]
		dir, file := path.Split(name)

		// if template name starts with an underscore it's a partial
		// partials share one template
		// views have their own template but have access to layouts
		if strings.HasPrefix(file, "_") {
			name = path.Join(dir, file[1:])

			t, err := partialsTemplate.New(name).Parse(string(content))
			if err != nil {
				return err
			}
			partialsTemplate = t
		} else {
			t, err := layoutsTemplate.Clone()
			if err != nil {
				return err
			}

			t, err = t.New(name).Parse(string(content))
			if err != nil {
				return err
			}

			viewTemplates[name] = t
		}

		return nil
	})

	if err != nil {
		return r, err
	}

	r.partialsTemplate = partialsTemplate
	r.viewTemplates = viewTemplates

	return r, nil
}

func (r *Renderer) View(w http.ResponseWriter, view string, data any) {
	b := r.bufPool.Get().(*bytes.Buffer)
	defer func() {
		b.Reset()
		r.bufPool.Put(b)
	}()

	if err := r.ExecuteView(b, view, data); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	r.SetContentType(w)
	io.Copy(w, b)
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
	b := r.bufPool.Get().(*bytes.Buffer)
	defer func() {
		b.Reset()
		r.bufPool.Put(b)
	}()

	if err := r.ExecuteLayout(b, partial, view, data); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	r.SetContentType(w)
	io.Copy(w, b)
}

func (r *Renderer) ExecuteLayout(w io.Writer, layout, view string, data any) error {
	tmpl, ok := r.viewTemplates[view]
	if !ok {
		return fmt.Errorf("render: view '%s' not found", view)
	}
	if err := tmpl.ExecuteTemplate(w, layout, data); err != nil {
		return errors.Wrapf(err, "render: ExecuteTemplate error, partial '%s' view '%s'", layout, view)
	}

	return nil
}

func (r *Renderer) Partial(w http.ResponseWriter, partial string, data any) {
	b := r.bufPool.Get().(*bytes.Buffer)
	defer func() {
		b.Reset()
		r.bufPool.Put(b)
	}()

	if err := r.ExecutePartial(b, partial, data); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	r.SetContentType(w)
	io.Copy(w, b)
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

func Ext(ext string) *Renderer {
	return defaultRenderer.Ext(ext)
}

func Dir(dir string) *Renderer {
	return defaultRenderer.Dir(dir)
}

func Funcs(funcs template.FuncMap) *Renderer {
	return defaultRenderer.Funcs(funcs)
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
