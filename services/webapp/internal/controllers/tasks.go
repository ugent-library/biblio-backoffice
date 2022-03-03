package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

type Tasks struct {
	Context
}

func NewTasks(c Context) *Tasks {
	return &Tasks{c}
}

func (c *Tasks) Status(w http.ResponseWriter, r *http.Request) {
	taskID := mux.Vars(r)["id"]
	resp, err := c.Engine.Temporal.QueryWorkflow(context.Background(), taskID, "", "state")
	if err != nil {
		log.Fatalln("Unable to query workflow", err)
	}
	var result string
	if err := resp.Get(&result); err != nil {
		log.Fatalln("Unable to decode query result", err)
	}

	c.Render.HTML(w, http.StatusOK, "task/_status", c.ViewData(r, struct {
		TaskID     string
		TaskResult string
	}{
		taskID,
		result,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
