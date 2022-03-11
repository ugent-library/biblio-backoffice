package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
	"go.temporal.io/sdk/client"
)

type Tasks struct {
	Context
}

func NewTasks(c Context) *Tasks {
	return &Tasks{c}
}

func (c *Tasks) Status(w http.ResponseWriter, r *http.Request) {
	taskID := mux.Vars(r)["id"]

	dw, err := c.Engine.Temporal.DescribeWorkflowExecution(context.Background(), taskID, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	taskState := models.TaskState{}

	// TODO using the constants in the enum package gives strange errors
	switch dw.GetWorkflowExecutionInfo().Status {
	case 1:
		taskState.Status = models.Running
	case 2:
		taskState.Status = models.Done
	case 3, 4, 5, 7:
		taskState.Status = models.Failed
	default:
		taskState.Status = models.Waiting
	}

	for _, a := range dw.GetPendingActivities() {
		client.NewValue(a.HeartbeatDetails).Get(&taskState.Progress)
	}

	// TODO move this to translations
	switch {
	case strings.HasPrefix(taskID, "orcid"):
		switch taskState.Status {
		case models.Waiting:
			taskState.Message = "Adding publications to your ORCID works"
		case models.Running:
			taskState.Message = fmt.Sprintf("Added %d publications to your ORCID works", taskState.Numerator)
		case models.Done:
			taskState.Message = "Finished adding publications to your ORCID works"
		case models.Failed:
			taskState.Message = "Adding publications to your ORCID works failed"
		}
	}

	c.Render.HTML(w, http.StatusOK, "task/_flash_message", c.ViewData(r, struct {
		TaskID    string
		TaskState models.TaskState
	}{
		taskID,
		taskState,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
