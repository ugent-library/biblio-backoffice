package handlers

import (
	"errors"
	"net/http"
	"sort"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views"
	dashboardviews "github.com/ugent-library/biblio-backoffice/views/dashboard"
)

func DashBoard(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	if c.UserRole == "curator" {
		// TODO port and render here as CuratorDashboard
		http.Redirect(w, r, c.PathTo("dashboard_publications", "type", "faculties").String(), http.StatusSeeOther)
	} else {
		dashboardviews.UserDashboard(c).Render(r.Context(), w)
	}
}

func DashBoardIcon(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	pHits, err := c.PublicationSearchIndex.Search(models.NewSearchArgs().
		WithPageSize(0).
		WithFilter("creator_id|author_id", c.User.ID).
		WithFilter("status", "private").
		WithFilter("locked", "false"))
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	if pHits.Total > 0 {
		views.DashboardIcon(c, true).Render(r.Context(), w)
		return
	}
	dHits, err := c.DatasetSearchIndex.Search(models.NewSearchArgs().
		WithPageSize(0).
		WithFilter("creator_id|author_id", c.User.ID).
		WithFilter("status", "private").
		WithFilter("locked", "false"))
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	if dHits.Total > 0 {
		views.DashboardIcon(c, true).Render(r.Context(), w)
		return
	}

	pHits, err = c.PublicationSearchIndex.Search(models.NewSearchArgs().
		WithPageSize(0).
		WithFilter("creator_id|author_id", c.User.ID).
		WithFilter("status", "returned").
		WithFilter("locked", "false"))
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	if pHits.Total > 0 {
		views.DashboardIcon(c, true).Render(r.Context(), w)
		return
	}
	dHits, err = c.DatasetSearchIndex.Search(models.NewSearchArgs().
		WithPageSize(0).
		WithFilter("creator_id|author_id", c.User.ID).
		WithFilter("status", "returned").
		WithFilter("locked", "false"))
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	if dHits.Total > 0 {
		views.DashboardIcon(c, true).Render(r.Context(), w)
		return
	}

	exists, err := c.Repo.PersonHasCandidateRecords(r.Context(), c.User.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	if exists {
		views.DashboardIcon(c, true).Render(r.Context(), w)
		return
	}

	views.DashboardIcon(c, false).Render(r.Context(), w)
}

func DraftsToComplete(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	pHits, err := c.PublicationSearchIndex.Search(models.NewSearchArgs().
		WithPageSize(0).
		WithFilter("creator_id|author_id", c.User.ID).
		WithFilter("status", "private").
		WithFilter("locked", "false"))
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	dHits, err := c.DatasetSearchIndex.Search(models.NewSearchArgs().
		WithPageSize(0).
		WithFilter("creator_id|author_id", c.User.ID).
		WithFilter("status", "private").
		WithFilter("locked", "false"))
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	dashboardviews.DraftsToComplete(c, pHits.Total, dHits.Total).Render(r.Context(), w)
}

func ActionRequired(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	pHits, err := c.PublicationSearchIndex.Search(models.NewSearchArgs().
		WithPageSize(0).
		WithFilter("creator_id|author_id", c.User.ID).
		WithFilter("status", "returned").
		WithFilter("locked", "false"))
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	dHits, err := c.DatasetSearchIndex.Search(models.NewSearchArgs().
		WithPageSize(0).
		WithFilter("creator_id|author_id", c.User.ID).
		WithFilter("status", "returned").
		WithFilter("locked", "false"))
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	dashboardviews.ActionRequired(c, pHits.Total, dHits.Total).Render(r.Context(), w)
}

func CandidateRecords(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	var total int
	var recs []*models.CandidateRecord
	var err error

	if c.FlagCandidateRecords() {
		searchArgs := models.NewSearchArgs().
			WithPageSize(4).
			WithFilter("status", "new").
			WithFilter("person_id", c.User.ID)

		total, recs, err = c.Repo.GetCandidateRecords(r.Context(), searchArgs)
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	dashboardviews.CandidateRecords(c, total, recs).Render(r.Context(), w)
}

func RecentActivity(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	var acts []views.Activity

	pHits, err := c.PublicationSearchIndex.Search(models.NewSearchArgs().
		WithPageSize(5).
		WithFilter("status", "private", "public", "returned").
		WithFilter("creator_id|author_id", c.User.ID))
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	for _, p := range pHits.Hits {
		prevP, err := c.Repo.GetPublicationSnapshotBefore(p.ID, *p.DateFrom)
		if err != nil && !errors.Is(err, models.ErrNotFound) {
			c.HandleError(w, r, err)
			return
		}

		acts = append(acts, GetPublicationActivity(c, p, prevP))
	}

	dHits, err := c.DatasetSearchIndex.Search(models.NewSearchArgs().
		WithPageSize(5).
		WithFilter("status", "private", "public", "returned").
		WithFilter("creator_id|author_id", c.User.ID))
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	for _, d := range dHits.Hits {
		prevD, err := c.Repo.GetDatasetSnapshotBefore(d.ID, *d.DateFrom)
		if err != nil && !errors.Is(err, models.ErrNotFound) {
			c.HandleError(w, r, err)
			return
		}

		acts = append(acts, GetDatasetActivity(c, d, prevD))
	}

	sort.Slice(acts, func(i, j int) bool {
		return acts[i].Datestamp.After(acts[j].Datestamp)
	})

	views.RecentActivity(c, acts).Render(r.Context(), w)
}

func GetPublicationActivity(c *ctx.Ctx, p *models.Publication, prevP *models.Publication) views.Activity {
	act := views.Activity{
		Object:    views.PublicationObject,
		User:      p.User,
		Datestamp: *p.DateUpdated,
		URL:       c.PathTo("publication", "id", p.ID).String(),
		Status:    p.Status,
		RecordID:  p.ID,
		Title:     p.Title,
	}

	if prevP == nil {
		act.Event = views.CreateEvent
	} else if p.Status == "public" && prevP.Status == "returned" {
		act.Event = views.RepublishEvent
	} else if p.Status == "public" && prevP.Status != "public" {
		act.Event = views.PublishEvent
	} else if p.Status == "returned" && prevP.Status != "returned" {
		act.Event = views.WithdrawEvent
	} else if p.Locked && !prevP.Locked {
		act.Event = views.LockEvent
	} else if !p.Locked && prevP.Locked {
		act.Event = views.UnlockEvent
	} else if p.Message != "" && p.Message != prevP.Message {
		act.Event = views.MessageEvent
	} else {
		act.Event = views.UpdateEvent
	}

	return act
}

func GetDatasetActivity(c *ctx.Ctx, d *models.Dataset, prevD *models.Dataset) views.Activity {
	act := views.Activity{
		Object:    views.DatasetObject,
		User:      d.User,
		Datestamp: *d.DateUpdated,
		URL:       c.PathTo("dataset", "id", d.ID).String(),
		Status:    d.Status,
		RecordID:  d.ID,
		Title:     d.Title,
	}

	if prevD == nil {
		act.Event = views.CreateEvent
	} else if d.Status == "public" && prevD.Status == "returned" {
		act.Event = views.RepublishEvent
	} else if d.Status == "public" && prevD.Status != "public" {
		act.Event = views.PublishEvent
	} else if d.Status == "returned" && prevD.Status != "returned" {
		act.Event = views.WithdrawEvent
	} else if d.Locked && !prevD.Locked {
		act.Event = views.LockEvent
	} else if !d.Locked && prevD.Locked {
		act.Event = views.UnlockEvent
	} else if d.Message != "" && d.Message != prevD.Message {
		act.Event = views.MessageEvent
	} else {
		act.Event = views.UpdateEvent
	}

	return act
}
