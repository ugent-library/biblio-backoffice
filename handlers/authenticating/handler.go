package authenticating

import (
	"fmt"
	"net/http"

	"slices"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
	"github.com/ugent-library/bind"
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
		c.Log.Errorw("authentication: OIDC client could not complete exchange:", "errors", err)
		c.HandleError(w, r, err)
		return
	}

	username := claims.GetString(h.usernameClaim)

	user, err := c.UserService.GetUserByUsername(username)
	if err != nil {
		c.Log.Warnw("authentication: No user with that name could be found:", "errors", err, "user", username)
		c.HandleError(w, r, err)
		return
	}

	session, err := c.SessionStore.Get(r, c.SessionName)
	if err != nil {
		c.Log.Errorw("authentication: session could not be retrieved:", "errors", err)
		c.HandleError(w, r, err)
		return
	}

	session.Values[handlers.UserIDKey] = user.ID
	if _, ok := session.Values[handlers.UserRoleKey]; !ok {
		session.Values[handlers.UserRoleKey] = "user"
	}

	if err := session.Save(r, w); err != nil {
		c.Log.Errorw("authentication: session could not be saved:", "errors", err)
		c.HandleError(w, r, err)
		return
	}

	http.Redirect(w, r, c.PathTo("home").String(), http.StatusFound)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	if err := h.auth.BeginAuth(w, r); err != nil {
		c.Log.Errorw("authentication: OIDC client could not begin exchange:", "errors", err)
		c.HandleError(w, r, err)
		return
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	session, err := c.SessionStore.Get(r, c.SessionName)
	if err != nil {
		c.Log.Errorw("authentication: session could not be retrieved:", "errors", err)
		c.HandleError(w, r, err)
		return
	}

	// only remember user role
	delete(session.Values, handlers.UserIDKey)
	delete(session.Values, handlers.OriginalUserIDKey)
	delete(session.Values, handlers.OriginalUserRoleKey)
	if err := session.Save(r, w); err != nil {
		c.Log.Errorw("authentication: session could not be saved:", "errors", err)
		c.HandleError(w, r, err)
		return
	}

	http.Redirect(w, r, c.PathTo("home").String(), http.StatusFound)
}

func UpdateRole(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	role := bind.PathValue(r, "role")

	if !slices.Contains(vocabularies.Map["user_roles"], role) {
		render.BadRequest(w, r, fmt.Errorf("%s is not a valid role", role))
		return
	}

	session, err := c.SessionStore.Get(r, c.SessionName)
	if err != nil {
		c.Log.Errorw("authentication: session could not be retrieved:", "errors", err)
		c.HandleError(w, r, err)
		return
	}

	session.Values[handlers.UserRoleKey] = role

	if err := session.Save(r, w); err != nil {
		c.Log.Errorw("authentication: session could not be saved:", "errors", err)
		c.HandleError(w, r, err)
		return
	}

	w.Header().Set("HX-Redirect", c.PathTo("dashboard").String())
}
