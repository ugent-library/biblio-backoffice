package publicationviewing

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	publicationviews "github.com/ugent-library/biblio-backoffice/views/publication"
)

var subNavs = []string{"description", "files", "contributors", "datasets", "activity"}

type YieldShowContributors struct {
	Context
	SubNavs      []string
	ActiveSubNav string
}

type YieldShowDatasets struct {
	Context
	SubNavs         []string
	ActiveSubNav    string
	RelatedDatasets []*models.Dataset
}

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

func (h *Handler) ShowContributors(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.View(w, "publication/show_contributors", YieldShowContributors{
		Context:      ctx,
		SubNavs:      subNavs,
		ActiveSubNav: "contributors",
	})
}

func ShowDatasets(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	datasets, err := c.Repo.GetVisiblePublicationDatasets(c.User, p)
	if err != nil {
		c.Log.Warn("show publication datasets: could not get publication datasets:", "errors", err, "publication", p.ID, "user", c.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	redirectURL := r.URL.Query().Get("redirect-url")
	if redirectURL == "" {
		redirectURL = c.PathTo("publications").String()
	}

	publicationviews.Datasets(c, p, datasets, redirectURL).Render(r.Context(), w)
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
