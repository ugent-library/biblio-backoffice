package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/leonelquinteros/gotext"
	"github.com/nics/ich"
	"github.com/oklog/ulid/v2"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/ffcontext"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/flash"
	"github.com/unrolled/secure"
	"go.uber.org/zap"
)

// TODO handlers should only have access to a url builder,
// the session and maybe the localizer
type BaseHandler struct {
	Router          *ich.Mux
	Logger          *zap.SugaredLogger
	SessionName     string
	SessionStore    sessions.Store
	UserService     backends.UserService
	Timezone        *time.Location
	Loc             *gotext.Locale
	BaseURL         *url.URL
	FrontendBaseUrl string
}

// also add fields to Yield method
type BaseContext struct {
	CurrentURL      *url.URL
	Flash           []flash.Flash
	Timezone        *time.Location
	Loc             *gotext.Locale
	User            *models.Person
	UserRole        string
	OriginalUser    *models.Person
	CSRFToken       string
	CSRFTag         template.HTML
	CSPNonce        string
	FrontendBaseUrl string
	flagContext     *ffcontext.EvaluationContext
}

func (c BaseContext) Yield(pairs ...any) map[string]any {
	yield := map[string]any{
		"CurrentURL":      c.CurrentURL,
		"Flash":           c.Flash,
		"Timezone":        c.Timezone,
		"Loc":             c.Loc,
		"User":            c.User,
		"UserRole":        c.UserRole,
		"OriginalUser":    c.OriginalUser,
		"CSRFToken":       c.CSRFToken,
		"CSRFTag":         c.CSRFTag,
		"FrontendBaseUrl": c.FrontendBaseUrl,
	}

	n := len(pairs)
	for i := 0; i < n; i += 2 {
		key := pairs[i].(string)
		val := pairs[i+1]
		yield[key] = val
	}

	return yield
}

// Note: keep in line with new ctx.Context
func (c BaseContext) getFlagContext() ffcontext.Context {
	if c.flagContext == nil {
		flagContext := ffcontext.NewEvaluationContext(c.User.Username)
		c.flagContext = &flagContext
	}
	return *c.flagContext
}

func (c BaseContext) FlagCandidateRecords() bool {
	flag, _ := ffclient.BoolVariation("candidate-records", c.getFlagContext(), false)
	return flag
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

	user, err := h.getUserFromSession(session, r, UserIDKey)
	if err != nil {
		return BaseContext{}, fmt.Errorf("could not get user from session: %w", err)
	}

	originalUser, err := h.getUserFromSession(session, r, OriginalUserIDKey)
	if err != nil {
		return BaseContext{}, fmt.Errorf("could not get original user from session: %w", err)
	}

	flash, err := h.getFlashFromCookies(r, w)
	if err != nil {
		return BaseContext{}, fmt.Errorf("could not get flash message from session: %w", err)
	}

	return BaseContext{
		CurrentURL:      r.URL,
		Flash:           flash,
		Timezone:        h.Timezone,
		Loc:             h.Loc,
		User:            user,
		UserRole:        h.getUserRoleFromSession(session),
		OriginalUser:    originalUser,
		CSRFToken:       csrf.Token(r),
		CSRFTag:         csrf.TemplateField(r),
		CSPNonce:        secure.CSPNonce(r.Context()),
		FrontendBaseUrl: h.FrontendBaseUrl,
	}, nil
}

func (h BaseHandler) AddFlash(r *http.Request, w http.ResponseWriter, f flash.Flash) error {
	j, err := json.Marshal(f)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     FlashCookiePrefix + ulid.Make().String(),
		Value:    base64.URLEncoding.EncodeToString(j),
		Expires:  time.Now().Add(3 * time.Minute),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	return nil
}

func (h BaseHandler) getFlashFromCookies(r *http.Request, w http.ResponseWriter) ([]flash.Flash, error) {
	flashes := []flash.Flash{}

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
			continue
		}

		f := flash.Flash{}
		if err = json.Unmarshal(j, &f); err == nil {
			flashes = append(flashes, f)
		}
	}

	return flashes, nil
}

func (h BaseHandler) getUserFromSession(session *sessions.Session, r *http.Request, sessionKey string) (*models.Person, error) {
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

func (h BaseHandler) getUserRoleFromSession(session *sessions.Session) string {
	role := session.Values[UserRoleKey]
	if role == nil {
		return ""
	}
	return role.(string)
}

func (h BaseHandler) PathFor(name string, pairs ...string) *url.URL {
	return h.Router.PathTo(name, pairs...)
}

func (h BaseHandler) URLFor(name string, pairs ...string) *url.URL {
	u := h.Router.PathTo(name, pairs...)
	u.Scheme = h.BaseURL.Scheme
	u.Host = h.BaseURL.Host
	return u
}

func (h BaseHandler) ActionError(w http.ResponseWriter, r *http.Request, ctx BaseContext, msg string, err error, ID string) {
	errID := ulid.Make().String()
	errMsg := fmt.Sprintf("[error: %s] %s", errID, msg)
	h.Logger.Errorw(errMsg, "errors", err, "publication", ID, "user", ctx.User.ID)
	h.ErrorModal(w, r, errID, ctx)
}
