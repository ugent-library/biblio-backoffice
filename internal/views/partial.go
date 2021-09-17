package views

import (
	"bytes"
	"html/template"

	"github.com/unrolled/render"
)

func RenderPartial(r *render.Render, tmpl string, data interface{}) (template.HTML, error) {
	buf := &bytes.Buffer{}
	var err error
	if tmpl := r.TemplateLookup(tmpl); tmpl != nil {
		err = tmpl.Execute(buf, data)
	}
	return template.HTML(buf.String()), err
}
