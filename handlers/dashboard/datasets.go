package dashboard

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
	"github.com/ugent-library/bind"
)

func CuratorDatasets(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	typ := bind.PathValue(r, "type") //TODO: bind via middleware
	var faculties []string

	var activeSubNav string

	switch typ {
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

	aSearcher := c.DatasetSearchIndex.WithScope("status", "private", "public", "returned")
	baseSearchUrl := c.PathTo("datasets")

	datasets, err := generateDatasetsDashboard(faculties, ptypes, aSearcher, baseSearchUrl, func(args *models.SearchArgs) *models.SearchArgs {
		return args
	})

	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	views.CuratorDashboardDatasets(c, &views.CuratorDashboardDatasetsArgs{
		ActiveSubNav: activeSubNav,
		PTypes:       locptypes,
		Datasets:     datasets,
		Faculties:    faculties,
	}).Render(r.Context(), w)
}

func generateDatasetsDashboard(faculties []string, ptypes []string, searcher backends.DatasetIndex, baseSearchUrl *url.URL, fn func(args *models.SearchArgs) *models.SearchArgs) (map[string]map[string][]string, error) {
	var datasets = make(map[string]map[string][]string)

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
		}
	}

	return datasets, nil
}
