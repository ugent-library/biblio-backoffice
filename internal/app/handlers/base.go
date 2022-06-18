package handlers

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/context"
)

type BaseContext struct {
	Locale       *locale.Locale
	User         *models.User
	OriginalUser *models.User
	CSRFToken    string
	CSRFTag      template.HTML
}

func NewBaseContext(base Base, r *http.Request) BaseContext {
	return BaseContext{
		Locale:       base.Localizer.GetLocale(r.Header.Get("Accept-Language")),
		User:         context.GetUser(r.Context()),
		OriginalUser: context.GetOriginalUser(r.Context()),
		CSRFToken:    csrf.Token(r),
		CSRFTag:      csrf.TemplateField(r),
	}
}

type Base struct {
	Localizer *locale.Localizer
}
