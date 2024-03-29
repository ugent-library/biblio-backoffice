package ctx

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/leonelquinteros/gotext"
	"github.com/nics/ich"
	"github.com/oklog/ulid/v2"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render/flash"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/mix"
	"github.com/ugent-library/zaphttp"
	"github.com/unrolled/secure"
	"go.uber.org/zap"
)

const (
	UserIDKey           = "user_id"
	OriginalUserIDKey   = "original_user_id"
	UserRoleKey         = "user_role"
	OriginalUserRoleKey = "original_user_role"
	FlashCookiePrefix   = "flash"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var ctxKey = contextKey("ctx")

func Get(r *http.Request) *Ctx {
	return r.Context().Value(ctxKey).(*Ctx)
}

func Set(config Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// setup ctx
			c := &Ctx{
				Config:    config,
				host:      r.Host,
				scheme:    r.URL.Scheme,
				Log:       zaphttp.Logger(r.Context()).Sugar(),
				Loc:       config.Loc,
				CSRFToken: csrf.Token(r),
				CSPNonce:  secure.CSPNonce(r.Context()),
			}
			if c.scheme == "" {
				c.scheme = "http"
			}

			r = r.WithContext(context.WithValue(r.Context(), ctxKey, c))

			// load user from session
			session, err := c.SessionStore.Get(r, c.SessionName)
			if err != nil {
				c.HandleError(w, r, err)
				return
			}
			user, err := c.getUserFromSession(r, session, UserIDKey)
			if err != nil {
				c.HandleError(w, r, fmt.Errorf("could not get user from session: %w", err))
				return
			}
			originalUser, err := c.getUserFromSession(r, session, OriginalUserIDKey)
			if err != nil {
				c.HandleError(w, r, fmt.Errorf("could not get original user from session: %w", err))
				return
			}

			c.User = user
			c.UserRole = c.getUserRoleFromSession(session)
			c.OriginalUser = originalUser

			// load flash from cookies
			f, err := c.getFlash(r, w)
			if err != nil {
				c.HandleError(w, r, err)
				return
			}
			c.Flash = f

			// handle request
			next.ServeHTTP(w, r)
		})
	}
}

type Config struct {
	*backends.Services
	Router        *ich.Mux
	Assets        mix.Manifest
	Timezone      *time.Location
	Loc           *gotext.Locale
	Env           string
	ErrorHandlers map[int]http.HandlerFunc
	SessionName   string
	SessionStore  sessions.Store
	BaseURL       *url.URL
	FrontendURL   string
	CSRFName      string
}

type Ctx struct {
	Config
	host         string
	scheme       string
	Log          *zap.SugaredLogger
	Loc          *gotext.Locale
	User         *models.Person
	UserRole     string
	OriginalUser *models.Person
	Flash        []flash.Flash
	CSRFToken    string
	CSPNonce     string
	Nav          string
}

func (c *Ctx) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == models.ErrNotFound {
		err = httperror.NotFound
	}

	var httpErr *httperror.Error
	if !errors.As(err, &httpErr) {
		httpErr = httperror.InternalServerError
	}

	if h, ok := c.ErrorHandlers[httpErr.StatusCode]; ok {
		h(w, r)
		return
	}

	c.Log.Error(err)

	http.Error(w, http.StatusText(httpErr.StatusCode), httpErr.StatusCode)
}

func (c *Ctx) PathTo(name string, pairs ...string) *url.URL {
	return c.Router.PathTo(name, pairs...)
}

func (c *Ctx) URLTo(name string, pairs ...string) *url.URL {
	u := c.Router.PathTo(name, pairs...)
	u.Scheme = c.BaseURL.Scheme
	u.Host = c.BaseURL.Host
	return u
}

func (c *Ctx) AssetPath(asset string) string {
	p, err := c.Assets.AssetPath(asset)
	if err != nil {
		panic(err)
	}
	return p
}

func (c *Ctx) PersistFlash(w http.ResponseWriter, f flash.Flash) {
	j, err := json.Marshal(f)
	if err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     FlashCookiePrefix + ulid.Make().String(),
		Value:    base64.URLEncoding.EncodeToString(j),
		Expires:  time.Now().Add(3 * time.Minute),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (c *Ctx) getFlash(r *http.Request, w http.ResponseWriter) ([]flash.Flash, error) {
	var flashes []flash.Flash

	for _, cookie := range r.Cookies() {
		if !strings.HasPrefix(cookie.Name, FlashCookiePrefix) {
			continue
		}

		// delete cookie
		http.SetCookie(w, &http.Cookie{
			Name:     cookie.Name,
			Value:    "",
			Expires:  time.Now(),
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})

		j, err := base64.URLEncoding.DecodeString(cookie.Value)
		if err != nil {
			return nil, err
		}

		f := flash.Flash{}
		if err = json.Unmarshal(j, &f); err != nil {
			return nil, err
		}
		flashes = append(flashes, f)
	}

	return flashes, nil
}

func (c *Ctx) getUserFromSession(r *http.Request, session *sessions.Session, sessionKey string) (*models.Person, error) {
	userID := session.Values[sessionKey]
	if userID == nil {
		return nil, nil
	}

	user, err := c.UserService.GetUser(userID.(string))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Ctx) getUserRoleFromSession(session *sessions.Session) string {
	role := session.Values[UserRoleKey]
	if role == nil {
		return ""
	}
	return role.(string)
}
