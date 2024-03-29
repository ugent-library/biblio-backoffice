package publicationediting

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/biblio-backoffice/displays"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/display"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/validation"
	"github.com/ugent-library/bind"
)

type BindConference struct {
	Name      string `form:"name"`
	Location  string `form:"location"`
	Organizer string `form:"organizer"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

type YieldConference struct {
	Context
	DisplayConference *display.Display
}

type YieldEditConference struct {
	Context
	Form     *form.Form
	Conflict bool
}

func (h *Handler) EditConference(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "show_modal", "publication/edit_conference", YieldEditConference{
		Context:  ctx,
		Form:     conferenceForm(ctx.Loc, ctx.Publication, nil),
		Conflict: false,
	})
}

func (h *Handler) UpdateConference(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindConference{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update publication conference: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	p := ctx.Publication

	p.ConferenceName = b.Name
	p.ConferenceLocation = b.Location
	p.ConferenceOrganizer = b.Organizer
	p.ConferenceStartDate = b.StartDate
	p.ConferenceEndDate = b.EndDate

	if validationErrs := p.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "publication/edit_conference", YieldEditConference{
			Context:  ctx,
			Form:     conferenceForm(ctx.Loc, p, validationErrs.(validation.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), p, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "publication/edit_conference", YieldEditConference{
			Context:  ctx,
			Form:     conferenceForm(ctx.Loc, p, nil),
			Conflict: true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("update publication conference: could not save the publication:", "error", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_conference", YieldConference{
		Context:           ctx,
		DisplayConference: displays.PublicationConference(ctx.User, ctx.Loc, p),
	})
}

func conferenceForm(loc *gotext.Locale, publication *models.Publication, errors validation.Errors) *form.Form {
	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(loc, errors)).
		AddSection(
			&form.Text{
				Name:  "name",
				Value: publication.ConferenceName,
				Label: loc.Get("builder.conference.name"),
				Cols:  9,
				Error: localize.ValidationErrorAt(
					loc,
					errors,
					"/conference_name",
				),
			},
			&form.Text{
				Name:  "location",
				Value: publication.ConferenceLocation,
				Label: loc.Get("builder.conference.location"),
				Cols:  9,
				Error: localize.ValidationErrorAt(
					loc,
					errors,
					"/conference_location",
				),
			},
			&form.Text{
				Name:  "organizer",
				Value: publication.ConferenceOrganizer,
				Label: loc.Get("builder.conference.organizer"),
				Cols:  9,
				Error: localize.ValidationErrorAt(
					loc,
					errors,
					"/conference_organizer",
				),
			},
			&form.Text{
				Name:  "start_date",
				Value: publication.ConferenceStartDate,
				Label: loc.Get("builder.conference.start_date"),
				Cols:  3,
				Help:  template.HTML(loc.Get("builder.conference.start_date.help")),
				Error: localize.ValidationErrorAt(
					loc,
					errors,
					"/conference_start_date",
				),
			},
			&form.Text{
				Name:  "end_date",
				Value: publication.ConferenceEndDate,
				Label: loc.Get("builder.conference.end_date"),
				Cols:  3,
				Help:  template.HTML(loc.Get("builder.conference.end_date.help")),
				Error: localize.ValidationErrorAt(
					loc,
					errors,
					"/conference_end_date",
				),
			},
		)
}
