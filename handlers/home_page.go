package handlers

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	views.HomePage(c).Render(r.Context(), w)
}

func ActionRequired(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	views.ActionRequired(c).Render(r.Context(), w)
}

func DraftsToComplete(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	pHits, err := c.PublicationSearchIndex.Search(models.NewSearchArgs().
		WithPageSize(0).
		WithFilter("creator_id", c.User.ID).
		WithFilter("status", "private"))
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	dHits, err := c.DatasetSearchIndex.Search(models.NewSearchArgs().
		WithPageSize(0).
		WithFilter("creator_id", c.User.ID).
		WithFilter("status", "private"))
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.DraftsToComplete(c, pHits.Total, dHits.Total).Render(r.Context(), w)
}
