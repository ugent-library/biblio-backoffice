package dataset

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
)

type EditContext struct {
	PathParams url.Values
	Locale     *locale.Locale
	Dataset    *models.Dataset
	CanEdit    bool
}

func (c EditContext) RenderYield(w http.ResponseWriter, tmpl string, yield interface{}) {
	render.Render(w, tmpl, struct {
		Locale  *locale.Locale
		CanEdit bool
		Yield   interface{}
	}{
		Locale:  c.Locale,
		CanEdit: c.CanEdit,
		Yield:   yield,
	})
}

type Controller struct {
	store backends.Store
}

func NewController(store backends.Store) *Controller {
	return &Controller{
		store: store,
	}
}

func (c *Controller) WithEditContext(fn func(http.ResponseWriter, *http.Request, EditContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := url.Values{}
		for k, v := range mux.Vars(r) {
			p.Set(k, v)
		}
		user := context.GetUser(r.Context())
		dataset := context.GetDataset(r.Context())
		fn(w, r, EditContext{
			PathParams: p,
			Locale:     locale.Get(r.Context()),
			Dataset:    dataset,
			CanEdit:    user.CanEditDataset(dataset),
		})
	}
}
