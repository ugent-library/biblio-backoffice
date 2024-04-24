package datasetediting

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	views "github.com/ugent-library/biblio-backoffice/views/dataset"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
)

type BindSuggestDepartments struct {
	Query string `query:"q"`
}

type BindDepartment struct {
	DepartmentID string `form:"department_id"`
}

type BindDeleteDepartment struct {
	DepartmentID string `path:"department_id"`
	SnapshotID   string `path:"snapshot_id"`
}

type YieldDepartments struct {
	Context
}

type YieldDeleteDepartment struct {
	Context
	DepartmentID string
}

func AddDepartment(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	hits, err := c.OrganizationSearchService.SuggestOrganizations("")
	if err != nil {
		c.Log.Errorw("add dataset department: could not suggest organization", "errors", err, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.AddDepartment(c, ctx.GetDataset(r), hits).Render(r.Context(), w)
}

func SuggestDepartments(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindSuggestDepartments{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("suggest dataset departments could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	hits, err := c.OrganizationSearchService.SuggestOrganizations(b.Query)
	if err != nil {
		c.Log.Errorw("add dataset department: could not suggest organization", "errors", err, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.SuggestDepartments(c, ctx.GetDataset(r), hits).Render(r.Context(), w)
}

func (h *Handler) CreateDepartment(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDepartment{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("create dataset department: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	org, err := h.OrganizationService.GetOrganization(b.DepartmentID)
	if err != nil {
		h.Logger.Errorw("create dataset department: could not find organization", "errors", err, "dataset", ctx.Dataset.ID, "department", b.DepartmentID, r, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	ctx.Dataset.AddOrganization(org)

	// TODO handle validation errors

	err = h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("dataset.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("create dataset department: Could not save the dataset:", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
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
		h.Logger.Warnw("confirm delete dataset department: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	// TODO why is this necessary for department id's containing an asterisk?
	depID, _ := url.QueryUnescape(b.DepartmentID)
	b.DepartmentID = depID

	if b.SnapshotID != ctx.Dataset.SnapshotID {
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("dataset.conflict_error_reload"),
		})
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
		h.Logger.Warnw("delete dataset department: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	// TODO why is this necessary for department id's containing an asterisk?
	depID, _ := url.QueryUnescape(b.DepartmentID)
	b.DepartmentID = depID

	ctx.Dataset.RemoveOrganization(b.DepartmentID)

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("dataset.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("delete dataset department: Could not save the dataset:", "error", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_departments", YieldDepartments{
		Context: ctx,
	})
}
