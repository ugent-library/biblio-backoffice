package publicationediting

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

	render.Render(w, "publication/add_department", YieldAddDepartment{
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

	render.Render(w, "publication/suggest_departments", YieldAddDepartment{
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

	department, err := h.OrganizationService.GetOrganization(b.DepartmentID)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	publicationDepartment := models.PublicationDepartment{ID: department.ID}

	for _, d := range department.Tree {
		publicationDepartment.Tree = append(publicationDepartment.Tree, models.PublicationDepartmentRef(d))
	}

	ctx.Publication.Department = append(ctx.Publication.Department, publicationDepartment)

	// TODO handle validation errors

	err = h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Render(w, "error_dialog", ctx.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "publication/refresh_departments", YieldDepartments{
		Context: ctx,
	})
}

func (h *Handler) ConfirmDeleteDepartment(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeleteDepartment{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.Render(w, "publication/confirm_delete_department", YieldDeleteDepartment{
		Context:  ctx,
		Position: b.Position,
	})
}

func (h *Handler) DeleteDepartment(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteDepartment
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	if err := ctx.Publication.RemoveDepartment(b.Position); err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	// TODO handle validation errors

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Render(w, "error_dialog", ctx.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "publication/refresh_departments", YieldDepartments{
		Context: ctx,
	})
}
