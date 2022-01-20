package controllers

import (
	"net/http"
)

type Home struct {
	Context
}

func NewHome(c Context) *Auth {
	return &Auth{c}
}

func (c *Auth) Home(w http.ResponseWriter, r *http.Request) {
	c.Render.HTML(w, http.StatusOK, "home/home", c.ViewData(r, struct {
		PageTitle string
	}{
		"Biblio",
	}))
}
