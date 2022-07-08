package publicationediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/displays"
	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
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
	Form *form.Form
}

func (h *Handler) EditConference(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.MustRenderLayout(w, "show_modal", "publication/edit_conference", YieldEditConference{
		Context: ctx,
		Form:    conferenceForm(ctx.Locale, ctx.Publication, nil),
	})
}

func (h *Handler) UpdateConference(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindConference{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	p := ctx.Publication

	p.Conference.Name = b.Name
	p.Conference.Location = b.Location
	p.Conference.Organizer = b.Organizer
	p.Conference.StartDate = b.StartDate
	p.Conference.EndDate = b.EndDate

	if validationErrs := p.Validate(); validationErrs != nil {
		form := conferenceForm(ctx.Locale, p, validationErrs.(validation.Errors))
		render.MustRenderLayout(w, "refresh_modal", "publication/edit_conference", YieldEditConference{
			Context: ctx,
			Form:    form,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), p)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Render(w, "error_dialog", ctx.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.MustRenderView(w, "publication/refresh_conference", YieldConference{
		Context:           ctx,
		DisplayConference: displays.PublicationConference(ctx.Locale, p.Conference),
	})
}

func conferenceForm(l *locale.Locale, p *models.Publication, errors validation.Errors) *form.Form {
	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(l, errors)).
		AddSection(
			&form.Text{
				Name:  "name",
				Value: p.Conference.Name,
				Label: l.T("builder.conference.name"),
				Cols:  9,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					"/conference/name",
				),
			},
			&form.Text{
				Name:  "location",
				Value: p.Conference.Location,
				Label: l.T("builder.conference.location"),
				Cols:  9,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					"/conference/location",
				),
			},
			&form.Text{
				Name:  "organizer",
				Value: p.Conference.Organizer,
				Label: l.T("builder.conference.organizer"),
				Cols:  9,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					"/conference/organizer",
				),
			},
			&form.Text{
				Name:  "start_date",
				Value: p.Conference.StartDate,
				Label: l.T("builder.conference.start_date"),
				Cols:  3,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					"/conference/start_date",
				),
			},
			&form.Text{
				Name:  "end_date",
				Value: p.Conference.EndDate,
				Label: l.T("builder.conference.end_date"),
				Cols:  3,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					"/conference/end_date",
				),
			},
		)
}
