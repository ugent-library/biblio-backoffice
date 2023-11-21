package displays

import (
	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/biblio-backoffice/helpers"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render/display"
)

func PublicationConference(user *models.Person, loc *gotext.Locale, p *models.Publication) *display.Display {
	return display.New().
		WithTheme("default").
		AddSection(
			&display.Text{
				Label: loc.Get("builder.conference.name"),
				Value: p.ConferenceName,
			},
			&display.Text{
				Label: loc.Get("builder.conference.location"),
				Value: p.ConferenceLocation,
			},
			&display.Text{
				Label: loc.Get("builder.conference.organizer"),
				Value: p.ConferenceOrganizer,
			},
			&display.Text{
				Label: loc.Get("builder.conference.date"),
				Value: helpers.FormatRange(p.ConferenceStartDate, p.ConferenceEndDate),
			},
		)
}
