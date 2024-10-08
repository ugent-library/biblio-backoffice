package publicationviewing

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
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

func ShowActivity(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	redirectURL := r.URL.Query().Get("redirect-url")
	if redirectURL == "" {
		redirectURL = c.PathTo("publications").String()
	}

	publicationviews.Activity(c, p, redirectURL).Render(r.Context(), w)
}
