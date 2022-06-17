package datasets

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
)

type YieldShow struct {
	Ctx        ViewContext
	PageTitle  string
	ActiveMenu string
	SearchArgs *models.SearchArgs
}

func (c *Controller) Show(w http.ResponseWriter, r *http.Request, ctx ViewContext) {
	// TODO bind search args
	searchArgs := models.NewSearchArgs()

	render.Wrap(w, "layouts/default", "dataset/show_page", YieldShow{
		Ctx:        ctx,
		PageTitle:  "Dataset - Biblio",
		ActiveMenu: "datasets",
		SearchArgs: searchArgs,
	})
}
