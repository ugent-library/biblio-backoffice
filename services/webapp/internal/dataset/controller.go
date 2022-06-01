package dataset

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
)

type Context struct {
	Locale  *locale.Locale
	Dataset *models.Dataset
	CanEdit bool
}

func (c Context) WithDataset(d *models.Dataset) Context {
	c.Dataset = d
	return c
}

type Controller struct {
	store backends.Store
}

func NewController(store backends.Store) *Controller {
	return &Controller{
		store: store,
	}
}

func (c *Controller) WithContext(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loc := locale.Get(r.Context())
		user := context.GetUser(r.Context())
		dataset := context.GetDataset(r.Context())

		fn(w, r, Context{
			Locale:  loc,
			Dataset: dataset,
			CanEdit: user.CanEditDataset(dataset),
		})
	}
}
