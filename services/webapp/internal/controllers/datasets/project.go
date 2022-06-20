package datasets

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
)

type BindActiveSearch struct {
	Search string `form:"search"`
}

type BindProject struct {
	ProjectID string `path:"project_id"`
	Position  int    `path:"position"`
}

type YieldProjects struct {
	Dataset *models.Dataset
}

type YieldChooseProject struct {
	Dataset *models.Dataset
	Hits    []models.Completion
}

type YieldConfirmRemoveProject struct {
	Dataset  *models.Dataset
	Position int
}

func (c *Controller) ChooseProject(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	// Get 20 random projects (no search, init state)
	hits, _ := c.projectSearchService.SuggestProjects("")

	ctx.RenderYield(w, "dataset/add_project", YieldChooseProject{
		Dataset: ctx.Dataset,
		Hits:    hits,
	})
}

func (c *Controller) ActiveSearchProject(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	b := BindActiveSearch{}
	if render.BadRequest(w, bind.RequestForm(r, &b)) {
		return
	}

	// Get 20 fresh results from the search query
	hits, _ := c.projectSearchService.SuggestProjects(b.Search)

	ctx.RenderYield(w, "dataset/add_project_hits", YieldChooseProject{
		Dataset: ctx.Dataset,
		Hits:    hits,
	})
}

func (c *Controller) AddProject(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	b := BindProject{}
	if render.BadRequest(w, bind.RequestPath(r, &b)) {
		return
	}

	project, getProjectErr := c.projectService.GetProject(b.ProjectID)
	if getProjectErr != nil {
		// Handle me
		return
	}

	ctx.Dataset.Project = append(ctx.Dataset.Project, models.DatasetProject{
		ID:   project.ID,
		Name: project.Title,
	})

	updateErr := c.store.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)
	// TODO handle conflict errors

	if !render.InternalServerError(w, updateErr) {
		ctx.RenderYield(w, "dataset/refresh_projects", YieldProjects{
			Dataset: ctx.Dataset,
		})
	}
}

func (c *Controller) ConfirmDeleteProject(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	b := BindProject{}
	if render.BadRequest(w, bind.RequestPath(r, &b)) {
		return
	}

	if _, err := ctx.Dataset.GetProject(b.Position); render.BadRequest(w, err) {
		return
	}

	ctx.RenderYield(w, "dataset/confirm_delete_project", YieldConfirmRemoveProject{
		Dataset:  ctx.Dataset,
		Position: b.Position,
	})
}

func (c *Controller) DeleteProject(w http.ResponseWriter, r *http.Request, ctx EditContext) {
	var b BindAbstract
	if render.BadRequest(w, bind.Request(r, &b)) {
		return
	}

	if render.BadRequest(w, ctx.Dataset.RemoveProject(b.Position)) {
		return
	}

	err := c.store.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)
	// TODO handle conflict errors

	if !render.InternalServerError(w, err) {
		ctx.RenderYield(w, "dataset/refresh_projects", YieldProjects{
			Dataset: ctx.Dataset,
		})
	}
}
