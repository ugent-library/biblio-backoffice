package views

import (
	"html/template"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-locale/locale"
	"github.com/unrolled/render"
)

type Data struct {
	renderer *render.Render
	request  *http.Request
	User     *models.User
	Locale   *locale.Locale
}

func NewData(renderer *render.Render, r *http.Request) Data {
	return Data{
		renderer: renderer,
		request:  r,
		User:     context.GetUser(r.Context()),
		Locale:   locale.Get(r.Context()),
	}
}

func (d Data) T(scope, key string, args ...interface{}) string {
	return d.Locale.Translate(scope, key, args...)
}

func (d Data) Partial(tmpl string, data interface{}) (template.HTML, error) {
	return RenderPartial(d.renderer, tmpl, data)
}

// func (d Data) PartialFallback(tmpl, tmplFallback string, data interface{}) (template.HTML, error) {
// 	return RenderPartialFallback(d.renderer, tmpl, tmplFallback, data)
// }

func (d Data) IsHTMXRequest() bool {
	return d.request.Header.Get("HX-Request") != ""
}

func (d Data) ActiveMenu() string {
	return context.GetActiveMenu(d.request.Context())
}
