package tasks

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/tasks"
)

type Handler struct {
	handlers.BaseHandler
	Tasks *tasks.Hub
}

type Context struct {
	handlers.BaseContext
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil {
			handlers.Unauthorized(w, r)
			return
		}

		fn(w, r, Context{
			BaseContext: ctx,
		})
	})
}

type BindStatus struct {
	ID string `path:"id"`
}

type YieldStatus struct {
	ID      string
	Status  tasks.Status
	Message string
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindStatus{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("tasks: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		handlers.BadRequest(w, r, err)
		return
	}

	status := h.Tasks.Status(b.ID)

	// TODO move this to translations
	var msg string
	switch {
	case strings.HasPrefix(b.ID, "orcid"):
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

	render.Partial(w, "task/status", YieldStatus{
		ID:      b.ID,
		Status:  status,
		Message: msg,
	})
}
