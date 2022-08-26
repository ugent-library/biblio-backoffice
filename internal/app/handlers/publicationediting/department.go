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
		h.Logger.Errorw("add publication department: could not suggest organization", "error", err)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "show_modal", "publication/add_department", YieldAddDepartment{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) SuggestDepartments(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindSuggestDepartments{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("suggest publication departments could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	hits, err := h.OrganizationSearchService.SuggestOrganizations(b.Query)
	if err != nil {
		h.Logger.Errorw("add publication department: could not suggest organization", "error", err)
		render.InternalServerError(w, r, err)
		return
	}

	render.Partial(w, "publication/suggest_departments", YieldAddDepartment{
		Context: ctx,
		Hits:    hits,
	})
}

func (h *Handler) CreateDepartment(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDepartment{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("create publication department: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	org, err := h.OrganizationService.GetOrganization(b.DepartmentID)
	if err != nil {
		h.Logger.Errorw("create publication department: could not find organization", "error", err, "request", r)
		render.InternalServerError(w, r, err)
		return
	}

	/*
		Note: AddDepartmentByOrg removes potential existing department
		and then adds the new one at the end
	*/
	ctx.Publication.AddDepartmentByOrg(org)

	// TODO handle validation errors

	err = h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		h.Logger.Warnf("create publication department: snapstore detected a conflicting publication:", "errors", errors.As(err, &conflict), "identifier", ctx.Publication.ID)
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		h.Logger.Errorf("create publication department: Could not save the publication:", "error", err, "identifier", ctx.Publication.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_departments", YieldDepartments{
		Context: ctx,
	})
}

func (h *Handler) ConfirmDeleteDepartment(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeleteDepartment{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("confirm delete publication department: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	render.Layout(w, "show_modal", "publication/confirm_delete_department", YieldDeleteDepartment{
		Context:      ctx,
		DepartmentID: b.DepartmentID,
	})
}

func (h *Handler) DeleteDepartment(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteDepartment
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete publication department: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	/*
		Ignore fact that department is removed in the mean time:
		conflict resolving will solve this
	*/
	ctx.Publication.RemoveDepartment(b.DepartmentID)

	// TODO handle validation errors

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		h.Logger.Warnf("delete publication department: snapstore detected a conflicting publication:", "errors", errors.As(err, &conflict), "identifier", ctx.Publication.ID)
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		h.Logger.Errorf("delete publication department: Could not save the publication:", "error", err, "identifier", ctx.Publication.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_departments", YieldDepartments{
		Context: ctx,
	})
}
