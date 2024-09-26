package dashboard

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	dashboardviews "github.com/ugent-library/biblio-backoffice/views/dashboard"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
)

type BindPublications struct {
	UYear string `query:"uyear" form:"uyear"`
	AYear string `query:"ayear" form:"ayear"`
}

func CuratorPublications(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	typ := bind.PathValue(r, "type") //TODO: bind via middleware
	var faculties []string

	var activeSubNav string

	socs := vocabularies.Map["faculties_socs"]
	core := vocabularies.Map["faculties_core"]

	switch typ {
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

	allUPublicationYears, err := allUPublicationYears(c.PublicationSearchIndex)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	allAPublicationYears, err := allAPublicationYears(c.PublicationSearchIndex)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	bindPublications := BindPublications{}
	if err := bind.Request(r, &bindPublications); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	// Publications with classification U
	uFacultyCols := append(aFacultyCols, []string{"UGent", "-"}...)

	uSearcher := c.PublicationSearchIndex
	baseSearchUrl := c.PathTo("publications")

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
		c.HandleError(w, r, err)
		return
	}

	// Publications with publication status "accepted"

	aSearcher := c.PublicationSearchIndex

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
		c.HandleError(w, r, err)
		return
	}

	dashboardviews.CuratorDashboardPublications(c, &dashboardviews.CuratorDashboardPublicationsArgs{
		Type:                 typ,
		ActiveSubNav:         activeSubNav,
		UPublications:        uPublications,
		APublications:        aPublications,
		UFaculties:           uFacultyCols,
		AFaculties:           aFacultyCols,
		UYear:                bindPublications.UYear,
		AYear:                bindPublications.AYear,
		AllUPublicationYears: allUPublicationYears,
		AllAPublicationYears: allAPublicationYears,
	}).Render(r.Context(), w)
}

func RefreshAPublications(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	typ := bind.PathValue(r, "type") //TODO: bind via middleware

	var faculties []string

	switch typ {
	case "socs":
		faculties = vocabularies.Map["faculties_socs"]
	default:
		faculties = vocabularies.Map["faculties_core"]
	}

	facultyCols := append([]string{"all"}, faculties...)

	ptypes := vocabularies.Map["publication_types"]
	ptypes = append([]string{"all"}, ptypes...)

	bindPublications := BindPublications{}
	if err := bind.Request(r, &bindPublications); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	baseSearchUrl := c.PathTo("publications")

	// Publications with publication status "accepted"
	publications, err := generatePublicationsDashboard(facultyCols, ptypes, c.PublicationSearchIndex, baseSearchUrl, func(fac string, args *models.SearchArgs) *models.SearchArgs {
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
		c.HandleError(w, r, err)
		return
	}

	w.Header().Add(
		"HX-Push-Url",
		c.URLTo(
			"dashboard_publications",
			"type", typ,
			"uyear", bindPublications.UYear,
			"ayear", bindPublications.AYear,
		).String(),
	)

	dashboardviews.CuratorDashboardTblPublications(c, facultyCols, publications).Render(r.Context(), w)
}

func RefreshUPublications(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	typ := bind.PathValue(r, "type") //TODO: bind via middleware
	var faculties []string

	switch typ {
	case "socs":
		faculties = vocabularies.Map["faculties_socs"]
	default:
		faculties = vocabularies.Map["faculties_core"]
	}

	facultyCols := append([]string{"all"}, faculties...)
	facultyCols = append(facultyCols, "UGent", "-")
	ptypes := append([]string{"all"}, vocabularies.Map["publication_types"]...)

	bindPublications := BindPublications{}
	if err := bind.Request(r, &bindPublications); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	// Publications with classification U
	baseSearchUrl := c.PathTo("publications")
	publications, err := generatePublicationsDashboard(facultyCols, ptypes, c.PublicationSearchIndex, baseSearchUrl, func(fac string, args *models.SearchArgs) *models.SearchArgs {
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
		c.HandleError(w, r, err)
		return
	}

	w.Header().Add(
		"HX-Push-Url",
		c.URLTo("dashboard_publications", "type", typ, "uyear", bindPublications.UYear, "ayear", bindPublications.AYear).String(),
	)

	dashboardviews.CuratorDashboardTblPublications(c, facultyCols, publications).Render(r.Context(), w)
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

func allUPublicationYears(idx backends.PublicationIndex) ([]string, error) {
	return getFilteredTokenValues(idx, "year", &models.SearchArgs{
		Filters: map[string][]string{
			"status":         {"public"},
			"classification": {"U"},
		},
	})
}

func allAPublicationYears(idx backends.PublicationIndex) ([]string, error) {
	return getFilteredTokenValues(idx, "year", &models.SearchArgs{
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
