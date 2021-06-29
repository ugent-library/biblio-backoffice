package controllers

import (
	"net/http"

	"github.com/unrolled/render"
)

type Publication struct {
	render *render.Render
}

func NewPublication(r *render.Render) *Publication {
	return &Publication{render: r}
}

func (c *Publication) List(w http.ResponseWriter, r *http.Request) {
	c.render.HTML(w, http.StatusOK, "publication/list", "")
}
