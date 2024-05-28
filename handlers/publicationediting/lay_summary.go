package publicationediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	publicationviews "github.com/ugent-library/biblio-backoffice/views/publication"
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

func AddLaySummary(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	views.ShowModal(publicationviews.EditLaySummaryDialog(c, p, &models.Text{}, -1, false, nil, true)).Render(r.Context(), w)
}

func CreateLaySummary(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	var b BindLaySummary
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	laySummary := &models.Text{
		Lang: b.Lang,
		Text: b.Text,
	}

	p.AddLaySummary(laySummary)

	idx := -1
	for i, l := range p.LaySummary {
		if l.ID == laySummary.ID {
			idx = i
		}
	}

	if validationErrs := p.Validate(); validationErrs != nil {
		views.ReplaceModal(publicationviews.EditLaySummaryDialog(c, p, laySummary, idx, false, validationErrs.(*okay.Errors), true)).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(publicationviews.EditLaySummaryDialog(c, p, laySummary, idx, true, nil, true)).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.CloseModalAndReplace(publicationviews.LaySummariesBodySelector, publicationviews.LaySummariesBody(c, p)).Render(r.Context(), w)
}

func EditLaySummary(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	var b BindLaySummary
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	laySummary := p.GetLaySummary(b.LaySummaryID)

	// TODO catch non-existing item in UI
	if laySummary == nil {
		c.Log.Warnf("edit publication lay summary: Could not fetch the lay summary:", "publication", p.ID, "lay_summary", b.LaySummaryID, "user", c.User.ID)
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	idx := -1
	for i, l := range p.LaySummary {
		if l.ID == laySummary.ID {
			idx = i
		}
	}

	views.ShowModal(publicationviews.EditLaySummaryDialog(c, p, laySummary, idx, false, nil, false)).Render(r.Context(), w)
}

func UpdateLaySummary(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindLaySummary{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	laySummary := p.GetLaySummary(b.LaySummaryID)

	if laySummary == nil {
		laySummary := &models.Text{
			Text: b.Text,
			Lang: b.Lang,
		}
		views.ReplaceModal(publicationviews.EditLaySummaryDialog(c, p, laySummary, -1, true, nil, false)).Render(r.Context(), w)
		return
	}

	laySummary.Text = b.Text
	laySummary.Lang = b.Lang

	p.SetLaySummary(laySummary)

	idx := -1
	for i, l := range p.LaySummary {
		if l.ID == laySummary.ID {
			idx = i
		}
	}

	if validationErrs := p.Validate(); validationErrs != nil {
		views.ReplaceModal(publicationviews.EditLaySummaryDialog(c, p, laySummary, idx, false, validationErrs.(*okay.Errors), false)).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(publicationviews.EditLaySummaryDialog(c, p, laySummary, idx, true, nil, false)).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.CloseModalAndReplace(publicationviews.LaySummariesBodySelector, publicationviews.LaySummariesBody(c, p)).Render(r.Context(), w)
}

func ConfirmDeleteLaySummary(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	var b BindDeleteLaySummary
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	if b.SnapshotID != p.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to remove this lay summary?",
		DeleteUrl:  c.PathTo("publication_delete_lay_summary", "id", p.ID, "lay_summary_id", b.LaySummaryID),
		SnapshotID: p.SnapshotID,
	}).Render(r.Context(), w)
}

func DeleteLaySummary(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	var b BindDeleteLaySummary
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	p.RemoveLaySummary(b.LaySummaryID)

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.CloseModalAndReplace(publicationviews.LaySummariesBodySelector, publicationviews.LaySummariesBody(c, p)).Render(r.Context(), w)
}
