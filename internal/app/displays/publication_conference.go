package displays

import (
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/helpers"
)

func PublicationConference(l *locale.Locale, c models.PublicationConference) *display.Display {
	return display.New().
		WithTheme("default").
		AddSection(
			&display.Text{
				Label: l.T("builder.conference.name"),
				Value: c.Name,
			},
			&display.Text{
				Label: l.T("builder.conference.location"),
				Value: c.Location,
			},
			&display.Text{
				Label: l.T("builder.conference.organizer"),
				Value: c.Organizer,
			},
			&display.Text{
				Label: l.T("builder.conference.date"),
				Value: helpers.FormatRange(c.StartDate, c.EndDate),
			},
		)
}
