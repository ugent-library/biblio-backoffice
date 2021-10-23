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

// func RenderPartialFallback(r *render.Render, tmpl, tmplFallback string, data interface{}) (template.HTML, error) {
// 	buf := &bytes.Buffer{}
// 	var err error
// 	t := r.TemplateLookup(tmpl)
// 	if t == nil {
// 		t = r.TemplateLookup(tmplFallback)
// 	}
// 	if t != nil {
// 		err = t.Execute(buf, data)
// 	}
// 	return template.HTML(buf.String()), err
// }
