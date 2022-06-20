package dataset

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
)

type EditContext struct {
	Locale  *locale.Locale
	Dataset *models.Dataset
	CanEdit bool
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
	store                backends.Repository
	projectSearchService backends.ProjectSearchService
	projectService       backends.ProjectService
}

func NewController(store backends.Repository, projectSearchService backends.ProjectSearchService, projectService backends.ProjectService) *Controller {
	return &Controller{
		store:                store,
		projectSearchService: projectSearchService,
		projectService:       projectService,
	}
}

func (c *Controller) WithEditContext(fn func(http.ResponseWriter, *http.Request, EditContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := context.GetUser(r.Context())
		dataset := context.GetDataset(r.Context())
		fn(w, r, EditContext{
			Locale:  locale.Get(r.Context()),
			Dataset: dataset,
			CanEdit: user.CanEditDataset(dataset),
		})
	}
}
