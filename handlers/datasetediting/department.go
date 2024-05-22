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

func AddDepartment(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	hits, err := c.OrganizationSearchService.SuggestOrganizations("")
	if err != nil {
		c.Log.Errorw("add dataset department: could not suggest organization", "errors", err, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	datasetviews.AddDepartment(c, ctx.GetDataset(r), hits).Render(r.Context(), w)
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

	datasetviews.SuggestDepartments(c, ctx.GetDataset(r), hits).Render(r.Context(), w)
}

func CreateDepartment(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	b := BindDepartment{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("create dataset department: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	org, err := c.OrganizationService.GetOrganization(b.DepartmentID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	d.AddOrganization(org)

	// TODO handle validation errors

	err = c.Repo.UpdateDataset(r.Header.Get("If-Match"), d, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("create dataset department: Could not save the dataset:", "errors", err, "dataset", d.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(datasetviews.DepartmentsBodySelector, datasetviews.DepartmentsBody(c, d)).Render(r.Context(), w)
}

func ConfirmDeleteDepartment(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindDeleteDepartment{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete dataset department: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	// TODO why is this necessary for department id's containing an asterisk?
	depID, _ := url.QueryUnescape(b.DepartmentID)
	b.DepartmentID = depID

	if b.SnapshotID != dataset.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to remove this department from the dataset?",
		DeleteUrl:  c.PathTo("dataset_delete_department", "id", dataset.ID, "department_id", b.DepartmentID),
		SnapshotID: dataset.SnapshotID,
	}).Render(r.Context(), w)

}

func DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	var b BindDeleteDepartment
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("delete dataset department: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	// TODO why is this necessary for department id's containing an asterisk?
	depID, _ := url.QueryUnescape(b.DepartmentID)
	b.DepartmentID = depID

	d.RemoveOrganization(b.DepartmentID)

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), d, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("delete dataset department: Could not save the dataset:", "error", err, "dataset", d.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(datasetviews.DepartmentsBodySelector, datasetviews.DepartmentsBody(c, d)).Render(r.Context(), w)
}
