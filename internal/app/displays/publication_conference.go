package displays

import (
	"github.com/ugent-library/biblio-backend/internal/app/helpers"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render/display"
)

func PublicationConference(user *models.User, l *locale.Locale, p *models.Publication) *display.Display {
	return display.New().
		WithTheme("default").
		AddSection(
			&display.Text{
				Label: l.T("builder.conference.name"),
				Value: p.ConferenceName,
			},
			&display.Text{
				Label: l.T("builder.conference.location"),
				Value: p.ConferenceLocation,
			},
			&display.Text{
				Label: l.T("builder.conference.organizer"),
				Value: p.ConferenceOrganizer,
			},
			&display.Text{
				Label: l.T("builder.conference.date"),
				Value: helpers.FormatRange(p.ConferenceStartDate, p.ConferenceEndDate),
			},
		)
}
