package controllers

import (
	"net/http"
)

type Home struct {
	Base
}

func NewHome(base Base) *Home {
	return &Home{base}
}

func (c *Home) Home(w http.ResponseWriter, r *http.Request) {
	c.Render.HTML(w, http.StatusOK, "home/home", c.ViewData(r, struct {
		PageTitle string
	}{
		"Biblio",
	}))
}
