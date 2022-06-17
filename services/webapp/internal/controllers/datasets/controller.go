package datasets

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
)

// TODO make base context
type ViewContext struct {
	// base context
	Locale       *locale.Locale
	User         *models.User
	OriginalUser *models.User
	CSRFToken    string
	CSRFTag      template.HTML
	// specific to context
	Dataset *models.Dataset
}
type EditContext struct {
	// base context
	Locale       *locale.Locale
	User         *models.User
	OriginalUser *models.User
	CSRFToken    string
	CSRFTag      template.HTML
	// specific to context
	Dataset *models.Dataset
}

type Controller struct {
	store backends.Repository
}

func NewController(store backends.Repository) *Controller {
	return &Controller{
		store: store,
	}
}

func (c *Controller) WithViewContext(fn func(http.ResponseWriter, *http.Request, ViewContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := context.GetUser(r.Context())
		dataset := context.GetDataset(r.Context())
		fn(w, r, ViewContext{
			// base context
			Locale:       locale.Get(r.Context()),
			User:         user,
			OriginalUser: context.GetOriginalUser(r.Context()),
			CSRFToken:    csrf.Token(r),
			CSRFTag:      csrf.TemplateField(r),
			// specific to context
			Dataset: dataset,
		})
	}
}
func (c *Controller) WithEditContext(fn func(http.ResponseWriter, *http.Request, EditContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := context.GetUser(r.Context())
		dataset := context.GetDataset(r.Context())
		fn(w, r, EditContext{
			// base context
			Locale:       locale.Get(r.Context()),
			User:         user,
			OriginalUser: context.GetOriginalUser(r.Context()),
			CSRFToken:    csrf.Token(r),
			CSRFTag:      csrf.TemplateField(r),
			// specific to context
			Dataset: dataset,
		})
	}
}
