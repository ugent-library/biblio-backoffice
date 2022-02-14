package controllers

import (
	"log"
	"net/http"

	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/unrolled/render"
)

type Users struct {
	Context
}

func NewUsers(c Context) *Users {
	return &Users{c}
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

	user, err := c.Engine.GetUserByUsername(username)
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
