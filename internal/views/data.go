package views

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-locale/locale"
)

type Data struct {
	request *http.Request
	User    *models.User
	Locale  *locale.Locale
}

func NewData(r *http.Request) Data {
	return Data{
		request: r,
		User:    context.GetUser(r.Context()),
		Locale:  locale.Get(r.Context()),
	}
}

func (d Data) T(scope, key string, args ...interface{}) string {
	return d.Locale.Translate(scope, key, args...)
}

func (d Data) IsHTMXRequest() bool {
	return d.request.Header.Get("HX-Request") != ""
}

func (d Data) ActiveMenu() string {
	return context.GetActiveMenu(d.request.Context())
}
