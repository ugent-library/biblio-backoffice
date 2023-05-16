package dashboard

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/internal/app/localize"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/render"
	"github.com/ugent-library/biblio-backoffice/internal/vocabularies"
)

type YieldPublications struct {
	Context
	PageTitle     string
	ActiveNav     string
	UPublications map[string]map[string][]string
	APublications map[string]map[string][]string
	UFaculties    []string
	AFaculties    []string
	PTypes        map[string]string
}

func (h *Handler) Publications(w http.ResponseWriter, r *http.Request, ctx Context) {
	var faculties []string

	var activeNav string

	socs := vocabularies.Map["faculties_socs"]
	core := vocabularies.Map["faculties_core"]
	all := vocabularies.Map["faculties"]

	switch ctx.Type {
	case "socs":
		faculties = socs
		activeNav = "dashboard_publications_socs"
	default:
		faculties = core
		activeNav = "dashboard_publications_faculties"
	}

	faculties = append([]string{"all"}, faculties...)

	ptypes := vocabularies.Map["publication_types"]
	ptypes = append([]string{"all"}, ptypes...)

	locptypes := localize.VocabularyTerms(ctx.Locale, "publication_types")
	locptypes["all"] = "All"

	// Publications with classification U
	ufaculties := faculties
	ufaculties = append(ufaculties, []string{"UGent", "-"}...)

	uSearcher := h.PublicationSearchService.NewIndex()
	baseSearchUrl := h.PathFor("publications")

	uPublications, err := generatePublicationsDashboard(ufaculties, ptypes, uSearcher, baseSearchUrl, func(fac string, args *models.SearchArgs) *models.SearchArgs {
		args.WithFilter("classification", "U")
		args.WithFilter("status", "public")

		switch fac {
		case "all":
			args.WithFilter("faculty", faculties...)
		case "-":
			args.WithFilter("!faculty", all...)
		case "UGent":
			args.WithFilter("department.id", "UGent")
		default:
			args.WithFilter("faculty", fac)
		}

		return args
	})

	if err != nil {
		h.Logger.Errorw("Dashboard: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// Publications with publication status "accepted"

	aSearcher := h.PublicationSearchService.NewIndex()

	aPublications, err := generatePublicationsDashboard(faculties, ptypes, aSearcher, baseSearchUrl, func(fac string, args *models.SearchArgs) *models.SearchArgs {
		args.WithFilter("publication_status", "accepted")
		args.WithFilter("status", "private", "public", "returned")

		switch fac {
		case "all":
			args.WithFilter("faculty", faculties...)
		case "-":
			args.WithFilter("!faculty", all...)
		default:
			args.WithFilter("faculty", fac)
		}

		return args
	})

	if err != nil {
		h.Logger.Errorw("Dashboard: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "layouts/default", "dashboard/pages/publications", YieldPublications{
		Context:       ctx,
		PageTitle:     "Dashboard - Publications - Biblio",
		ActiveNav:     activeNav,
		UPublications: uPublications,
		APublications: aPublications,
		UFaculties:    ufaculties,
		AFaculties:    faculties,
		PTypes:        locptypes,
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
