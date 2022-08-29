package authenticating

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/go-oidc/oidc"
)

type Handler struct {
	handlers.BaseHandler
	OIDCClient *oidc.Client
}

type Context struct {
	handlers.BaseContext
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		fn(w, r, Context{
			BaseContext: ctx,
		})
	})
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request, ctx Context) {
	code := r.URL.Query().Get("code")
	claims := &oidc.Claims{}
	if err := h.OIDCClient.Exchange(code, claims); err != nil {
		h.Logger.Errorw("authentication: OIDC client could not complete exchange:", "errors", err)
		render.InternalServerError(w, r, err)
		return
	}

	user, err := h.UserService.GetUserByUsername(claims.PreferredUsername)
	if err != nil {
		h.Logger.Warnw("authentication: No user with that name could be found:", "errors", err, "user", claims.PreferredUsername)
		render.NotFoundError(w, r, err)
		return
	}

	session, err := h.SessionStore.Get(r, h.SessionName)
	if err != nil {
		h.Logger.Errorw("authentication: session could not be retrieved:", "errors", err)
		render.InternalServerError(w, r, err)
		return
	}

	session.Values[handlers.UserSessionKey] = user.ID
	if err := session.Save(r, w); err != nil {
		h.Logger.Errorw("authentication: session could not be saved:", "errors", err)
		render.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, h.PathFor("home").String(), http.StatusFound)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request, ctx Context) {
	http.Redirect(w, r, h.OIDCClient.AuthorizationURL(), http.StatusFound)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request, ctx Context) {
	session, err := h.SessionStore.Get(r, h.SessionName)
	if err != nil {
		h.Logger.Errorw("authentication: session could not be retrieved:", "errors", err)
		render.InternalServerError(w, r, err)
		return
	}

	delete(session.Values, handlers.UserSessionKey)
	delete(session.Values, handlers.OriginalUserSessionKey)
	if err := session.Save(r, w); err != nil {
		h.Logger.Errorw("authentication: session could not be saved:", "errors", err)
		render.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, h.PathFor("home").String(), http.StatusFound)
}
