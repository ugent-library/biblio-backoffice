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
	CanEdit bool
}

type Controller struct {
	store        backends.Store
	abstractView render.Partial
}

func NewController(store backends.Store) *Controller {
	return &Controller{
		store: store,
		abstractView: render.NewPartial(
			"dataset/_add_abstract",
			"dataset/_create_abstract",
			"dataset/_abstracts_table",
		),
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
