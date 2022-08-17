package datasetediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
)

type BindSuggestDepartments struct {
	Query string `query:"q"`
}

type BindDepartment struct {
	DepartmentID string `form:"department_id"`
}

type BindDeleteDepartment struct {
	DepartmentID string `path:"department_id"`
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
	DepartmentID string
}

func (h *Handler) AddDepartment(w http.ResponseWriter, r *http.Request, ctx Context) {
	hits, err := h.OrganizationSearchService.SuggestOrganizations("")
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "show_modal", "dataset/add_department", YieldAddDepartment{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) SuggestDepartments(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindSuggestDepartments{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	hits, err := h.OrganizationSearchService.SuggestOrganizations(b.Query)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Partial(w, "dataset/suggest_departments", YieldAddDepartment{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) CreateDepartment(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDepartment{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	org, err := h.OrganizationService.GetOrganization(b.DepartmentID)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	ctx.Dataset.AddDepartmentByOrg(org)

	// TODO handle validation errors

	err = h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("dataset.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_departments", YieldDepartments{
		Context: ctx,
	})
}

func (h *Handler) ConfirmDeleteDepartment(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeleteDepartment{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.Layout(w, "show_modal", "dataset/confirm_delete_department", YieldDeleteDepartment{
		Context:      ctx,
		DepartmentID: b.DepartmentID,
	})
}

func (h *Handler) DeleteDepartment(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteDepartment
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	/*
		ignore possibility that department is already removed:
		conflict resolving will solve this anyway
	*/
	ctx.Dataset.RemoveDepartment(b.DepartmentID)

	// TODO handle validation errors

	err := h.Repository.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "show_modal", "error_dialog", ctx.Locale.T("dataset.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_departments", YieldDepartments{
		Context: ctx,
	})
}
