package datasetediting

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	datasetviews "github.com/ugent-library/biblio-backoffice/views/dataset"
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
		c.Log.Errorw("add dataset project: could not suggest projects:", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	datasetviews.AddProject(c, ctx.GetDataset(r), hits).Render(r.Context(), w)

}

func SuggestProjects(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindSuggestProjects{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("suggest dataset project: could not bind request arguments:", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	hits, err := c.ProjectSearchService.SuggestProjects(b.Query)
	if err != nil {
		c.Log.Errorw("suggest dataset project: could not suggest projects:", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	datasetviews.SuggestProjects(c, ctx.GetDataset(r), hits).Render(r.Context(), w)
}

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindProject{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("create dataset project: could not bind request arguments:", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	project, err := h.ProjectService.GetProject(b.ProjectID)
	if err != nil {
		h.Logger.Errorw("create dataset project: could not get project:", "errors", err, "dataset", ctx.Dataset.ID, "project", b.ProjectID, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	ctx.Dataset.AddProject(project)

	// TODO handle validation errors

	err = h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(ctx.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		h.Logger.Errorf("create dataset project: Could not save the dataset:", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_projects", YieldProjects{
		Context: ctx,
	})
}

func ConfirmDeleteProject(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindDeleteProject{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete dataset project: could not bind request arguments:", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	if b.SnapshotID != dataset.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	projectID, _ := url.PathUnescape(b.ProjectID)

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to remove this project from the dataset?",
		DeleteUrl:  c.PathTo("dataset_delete_project", "id", dataset.ID, "project_id", projectID),
		SnapshotID: dataset.SnapshotID,
	}).Render(r.Context(), w)
}

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteProject
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete dataset project: could not bind request arguments:", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	projectID, _ := url.PathUnescape(b.ProjectID)

	ctx.Dataset.RemoveProject(projectID)

	// TODO handle validation errors

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(ctx.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		h.Logger.Errorf("delete dataset project: Could not save the dataset:", "error", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_projects", YieldProjects{
		Context: ctx,
	})
}
