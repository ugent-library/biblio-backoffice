package presenters

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/unrolled/render"
)

type PartData struct {
	Label        string
	Text         string
	Tooltip      string
	PrefilledDOI bool
	Required     bool
}

type Publication struct {
	*engine.Publication
	Render *render.Render
}

//
// Render each collapsible card
//

// Render the "Publication details" section under the Description tab

func (p *Publication) RenderDetails() template.HTML {
	if data, ok := p.Data()["type"]; ok {
		if val, ok := data.(string); ok {
			tpl := fmt.Sprintf("publication/details/_%s", val)
			return p.renderPartial(tpl, p)
		}
	}

	return template.HTML("")
}

func (p *Publication) RenderConference() template.HTML {
	if data, ok := p.Data()["type"]; ok {
		if val, ok := data.(string); ok {
			tpl := fmt.Sprintf("publication/conference/_%s", val)
			return p.renderPartial(tpl, p)
		}
	}

	return template.HTML("")
}

//
// Render each field
//

// Render the "Publication type" field

func (p *Publication) RenderType() template.HTML {
	text := "-"
	if data, ok := p.Data()["type"]; ok {
		if val, ok := data.(string); ok {
			text = val
		}
	}

	return p.renderPartial("part/_text", &PartData{
		Label: "Publication type",
		Text:  text,
	})
}

// Render the "DOI" field

func (p *Publication) RenderDOI() template.HTML {
	text := "-"
	if data, ok := p.Data()["doi"]; ok {
		if val, ok := data.(string); ok {
			text = val
		}
	}

	return p.renderPartial("part/_text", &PartData{
		Label: "DOI",
		Text:  text,
	})
}

// Render the "ISSN/ISBN" field

func (p *Publication) RenderISXN() template.HTML {
	text := "-"

	if data, ok := p.Data()["issn"]; ok {
		if val, ok := data.(string); ok {
			text = val
		}
	}

	return p.renderPartial("part/_text", &PartData{
		Label: "ISSN/ISBN",
		Text:  text,
	})
}

// Render the "Title" field

func (p *Publication) RenderTitle() template.HTML {
	text := "-"

	if data, ok := p.Data()["title"]; ok {
		if val, ok := data.(string); ok {
			text = val
		}
	}

	return p.renderPartial("part/_text", &PartData{
		Label:        "Title",
		Text:         text,
		PrefilledDOI: true,
		Required:     true,
	})
}

// Render the "Alternative Title" field

func (p *Publication) RenderAlternativeTitle() template.HTML {
	text := "-"

	if data, ok := p.Data()["alternative_title"]; ok {
		if val, ok := data.(string); ok {
			text = val
		}
	}

	return p.renderPartial("part/_text", &PartData{
		Label: "Alternative title",
		Text:  text,
	})
}

func (p *Publication) renderPartial(name string, vars interface{}) template.HTML {
	buf := &bytes.Buffer{}
	tmpl := p.Render.TemplateLookup(name)
	tmpl.Execute(buf, vars)
	return template.HTML(buf.String())
}
