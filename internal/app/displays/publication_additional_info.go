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
		).
		AddSection(
			&display.Text{
				Label: l.T("builder.message"),
				Value: p.Message,
			},
		)
}
