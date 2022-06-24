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
		render.BadRequest(w, r, errors.New("already impersonating"))
	}

	if !ctx.User.CanImpersonateUser() {
		render.Unauthorized(w, r)
		return
	}

	render.Render(w, "impersonation/add", ctx)
}

func (h *Handler) CreateImpersonation(w http.ResponseWriter, r *http.Request, ctx Context) {
	if ctx.OriginalUser != nil {
		render.BadRequest(w, r, errors.New("already impersonating"))
	}

	if !ctx.User.CanImpersonateUser() {
		render.Unauthorized(w, r)
		return
	}

	b := BindImpersonation{}
	if err := bind.Request(r, &b); err != nil {
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
		render.InternalServerError(w, r, err)
		return
	}

	session.Values[handlers.OriginalUserSessionKey] = ctx.User.ID
	session.Values[handlers.UserSessionKey] = user.ID

	if err = session.Save(r, w); err != nil {
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
		render.InternalServerError(w, r, err)
		return
	}

	if origUserID := session.Values[handlers.OriginalUserSessionKey]; origUserID != nil {
		delete(session.Values, handlers.OriginalUserSessionKey)
		session.Values[handlers.UserSessionKey] = origUserID

		if err = session.Save(r, w); err != nil {
			render.InternalServerError(w, r, err)
			return
		}
	}

	http.Redirect(w, r, h.PathFor("home").String(), http.StatusFound)
}
