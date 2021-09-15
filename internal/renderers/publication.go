package renderers

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
)

type textData struct {
	Text     string
	Label    string
	Required bool
	Tooltip  string
}

type listData struct {
	List     []string
	Label    string
	Required bool
	Tooltip  string
}

type Publication struct {
	Render *render.Render
}

//
// Render each collapsible card
//

// Render the "Publication details" section under the Description tab

func (r *Publication) Details(p *models.Publication) template.HTML {
	return r.Partial(fmt.Sprintf("publication/details/_%s", p.Type), p)
}

func (r *Publication) Conference(p *models.Publication) template.HTML {
	return r.Partial(fmt.Sprintf("publication/conference/_%s", p.Type), p.Conference)
}

func (r *Publication) Abstract(p *models.Publication) template.HTML {
	return r.Partial(fmt.Sprintf("publication/abstract/_%s", p.Type), p.Abstract)
}

//
// Render each field
//

// Render the "ISSN/ISBN" field

func (r *Publication) ISXN(p *models.Publication, label string, required bool) template.HTML {
	var list []string
	for _, val := range p.ISSN {
		list = append(list, fmt.Sprintf("ISSN: %s", val))
	}
	for _, val := range p.ISBN {
		list = append(list, fmt.Sprintf("ISBN: %s", val))
	}

	return r.Partial("part/_list", &listData{
		Label: "ISSN/ISBN",
		List:  list,
	})
}

func (r *Publication) Text(text, label string, required bool) template.HTML {
	return r.Partial("part/_text", &textData{
		Label:    label,
		Text:     text,
		Required: required,
	})
}

func (r *Publication) List(list []string, label string, required bool) template.HTML {
	return r.Partial("part/_list", &listData{
		Label:    label,
		List:     list,
		Required: required,
	})
}

func (r *Publication) Range(start, end, label string, required bool) template.HTML {
	var text string
	if len(start) > 0 && len(end) > 0 && start == end {
		text = start
	} else if len(start) > 0 && len(end) > 0 {
		text = fmt.Sprintf("%s - %s", start, end)
	} else if len(start) > 0 {
		text = fmt.Sprintf("%s -", start)
	} else if len(end) > 0 {
		text = fmt.Sprintf("- %s", end)
	}
	return r.Partial("part/_text", &textData{
		Label:    label,
		Text:     text,
		Required: required,
	})
}

func (r *Publication) Partial(partial string, vars interface{}) template.HTML {
	buf := &bytes.Buffer{}
	if tmpl := r.Render.TemplateLookup(partial); tmpl != nil {
		tmpl.Execute(buf, vars)
	}
	return template.HTML(buf.String())
}
