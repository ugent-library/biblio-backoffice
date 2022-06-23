package datasetviewing

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/bind"
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
	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		render.BadRequest(w, r, err)
		return
	}

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
	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.Render(w, "dataset/show_description", YieldShowDescription{
		Context:      ctx,
		ActiveSubNav: "description",
		SearchArgs:   searchArgs,
		DetailFields: detailFields(ctx),
	})
}

func (h *Handler) ShowContributors(w http.ResponseWriter, r *http.Request, ctx Context) {
	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.Render(w, "dataset/show_contributors", YieldShowContributors{
		Context:      ctx,
		ActiveSubNav: "contributors",
		SearchArgs:   searchArgs,
	})
}

func (h *Handler) ShowPublications(w http.ResponseWriter, r *http.Request, ctx Context) {
	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	relatedPublications, err := h.Repository.GetDatasetPublications(ctx.Dataset)
	if err != nil {
		render.InternalServerError(w, r, err)
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
					Value:    ctx.Dataset.Title,
					Required: true,
				},
				&fields.Text{
					Label:         ctx.T("builder.doi"),
					Value:         ctx.Dataset.DOI,
					Required:      true,
					ValueTemplate: "format/doi",
				},
				&fields.Text{
					Label:         ctx.T("builder.url"),
					Value:         ctx.Dataset.URL,
					ValueTemplate: "format/link",
				},
			},
		},
		{
			Theme: "default",
			Fields: []fields.Field{
				&fields.Text{
					Label:    ctx.T("builder.publisher"),
					Value:    ctx.Dataset.Publisher,
					Required: true,
				},
				&fields.Text{
					Label:    ctx.T("builder.year"),
					Value:    ctx.Dataset.Year,
					Required: true,
				},
			},
		},
		{
			Theme: "default",
			Fields: []fields.Field{
				&fields.Text{
					Label:    ctx.T("builder.format"),
					Values:   ctx.Dataset.Format,
					List:     true,
					Required: true,
				},
				&fields.Text{
					Label:         ctx.T("builder.keyword"),
					Values:        ctx.Dataset.Keyword,
					ValueTemplate: "format/badge",
				},
			},
		},
		{
			Theme: "default",
			Fields: []fields.Field{
				&fields.Text{
					Label:    ctx.T("builder.license"),
					Value:    ctx.TS("cc_licenses", ctx.Dataset.License),
					Required: true,
				},
				&fields.Text{
					Label: ctx.T("builder.other_license"),
					Value: ctx.Dataset.OtherLicense,
				},
				&fields.Text{
					Label:    ctx.T("builder.access_level"),
					Value:    ctx.TS("access_levels", ctx.Dataset.AccessLevel),
					Required: true,
				},
				&fields.Text{
					Label: ctx.T("builder.embargo"),
					Value: ctx.Dataset.Embargo,
				},
				&fields.Text{
					Label: ctx.T("builder.embargo_to"),
					Value: ctx.TS("access_levels", ctx.Dataset.EmbargoTo),
				},
			},
		},
	}
}
