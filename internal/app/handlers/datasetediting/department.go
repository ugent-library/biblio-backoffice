package datasetediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
)

type BindDepartmentSuggestions struct {
	Query string `query:"q"`
}

type BindDepartment struct {
	DepartmentID string `form:"department_id"`
}

type BindDeleteDepartment struct {
	Position int `path:"position"`
}

type YieldDepartments struct {
	Context
}

type YieldAddDepartment struct {
	Context
	Hits []models.Completion
}

type YieldDeleteDepartment struct {
	Context
	Position int
}

func (h *Handler) AddDepartment(w http.ResponseWriter, r *http.Request, ctx Context) {
	hits, err := h.OrganizationSearchService.SuggestOrganizations("")
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "dataset/add_department", YieldAddDepartment{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) DepartmentSuggestions(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDepartmentSuggestions{}
	if err := bind.RequestQuery(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	hits, err := h.OrganizationSearchService.SuggestOrganizations(b.Query)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "dataset/department_suggestions", YieldAddDepartment{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) CreateDepartment(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDepartment{}
	if err := bind.RequestForm(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	department, err := h.OrganizationService.GetOrganization(b.DepartmentID)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	datasetDepartment := models.DatasetDepartment{ID: department.ID}

	for _, d := range department.Tree {
		datasetDepartment.Tree = append(datasetDepartment.Tree, models.DatasetDepartmentRef(d))
	}

	ctx.Dataset.Department = append(ctx.Dataset.Department, datasetDepartment)

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

	render.Render(w, "dataset/refresh_departments", YieldDepartments{
		Context: ctx,
	})
}

func (h *Handler) ConfirmDeleteDepartment(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeleteDepartment{}
	if err := bind.RequestPath(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	if _, err := ctx.Dataset.GetDepartment(b.Position); err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "dataset/confirm_delete_department", YieldDeleteDepartment{
		Context:  ctx,
		Position: b.Position,
	})
}

func (h *Handler) DeleteDepartment(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindAbstract
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	if err := ctx.Dataset.RemoveDepartment(b.Position); err != nil {
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

	render.Render(w, "dataset/refresh_departments", YieldDepartments{
		Context: ctx,
	})
}
