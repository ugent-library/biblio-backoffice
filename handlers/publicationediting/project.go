package publicationediting

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/ctx"
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

func CreateProject(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindProject{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("create publication project: could not bind request arguments:", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	project, err := c.ProjectService.GetProject(b.ProjectID)
	if err != nil {
		c.Log.Errorw("create publication project: could not get project:", "errors", err, "publication", p.ID, "project", b.ProjectID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}
	p.AddProject(project)

	// TODO handle validation errors

	err = c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("create publication project: Could not save the publication:", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(publicationviews.ProjectsBodySelector, publicationviews.ProjectsBody(c, p)).Render(r.Context(), w)
}

func ConfirmDeleteProject(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindDeleteProject{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete publication project: could not bind request arguments:", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	if b.SnapshotID != p.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	projectID, _ := url.PathUnescape(b.ProjectID)

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to remove this project from the publication?",
		DeleteUrl:  c.PathTo("publication_delete_project", "id", p.ID, "project_id", projectID),
		SnapshotID: p.SnapshotID,
	}).Render(r.Context(), w)
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	var b BindDeleteProject
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("delete publication project: could not bind request arguments:", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	projectID, _ := url.PathUnescape(b.ProjectID)

	p.RemoveProject(projectID)

	// TODO handle validation errors

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("delete publication project: Could not save the publication:", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(publicationviews.ProjectsBodySelector, publicationviews.ProjectsBody(c, p)).Render(r.Context(), w)
}
