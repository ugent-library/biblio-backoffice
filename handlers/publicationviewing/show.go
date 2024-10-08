package publicationviewing

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views"
	publicationviews "github.com/ugent-library/biblio-backoffice/views/publication"
)

func Show(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	redirectURL := r.URL.Query().Get("redirect-url")
	if redirectURL == "" {
		redirectURL = c.PathTo("publications").String()
	}

	subNav := r.URL.Query().Get("show")
	if subNav == "" {
		subNav = "description"
	}
	c.SubNav = subNav

	publicationviews.Show(c, p, redirectURL).Render(r.Context(), w)
}

func ShowDescription(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	redirectURL := r.URL.Query().Get("redirect-url")
	if redirectURL == "" {
		redirectURL = c.PathTo("publications").String()
	}

	publicationviews.Description(c, p, redirectURL).Render(r.Context(), w)
}

func ShowFiles(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	redirectURL := r.URL.Query().Get("redirect-url")
	if redirectURL == "" {
		redirectURL = c.PathTo("publications").String()
	}

	publicationviews.Files(c, p, redirectURL).Render(r.Context(), w)
}

func ShowContributors(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	redirectURL := r.URL.Query().Get("redirect-url")
	if redirectURL == "" {
		redirectURL = c.PathTo("publications").String()
	}

	publicationviews.Contributors(c, p, redirectURL).Render(r.Context(), w)
}

func ShowDatasets(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	datasets, err := c.Repo.GetPublicationDatasets(p)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	publicationviews.Datasets(c, p, datasets).Render(r.Context(), w)
}

func BiblioMessages(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	publicationviews.Messages(c, publicationviews.MessagesArgs{
		Publication: p,
	}).Render(r.Context(), w)
}

func RecentActivity(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	var (
		acts         []views.Activity
		nextSnapshot *models.Publication
	)

	err := c.Repo.PublicationHistory(p.ID, func(snapshot *models.Publication) bool {
		if nextSnapshot != nil {
			acts = append(acts, handlers.GetPublicationActivity(c, nextSnapshot, snapshot))
		}

		nextSnapshot = snapshot

		return true
	})
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	acts = append(acts, handlers.GetPublicationActivity(c, nextSnapshot, nil))

	publicationviews.RecentActivity(c, acts, p).Render(r.Context(), w)
}
