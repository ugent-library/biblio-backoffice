package views

import (
	"bytes"
	"html/template"

	"github.com/unrolled/render"
)

func RenderPartial(r *render.Render, tmpl string, data interface{}) (template.HTML, error) {
	buf := &bytes.Buffer{}
	var err error
	if t := r.TemplateLookup(tmpl); t != nil {
		err = t.Execute(buf, data)
	}
	return template.HTML(buf.String()), err
}
