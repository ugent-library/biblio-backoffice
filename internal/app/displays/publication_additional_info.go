package displays

import (
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render/display"
)

func PublicationAdditionalInfo(l *locale.Locale, p *models.Publication) *display.Display {
	return display.New().
		WithTheme("default").
		AddSection(
			&display.Text{
				Label:  l.T("builder.research_field"),
				List:   true,
				Values: p.ResearchField,
			},
			&display.Text{
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
