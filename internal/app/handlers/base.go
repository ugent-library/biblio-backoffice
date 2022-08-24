package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/flash"
)

func init() {
	// register flash.Flash as a gob Type to make SecureCookieStore happy
	// see https://github.com/gin-contrib/sessions/issues/39
	gob.Register(flash.Flash{})
}

// TODO handlers should only have access to a url builder,
// the session and maybe the localizer
type BaseHandler struct {
	Router       *mux.Router
	Logger       backends.Logger
	SessionName  string
	SessionStore sessions.Store
	UserService  backends.UserService
	Localizer    *locale.Localizer
}

func (h BaseHandler) Wrap(fn func(http.ResponseWriter, *http.Request, BaseContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := h.NewContext(r, w)
		if err != nil {
			h.Logger.Errorw("could not create new context.", err)
			render.InternalServerError(w, r, err)
			return
		}
		fn(w, r, ctx)
	}
}

func (h BaseHandler) NewContext(r *http.Request, w http.ResponseWriter) (BaseContext, error) {
	session, err := h.SessionStore.Get(r, h.SessionName)
	if err != nil {
		return BaseContext{}, err
	}

	user, err := h.getUserFromSession(session, r, UserSessionKey)
	if err != nil {
		return BaseContext{}, fmt.Errorf("could not get user from session: %w", err)
	}

	originalUser, err := h.getUserFromSession(session, r, OriginalUserSessionKey)
	if err != nil {
		return BaseContext{}, fmt.Errorf("could not get original user from session: %w", err)
	}

	flash, err := h.getFlashFromSession(session, r, w)
	if err != nil {
		return BaseContext{}, fmt.Errorf("could not get flash message from session: %w", err)
	}

	return BaseContext{
		CurrentURL:   r.URL,
		Flash:        flash,
		Locale:       h.Localizer.GetLocale(r.Header.Get("Accept-Language")),
		User:         user,
		OriginalUser: originalUser,
		CSRFToken:    csrf.Token(r),
		CSRFTag:      csrf.TemplateField(r),
	}, nil
}

func (h BaseHandler) AddSessionFlash(r *http.Request, w http.ResponseWriter, f flash.Flash) error {
	session, err := h.SessionStore.Get(r, h.SessionName)
	if err != nil {
		return fmt.Errorf("could not get session from store: %w", err)
	}

	session.AddFlash(f, FlashSessionKey)

	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("could not save data to session: %w", err)
	}

	return nil
}

func (h BaseHandler) getFlashFromSession(session *sessions.Session, r *http.Request, w http.ResponseWriter) ([]flash.Flash, error) {
	sessionFlashes := session.Flashes(FlashSessionKey)

	if err := session.Save(r, w); err != nil {
		return []flash.Flash{}, fmt.Errorf("could not save data to session: %w", err)
	}

	flashes := []flash.Flash{}
	for _, f := range sessionFlashes {
		flashes = append(flashes, f.(flash.Flash))
	}

	return flashes, nil
}

func (h BaseHandler) getUserFromSession(session *sessions.Session, r *http.Request, sessionKey string) (*models.User, error) {
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
			h.Logger.Panic("Could not reverse route %s: %w", name, err)
		}
		return u
	}
	h.Logger.Panicf("Could not find route named %s", name)
	return nil

}

func (h BaseHandler) URLFor(name string, vars ...string) *url.URL {
	if route := h.Router.Get(name); route != nil {
		u, err := route.URL(vars...)
		if err != nil {
			h.Logger.Panic("Could not reverse route %s: %w", name, err)
		}
		return u
	}
	h.Logger.Panic("Could not find route named %s", name)
	return nil
}

type BaseContext struct {
	CurrentURL   *url.URL
	Flash        []flash.Flash
	Locale       *locale.Locale
	User         *models.User
	OriginalUser *models.User
	CSRFToken    string
	CSRFTag      template.HTML
}
