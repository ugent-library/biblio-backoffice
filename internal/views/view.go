package views

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/unrolled/render"
)

type Renderer interface {
	Render(http.ResponseWriter, *http.Request, int, interface{})
	RenderPartial(*http.Request, interface{}) template.HTML
}

type View struct {
	render   *render.Render
	template string
}

func NewView(r *render.Render, t string) *View {
	return &View{r, t}
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	v.render.HTML(w, status, v.template, newContext(r, data))
}

func (v *View) RenderPartial(r *http.Request, data interface{}) template.HTML {
	buf := &bytes.Buffer{}
	if tmpl := v.render.TemplateLookup(v.template); tmpl != nil {
		tmpl.Execute(buf, newContext(r, data))
	}
	return template.HTML(buf.String())
}
