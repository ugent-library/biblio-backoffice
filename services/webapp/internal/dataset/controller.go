package dataset

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/render"
	"github.com/ugent-library/go-locale/locale"
)

type Context struct {
	Locale  *locale.Locale
	Dataset *models.Dataset
}

type Controller struct {
	store              backends.Store
	addAbstractPartial render.Partial
}

func NewController(store backends.Store) *Controller {
	return &Controller{
		store:              store,
		addAbstractPartial: render.NewPartial("add-abstract", "dataset/_add_abstract"),
	}
}

func (c *Controller) WithContext(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loc := locale.Get(r.Context())
		dataset := context.GetDataset(r.Context())

		fn(w, r, Context{
			Locale:  loc,
			Dataset: dataset,
		})
	}
}
