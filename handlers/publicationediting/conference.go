package publicationediting

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	publicationviews "github.com/ugent-library/biblio-backoffice/views/publication"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

type BindConference struct {
	Name      string `form:"name"`
	Location  string `form:"location"`
	Organizer string `form:"organizer"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

func EditConference(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)
	views.ShowModal(publicationviews.EditConferenceDialog(c, p, false, nil)).Render(r.Context(), w)
}

func UpdateConference(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindConference{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	p.ConferenceName = b.Name
	p.ConferenceLocation = b.Location
	p.ConferenceOrganizer = b.Organizer
	p.ConferenceStartDate = b.StartDate
	p.ConferenceEndDate = b.EndDate

	if validationErrs := p.Validate(); validationErrs != nil {
		views.ReplaceModal(publicationviews.EditConferenceDialog(c, p, false, validationErrs.(*okay.Errors))).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(publicationviews.EditConferenceDialog(c, p, true, nil)).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, httperror.InternalServerError.Wrap(fmt.Errorf("could not save the publication: %w", err)))
		return
	}

	views.CloseModalAndReplace(publicationviews.ConferenceBodySelector, publicationviews.ConferenceBody(c, p)).Render(r.Context(), w)
}
