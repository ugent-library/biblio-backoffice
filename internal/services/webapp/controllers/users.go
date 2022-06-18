package controllers

import (
	"log"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/context"
	"github.com/unrolled/render"
)

type Users struct {
	Base
	userService backends.UserService
}

func NewUsers(base Base, userService backends.UserService) *Users {
	return &Users{
		Base:        base,
		userService: userService,
	}
}

func (c *Users) ImpersonateChoose(w http.ResponseWriter, r *http.Request) {
	c.Render.HTML(w, http.StatusOK, "user/_impersonate_choose",
		c.ViewData(r, nil),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *Users) Impersonate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")

	user, err := c.userService.GetUserByUsername(username)
	if err != nil {
		log.Printf("impersonate get user error: %s", err)
		// TODO
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !context.GetUser(r.Context()).CanImpersonateUser() {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	session, _ := c.Session(r)
	session.Values["original_user_id"] = context.GetUser(r.Context()).ID
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	redirectURL, _ := c.Router.Get("publications").URLPath()
	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}

func (c *Users) ImpersonateRemove(w http.ResponseWriter, r *http.Request) {
	session, _ := c.Session(r)
	if origUserID := session.Values["original_user_id"]; origUserID != nil {
		delete(session.Values, "original_user_id")
		session.Values["user_id"] = origUserID
		session.Save(r, w)
	}

	redirectURL, _ := c.Router.Get("publications").URLPath()
	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}
