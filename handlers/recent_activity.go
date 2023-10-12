package handlers

import (
	"net/http"
	"sort"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views"
)

func RecentActivity(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	var acts []views.Activity

	pHits, err := c.PublicationSearchIndex.Search(models.NewSearchArgs().
		WithPageSize(5).
		WithFilter("creator_id|author_id", c.User.ID))
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	for _, p := range pHits.Hits {
		prevP, err := c.Repo.GetPublicationSnapshotBefore(p.ID, *p.DateFrom)
		if err != nil && err != models.ErrNotFound {
			c.HandleError(w, r, err)
			return
		}
		act := views.Activity{
			Object:    views.PublicationObject,
			User:      p.User,
			Datestamp: *p.DateUpdated,
			URL:       c.PathTo("publication", "id", p.ID).String(),
			Status:    p.Status,
			Title:     p.Title,
		}
		if prevP == nil {
			act.Event = views.CreateEvent
		} else if p.Status == "deleted" && prevP.Status != "deleted" {
			act.Event = views.DeleteEvent
			act.Status = prevP.Status
		} else if p.Status == "public" && prevP.Status == "returned" {
			act.Event = views.RepublishEvent
		} else if p.Status == "public" && prevP.Status != "public" {
			act.Event = views.PublishEvent
		} else if p.Status == "returned" && prevP.Status != "returned" {
			act.Event = views.WithdrawEvent
		} else if p.Locked && !prevP.Locked {
			act.Event = views.LockEvent
		} else {
			act.Event = views.UpdateEvent
		}
		acts = append(acts, act)
	}

	dHits, err := c.DatasetSearchIndex.Search(models.NewSearchArgs().
		WithPageSize(5).
		WithFilter("creator_id|author_id", c.User.ID))
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	for _, d := range dHits.Hits {
		prevd, err := c.Repo.GetDatasetSnapshotBefore(d.ID, *d.DateFrom)
		if err != nil && err != models.ErrNotFound {
			c.HandleError(w, r, err)
			return
		}
		act := views.Activity{
			Object:    views.DatasetObject,
			User:      d.User,
			Datestamp: *d.DateUpdated,
			URL:       c.PathTo("dataset", "id", d.ID).String(),
			Status:    d.Status,
			Title:     d.Title,
		}
		if prevd == nil {
			act.Event = views.CreateEvent
		} else if d.Status == "deleted" && prevd.Status != "deleted" {
			act.Event = views.DeleteEvent
			act.Status = prevd.Status
		} else if d.Status == "public" && prevd.Status == "returned" {
			act.Event = views.RepublishEvent
		} else if d.Status == "public" && prevd.Status != "public" {
			act.Event = views.PublishEvent
		} else if d.Status == "returned" && prevd.Status != "returned" {
			act.Event = views.WithdrawEvent
		} else if d.Locked && !prevd.Locked {
			act.Event = views.LockEvent
		} else {
			act.Event = views.UpdateEvent
		}
		acts = append(acts, act)
	}

	sort.Slice(acts, func(i, j int) bool {
		return acts[i].Datestamp.After(acts[j].Datestamp)
	})

	views.RecentActivity(c, acts).Render(r.Context(), w)
}
