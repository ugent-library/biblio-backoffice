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
		h.Logger.Warnw("create impersonation: could not bind request arguments", "error", err, "request", r)
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
		h.Logger.Errorw("create impersonation: session could not be retrieved:", "error", err)
		render.InternalServerError(w, r, err)
		return
	}

	session.Values[handlers.OriginalUserSessionKey] = ctx.User.ID
	session.Values[handlers.UserSessionKey] = user.ID

	if err = session.Save(r, w); err != nil {
		h.Logger.Errorw("create impersonation: session could not be saved:", "error", err)
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
		h.Logger.Errorw("delete impersonation: session could not be retrieved:", "error", err)
		render.InternalServerError(w, r, err)
		return
	}

	if origUserID := session.Values[handlers.OriginalUserSessionKey]; origUserID != nil {
		delete(session.Values, handlers.OriginalUserSessionKey)
		session.Values[handlers.UserSessionKey] = origUserID

		if err = session.Save(r, w); err != nil {
			h.Logger.Errorw("delete impersonation: session could not be saved:", "error", err)
			render.InternalServerError(w, r, err)
			return
		}
	}

	http.Redirect(w, r, h.PathFor("home").String(), http.StatusFound)
}
