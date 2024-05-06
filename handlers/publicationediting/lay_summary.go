package publicationediting

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

type BindLaySummary struct {
	LaySummaryID string `path:"lay_summary_id"`
	Text         string `form:"text"`
	Lang         string `form:"lang"`
}

type BindDeleteLaySummary struct {
	LaySummaryID string `path:"lay_summary_id"`
	SnapshotID   string `path:"snapshot_id"`
}

type YieldLaySummaries struct {
	Context
}
type YieldAddLaySummary struct {
	Context
	Form     *form.Form
	Conflict bool
}
type YieldEditLaySummary struct {
	Context
	LaySummaryID string
	Form         *form.Form
	Conflict     bool
}

func (h *Handler) AddLaySummary(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "show_modal", "publication/add_lay_summary", YieldAddLaySummary{
		Context:  ctx,
		Form:     laySummaryForm(ctx.Loc, ctx.Publication, &models.Text{}, nil),
		Conflict: false,
	})
}

func (h *Handler) CreateLaySummary(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindLaySummary
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("create publication lay summary: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	laySummary := models.Text{
		Lang: b.Lang,
		Text: b.Text,
	}

	ctx.Publication.AddLaySummary(&laySummary)

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		h.Logger.Warnw("create publication lay summary: could not validate contributor:", "errors", validationErrs, "identifier", ctx.Publication.ID)
		render.Layout(w, "refresh_modal", "publication/add_lay_summary", YieldAddLaySummary{
			Context:  ctx,
			Form:     laySummaryForm(ctx.Loc, ctx.Publication, &laySummary, validationErrs.(*okay.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "publication/add_lay_summary", YieldAddLaySummary{
			Context:  ctx,
			Form:     laySummaryForm(ctx.Loc, ctx.Publication, &laySummary, nil),
			Conflict: true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("create publication lay summary: Could not save the publication:", "error", err, "identifier", ctx.Publication.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_lay_summaries", YieldLaySummaries{
		Context: ctx,
	})
}

func (h *Handler) EditLaySummary(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindLaySummary
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("edit publication lay summary: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	laySummary := ctx.Publication.GetLaySummary(b.LaySummaryID)

	// TODO catch non-existing item in UI
	if laySummary == nil {
		h.Logger.Warnf("edit publication lay summary: Could not fetch the lay summary:", "publication", ctx.Publication.ID, "abstract", b.LaySummaryID, "user", ctx.User.ID)
		views.ShowModal(views.ErrorDialog(ctx.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	render.Layout(w, "show_modal", "publication/edit_lay_summary", YieldEditLaySummary{
		Context:      ctx,
		LaySummaryID: b.LaySummaryID,
		Form:         laySummaryForm(ctx.Loc, ctx.Publication, laySummary, nil),
	})
}

func (h *Handler) UpdateLaySummary(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLaySummary{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update publication lay summary: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	laySummary := ctx.Publication.GetLaySummary(b.LaySummaryID)

	if laySummary == nil {
		laySummary := &models.Text{
			Text: b.Text,
			Lang: b.Lang,
		}
		render.Layout(w, "refresh_modal", "publication/edit_lay_summary", YieldEditLaySummary{
			Context:      ctx,
			LaySummaryID: b.LaySummaryID,
			Form:         laySummaryForm(ctx.Loc, ctx.Publication, laySummary, nil),
			Conflict:     true,
		})
		return
	}

	laySummary.Text = b.Text
	laySummary.Lang = b.Lang

	ctx.Publication.SetLaySummary(laySummary)

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "publication/edit_lay_summary", YieldEditLaySummary{
			Context:      ctx,
			LaySummaryID: b.LaySummaryID,
			Form:         laySummaryForm(ctx.Loc, ctx.Publication, laySummary, validationErrs.(*okay.Errors)),
			Conflict:     false,
		})
		return
	}

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "publication/edit_lay_summary", YieldEditLaySummary{
			Context:      ctx,
			LaySummaryID: b.LaySummaryID,
			Form:         laySummaryForm(ctx.Loc, ctx.Publication, laySummary, nil),
			Conflict:     true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("update publication lay summary: Could not save the publication:", "error", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_lay_summaries", YieldLaySummaries{
		Context: ctx,
	})
}

func ConfirmDeleteLaySummary(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	var b BindDeleteLaySummary
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete publication lay summary: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	if b.SnapshotID != publication.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to remove this lay summary?",
		DeleteUrl:  c.PathTo("publication_delete_lay_summary", "id", publication.ID, "lay_summary_id", b.LaySummaryID),
		SnapshotID: publication.SnapshotID,
	}).Render(r.Context(), w)
}

func (h *Handler) DeleteLaySummary(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteLaySummary
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete publication lay summary: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	ctx.Publication.RemoveLaySummary(b.LaySummaryID)

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(ctx.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		h.Logger.Errorf("delete publication lay summary: Could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_lay_summaries", YieldLaySummaries{
		Context: ctx,
	})
}

func laySummaryForm(loc *gotext.Locale, publication *models.Publication, laySummary *models.Text, errors *okay.Errors) *form.Form {
	idx := -1
	for i, ls := range publication.LaySummary {
		if ls.ID == laySummary.ID {
			idx = i
			break
		}
	}
	return form.New().
		WithTheme("cols").
		WithErrors(localize.ValidationErrors(loc, errors)).
		AddSection(
			&form.TextArea{
				Name:  "text",
				Value: laySummary.Text,
				Label: loc.Get("builder.lay_summary.text"),
				Cols:  12,
				Rows:  6,
				Error: localize.ValidationErrorAt(loc, errors, fmt.Sprintf("/lay_summary/%d/text", idx)),
			},
			&form.Select{
				Name:    "lang",
				Value:   laySummary.Lang,
				Label:   loc.Get("builder.lay_summary.lang"),
				Options: localize.LanguageSelectOptions(),
				Cols:    12,
				Error:   localize.ValidationErrorAt(loc, errors, fmt.Sprintf("/lay_summary/%d/lang", idx)),
			},
		)
}
