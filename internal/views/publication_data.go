package views

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
)

type PublicationData struct {
	Data
	render      *render.Render
	Publication *models.Publication
}

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

func NewPublicationData(r *http.Request, render *render.Render, p *models.Publication) PublicationData {
	return PublicationData{Data: NewData(r), render: render, Publication: p}
}

func (d PublicationData) RenderDetails() (template.HTML, error) {
	return RenderPartial(d.render, fmt.Sprintf("publication/details/_%s", d.Publication.Type), d)
}

func (d PublicationData) RenderConference() (template.HTML, error) {
	return RenderPartial(d.render, fmt.Sprintf("publication/conference/_%s", d.Publication.Type), d)
}

func (d PublicationData) RenderAbstract() (template.HTML, error) {
	return RenderPartial(d.render, fmt.Sprintf("publication/abstract/_%s", d.Publication.Type), d)
}

func (d PublicationData) RenderLinks() (template.HTML, error) {
	return RenderPartial(d.render, fmt.Sprintf("publication/links/_%s", d.Publication.Type), d)
}

func (d PublicationData) RenderAdditionalInfo() (template.HTML, error) {
	return RenderPartial(d.render, fmt.Sprintf("publication/additional_info/_%s", d.Publication.Type), d)
}

func (d PublicationData) RenderAuthors() (template.HTML, error) {
	return RenderPartial(d.render, fmt.Sprintf("publication/authors/_%s", d.Publication.Type), d)
}

func (d PublicationData) RenderEditors() (template.HTML, error) {
	return RenderPartial(d.render, fmt.Sprintf("publication/editors/_%s", d.Publication.Type), d)
}

func (d PublicationData) RenderSupervisors() (template.HTML, error) {
	return RenderPartial(d.render, fmt.Sprintf("publication/supervisors/_%s", d.Publication.Type), d)
}

func (d PublicationData) RenderText(text, label string, required bool) (template.HTML, error) {
	return RenderPartial(d.render, "part/_text", &textData{
		Label:    label,
		Text:     text,
		Required: required,
	})
}

func (d PublicationData) RenderList(list []string, label string, required bool) (template.HTML, error) {
	return RenderPartial(d.render, "part/_list", &listData{
		Label:    label,
		List:     list,
		Required: required,
	})
}

func (d PublicationData) RenderRange(start, end, label string, required bool) (template.HTML, error) {
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

	return RenderPartial(d.render, "part/_text", &textData{
		Label:    label,
		Text:     text,
		Required: required,
	})
}

func (d PublicationData) RenderBadgeList(list []string, label string, required bool) (template.HTML, error) {
	return RenderPartial(d.render, "part/_badge_list", &listData{
		Label:    label,
		List:     list,
		Required: required,
	})
}
