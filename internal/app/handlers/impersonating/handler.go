package impersonating

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/render"
)

type Handler struct {
	handlers.BaseHandler
}

type Context struct {
	handlers.BaseContext
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil {
			render.Unauthorized(w, r)
			return
		}

		fn(w, r, Context{
			BaseContext: ctx,
		})
	})
}

type BindImpersonation struct {
	Username string `form:"username"`
}

func (h *Handler) AddImpersonation(w http.ResponseWriter, r *http.Request, ctx Context) {
	if ctx.OriginalUser != nil {
		h.Logger.Warn("add impersonation: already impersonating", "user", ctx.OriginalUser.ID)
		render.BadRequest(w, r, errors.New("already impersonating"))
	}

	if !ctx.User.CanImpersonateUser() {
		h.Logger.Warn("add impersonation: user does not have permission to impersonate", "user", ctx.User.ID)
		render.Unauthorized(w, r)
		return
	}

	render.Layout(w, "show_modal", "impersonation/add", ctx)
}

func (h *Handler) CreateImpersonation(w http.ResponseWriter, r *http.Request, ctx Context) {
	if ctx.OriginalUser != nil {
		h.Logger.Warn("create impersonation: already impersonating", "user", ctx.OriginalUser.ID)
		render.BadRequest(w, r, errors.New("already impersonating"))
	}

	if !ctx.User.CanImpersonateUser() {
		h.Logger.Warn("create impersonation: user does not have permission to impersonate", "user", ctx.User.ID)
		render.Unauthorized(w, r)
		return
	}

	b := BindImpersonation{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("create impersonation: could not bind request arguments", "errors", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	user, err := h.UserService.GetUserByUsername(b.Username)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	// TODO handle user not found

	session, err := h.SessionStore.Get(r, h.SessionName)
	if err != nil {
		h.Logger.Errorw("create impersonation: session could not be retrieved:", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	session.Values[handlers.OriginalUserIDKey] = ctx.User.ID
	session.Values[handlers.OriginalUserRoleKey] = ctx.UserRole
	session.Values[handlers.UserIDKey] = user.ID
	session.Values[handlers.UserRoleKey] = "user"

	if err = session.Save(r, w); err != nil {
		h.Logger.Errorw("create impersonation: session could not be saved:", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, h.PathFor("home").String(), http.StatusFound)
}

func (h *Handler) DeleteImpersonation(w http.ResponseWriter, r *http.Request, ctx Context) {
	if ctx.OriginalUser == nil {
		render.BadRequest(w, r, errors.New("no impersonation"))
	}

	session, err := h.SessionStore.Get(r, h.SessionName)
	if err != nil {
		h.Logger.Errorw("delete impersonation: session could not be retrieved:", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	session.Values[handlers.UserIDKey] = session.Values[handlers.OriginalUserIDKey]
	session.Values[handlers.UserRoleKey] = session.Values[handlers.OriginalUserRoleKey]
	delete(session.Values, handlers.OriginalUserIDKey)
	delete(session.Values, handlers.OriginalUserRoleKey)

	if err = session.Save(r, w); err != nil {
		h.Logger.Errorw("delete impersonation: session could not be saved:", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, h.PathFor("home").String(), http.StatusFound)
}
