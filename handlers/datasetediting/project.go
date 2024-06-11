package datasetediting

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/ctx"
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

func AddProject(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	hits, err := c.ProjectSearchService.SuggestProjects("")
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.ShowModal(datasetviews.AddProject(c, ctx.GetDataset(r), hits)).Render(r.Context(), w)
}

func SuggestProjects(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindSuggestProjects{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	hits, err := c.ProjectSearchService.SuggestProjects(b.Query)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	datasetviews.SuggestProjects(c, ctx.GetDataset(r), hits).Render(r.Context(), w)
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	b := BindProject{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	project, err := c.ProjectService.GetProject(b.ProjectID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	d.AddProject(project)

	// TODO handle validation errors

	err = c.Repo.UpdateDataset(r.Header.Get("If-Match"), d, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.CloseModalAndReplace(datasetviews.ProjectsBodySelector, datasetviews.ProjectsBody(c, d)).Render(r.Context(), w)
}

func ConfirmDeleteProject(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	b := BindDeleteProject{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	if b.SnapshotID != d.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	projectID, _ := url.PathUnescape(b.ProjectID)

	views.ConfirmDeleteDialog(views.ConfirmDeleteDialogArgs{
		Context:    c,
		Question:   "Are you sure you want to remove this project from the dataset?",
		DeleteUrl:  c.PathTo("dataset_delete_project", "id", d.ID, "project_id", projectID),
		SnapshotID: d.SnapshotID,
	}).Render(r.Context(), w)
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	var b BindDeleteProject
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	projectID, _ := url.PathUnescape(b.ProjectID)

	d.RemoveProject(projectID)

	// TODO handle validation errors

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), d, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.CloseModalAndReplace(datasetviews.ProjectsBodySelector, datasetviews.ProjectsBody(c, d)).Render(r.Context(), w)
}
