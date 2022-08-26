package publicationediting

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type BindAddContributor struct {
	Role       string   `path:"role"`
	CreditRole []string `query:"credit_role"`
	FirstName  string   `query:"first_name"`
	LastName   string   `query:"last_name"`
}

type BindSuggestContributors struct {
	Role       string   `path:"role"`
	Position   int      `path:"position"`
	CreditRole []string `query:"credit_role"`
	FirstName  string   `query:"first_name"`
	LastName   string   `query:"last_name"`
}

type BindConfirmContributor struct {
	Role       string   `path:"role"`
	Position   int      `path:"position"`
	ID         string   `form:"id"`
	CreditRole []string `form:"credit_role"`
}

type BindUnconfirmContributor struct {
	Role       string   `path:"role"`
	Position   int      `path:"position"`
	CreditRole []string `form:"credit_role"`
	FirstName  string   `form:"first_name"`
	LastName   string   `form:"last_name"`
}

type BindCreateContributor struct {
	Role       string   `path:"role"`
	ID         string   `form:"id"`
	CreditRole []string `form:"credit_role"`
	FirstName  string   `form:"first_name"`
	LastName   string   `form:"last_name"`
}

type BindEditContributor struct {
	Role     string `path:"role"`
	Position int    `path:"position"`
}

type BindUpdateContributor struct {
	Role       string   `path:"role"`
	Position   int      `path:"position"`
	ID         string   `form:"id"`
	CreditRole []string `form:"credit_role"`
	FirstName  string   `form:"first_name"`
	LastName   string   `form:"last_name"`
}

type BindDeleteContributor struct {
	Role     string `path:"role"`
	Position int    `path:"position"`
}

type BindOrderContributors struct {
	Role      string `path:"role"`
	Positions []int  `form:"position"`
}

type YieldContributors struct {
	Context
	Role string
}

type YieldContributorForm struct {
	Context
	Role        string
	Position    int
	Contributor *models.Contributor
	Form        *form.Form
}

type YieldSuggestContributors struct {
	Context
	Role        string
	Position    int
	Contributor *models.Contributor
	Hits        []models.Person
}

type YieldDeleteContributor struct {
	Context
	Role     string
	Position int
}

func (h *Handler) AddContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAddContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("add publication contributor: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	position := len(ctx.Publication.Contributors(b.Role))

	c := &models.Contributor{
		CreditRole: b.CreditRole,
		FirstName:  b.FirstName,
		LastName:   b.LastName,
	}

	f := contributorForm(ctx, b.Role, position, c, nil)

	render.Layout(w, "show_modal", "publication/add_contributor", YieldContributorForm{
		Context:     ctx,
		Role:        b.Role,
		Position:    position,
		Contributor: c,
		Form:        f,
	})
}

