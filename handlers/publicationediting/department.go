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

func AddDepartment(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	hits, err := c.OrganizationSearchService.SuggestOrganizations("")
	if err != nil {
		c.Log.Errorw("add publication department: could not suggest organization", "errors", err, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	publicationviews.AddDepartment(c, ctx.GetPublication(r), hits).Render(r.Context(), w)
}

func SuggestDepartments(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindSuggestDepartments{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("suggest publication departments could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	hits, err := c.OrganizationSearchService.SuggestOrganizations(b.Query)
	if err != nil {
		c.Log.Errorw("add publication department: could not suggest organization", "errors", err, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	publicationviews.SuggestDepartments(c, ctx.GetPublication(r), hits).Render(r.Context(), w)
}

func CreateDepartment(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindDepartment{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("create publication department: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	org, err := c.OrganizationService.GetOrganization(b.DepartmentID)
	if err != nil {
		c.Log.Errorw("create publication department: could not find organization", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	p.AddOrganization(org)

	// TODO handle validation errors

	err = c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("create publication department: Could not save the publication:", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(publicationviews.DepartmentsBodySelector, publicationviews.DepartmentsBody(c, p)).Render(r.Context(), w)
}

func ConfirmDeleteDepartment(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindDeleteDepartment{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete publication department: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	// TODO why is this necessary for department id's containing an asterisk?
	depID, _ := url.QueryUnescape(b.DepartmentID)
	b.DepartmentID = depID

	if b.SnapshotID != p.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to remove this department from the publication?",
		DeleteUrl:  c.PathTo("publication_delete_department", "id", p.ID, "department_id", b.DepartmentID),
		SnapshotID: b.SnapshotID,
	}).Render(r.Context(), w)
}

func DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	var b BindDeleteDepartment
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("delete publication department: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	// TODO why is this necessary for department id's containing an asterisk?
	depID, _ := url.QueryUnescape(b.DepartmentID)
	b.DepartmentID = depID

	p.RemoveOrganization(b.DepartmentID)

	// TODO handle validation errors

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("delete publication department: Could not save the publication:", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(publicationviews.DepartmentsBodySelector, publicationviews.DepartmentsBody(c, p)).Render(r.Context(), w)
}
