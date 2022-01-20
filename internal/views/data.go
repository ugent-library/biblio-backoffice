package views

import (
	"html/template"
	"time"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-locale/locale"
)

type Flash struct {
	Type         string
	Message      string
	DismissAfter time.Duration
}

type Data struct {
	Mode              string
	RenderPartialFunc func(string, interface{}) (template.HTML, error)
	User              *models.User
	OriginalUser      *models.User
	Locale            *locale.Locale
	Flash             []Flash
	CSRFToken         string
	CSRFTag           template.HTML
	ActiveMenu        string
	IsHTMXRequest     bool // TODO get rid of this
	Data              interface{}
}

func (d *Data) ViewData(data interface{}) *Data {
	vd := *d
	vd.Data = data
	return &vd
}

func (d *Data) D() interface{} {
	return d.Data
}

func (d *Data) T(scope, key string, args ...interface{}) string {
	return d.Locale.Translate(scope, key, args...)
}

func (d *Data) RenderPartial(tmpl string, data interface{}) (template.HTML, error) {
	return d.RenderPartialFunc(tmpl, data)
}
