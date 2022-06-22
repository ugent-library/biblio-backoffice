package datasetediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
)

type BindProjectSuggestions struct {
	Query string `query:"q"`
}
type BindProject struct {
	ProjectID string `form:"project_id"`
}
type BindDeleteProject struct {
	Position int `path:"position"`
}

type YieldProjects struct {
	Context
}
type YieldAddProject struct {
	Context
	Hits []models.Completion
}
type YieldDeleteProject struct {
	Context
	Position int
}

func (h *Handler) AddProject(w http.ResponseWriter, r *http.Request, ctx Context) {
	hits, err := h.ProjectSearchService.SuggestProjects("")
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "dataset/add_project", YieldAddProject{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) ProjectSuggestions(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindProjectSuggestions{}
	if err := bind.RequestQuery(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	hits, err := h.ProjectSearchService.SuggestProjects(b.Query)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "dataset/project_suggestions", YieldAddProject{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindProject{}
	if err := bind.RequestForm(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	project, err := h.ProjectService.GetProject(b.ProjectID)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	ctx.Dataset.Project = append(ctx.Dataset.Project, models.DatasetProject{
		ID:   project.ID,
		Name: project.Title,
	})

	// TODO handle validation errors

	err = h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Render(w, "error_dialog", ctx.T("dataset.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "dataset/refresh_projects", YieldProjects{
		Context: ctx,
	})
}

func (h *Handler) ConfirmDeleteProject(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeleteProject{}
	if err := bind.RequestPath(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	if _, err := ctx.Dataset.GetProject(b.Position); err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "dataset/confirm_delete_project", YieldDeleteProject{
		Context:  ctx,
		Position: b.Position,
	})
}

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindAbstract
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	if err := ctx.Dataset.RemoveProject(b.Position); err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	// TODO handle validation errors

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Render(w, "error_dialog", ctx.T("dataset.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "dataset/refresh_projects", YieldProjects{
		Context: ctx,
	})
}
