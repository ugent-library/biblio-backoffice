package authenticating

import (
	"fmt"
	"net/http"

	"slices"

	"github.com/gorilla/sessions"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/oidc"
)

type AuthHandler struct {
	auth          *oidc.Auth
	usernameClaim string
}

func NewAuthHandler(auth *oidc.Auth, usernameClaim string) *AuthHandler {
	return &AuthHandler{
		auth:          auth,
		usernameClaim: usernameClaim,
	}
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	claims := &oidc.Claims{}
	if err := h.auth.CompleteAuth(w, r, &claims); err != nil {
		c.HandleError(w, r, err)
		return
	}

	username := claims.GetString(h.usernameClaim)

	user, err := c.UserService.GetUserByUsername(username)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	session, err := c.SessionStore.Get(r, c.SessionName)
	if err != nil {
		c.HandleError(w, r, fmt.Errorf("session could not be retrieved: %w", err))
		return
	}

	session.Values[ctx.UserIDKey] = user.ID
	if user.CanCurate() {
		session.Values[ctx.UserRoleKey] = "curator"
	} else {
		session.Values[ctx.UserRoleKey] = "user"
	}

	if err := session.Save(r, w); err != nil {
		c.HandleError(w, r, fmt.Errorf("session could not be saved: %w", err))
		return
	}

	http.Redirect(w, r, c.PathTo("home").String(), http.StatusFound)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	if err := h.auth.BeginAuth(w, r); err != nil {
		c.HandleError(w, r, err)
		return
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	session, err := c.SessionStore.Get(r, c.SessionName)
	if err != nil {
		c.HandleError(w, r, fmt.Errorf("session could not be retrieved: %w", err))
		return
	}

	if err := ClearSession(w, r, session); err != nil {
		c.HandleError(w, r, fmt.Errorf("session could not be saved: %w", err))
		return
	}

	http.Redirect(w, r, c.PathTo("home").String(), http.StatusFound)
}

func UpdateRole(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	role := bind.PathValue(r, "role")

	if !slices.Contains(vocabularies.Map["user_roles"], role) {
		c.HandleError(w, r, httperror.BadRequest.Wrap(fmt.Errorf("%s is not a valid role", role)))
		return
	}

	session, err := c.SessionStore.Get(r, c.SessionName)
	if err != nil {
		c.HandleError(w, r, fmt.Errorf("session could not be retrieved: %w", err))
		return
	}

	session.Values[ctx.UserRoleKey] = role

	if err := session.Save(r, w); err != nil {
		c.HandleError(w, r, fmt.Errorf("session could not be saved: %w", err))
		return
	}

	w.Header().Set("HX-Redirect", c.PathTo("dashboard").String())
}

func ClearSession(w http.ResponseWriter, r *http.Request, session *sessions.Session) error {
	delete(session.Values, ctx.UserIDKey)
	delete(session.Values, ctx.OriginalUserIDKey)
	delete(session.Values, ctx.OriginalUserRoleKey)
	delete(session.Values, ctx.UserRoleKey)
	return session.Save(r, w)
}
