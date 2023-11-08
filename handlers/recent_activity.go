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
		WithFilter("status", "private", "public", "returned").
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
		if prevP != nil && p.Message != "" && p.Message != prevP.Message {
			acts = append(acts, views.Activity{
				Event:     views.MessageEvent,
				Object:    views.PublicationObject,
				User:      p.User,
				Datestamp: *p.DateUpdated,
				URL:       c.PathTo("publication", "id", p.ID).String(),
				Status:    p.Status,
				Title:     p.Title,
			})
		}
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
		} else {
			act.Event = views.UpdateEvent
		}
		acts = append(acts, act)
		if prevD != nil && d.Message != "" && d.Message != prevD.Message {
			acts = append(acts, views.Activity{
				Event:     views.MessageEvent,
				Object:    views.DatasetObject,
				User:      d.User,
				Datestamp: *d.DateUpdated,
				URL:       c.PathTo("dataset", "id", d.ID).String(),
				Status:    d.Status,
				Title:     d.Title,
			})
		}
	}

	sort.Slice(acts, func(i, j int) bool {
		return acts[i].Datestamp.After(acts[j].Datestamp)
	})

	views.RecentActivity(c, acts).Render(r.Context(), w)
}
