package handlers

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views"
)

func ActionRequired(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	pHits, err := c.PublicationSearchIndex.Search(models.NewSearchArgs().
		WithPageSize(0).
		WithFilter("creator_id|author_id", c.User.ID).
		WithFilter("status", "returned"))
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	dHits, err := c.DatasetSearchIndex.Search(models.NewSearchArgs().
		WithPageSize(0).
		WithFilter("creator_id|author_id", c.User.ID).
		WithFilter("status", "returned"))
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.ActionRequired(c, pHits.Total, dHits.Total).Render(r.Context(), w)
}
