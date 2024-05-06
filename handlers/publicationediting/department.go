package publicationediting

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/render"
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

func (h *Handler) CreateDepartment(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDepartment{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("create publication department: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	org, err := h.OrganizationService.GetOrganization(b.DepartmentID)
	if err != nil {
		h.Logger.Errorw("create publication department: could not find organization", "errors", err, "request", r, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	/*
		Note: AddDepartmentByOrg removes potential existing department
		and then adds the new one at the end
	*/
	ctx.Publication.AddOrganization(org)

	// TODO handle validation errors

	err = h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(ctx.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		h.Logger.Errorf("create publication department: Could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_departments", YieldDepartments{
		Context: ctx,
	})
}

func ConfirmDeleteDepartment(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	b := BindDeleteDepartment{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete publication department: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	// TODO why is this necessary for department id's containing an asterisk?
	depID, _ := url.QueryUnescape(b.DepartmentID)
	b.DepartmentID = depID

	if b.SnapshotID != publication.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to remove this department from the publication?",
		DeleteUrl:  c.PathTo("publication_delete_department", "id", publication.ID, "department_id", b.DepartmentID),
		SnapshotID: b.SnapshotID,
	}).Render(r.Context(), w)
}

func (h *Handler) DeleteDepartment(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteDepartment
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete publication department: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	// TODO why is this necessary for department id's containing an asterisk?
	depID, _ := url.QueryUnescape(b.DepartmentID)
	b.DepartmentID = depID

	ctx.Publication.RemoveOrganization(b.DepartmentID)

	// TODO handle validation errors

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(ctx.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		h.Logger.Errorf("delete publication department: Could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_departments", YieldDepartments{
		Context: ctx,
	})
}
