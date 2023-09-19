package displays

import (
	"github.com/ugent-library/biblio-backoffice/internal/locale"
	"github.com/ugent-library/biblio-backoffice/internal/render/display"
	"github.com/ugent-library/biblio-backoffice/models"
)

func PublicationAdditionalInfo(user *models.User, l *locale.Locale, p *models.Publication) *display.Display {
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
