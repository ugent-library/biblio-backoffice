package displays

import (
	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render/display"
)

func PublicationAdditionalInfo(user *models.Person, loc *gotext.Locale, p *models.Publication) *display.Display {
	return display.New().
		WithTheme("default").
		AddSection(
			&display.List{
				Label:  loc.Get("builder.research_field"),
				Values: p.ResearchField,
			},
			&display.List{
				Inline:        true,
				Label:         loc.Get("builder.keyword"),
				Values:        p.Keyword,
				ValueTemplate: "format/badge",
			},
			&display.Text{
				Label: loc.Get("builder.additional_info"),
				Value: p.AdditionalInfo,
			},
		)
}
