package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
)

// TODO handlers should only have access to a url builder,
// the session and maybe the localizer
type BaseHandler struct {
	Router       *mux.Router
	SessionName  string
	SessionStore sessions.Store
	UserService  backends.UserService
	Localizer    *locale.Localizer
}

func (h BaseHandler) Wrap(fn func(http.ResponseWriter, *http.Request, BaseContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := h.NewContext(r)
		if err != nil {
			render.InternalServerError(w, r, err)
			return
		}
		fn(w, r, ctx)
	}
}

func (h BaseHandler) NewContext(r *http.Request) (BaseContext, error) {
	user, err := h.getUserFromSession(r, UserSessionKey)
	if err != nil {
		return BaseContext{}, err
	}

	originalUser, err := h.getUserFromSession(r, OriginalUserSessionKey)
	if err != nil {
		return BaseContext{}, err
	}

	return BaseContext{
		Locale:       h.Localizer.GetLocale(r.Header.Get("Accept-Language")),
		User:         user,
		OriginalUser: originalUser,
		CSRFToken:    csrf.Token(r),
		CSRFTag:      csrf.TemplateField(r),
	}, nil
}

func (h BaseHandler) getUserFromSession(r *http.Request, sessionKey string) (*models.User, error) {
	session, err := h.SessionStore.Get(r, h.SessionName)
	if err != nil {
		return nil, err
	}
	userID := session.Values[sessionKey]
	if userID == nil {
		return nil, nil
	}

	user, err := h.UserService.GetUser(userID.(string))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h BaseHandler) PathFor(name string, vars ...string) *url.URL {
	if route := h.Router.Get(name); route != nil {
		u, err := route.URLPath(vars...)
		if err != nil {
			log.Panic(fmt.Errorf("can't reverse route %s: %w", name, err))
		}
		return u
	}
	log.Panic(fmt.Errorf("route %s not found", name))
	return nil

}

func (h BaseHandler) URLFor(name string, vars ...string) *url.URL {
	if route := h.Router.Get(name); route != nil {
		u, err := route.URL(vars...)
		if err != nil {
			log.Panic(fmt.Errorf("can't reverse route %s: %w", name, err))
		}
		return u
	}
	log.Panic(fmt.Errorf("route %s not found", name))
	return nil
}

type BaseContext struct {
	Locale       *locale.Locale
	User         *models.User
	OriginalUser *models.User
	CSRFToken    string
	CSRFTag      template.HTML
}

func (c BaseContext) T(key string, args ...any) string {
	return c.Locale.T(key, args...)
}

func (c BaseContext) TS(scope, key string, args ...any) string {
	return c.Locale.TS(scope, key, args...)
}
