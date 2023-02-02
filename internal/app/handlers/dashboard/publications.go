package dashboard

import (
	"fmt"
	"log"
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
	Faculties     []string
	PTypes        map[string]string
}

func (h *Handler) Publications(w http.ResponseWriter, r *http.Request, ctx Context) {
	var faculties []string

	var activeNav string

	switch ctx.Type {
	case "socs":
		faculties = vocabularies.Map["faculties_socs"]
		activeNav = "dashboard_publications_socs"
	default:
		faculties = vocabularies.Map["faculties_core"]
		activeNav = "dashboard_publications_faculties"
	}

	ptypes := vocabularies.Map["publication_types"]

	faculties = append([]string{"all"}, faculties...)
	ptypes = append([]string{"all"}, ptypes...)

	locptypes := localize.VocabularyTerms(ctx.Locale, "publication_types")
	locptypes["all"] = "All"

	// Publications with classification U

	uSearcher := h.PublicationSearchService.WithScope("status", "public", "returned")
	baseSearchUrl := h.PathFor("publications")

	uPublications, err := generatePublicationsDashboard(faculties, ptypes, uSearcher, baseSearchUrl, func(args *models.SearchArgs) *models.SearchArgs {
		args.WithFilter("classification", "U")
		args.WithFilter("status", "public")
		return args
	})

	if err != nil {
		h.Logger.Errorw("Dashboard: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// Publications with publication status "accepted"

	aSearcher := h.PublicationSearchService.WithScope("status", "private", "public", "returned")

	aPublications, err := generatePublicationsDashboard(faculties, ptypes, aSearcher, baseSearchUrl, func(args *models.SearchArgs) *models.SearchArgs {
		args.WithFilter("publication_status", "accepted")
		return args
	})

	if err != nil {
		h.Logger.Errorw("Dashboard: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	log.Println(activeNav)

	render.Layout(w, "layouts/default", "dashboard/pages/publications", YieldPublications{
		Context:       ctx,
		PageTitle:     "Dashboard - Publications - Biblio",
		ActiveNav:     activeNav,
		UPublications: uPublications,
		APublications: aPublications,
		Faculties:     faculties,
		PTypes:        locptypes,
	})
}

func generatePublicationsDashboard(faculties []string, ptypes []string, searcher backends.PublicationSearchService, baseSearchUrl *url.URL, fn func(args *models.SearchArgs) *models.SearchArgs) (map[string]map[string][]string, error) {
	var publications = make(map[string]map[string][]string)

	// pool := pond.New(100, 300)
	// defer pool.StopAndWait()
	// group := pool.Group()

	for _, fac := range faculties {
		publications[fac] = map[string][]string{}

		for _, ptype := range ptypes {
			searchUrl := *baseSearchUrl
			searchArgs := models.NewSearchArgs()
			queryVals := searchUrl.Query()

			if fac != "all" {
				searchArgs.WithFilter("faculty", fac)
			} else {
				searchArgs.WithFilter("faculty", faculties...)
			}

			if ptype != "all" {
				searchArgs.WithFilter("type", ptype)
			}

			searchArgs = fn(searchArgs)

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

			// f := fac
			// p := ptype

			// var lock sync.Mutex
			// group.Submit(func(f string, pt string, p map[string]map[string][]string, searchUrl string) func() {
			// 	return func() {
			// 		lock.Lock()
			// 		hits, err := searcher.Search(searchArgs)
			// 		if err != nil {
			// 			p[f][pt] = []string{"Error", ""}
			// 		} else {
			// 			p[f][pt] = []string{fmt.Sprint(hits.Total), searchUrl}
			// 		}
			// 		lock.Unlock()
			// 	}
			// }(f, p, publications, searchUrl.String()))
		}
	}

	// group.Wait()

	return publications, nil
}
