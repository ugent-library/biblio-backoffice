package publicationsearching

import (
	"fmt"
	"net/http"

	"github.com/cshum/imagor/imagorpath"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
)

var (
	userScopes = []string{"all", "contributed", "created"}
)

type YieldSearch struct {
	Context
	PageTitle string
	ActiveNav string
	Scopes    []string
	Hits      *models.PublicationHits
}

type YieldHit struct {
	Context
	Publication *models.Publication
}

func (y YieldSearch) YieldHit(d *models.Publication) YieldHit {
	return YieldHit{y.Context, d}
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request, ctx Context) {
	if ctx.SearchArgs.FilterFor("scope") == "" {
		ctx.SearchArgs.WithFilter("scope", "all")
	}

	searcher := h.PublicationSearchService.WithScope("status", "private", "public")
	args := ctx.SearchArgs.Clone()

	switch args.FilterFor("scope") {
	case "created":
		searcher = searcher.WithScope("creator_id", ctx.User.ID)
	case "contributed":
		searcher = searcher.WithScope("author.id", ctx.User.ID)
	case "all":
		searcher = searcher.WithScope("creator_id|author.id", ctx.User.ID)
	default:
		render.BadRequest(w, r, fmt.Errorf("unknown scope %s", args.FilterFor("scope")))
		return
	}
	delete(args.Filters, "scope")

	hits, err := searcher.IncludeFacets(true).Search(args)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	for _, p := range hits.Hits {
		h.addThumbnailURLs(p)
	}

	render.Layout(w, "layouts/default", "publication/pages/search", YieldSearch{
		Context:   ctx,
		PageTitle: "Overview - Publications - Biblio",
		ActiveNav: "publications",
		Scopes:    userScopes,
		Hits:      hits,
	})
}

func (h *Handler) CurationSearch(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCuratePublications() {
		render.Forbidden(w, r)
		return
	}

	searcher := h.PublicationSearchService.WithScope("status", "private", "public")
	hits, err := searcher.IncludeFacets(true).Search(ctx.SearchArgs)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	for _, p := range hits.Hits {
		h.addThumbnailURLs(p)
	}

	render.Layout(w, "layouts/default", "publication/pages/search", YieldSearch{
		Context:   ctx,
		PageTitle: "Overview - Publications - Biblio",
		ActiveNav: "publications",
		Hits:      hits,
	})
}

// TODO clean this up
func (h *Handler) addThumbnailURLs(p *models.Publication) {
	var u string
	for _, f := range p.File {
		if f.ContentType == "application/pdf" && f.Size <= 25000000 {
			params := imagorpath.Params{
				Image:  h.FileStore.RelativeFilePath(f.SHA256),
				FitIn:  true,
				Width:  156,
				Height: 156,
			}
			p := imagorpath.Generate(params, imagorpath.NewDefaultSigner(viper.GetString("imagor-secret")))
			u = viper.GetString("imagor-url") + "/" + p
			f.ThumbnailURL = u
		}
	}
}
