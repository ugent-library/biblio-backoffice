package authenticating

import (
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/internal/app/handlers"
	"github.com/ugent-library/biblio-backoffice/internal/bind"
	"github.com/ugent-library/biblio-backoffice/internal/render"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
	"github.com/ugent-library/biblio-backoffice/internal/vocabularies"
	"github.com/ugent-library/oidc"
)

type Handler struct {
	handlers.BaseHandler
	OIDCAuth *oidc.Auth
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
	claims := &oidc.Claims{}
	if err := h.OIDCAuth.CompleteAuth(w, r, &claims); err != nil {
		h.Logger.Errorw("authentication: OIDC client could not complete exchange:", "errors", err)
		h.InternalServerError(w, r, ctx.BaseContext)
		return
	}

	user, err := h.UserService.GetUserByUsername(claims.PreferredUsername)
	if err != nil {
		h.Logger.Warnw("authentication: No user with that name could be found:", "errors", err, "user", claims.PreferredUsername)
		h.InternalServerError(w, r, ctx.BaseContext)
		return
	}

	session, err := h.SessionStore.Get(r, h.SessionName)
	if err != nil {
		h.Logger.Errorw("authentication: session could not be retrieved:", "errors", err)
		h.InternalServerError(w, r, ctx.BaseContext)
		return
	}

	session.Values[handlers.UserIDKey] = user.ID
	if _, ok := session.Values[handlers.UserRoleKey]; !ok {
		session.Values[handlers.UserRoleKey] = "user"
	}

	if err := session.Save(r, w); err != nil {
		h.Logger.Errorw("authentication: session could not be saved:", "errors", err)
		h.InternalServerError(w, r, ctx.BaseContext)
		return
	}

	http.Redirect(w, r, h.PathFor("home").String(), http.StatusFound)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request, ctx Context) {
	if err := h.OIDCAuth.BeginAuth(w, r); err != nil {
		h.Logger.Errorw("authentication: OIDC client could not begin exchange:", "errors", err)
		h.InternalServerError(w, r, ctx.BaseContext)
		return
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request, ctx Context) {
	session, err := h.SessionStore.Get(r, h.SessionName)
	if err != nil {
		h.Logger.Errorw("authentication: session could not be retrieved:", "errors", err)
		render.InternalServerError(w, r, err)
		return
	}

	// only remember user role
	delete(session.Values, handlers.UserIDKey)
	delete(session.Values, handlers.OriginalUserIDKey)
	delete(session.Values, handlers.OriginalUserRoleKey)
	if err := session.Save(r, w); err != nil {
		h.Logger.Errorw("authentication: session could not be saved:", "errors", err)
		h.InternalServerError(w, r, ctx.BaseContext)
		return
	}

	http.Redirect(w, r, h.PathFor("home").String(), http.StatusFound)
}

func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request, ctx Context) {
	if ctx.User == nil || !ctx.User.CanCurate() {
		render.Unauthorized(w, r)
		return
	}

	role := bind.PathValues(r).Get("role")

	if !validation.InArray(vocabularies.Map["user_roles"], role) {
		render.BadRequest(w, r, fmt.Errorf("%s is not a valid role", role))
		return
	}

	session, err := h.SessionStore.Get(r, h.SessionName)
	if err != nil {
		h.Logger.Errorw("authentication: session could not be retrieved:", "errors", err)
		h.InternalServerError(w, r, ctx.BaseContext)
		return
	}

	session.Values[handlers.UserRoleKey] = role

	if err := session.Save(r, w); err != nil {
		h.Logger.Errorw("authentication: session could not be saved:", "errors", err)
		h.InternalServerError(w, r, ctx.BaseContext)
		return
	}

	w.Header().Set("HX-Redirect", h.PathFor("publications").String())
}
