package datasetviewing

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/context"
)

type Handler struct {
	handlers.Base
	Repo backends.Repository
}

type Context struct {
	handlers.BaseContext
	Dataset *models.Dataset
}

// TODO check edit rights
func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := context.GetDataset(r.Context())
		fn(w, r, Context{
			BaseContext: handlers.NewBaseContext(h.Base, r),
			Dataset:     d,
		})
	}
}
