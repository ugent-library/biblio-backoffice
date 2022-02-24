package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/task"
	"github.com/unrolled/render"
)

type Tasks struct {
	Context
}

func NewTasks(c Context) *Tasks {
	return &Tasks{c}
}

func (c *Tasks) Status(w http.ResponseWriter, r *http.Request) {
	t, _ := c.Engine.Tasks.Get(mux.Vars(r)["id"])
	c.Render.HTML(w, http.StatusOK, "task/_status", c.ViewData(r, struct {
		Task task.Task
	}{
		t,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
