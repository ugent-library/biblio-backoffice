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

type listData struct {
	Label    string
	List     []string
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
	return r.partial("part/_text", &textData{
		Label: "DOI",
		Text:  p.DOI,
	})
}

// Render the "ISSN/ISBN" field

func (r *Publication) ISXN(p *models.Publication) template.HTML {
	var list []string
	for _, val := range p.ISSN {
		list = append(list, fmt.Sprintf("ISSN: %s", val))
	}
	for _, val := range p.ISBN {
		list = append(list, fmt.Sprintf("ISBN: %s", val))
	}

	return r.partial("part/_list", &listData{
		Label: "ISSN/ISBN",
		List:  list,
	})
}

// Render the "Title" field

func (r *Publication) Title(p *models.Publication) template.HTML {
	return r.partial("part/_text", &textData{
		Label:    "Title",
		Text:     p.Title,
		Required: true,
	})
}

// Render the "Alternative Title" field

func (r *Publication) AlternativeTitle(p *models.Publication) template.HTML {
	return r.partial("part/_list", &listData{
		Label: "Alternative title",
		List:  p.AlternativeTitle,
	})
}

func (r *Publication) partial(partial string, vars interface{}) template.HTML {
	buf := &bytes.Buffer{}
	if tmpl := r.Render.TemplateLookup(partial); tmpl != nil {
		tmpl.Execute(buf, vars)
	}
	return template.HTML(buf.String())
}
