package dashboard

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
)

type YieldDatasets struct {
	Context
	PageTitle    string
	ActiveNav    string
	ActiveSubNav string
	Datasets     map[string]map[string][]string
	Faculties    []string
	PTypes       map[string]string
}

func (h *Handler) Datasets(w http.ResponseWriter, r *http.Request, ctx Context) {
	var faculties []string

	var activeSubNav string

	switch ctx.Type {
	case "socs":
		faculties = vocabularies.Map["faculties_socs"]
		activeSubNav = "dashboard_datasets_socs"
	default:
		faculties = vocabularies.Map["faculties_core"]
		activeSubNav = "dashboard_datasets_faculties"
	}

	faculties = append([]string{"all"}, faculties...)
	ptypes := []string{"all"}

	locptypes := make(map[string]string)
	locptypes["all"] = "All"

	aSearcher := h.DatasetSearchIndex.WithScope("status", "private", "public", "returned")
	baseSearchUrl := h.PathFor("datasets")

	datasets, err := generateDatasetsDashboard(faculties, ptypes, aSearcher, baseSearchUrl, func(args *models.SearchArgs) *models.SearchArgs {
		return args
	})

	if err != nil {
		h.Logger.Errorw("Dashboard: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "layouts/default", "dashboard/pages/datasets", YieldDatasets{
		Context:      ctx,
		PageTitle:    "Dashboard - Datasets - Biblio",
		ActiveNav:    "dashboard",
		ActiveSubNav: activeSubNav,
		PTypes:       locptypes,
		Datasets:     datasets,
		Faculties:    faculties,
	})
}

func generateDatasetsDashboard(faculties []string, ptypes []string, searcher backends.DatasetIndex, baseSearchUrl *url.URL, fn func(args *models.SearchArgs) *models.SearchArgs) (map[string]map[string][]string, error) {
	var datasets = make(map[string]map[string][]string)

	// pool := pond.New(100, 300)
	// defer pool.StopAndWait()
	// group := pool.Group()

	for _, fac := range faculties {
		datasets[fac] = map[string][]string{}

		for _, ptype := range ptypes {
			searchUrl := *baseSearchUrl
			searchArgs := models.NewSearchArgs()
			queryVals := searchUrl.Query()

			if fac != "all" {
				searchArgs.WithFilter("faculty_id", fac)
			} else {
				searchArgs.WithFilter("faculty_id", faculties...)
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
				datasets[fac][ptype] = []string{"Error", ""}
			} else {
				datasets[fac][ptype] = []string{fmt.Sprint(hits.Total), searchUrl.String()}
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

	return datasets, nil
}
