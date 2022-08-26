package datasetediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
)

type BindSuggestProjects struct {
	Query string `query:"q"`
}
type BindProject struct {
	ProjectID string `form:"project_id"`
}
type BindDeleteProject struct {
	ProjectID string `path:"project_id"`
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
	ProjectID string
}

func (h *Handler) AddProject(w http.ResponseWriter, r *http.Request, ctx Context) {
	hits, err := h.ProjectSearchService.SuggestProjects("")
	if err != nil {
		h.Logger.Errorw("add dataset project: could not suggest projects:", "error", err, "request", r)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "show_modal", "dataset/add_project", YieldAddProject{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) SuggestProjects(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindSuggestProjects{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("suggest dataset project: could not bind request arguments:", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	hits, err := h.ProjectSearchService.SuggestProjects(b.Query)
	if err != nil {
		h.Logger.Errorw("suggest dataset project: could not suggest projects:", "error", err, "request", r)
		render.InternalServerError(w, r, err)
		return
	}

	render.Partial(w, "dataset/suggest_projects", YieldAddProject{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindProject{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("create dataset project: could not bind request arguments:", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	project, err := h.ProjectService.GetProject(b.ProjectID)
	if err != nil {
		h.Logger.Errorw("create dataset project: could not get project:", "error", err, "dataset", ctx.Dataset.ID, "identifier", b.ProjectID)
		render.BadRequest(w, r, err)
		return
	}

	/*
		Note: AddProject removes potential existing project
		and adds a new one at the end
	*/
	ctx.Dataset.AddProject(&models.DatasetProject{
		ID:   project.ID,
		Name: project.Title,
	})

	// TODO handle validation errors

	err = h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		h.Logger.Warnf("create dataset project: snapstore detected a conflicting dataset:", "errors", errors.As(err, &conflict), "identifier", ctx.Dataset.ID)
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("dataset.conflict_error"))
		return
	}

	if err != nil {
		h.Logger.Errorf("create dataset project: Could not save the dataset:", "error", err, "identifier", ctx.Dataset.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_projects", YieldProjects{
		Context: ctx,
	})
}

func (h *Handler) ConfirmDeleteProject(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeleteProject{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("confirm delete dataset project: could not bind request arguments:", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	render.Layout(w, "show_modal", "dataset/confirm_delete_project", YieldDeleteProject{
		Context:   ctx,
		ProjectID: b.ProjectID,
	})
}

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteProject
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete dataset project: could not bind request arguments:", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	/*
		ignore possibility that project is already removed:
		conflict resolving will solve this anyway
	*/
	ctx.Dataset.RemoveProject(b.ProjectID)

	// TODO handle validation errors

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		h.Logger.Warnf("delete dataset project: snapstore detected a conflicting dataset:", "errors", errors.As(err, &conflict), "identifier", ctx.Dataset.ID)
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("dataset.conflict_error"))
		return
	}

	if err != nil {
		h.Logger.Errorf("delete dataset project: Could not save the dataset:", "error", err, "identifier", ctx.Dataset.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_projects", YieldProjects{
		Context: ctx,
	})
}
