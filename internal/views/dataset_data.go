package views

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
)

type DatasetData struct {
	Data
	render  *render.Render
	Dataset *models.Dataset
}

func NewDatasetData(r *http.Request, render *render.Render, d *models.Dataset) DatasetData {
	return DatasetData{Data: NewData(r), render: render, Dataset: d}
}

func (d DatasetData) RenderDetails() (template.HTML, error) {
	return RenderPartial(d.render, fmt.Sprintf("publication/details/_%s", d.Dataset.Type), d)
}

func (d DatasetData) RenderAbstract() (template.HTML, error) {
	return RenderPartial(d.render, fmt.Sprintf("publication/abstracts/_%s", d.Dataset.Type), d)
}

func (d DatasetData) RenderLinks() (template.HTML, error) {
	return RenderPartial(d.render, fmt.Sprintf("publication/links/_%s", d.Dataset.Type), d)
}

func (d DatasetData) RenderAdditionalInfo() (template.HTML, error) {
	return RenderPartial(d.render, fmt.Sprintf("publication/additional_info/_%s", d.Dataset.Type), d)
}

func (d DatasetData) RenderAuthors() (template.HTML, error) {
	return RenderPartial(d.render, fmt.Sprintf("publication/authors/_%s", d.Dataset.Type), d)
}

func (d DatasetData) RenderEditors() (template.HTML, error) {
	return RenderPartial(d.render, fmt.Sprintf("publication/editors/_%s", d.Dataset.Type), d)
}

func (d DatasetData) RenderSupervisors() (template.HTML, error) {
	return RenderPartial(d.render, fmt.Sprintf("publication/supervisors/_%s", d.Dataset.Type), d)
}

func (d DatasetData) RenderText(text, label string, required bool) (template.HTML, error) {
	return RenderPartial(d.render, "part/_text", &textData{
		Label:    label,
		Text:     text,
		Required: required,
	})
}

func (d DatasetData) RenderList(list []string, label string, required bool) (template.HTML, error) {
	return RenderPartial(d.render, "part/_list", &listData{
		Label:    label,
		List:     list,
		Required: required,
	})
}

func (d DatasetData) RenderRange(start, end, label string, required bool) (template.HTML, error) {
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

func (d DatasetData) RenderBadgeList(list []string, label string, required bool) (template.HTML, error) {
	return RenderPartial(d.render, "part/_badge_list", &listData{
		Label:    label,
		List:     list,
		Required: required,
	})
}
