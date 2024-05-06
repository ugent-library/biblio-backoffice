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

func (h *Handler) ShowDatasets(w http.ResponseWriter, r *http.Request, ctx Context) {
	relatedDatasets, err := h.Repo.GetVisiblePublicationDatasets(ctx.User, ctx.Publication)
	if err != nil {
		h.Logger.Warn("show publication datasets: could not get publication datasets:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/show_datasets", YieldShowDatasets{
		Context:         ctx,
		SubNavs:         subNavs,
		ActiveSubNav:    "datasets",
		RelatedDatasets: relatedDatasets,
	})
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
