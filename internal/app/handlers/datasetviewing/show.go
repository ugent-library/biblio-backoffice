package datasetviewing

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
)

type YieldShow struct {
	Context
	PageTitle  string
	ActiveMenu string
	SearchArgs *models.SearchArgs
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request, ctx Context) {
	// TODO bind search args
	searchArgs := models.NewSearchArgs()

	render.Wrap(w, "layouts/default", "dataset/show_page", YieldShow{
		Context:    ctx,
		PageTitle:  "Dataset - Biblio",
		ActiveMenu: "datasets",
		SearchArgs: searchArgs,
	})
}
