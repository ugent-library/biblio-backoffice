package handlers

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
)

type Base struct {
	SessionName  string
	SessionStore sessions.Store
	UserService  backends.UserService
	Localizer    *locale.Localizer
}

func (b Base) Wrap(fn func(http.ResponseWriter, *http.Request, BaseContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := b.NewContext(r)
		if err != nil {
			render.InternalServerError(w, r, err)
			return
		}
		fn(w, r, ctx)
	}
}

func (b Base) NewContext(r *http.Request) (BaseContext, error) {
	user, err := b.getUserFromSession(r, "user_id")
	if err != nil {
		return BaseContext{}, err
	}

	originalUser, err := b.getUserFromSession(r, "original_user_id")
	if err != nil {
		return BaseContext{}, err
	}

	return BaseContext{
		Locale:       b.Localizer.GetLocale(r.Header.Get("Accept-Language")),
		User:         user,
		OriginalUser: originalUser,
		CSRFToken:    csrf.Token(r),
		CSRFTag:      csrf.TemplateField(r),
	}, nil
}

func (b Base) getUserFromSession(r *http.Request, sessionKey string) (*models.User, error) {
	session, err := b.SessionStore.Get(r, b.SessionName)
	if err != nil {
		return nil, err
	}
	userID := session.Values[sessionKey]
	if userID == nil {
		return nil, nil
	}

	user, err := b.UserService.GetUser(userID.(string))
	if err != nil {
		return nil, err
	}

	return user, nil
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
