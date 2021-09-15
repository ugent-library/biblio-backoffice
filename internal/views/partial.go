package views

import (
	"bytes"
	"html/template"

	"github.com/unrolled/render"
)

func RenderPartial(r *render.Render, tmpl string, data interface{}) template.HTML {
	buf := &bytes.Buffer{}
	if tmpl := r.TemplateLookup(tmpl); tmpl != nil {
		tmpl.Execute(buf, data)
	}
	return template.HTML(buf.String())
}
