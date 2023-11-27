package dashboard

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
	"github.com/ugent-library/bind"
)

type YieldPublications struct {
	Context
	PageTitle            string
	ActiveNav            string
	ActiveSubNav         string
	UPublications        map[string]map[string][]string
	APublications        map[string]map[string][]string
	UFaculties           []string
	AFaculties           []string
	PTypes               map[string]string
	UYear                string
	AYear                string
	AllUPublicationYears []string
	AllAPublicationYears []string
}

type BindPublications struct {
	UYear string `query:"uyear" form:"uyear"`
	AYear string `query:"ayear" form:"ayear"`
}

type YieldTblPublications struct {
	Context
	Publications map[string]map[string][]string
	Faculties    []string
	PTypes       map[string]string
}

func (h *Handler) Publications(w http.ResponseWriter, r *http.Request, ctx Context) {
	var faculties []string

	var activeSubNav string

	socs := vocabularies.Map["faculties_socs"]
	core := vocabularies.Map["faculties_core"]

	switch ctx.Type {
	case "socs":
		faculties = socs
		activeSubNav = "dashboard_publications_socs"
	default:
		faculties = core
		activeSubNav = "dashboard_publications_faculties"
	}

	aFacultyCols := append([]string{"all"}, faculties...)

	ptypes := vocabularies.Map["publication_types"]
	ptypes = append([]string{"all"}, ptypes...)

	locptypes := localize.VocabularyTerms(ctx.Loc, "publication_types")
	locptypes["all"] = "All"

	allUPublicationYears, err := h.allUPublicationYears()
	if err != nil {
		h.Logger.Errorw("Dashboard: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}
	allAPublicationYears, err := h.allAPublicationYears()
	if err != nil {
		h.Logger.Errorw("Dashboard: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}
	bindPublications := BindPublications{}
	if err := bind.Request(r, &bindPublications); err != nil {
		h.Logger.Warnw("publication dashboard could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	// Publications with classification U
	uFacultyCols := append(aFacultyCols, []string{"UGent", "-"}...)

	uSearcher := h.PublicationSearchIndex
	baseSearchUrl := h.PathFor("publications")

	uPublications, err := generatePublicationsDashboard(uFacultyCols, ptypes, uSearcher, baseSearchUrl, func(fac string, args *models.SearchArgs) *models.SearchArgs {
		args.WithFilter("classification", "U")
		args.WithFilter("status", "public")
		if bindPublications.UYear != "" {
			args.WithFilter("year", bindPublications.UYear)
		}

		switch fac {
		case "all":
			args.WithFilter("faculty_id", faculties...)
		case "-":
			args.WithFilter("faculty_id", backends.MissingValue)
		case "UGent":
			args.WithFilter("organization_id", "UGent")
		default:
			args.WithFilter("faculty_id", fac)
		}

		return args
	})

	if err != nil {
		h.Logger.Errorw("Dashboard: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// Publications with publication status "accepted"

	aSearcher := h.PublicationSearchIndex

	aPublications, err := generatePublicationsDashboard(aFacultyCols, ptypes, aSearcher, baseSearchUrl, func(fac string, args *models.SearchArgs) *models.SearchArgs {
		args.WithFilter("publication_status", "accepted")
		args.WithFilter("status", "private", "public", "returned")
		if bindPublications.AYear != "" {
			args.WithFilter("year", bindPublications.AYear)
		}

		switch fac {
		case "all":
			args.WithFilter("faculty_id", faculties...)
		case "-":
			args.WithFilter("faculty_id", backends.MissingValue)
		default:
			args.WithFilter("faculty_id", fac)
		}

		return args
	})

	if err != nil {
		h.Logger.Errorw("Dashboard: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "layouts/default", "dashboard/pages/publications", YieldPublications{
		Context:              ctx,
		PageTitle:            "Dashboard - Publications - Biblio",
		ActiveNav:            "dashboard",
		ActiveSubNav:         activeSubNav,
		UPublications:        uPublications,
		APublications:        aPublications,
		UFaculties:           uFacultyCols,
		AFaculties:           aFacultyCols,
		PTypes:               locptypes,
		UYear:                bindPublications.UYear,
		AYear:                bindPublications.AYear,
		AllUPublicationYears: allUPublicationYears,
		AllAPublicationYears: allAPublicationYears,
	})
}

func (h *Handler) RefreshAPublications(w http.ResponseWriter, r *http.Request, ctx Context) {
	var faculties []string

	switch ctx.Type {
	case "socs":
		faculties = vocabularies.Map["faculties_socs"]
	default:
		faculties = vocabularies.Map["faculties_core"]
	}

	facultyCols := append([]string{"all"}, faculties...)

	ptypes := vocabularies.Map["publication_types"]
	ptypes = append([]string{"all"}, ptypes...)

	locptypes := localize.VocabularyTerms(ctx.Loc, "publication_types")
	locptypes["all"] = "All"

	bindPublications := BindPublications{}
	if err := bind.Request(r, &bindPublications); err != nil {
		h.Logger.Warnw("publication dashboard could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	baseSearchUrl := h.PathFor("publications")

	// Publications with publication status "accepted"
	publications, err := generatePublicationsDashboard(facultyCols, ptypes, h.PublicationSearchIndex, baseSearchUrl, func(fac string, args *models.SearchArgs) *models.SearchArgs {
		args.WithFilter("publication_status", "accepted")
		args.WithFilter("status", "private", "public", "returned")
		if bindPublications.AYear != "" {
			args.WithFilter("year", bindPublications.AYear)
		}

		switch fac {
		case "all":
			args.WithFilter("faculty_id", faculties...)
		case "-":
			args.WithFilter("faculty_id", backends.MissingValue)
		default:
			args.WithFilter("faculty_id", fac)
		}

		return args
	})
	if err != nil {
		h.Logger.Errorw("Dashboard: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	w.Header().Add(
		"HX-Push-Url",
		h.URLFor(
			"dashboard_publications",
			"type", ctx.Type,
			"uyear", bindPublications.UYear,
			"ayear", bindPublications.AYear,
		).String(),
	)

	render.Partial(w, "dashboard/partials/tbl_publications", YieldTblPublications{
		Context:      ctx,
		Publications: publications,
		Faculties:    facultyCols,
		PTypes:       locptypes,
	})
}

func (h *Handler) RefreshUPublications(w http.ResponseWriter, r *http.Request, ctx Context) {
	var faculties []string

	switch ctx.Type {
	case "socs":
		faculties = vocabularies.Map["faculties_socs"]
	default:
		faculties = vocabularies.Map["faculties_core"]
	}

	facultyCols := append([]string{"all"}, faculties...)
	facultyCols = append(facultyCols, "UGent", "-")
	ptypes := append([]string{"all"}, vocabularies.Map["publication_types"]...)

	locptypes := localize.VocabularyTerms(ctx.Loc, "publication_types")
	locptypes["all"] = "All"

	bindPublications := BindPublications{}
	if err := bind.Request(r, &bindPublications); err != nil {
		h.Logger.Warnw("publication dashboard could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	// Publications with classification U
	baseSearchUrl := h.PathFor("publications")
	publications, err := generatePublicationsDashboard(facultyCols, ptypes, h.PublicationSearchIndex, baseSearchUrl, func(fac string, args *models.SearchArgs) *models.SearchArgs {
		args.WithFilter("classification", "U")
		args.WithFilter("status", "public")
		if bindPublications.UYear != "" {
			args.WithFilter("year", bindPublications.UYear)
		}

		switch fac {
		case "all":
			args.WithFilter("faculty_id", faculties...)
		case "-":
			args.WithFilter("faculty_id", backends.MissingValue)
		case "UGent":
			args.WithFilter("organization_id", "UGent")
		default:
			args.WithFilter("faculty_id", fac)
		}

		return args
	})
	if err != nil {
		h.Logger.Errorw("Dashboard: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	w.Header().Add(
		"HX-Push-Url",
		h.URLFor(
			"dashboard_publications",
			"type", ctx.Type,
			"uyear", bindPublications.UYear,
			"ayear", bindPublications.AYear,
		).String(),
	)

	render.Partial(w, "dashboard/partials/tbl_publications", YieldTblPublications{
		Context:      ctx,
		Publications: publications,
		Faculties:    facultyCols,
		PTypes:       locptypes,
	})
}

func generatePublicationsDashboard(faculties []string, ptypes []string, searcher backends.PublicationIndex, baseSearchUrl *url.URL, fn func(fac string, args *models.SearchArgs) *models.SearchArgs) (map[string]map[string][]string, error) {
	var publications = make(map[string]map[string][]string)

	for _, fac := range faculties {
		publications[fac] = map[string][]string{}

		for _, ptype := range ptypes {
			searchUrl := *baseSearchUrl
			searchArgs := models.NewSearchArgs()
			queryVals := searchUrl.Query()

			if ptype != "all" {
				searchArgs.WithFilter("type", ptype)
			}

			searchArgs = fn(fac, searchArgs)

			for f, varr := range searchArgs.Filters {
				for _, v := range varr {
					queryVals.Add(fmt.Sprintf("f[%s]", f), v)
				}
			}

			searchArgs.PageSize = 0
			searchArgs.Page = 1

			searchUrl.RawQuery = queryVals.Encode()

			hits, err := searcher.Search(searchArgs)
			if err != nil {
				publications[fac][ptype] = []string{"Error", ""}
			} else {
				publications[fac][ptype] = []string{fmt.Sprint(hits.Total), searchUrl.String()}
			}
		}
	}

	return publications, nil
}

func (h *Handler) allUPublicationYears() ([]string, error) {
	return getFilteredTokenValues(h.PublicationSearchIndex, "year", &models.SearchArgs{
		Filters: map[string][]string{
			"status":         {"public"},
			"classification": {"U"},
		},
	})
}

func (h *Handler) allAPublicationYears() ([]string, error) {
	return getFilteredTokenValues(h.PublicationSearchIndex, "year", &models.SearchArgs{
		Filters: map[string][]string{
			"status":             {"private", "public", "returned"},
			"publication_status": {"accepted"},
		},
	})
}

func getFilteredTokenValues(idx backends.PublicationIndex, field string, baseSearchArgs *models.SearchArgs) ([]string, error) {
	tokens := make([]string, 0)

	searchArgs := &models.SearchArgs{
		PageSize: 0,
		Page:     1,
		Facets:   []string{field},
		Filters:  map[string][]string{},
	}

	for f, values := range baseSearchArgs.Filters {
		searchArgs.WithFilter(f, values...)
	}

	hits, err := idx.Search(searchArgs)
	if err != nil {
		return nil, err
	}

	fieldFacets, ok := hits.Facets[field]
	if !ok {
		return tokens, nil
	}

	for _, fieldFacet := range fieldFacets {
		tokens = append(tokens, fieldFacet.Value)
	}

	return tokens, nil
}
