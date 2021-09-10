package presenters

import (
	"bytes"
	"html/template"

	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/unrolled/render"
)

type PartData struct {
	Label    string
	Text     string
	Tooltip  string
	Required bool
}

type Publication struct {
	*engine.Publication
	Render *render.Render
}

func (p *Publication) RenderDetails() template.HTML {
	partial := p.Data()["type"].(string)
	return p.renderPartial(partial, p)
}

func (p *Publication) RenderType() template.HTML {
	return p.renderPartial("part/_text", &PartData{
		Label: "Publication type",
		Text:  p.Data()["type"].(string),
	})
}

func (p *Publication) RenderDOI() template.HTML {
	return p.renderPartial("part/_text", &PartData{
		Label: "DOI",
		Text:  p.Data()["doi"].(string),
	})
}

func (p *Publication) renderPartial(name string, vars interface{}) template.HTML {
	buf := &bytes.Buffer{}
	tmpl := p.Render.TemplateLookup(name)
	tmpl.Execute(buf, vars)
	return template.HTML(buf.String())
}
