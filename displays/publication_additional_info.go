package displays

import (
	"github.com/ugent-library/biblio-backoffice/locale"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render/display"
)

func PublicationAdditionalInfo(user *models.Person, l *locale.Locale, p *models.Publication) *display.Display {
	return display.New().
		WithTheme("default").
		AddSection(
			&display.List{
				Label:  l.T("builder.research_field"),
				Values: p.ResearchField,
			},
			&display.List{
				Inline:        true,
				Label:         l.T("builder.keyword"),
				Values:        p.Keyword,
				ValueTemplate: "format/badge",
			},
			&display.Text{
				Label: l.T("builder.additional_info"),
				Value: p.AdditionalInfo,
			},
		)
}
