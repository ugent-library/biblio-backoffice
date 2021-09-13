package renderers

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
)

type textData struct {
	Label    string
	Text     string
	Tooltip  string
	Required bool
}

type Publication struct {
	Render *render.Render
}

//
// Render each collapsible card
//

// Render the "Publication details" section under the Description tab

func (r *Publication) Details(p *models.Publication) template.HTML {
	return r.partial(fmt.Sprintf("publication/details/_%s", p.Type), p)
}

func (r *Publication) Conference(p *models.Publication) template.HTML {
	return r.partial(fmt.Sprintf("publication/conference/_%s", p.Type), p)
}

//
// Render each field
//

// Render the "Publication type" field

func (r *Publication) Type(p *models.Publication) template.HTML {
	return r.partial("part/_text", &textData{
		Label: "Publication type",
		Text:  p.Type,
	})
}

// Render the "DOI" field

func (r *Publication) DOI(p *models.Publication) template.HTML {
	text := p.DOI
	if text == "" {
		text = "-"
	}

	return r.partial("part/_text", &textData{
		Label: "DOI",
		Text:  text,
	})
}

// Render the "ISSN/ISBN" field

func (r *Publication) ISXN(p *models.Publication) template.HTML {
	text := "-"
	if len(p.ISSN) > 0 {
		text = p.ISSN[0]
	}

	return r.partial("part/_text", &textData{
		Label: "ISSN/ISBN",
		Text:  text,
	})
}

// Render the "Title" field

func (r *Publication) Title(p *models.Publication) template.HTML {
	text := p.Title
	if text == "" {
		text = "-"
	}

	return r.partial("part/_text", &textData{
		Label:    "Title",
		Text:     text,
		Required: true,
	})
}

// Render the "Alternative Title" field

func (r *Publication) AlternativeTitle(p *models.Publication) template.HTML {
	text := "-"
	if len(p.AlternativeTitle) > 0 {
		text = p.AlternativeTitle[0]
	}

	return r.partial("part/_text", &textData{
		Label: "Alternative title",
		Text:  text,
	})
}

func (r *Publication) partial(partial string, vars interface{}) template.HTML {
	buf := &bytes.Buffer{}
	if tmpl := r.Render.TemplateLookup(partial); tmpl != nil {
		tmpl.Execute(buf, vars)
	}
	return template.HTML(buf.String())
}
