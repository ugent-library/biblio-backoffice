package publicationediting

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	publicationviews "github.com/ugent-library/biblio-backoffice/views/publication"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
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

func AddProject(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	hits, err := c.ProjectSearchService.SuggestProjects("")
	if err != nil {
		c.Log.Errorw("add publication project: could not suggest projects:", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	publicationviews.AddProject(c, ctx.GetPublication(r), hits).Render(r.Context(), w)
}

func SuggestProjects(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindSuggestProjects{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("suggest publication project: could not bind request arguments:", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	hits, err := c.ProjectSearchService.SuggestProjects(b.Query)
	if err != nil {
		c.Log.Errorw("suggest publication project: could not suggest projects:", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	publicationviews.SuggestProjects(c, ctx.GetPublication(r), hits).Render(r.Context(), w)
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

func ConfirmDeleteProject(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	b := BindDeleteProject{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete publication project: could not bind request arguments:", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	if b.SnapshotID != publication.SnapshotID {
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: c.Loc.Get("publication.conflict_error_reload"),
		})
		return
	}

	projectID, _ := url.PathUnescape(b.ProjectID)

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to remove this project from the publication?",
		DeleteUrl:  c.PathTo("publication_delete_project", "id", publication.ID, "project_id", projectID),
		SnapshotID: publication.SnapshotID,
	}).Render(r.Context(), w)
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
