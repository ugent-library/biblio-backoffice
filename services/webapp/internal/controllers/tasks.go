package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/tasks"
	"github.com/unrolled/render"
)

type Tasks struct {
	Base
}

func NewTasks(c Base) *Tasks {
	return &Tasks{c}
}

func (c *Tasks) Status(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	status := c.Services.Tasks.Status(id)

	var msg string

	// TODO move this to translations
	switch {
	case strings.HasPrefix(id, "orcid"):
		switch {
		case status.Waiting():
			msg = "Adding publications to your ORCID works"
		case status.Running():
			msg = fmt.Sprintf("Added %d publications to your ORCID works", status.Progress.Numerator)
		case status.Done():
			msg = "Finished adding publications to your ORCID works"
		case status.Failed():
			msg = "Adding publications to your ORCID works failed"
		}
	}

	c.Render.HTML(w, http.StatusOK, "task/_flash_message", c.ViewData(r, struct {
		ID      string
		Status  tasks.Status
		Message string
	}{
		id,
		status,
		msg,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
