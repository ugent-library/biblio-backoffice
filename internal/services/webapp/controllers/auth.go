package controllers

import (
	"log"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/go-oidc/oidc"
)

type Auth struct {
	Base
	oidcClient  *oidc.Client
	userService backends.UserService
}

func NewAuth(base Base, oidcClient *oidc.Client, userService backends.UserService) *Auth {
	return &Auth{
		Base:        base,
		oidcClient:  oidcClient,
		userService: userService,
	}
}

func (c *Auth) Callback(w http.ResponseWriter, r *http.Request) {
	claims := &oidc.Claims{}
	err := c.oidcClient.Exchange(r.URL.Query().Get("code"), claims)
	if err != nil {
		log.Printf("oidc error: %s", err)
		// TODO
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := c.Session(r)
	if err != nil {
		log.Printf("session error: %s", err)
		// TODO
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := c.userService.GetUserByUsername(claims.PreferredUsername)
	if err != nil {
		log.Printf("get user error: %s", err)
		// TODO
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user_id"] = user.ID
	if err = session.Save(r, w); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURL, _ := c.Router.Get("home").URLPath()
	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}

func (c *Auth) Login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, c.oidcClient.AuthorizationURL(), http.StatusFound)
}

func (c *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := c.Session(r)
	if err != nil {
		// TODO
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	delete(session.Values, "user_id")
	delete(session.Values, "original_user_id")
	session.Save(r, w)

	redirectURL, _ := c.Router.Get("home").URLPath()
	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}
