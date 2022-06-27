package displays

import (
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/helpers"
)

func DisplayConference(l *locale.Locale, p *models.Publication) *display.Display {
	return display.New().
		WithTheme("default").
		AddSection(
			&display.Text{
				Label: l.T("builder.conference.name"),
				Value: p.Conference.Name,
			},
			&display.Text{
				Label: l.T("builder.conference.location"),
				Value: p.Conference.Location,
			},
			&display.Text{
				Label: l.T("builder.conference.organizer"),
				Value: p.Conference.Organizer,
			},
			&display.Text{
				Label: l.T("builder.conference.date"),
				Value: helpers.FormatRange(p.Conference.StartDate, p.Conference.EndDate),
			},
		)
}
