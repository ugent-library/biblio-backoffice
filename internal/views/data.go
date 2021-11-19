package views

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-locale/locale"
	"github.com/unrolled/render"
)

type Flash struct {
	Type         string
	Message      string
	DismissAfter time.Duration
}

type Data struct {
	renderer *render.Render
	request  *http.Request
	User     *models.User
	Locale   *locale.Locale
	Flash    []Flash
	Data     interface{}
}

func NewData(renderer *render.Render, r *http.Request, data interface{}, flashes ...Flash) *Data {
	return &Data{
		renderer: renderer,
		request:  r,
		User:     context.GetUser(r.Context()),
		Locale:   locale.Get(r.Context()),
		Data:     data,
		Flash:    flashes,
	}
}

func (d *Data) NewData(data interface{}) *Data {
	return NewData(d.renderer, d.request, data)
}

func (d *Data) D() interface{} {
	return d.Data
}

func (d *Data) T(scope, key string, args ...interface{}) string {
	return d.Locale.Translate(scope, key, args...)
}

func (d *Data) Partial(tmpl string, data interface{}) (template.HTML, error) {
	return RenderPartial(d.renderer, tmpl, data)
}

func (d *Data) CSRFToken() string {
	return csrf.Token(d.request)
}

func (d *Data) CSRFTag() template.HTML {
	return csrf.TemplateField(d.request)
}

func (d *Data) OriginalUser() *models.User {
	return context.GetOriginalUser(d.request.Context())
}

func (d *Data) IsHTMXRequest() bool {
	return d.request.Header.Get("HX-Request") != ""
}

func (d *Data) ActiveMenu() string {
	return context.GetActiveMenu(d.request.Context())
}
