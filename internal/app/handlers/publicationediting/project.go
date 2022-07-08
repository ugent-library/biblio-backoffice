package publicationediting

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

	render.MustRenderLayout(w, "show_modal", "publication/add_project", YieldAddProject{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) SuggestProjects(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindSuggestProjects{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	hits, err := h.ProjectSearchService.SuggestProjects(b.Query)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.MustRenderPartial(w, "publication/suggest_projects", YieldAddProject{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindProject{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	project, err := h.ProjectService.GetProject(b.ProjectID)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	ctx.Publication.Project = append(ctx.Publication.Project, models.PublicationProject{
		ID:   project.ID,
		Name: project.Title,
	})

	// TODO handle validation errors

	err = h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.MustRenderLayout(w, "refresh_modal", "error_dialog", ctx.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.MustRenderView(w, "publication/refresh_projects", YieldProjects{
		Context: ctx,
	})
}

func (h *Handler) ConfirmDeleteProject(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeleteProject{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.MustRenderLayout(w, "show_modal", "publication/confirm_delete_project", YieldDeleteProject{
		Context:  ctx,
		Position: b.Position,
	})
}

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteProject
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	if err := ctx.Publication.RemoveProject(b.Position); err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	// TODO handle validation errors

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.MustRenderLayout(w, "refresh_modal", "error_dialog", ctx.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.MustRenderView(w, "publication/refresh_projects", YieldProjects{
		Context: ctx,
	})
}
