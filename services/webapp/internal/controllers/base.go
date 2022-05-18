package controllers

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-oidc/oidc"
	"github.com/unrolled/render"
)

type Base struct {
	Engine       *engine.Engine
	Mode         string
	BaseURL      *url.URL
	Router       *mux.Router
	Render       *render.Render
	Localizer    *locale.Localizer
	SessionName  string
	SessionStore sessions.Store
	OIDC         *oidc.Client
}

func (c *Base) Session(r *http.Request) (*sessions.Session, error) {
	return c.SessionStore.Get(r, c.SessionName)
}

func (c *Base) RenderPartial(tmpl string, data interface{}) (template.HTML, error) {
	buf := &bytes.Buffer{}
	var err error
	if t := c.Render.TemplateLookup(tmpl); t != nil {
		err = t.Execute(buf, data)
	}
	return template.HTML(buf.String()), err
}

func (c *Base) ViewData(r *http.Request, data interface{}, flash ...views.Flash) *views.Data {
	return &views.Data{
		Mode:              c.Mode,
		RenderPartialFunc: c.RenderPartial,
		Locale:            locale.Get(r.Context()),
		Vocabularies:      vocabularies.Map,
		User:              context.GetUser(r.Context()),
		OriginalUser:      context.GetOriginalUser(r.Context()),
		Data:              data,
		Flash:             flash,
		CSRFToken:         csrf.Token(r),
		CSRFTag:           csrf.TemplateField(r),
		ActiveMenu:        context.GetActiveMenu(r.Context()),
		IsHTMXRequest:     r.Header.Get("HX-Request") != "",
	}
}
