package presenters

import (
	"bytes"
	"fmt"
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

//
// Render each collapsible card
//

// Render the "Publication details" section under the Description tab

func (p *Publication) RenderDetails() template.HTML {
	if val, ok := p.Data()["type"]; ok {
		tpl := fmt.Sprintf("publication/details/_%s", val.(string))
		return p.renderPartial(tpl, p, "publication/details/_default")
	}

	return template.HTML("")
}

func (p *Publication) RenderConference() template.HTML {
	if val, ok := p.Data()["type"]; ok {
		tpl := fmt.Sprintf("publication/conference/_%s", val.(string))
		return p.renderPartial(tpl, p, "publication/conference/_default")
	}

	return template.HTML("")
}

//
// Render each field
//

// Render the "Publication type" field

func (p *Publication) RenderType() template.HTML {
	text := "-"
	if val, ok := p.Data()["type"]; ok {
		text = val.(string)
	}

	return p.renderPartial("part/_text", &PartData{
		Label: "Publication type",
		Text:  text,
	})
}

// Render the "DOI" field

func (p *Publication) RenderDOI() template.HTML {
	text := "-"
	if val, ok := p.Data()["doi"]; ok {
		text = val.(string)
	}

	return p.renderPartial("part/_text", &PartData{
		Label: "DOI",
		Text:  text,
	})
}

// Render the "ISSN/ISBN" field

func (p *Publication) RenderISXN() template.HTML {
	text := "-"
	if val, ok := p.Data()["issn"]; ok {
		text = val.(string)
	}

	return p.renderPartial("part/_text", &PartData{
		Label: "ISSN/ISBN",
		Text:  text,
	})
}

// Render the "Title" field

func (p *Publication) RenderTitle() template.HTML {
	text := "-"
	if val, ok := p.Data()["title"]; ok {
		text = val.(string)
	}

	return p.renderPartial("part/_text", &PartData{
		Label:    "Title",
		Text:     text,
		Required: true,
	})
}

// Render the "Alternative Title" field

func (p *Publication) RenderAlternativeTitle() template.HTML {
	text := "-"
	if val, ok := p.Data()["alternative_title"]; ok {
		text = val.(string)
	}

	return p.renderPartial("part/_text", &PartData{
		Label: "Alternative title",
		Text:  text,
	})
}

func (p *Publication) renderPartial(partial string, vars interface{}, fallbacks ...string) template.HTML {
	buf := &bytes.Buffer{}
	tmpl := p.Render.TemplateLookup(partial)
	// look for fallback partial if template doeesn't exist
	if tmpl == nil {
		for _, f := range fallbacks {
			if t := p.Render.TemplateLookup(f); t != nil {
				tmpl = t
				break
			}
		}
	}
	if tmpl != nil {
		tmpl.Execute(buf, vars)
	}
	return template.HTML(buf.String())
}