func (h *Handler) SuggestContributors(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindSuggestContributors{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("suggest publication contributor: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	hits, err := h.PersonSearchService.SuggestPeople(b.FirstName + " " + b.LastName)
	if err != nil {
		h.Logger.Errorw("suggest publication contributor: could not suggest people", "error", err, "request", r)
		render.InternalServerError(w, r, err)
		return
	}

	c := &models.Contributor{
		CreditRole: b.CreditRole,
		FirstName:  b.FirstName,
		LastName:   b.LastName,
	}

	render.Layout(w, "refresh_modal", "publication/suggest_contributors", YieldSuggestContributors{
		Context:     ctx,
		Role:        b.Role,
		Position:    b.Position,
		Contributor: c,
		Hits:        hits,
	})
}

func (h *Handler) ConfirmContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindConfirmContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("confirm publication contributor: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	p, err := h.PersonService.GetPerson(b.ID)
	if err != nil {
		h.Logger.Errorw("confirm publication contributor: could not find person", "error", err, "identifier", b.ID)
		render.InternalServerError(w, r, err)
		return
	}

	c := &models.Contributor{
		ID:         p.ID,
		CreditRole: b.CreditRole,
		FirstName:  p.FirstName,
		LastName:   p.LastName,
	}

	f := contributorForm(ctx, b.Role, b.Position, c, nil)

	var tmpl string
	if len(ctx.Publication.Contributors(b.Role)) > b.Position {
		tmpl = "publication/edit_contributor"
	} else {
		tmpl = "publication/add_contributor"
	}

	render.Layout(w, "refresh_modal", tmpl, YieldContributorForm{
		Context:     ctx,
		Role:        b.Role,
		Position:    b.Position,
		Contributor: c,
		Form:        f,
	})
}

func (h *Handler) UnconfirmContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindUnconfirmContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("unconfirm publication contributor: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	c := &models.Contributor{
		CreditRole: b.CreditRole,
		FirstName:  b.FirstName,
		LastName:   b.LastName,
	}

	f := contributorForm(ctx, b.Role, b.Position, c, nil)

	var tmpl string
	if len(ctx.Publication.Contributors(b.Role)) > b.Position {
		tmpl = "publication/edit_contributor"
	} else {
		tmpl = "publication/add_contributor"
	}

	render.Layout(w, "refresh_modal", tmpl, YieldContributorForm{
		Context:     ctx,
		Role:        b.Role,
		Position:    b.Position,
		Contributor: c,
		Form:        f,
	})
}

func (h *Handler) CreateContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindCreateContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("create publication contributor: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	position := len(ctx.Publication.Contributors(b.Role))

	c := &models.Contributor{CreditRole: b.CreditRole}
	if b.ID != "" {
		p, err := h.PersonService.GetPerson(b.ID)
		if err != nil {
			render.InternalServerError(w, r, err)
			return
		}
		c.ID = p.ID
		c.FirstName = p.FirstName
		c.LastName = p.LastName
	} else {
		c.FirstName = b.FirstName
		c.LastName = b.LastName
	}

	ctx.Publication.AddContributor(b.Role, c)

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		h.Logger.Warnw("create publication contributor: could not validate contributor:", "errors", validationErrs, "identifier", ctx.Publication.ID)
		f := contributorForm(ctx, b.Role, position, c, validationErrs.(validation.Errors))
		render.Layout(w, "refresh_modal", "publication/add_contributor", YieldContributorForm{
			Context:     ctx,
			Role:        b.Role,
			Position:    position,
			Contributor: c,
			Form:        f,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		h.Logger.Warnf("create publication contributor: snapstore detected a conflicting publication:", "errors", errors.As(err, &conflict), "identifier", ctx.Publication.ID)
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		h.Logger.Errorf("create publication contributor: Could not save the publication:", "error", err, "identifier", ctx.Publication.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_contributors", YieldContributors{
		Context: ctx,
		Role:    b.Role,
	})
}

func (h *Handler) EditContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindEditContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("edit publication contributor: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	c, err := ctx.Publication.GetContributor(b.Role, b.Position)
	if err != nil {
		h.Logger.Errorw("edit publication contributor: could not get the contributor", "error", err, "publication", ctx.Publication.ID)
		render.InternalServerError(w, r, err)
		return
	}

	f := contributorForm(ctx, b.Role, b.Position, c, nil)

	render.Layout(w, "show_modal", "publication/edit_contributor", YieldContributorForm{
		Context:     ctx,
		Role:        b.Role,
		Position:    b.Position,
		Contributor: c,
		Form:        f,
	})
}

func (h *Handler) UpdateContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindUpdateContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update publication contributor: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	c := &models.Contributor{CreditRole: b.CreditRole}
	if b.ID != "" {
		p, err := h.PersonService.GetPerson(b.ID)
		if err != nil {
			h.Logger.Errorw("update publication contributor: could not fetch person", "error", err, "personid", b.ID, "publication", ctx.Publication.ID)
			render.InternalServerError(w, r, err)
			return
		}
		c.ID = p.ID
		c.FirstName = p.FirstName
		c.LastName = p.LastName
	} else {
		c.FirstName = b.FirstName
		c.LastName = b.LastName
	}

	if err := ctx.Publication.SetContributor(b.Role, b.Position, c); err != nil {
		h.Logger.Errorw("update publication contributor: could not set the contributor", "error", err, "publication", ctx.Publication.ID)
		render.InternalServerError(w, r, err)
		return
	}

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		h.Logger.Warnw("update publication contributor: could not validate contributor:", "errors", validationErrs, "identifier", ctx.Publication.ID)
		f := contributorForm(ctx, b.Role, b.Position, c, validationErrs.(validation.Errors))
		render.Layout(w, "refresh_modal", "publication/edit_contributor", YieldContributorForm{
			Context:     ctx,
			Role:        b.Role,
			Position:    b.Position,
			Contributor: c,
			Form:        f,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		h.Logger.Warnf("update publication contributor: snapstore detected a conflicting publication:", "errors", errors.As(err, &conflict), "identifier", ctx.Publication.ID)
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		h.Logger.Errorf("update publication contributor: Could not save the publication:", "error", err, "identifier", ctx.Publication.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_contributors", YieldContributors{
		Context: ctx,
		Role:    b.Role,
	})
}

func (h *Handler) ConfirmDeleteContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeleteContributor{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("confirm delete publication contributor: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	render.Layout(w, "show_modal", "publication/confirm_delete_contributor", YieldDeleteContributor{
		Context:  ctx,
		Role:     b.Role,
		Position: b.Position,
	})
}

func (h *Handler) DeleteContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeleteContributor{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete publication contributor: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	if err := ctx.Publication.RemoveContributor(b.Role, b.Position); err != nil {
		h.Logger.Warnw("delete publication contributor: could not remove contributor", "error", err, "publication", ctx.Publication.ID)
		render.InternalServerError(w, r, err)
		return
	}

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		h.Logger.Warnw("delete publication contributor: could not validate contributor:", "errors", validationErrs, "identifier", ctx.Publication.ID)
		errors := form.Errors(localize.ValidationErrors(ctx.Locale, validationErrs.(validation.Errors)))
		render.Layout(w, "refresh_modal", "form_errors_dialog", struct {
			Title  string
			Errors form.Errors
		}{
			Title:  "Can't delete this contributor due to the following errors",
			Errors: errors,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		h.Logger.Warnf("delete publication contributor: snapstore detected a conflicting publication:", "errors", errors.As(err, &conflict), "identifier", ctx.Publication.ID)
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		h.Logger.Errorf("delete publication contributor: Could not save the publication:", "error", err, "identifier", ctx.Publication.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_contributors", YieldContributors{
		Context: ctx,
		Role:    b.Role,
	})
}

func (h *Handler) OrderContributors(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindOrderContributors{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("order publication contributors: could not bind request arguments", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	contributors := ctx.Publication.Contributors(b.Role)
	if len(b.Positions) != len(contributors) {
		err := fmt.Errorf("positions don't match number of contributors")
		h.Logger.Warnw("order publication contributors: could not order contributors", "error", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}
	newContributors := make([]*models.Contributor, len(contributors))
	for i, pos := range b.Positions {
		newContributors[i] = contributors[pos]
	}
	ctx.Publication.SetContributors(b.Role, newContributors)

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		h.Logger.Warnf("order publication contributors: snapstore detected a conflicting publication:", "errors", errors.As(err, &conflict), "identifier", ctx.Publication.ID)
		render.Layout(w, "show_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		h.Logger.Errorf("order publication contributors: Could not save the publication:", "error", err, "identifier", ctx.Publication.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_contributors", YieldContributors{
		Context: ctx,
		Role:    b.Role,
	})
}

func contributorForm(ctx Context, role string, position int, c *models.Contributor, errors validation.Errors) *form.Form {
	f := form.New().
		WithTheme("cols").
		WithErrors(localize.ValidationErrors(ctx.Locale, errors)).
		AddSection(
			&form.Hidden{
				Name:  "id",
				Value: c.ID,
			},
			&form.Text{
				Name:     "first_name",
				Value:    c.FirstName,
				Label:    "First name",
				Readonly: c.ID != "",
				Error:    localize.ValidationErrorAt(ctx.Locale, errors, fmt.Sprintf("/%s/%d/first_name", role, position)),
			},
			&form.Text{
				Name:     "last_name",
				Value:    c.LastName,
				Label:    "Last name",
				Readonly: c.ID != "",
				Error:    localize.ValidationErrorAt(ctx.Locale, errors, fmt.Sprintf("/%s/%d/last_name", role, position)),
			},
		)

	if role == "author" {
		f.AddSection(&form.SelectRepeat{
			Name:        "credit_role",
			Label:       "Credit roles",
			Options:     localize.VocabularySelectOptions(ctx.Locale, "credit_roles"),
			Values:      c.CreditRole,
			EmptyOption: true,
		})
	}

	return f
}
