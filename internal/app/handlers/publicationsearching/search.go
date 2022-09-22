package publicationsearching

import (
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"

	"github.com/ugent-library/biblio-backend/internal/vocabularies"
)

var (
	userScopes = []string{"all", "contributed", "created"}
)

type YieldSearch struct {
	Context
	PageTitle   string
	ActiveNav   string
	Scopes      []string
	Hits        *models.PublicationHits
	ActionItems []*models.ActionItem
}

type YieldHit struct {
	Context
	Publication *models.Publication
}

func (y YieldSearch) YieldHit(d *models.Publication) YieldHit {
	return YieldHit{y.Context, d}
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request, ctx Context) {
	ctx.SearchArgs.WithFacets(vocabularies.Map["publication_facets"]...)
	if ctx.SearchArgs.FilterFor("scope") == "" {
		ctx.SearchArgs.WithFilter("scope", "all")
	}

	searcher := h.PublicationSearchService.WithScope("status", "private", "public")
	args := ctx.SearchArgs.Clone()

	switch args.FilterFor("scope") {
	case "created":
		searcher = searcher.WithScope("creator.id", ctx.User.ID)
	case "contributed":
		searcher = searcher.WithScope("author.id", ctx.User.ID)
	case "all":
		searcher = searcher.WithScope("creator.id|author.id", ctx.User.ID)
	default:
		errorUnkownScope := fmt.Errorf("unknown scope: %s", args.FilterFor("scope"))
		h.Logger.Warnw("publication search: could not create searcher with passed filters", "errors", errorUnkownScope, "user", ctx.User.ID)
		render.BadRequest(w, r, errorUnkownScope)
		return
	}
	delete(args.Filters, "scope")

	hits, err := searcher.Search(args)
	if err != nil {
		h.Logger.Errorw("publication search: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "layouts/default", "publication/pages/search", YieldSearch{
		Context:     ctx,
		PageTitle:   "Overview - Publications - Biblio",
		ActiveNav:   "publications",
		Scopes:      userScopes,
		Hits:        hits,
		ActionItems: h.getSearchActions(ctx),
	})
}

func (h *Handler) CurationSearch(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCuratePublications() {
		render.Forbidden(w, r)
		return
	}

	ctx.SearchArgs.WithFacets(vocabularies.Map["publication_curation_facets"]...)

	searcher := h.PublicationSearchService.WithScope("status", "private", "public")
	hits, err := searcher.Search(ctx.SearchArgs)
	if err != nil {
		h.Logger.Errorw("publication search: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "layouts/default", "publication/pages/search", YieldSearch{
		Context:     ctx,
		PageTitle:   "Overview - Publications - Biblio",
		ActiveNav:   "publications",
		Hits:        hits,
		ActionItems: h.getCurationSearchActions(ctx),
	})
}

func (h *Handler) getCurationSearchActions(ctx Context) []*models.ActionItem {
	actionItems := make([]*models.ActionItem, 0)
	if oa := h.getOrcidAction(ctx); oa != nil {
		actionItems = append(actionItems, oa)
	}
	u := h.PathFor("export_curation_publications", "format", "xlsx")
	q, _ := bind.EncodeQuery(ctx.SearchArgs)
	u.RawQuery = q.Encode()
	actionItems = append(actionItems, &models.ActionItem{
		Label:    ctx.Locale.T("export_to.xlsx"),
		URL:      u,
		Template: "actions/export",
	})
	return actionItems
}

func (h *Handler) getSearchActions(ctx Context) []*models.ActionItem {
	actionItems := make([]*models.ActionItem, 0)
	if oa := h.getOrcidAction(ctx); oa != nil {
		actionItems = append(actionItems, oa)
	}
	return actionItems
}

func (h *Handler) getOrcidAction(ctx Context) *models.ActionItem {
	if ctx.User.ORCID == "" || ctx.User.ORCIDToken == "" {
		return nil
	}
	return &models.ActionItem{
		Label:    "Send my publications to ORCID",
		URL:      h.PathFor("publication_orcid_add_all"),
		Template: "actions/publication_orcid_add_all",
	}
}
