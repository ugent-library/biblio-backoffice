package controllers

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/views"
)

type Home struct {
	Context
}

func NewHome(c Context) *Auth {
	return &Auth{c}
}

func (c *Auth) Home(w http.ResponseWriter, r *http.Request) {
	c.Render.HTML(w, http.StatusOK, "home/home", views.NewData(c.Render, r, struct {
		PageTitle string
	}{
		"Biblio",
	}))
}
