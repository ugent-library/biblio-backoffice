package datasetviewing

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/fields"
)

type YieldShow struct {
	Context
	PageTitle    string
	ActiveNav    string
	ActiveSubNav string
	SearchArgs   *models.SearchArgs
}

type YieldShowDescription struct {
	Context
	ActiveSubNav string
	SearchArgs   *models.SearchArgs
	DetailFields []*fields.Fields
}

type YieldShowContributors struct {
	Context
	ActiveSubNav string
	SearchArgs   *models.SearchArgs
}

type YieldShowPublications struct {
	Context
	ActiveSubNav        string
	SearchArgs          *models.SearchArgs
	RelatedPublications []*models.Publication
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request, ctx Context) {
	// TODO bind search args
	searchArgs := models.NewSearchArgs()

	// TODO bind and validate
	activeSubNav := "description"
	if r.URL.Query().Get("show") == "contributors" || r.URL.Query().Get("show") == "publications" {
		activeSubNav = r.URL.Query().Get("show")
	}

	render.Wrap(w, "layouts/default", "dataset/show_page", YieldShow{
		Context:      ctx,
		PageTitle:    ctx.T("dataset.page.show.title"),
		ActiveNav:    "datasets",
		ActiveSubNav: activeSubNav,
		SearchArgs:   searchArgs,
	})
}

func (h *Handler) ShowDescription(w http.ResponseWriter, r *http.Request, ctx Context) {
	// TODO bind search args
	searchArgs := models.NewSearchArgs()

	render.Render(w, "dataset/show_description", YieldShowDescription{
		Context:      ctx,
		ActiveSubNav: "description",
		SearchArgs:   searchArgs,
		DetailFields: detailFields(ctx),
	})
}

func (h *Handler) ShowContributors(w http.ResponseWriter, r *http.Request, ctx Context) {
	// TODO bind search args
	searchArgs := models.NewSearchArgs()

	render.Render(w, "dataset/show_contributors", YieldShowContributors{
		Context:      ctx,
		ActiveSubNav: "contributors",
		SearchArgs:   searchArgs,
	})
}

func (h *Handler) ShowPublications(w http.ResponseWriter, r *http.Request, ctx Context) {
	// TODO bind search args
	searchArgs := models.NewSearchArgs()

	relatedPublications, err := h.Repo.GetDatasetPublications(ctx.Dataset)
	if render.InternalServerError(w, err) {
		return
	}

	render.Render(w, "dataset/show_publications", YieldShowPublications{
		Context:             ctx,
		ActiveSubNav:        "publications",
		SearchArgs:          searchArgs,
		RelatedPublications: relatedPublications,
	})
}

func detailFields(ctx Context) []*fields.Fields {
	return []*fields.Fields{
		{
			Theme: "default",
			Fields: []fields.Field{
				&fields.Text{
					Label:    ctx.T("builder.title"),
					Values:   []string{ctx.Dataset.Title},
					Required: true,
				},
				&fields.Text{
					Label:         ctx.T("builder.doi"),
					Values:        []string{ctx.Dataset.DOI},
					Required:      true,
					ValueTemplate: "format/doi",
				},
			},
		},
	}
}
