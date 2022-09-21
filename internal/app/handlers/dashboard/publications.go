package dashboard

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
)

type YieldPublications struct {
	Context
	PageTitle     string
	ActiveNav     string
	UPublications map[string]map[string]int
	APublications map[string]map[string]int
	Faculties     []string
	PTypes        map[string]string
}

func (h *Handler) Publications(w http.ResponseWriter, r *http.Request, ctx Context) {
	faculties := vocabularies.Map["faculties"]
	ptypes := vocabularies.Map["publication_types"]

	faculties = append([]string{"all"}, faculties...)
	ptypes = append([]string{"all"}, ptypes...)

	locptypes := localize.VocabularyTerms(ctx.Locale, "publication_types")
	locptypes["all"] = "All"

	// Publications with classification U

	var searcher = h.PublicationSearchService.WithScope("status", "public")

	err, uPublications := generateDashboard(faculties, ptypes, searcher, func(args *models.SearchArgs) *models.SearchArgs {
		args.WithFilter("classification", "U")
		return args
	})

	if err != nil {
		h.Logger.Errorw("publication search: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	// Publications with publication status "accepted"

	searcher = h.PublicationSearchService

	err, aPublications := generateDashboard(faculties, ptypes, searcher, func(args *models.SearchArgs) *models.SearchArgs {
		args.WithFilter("publication_status", "accepted")
		return args
	})

	if err != nil {
		h.Logger.Errorw("publication search: could not execute search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "layouts/default", "dashboard/pages/publications", YieldPublications{
		Context:       ctx,
		PageTitle:     "Dashboard - Publications - Biblio",
		ActiveNav:     "dashboard",
		UPublications: uPublications,
		APublications: aPublications,
		Faculties:     faculties,
		PTypes:        locptypes,
	})
}

func generateDashboard(faculties []string, ptypes []string, searcher backends.PublicationSearchService, fn func(args *models.SearchArgs) *models.SearchArgs) (error, map[string]map[string]int) {
	var publications = make(map[string]map[string]int)

	for _, fac := range faculties {
		publications[fac] = map[string]int{}

		for _, ptype := range ptypes {
			searchArgs := models.NewSearchArgs()

			if fac != "all" {
				searchArgs.WithFilter("faculty", fac)
			}

			if ptype != "all" {
				searchArgs.WithFilter("type", ptype)
			}

			searchArgs = fn(searchArgs)

			hits, err := searcher.Search(searchArgs)
			if err != nil {
				return err, nil
			}

			publications[fac][ptype] = hits.Total
		}
	}

	return nil, publications
}
