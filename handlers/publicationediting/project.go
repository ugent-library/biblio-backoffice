package publicationediting

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/bind"
)

type BindSuggestProjects struct {
	Query string `query:"q"`
}
type BindProject struct {
	ProjectID string `form:"project_id"`
}
type BindDeleteProject struct {
	ProjectID  string `path:"project_id"`
	SnapshotID string `path:"snapshot_id"`
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
		h.Logger.Errorw("add publication project: could not suggest projects:", "errors", err, "request", r, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "show_modal", "publication/add_project", YieldAddProject{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) SuggestProjects(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindSuggestProjects{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("suggest publication project: could not bind request arguments:", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	hits, err := h.ProjectSearchService.SuggestProjects(b.Query)
	if err != nil {
		h.Logger.Errorw("suggest publication project: could not suggest projects:", "errors", err, "request", r, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Partial(w, "publication/suggest_projects", YieldAddProject{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindProject{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("create publication project: could not bind request arguments:", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	project, err := h.ProjectService.GetProject(b.ProjectID)
	if err != nil {
		h.Logger.Errorw("create publication project: could not get project:", "errors", err, "publication", ctx.Publication.ID, "project", b.ProjectID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}
	ctx.Publication.AddProject(project)

	// TODO handle validation errors

	err = h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("publication.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("create publication project: Could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_projects", YieldProjects{
		Context: ctx,
	})
}

func (h *Handler) ConfirmDeleteProject(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeleteProject{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("confirm delete publication project: could not bind request arguments:", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	if b.SnapshotID != ctx.Publication.SnapshotID {
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("publication.conflict_error_reload"),
		})
		return
	}

	projectID, _ := url.PathUnescape(b.ProjectID)

	render.Layout(w, "show_modal", "publication/confirm_delete_project", YieldDeleteProject{
		Context:   ctx,
		ProjectID: projectID,
	})
}

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteProject
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete publication project: could not bind request arguments:", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	projectID, _ := url.PathUnescape(b.ProjectID)

	ctx.Publication.RemoveProject(projectID)

	// TODO handle validation errors

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("publication.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("delete publication project: Could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_projects", YieldProjects{
		Context: ctx,
	})
}
