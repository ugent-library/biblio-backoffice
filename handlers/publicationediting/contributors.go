package publicationediting

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	publicationviews "github.com/ugent-library/biblio-backoffice/views/publication"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

type BindAddContributor struct {
	Role      string `path:"role"`
	FirstName string `query:"first_name"`
	LastName  string `query:"last_name"`
}

type BindAddContributorSuggest struct {
	Role       string   `path:"role"`
	CreditRole []string `query:"credit_role"`
	FirstName  string   `query:"first_name"`
	LastName   string   `query:"last_name"`
}

type BindConfirmCreateContributor struct {
	Role      string `path:"role"`
	ID        string `query:"id"`
	FirstName string `query:"first_name"`
	LastName  string `query:"last_name"`
}

type BindCreateContributor struct {
	Role       string   `path:"role"`
	ID         string   `form:"id"`
	CreditRole []string `form:"credit_role"`
	FirstName  string   `form:"first_name"`
	LastName   string   `form:"last_name"`
	AddNext    bool     `form:"add_next"`
}

type BindEditContributor struct {
	Role      string `path:"role"`
	Position  int    `path:"position"`
	FirstName string `query:"first_name"`
	LastName  string `query:"last_name"`
}

type BindEditContributorSuggest struct {
	Role       string   `path:"role"`
	Position   int      `path:"position"`
	CreditRole []string `query:"credit_role"`
	FirstName  string   `query:"first_name"`
	LastName   string   `query:"last_name"`
}

type BindConfirmUpdateContributor struct {
	Role      string `path:"role"`
	Position  int    `path:"position"`
	ID        string `query:"id"`
	FirstName string `query:"first_name"`
	LastName  string `query:"last_name"`
}

type BindUpdateContributor struct {
	Role       string   `path:"role"`
	Position   int      `path:"position"`
	ID         string   `form:"id"`
	CreditRole []string `form:"credit_role"`
	FirstName  string   `form:"first_name"`
	LastName   string   `form:"last_name"`
	EditNext   bool     `form:"edit_next"`
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

type YieldAddContributor struct {
	Context
	Role        string
	Contributor *models.Contributor
	Form        *form.Form
	Hits        []*models.Contributor
}

type YieldAddContributorSuggest struct {
	Context
	Role        string
	Contributor *models.Contributor
	Hits        []*models.Contributor
}

type YieldConfirmCreateContributor struct {
	Context
	Role        string
	Contributor *models.Contributor
	Active      bool
	Form        *form.Form
}

type YieldEditContributor struct {
	Context
	Role        string
	Position    int
	Contributor *models.Contributor
	FirstName   string
	LastName    string
	Active      bool
	Hits        []*models.Contributor
	Form        *form.Form
}

type YieldEditContributorSuggest struct {
	Context
	Role        string
	Position    int
	Contributor *models.Contributor
	FirstName   string
	LastName    string
	Hits        []*models.Contributor
}

type YieldConfirmUpdateContributor struct {
	Context
	Role        string
	Position    int
	Contributor *models.Contributor
	Active      bool
	Form        *form.Form
	EditNext    bool
}

func (h *Handler) AddContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAddContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("add publication contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	var hits []*models.Contributor

	if b.FirstName != "" || b.LastName != "" {
		people, err := h.PersonSearchService.SuggestPeople(b.FirstName + " " + b.LastName)
		if err != nil {
			h.Logger.Errorw("suggest publication contributor: could not suggest people", "errors", err, "request", r, "user", ctx.User.ID)
			render.InternalServerError(w, r, err)
			return
		}
		hits = make([]*models.Contributor, len(people))
		for i, person := range people {
			hits[i] = models.ContributorFromPerson(person)
		}
	}

	c := models.ContributorFromFirstLastName(b.FirstName, b.LastName)

	suggestURL := h.PathFor("publication_add_contributor_suggest", "id", ctx.Publication.ID, "role", b.Role).String()

	render.Layout(w, "show_modal", "publication/add_contributor", YieldAddContributor{
		Context:     ctx,
		Role:        b.Role,
		Contributor: c,
		Form:        contributorForm(ctx, c, suggestURL),
		Hits:        hits,
	})
}

func (h *Handler) AddContributorSuggest(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAddContributorSuggest{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("suggest publication contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	var hits []*models.Contributor

	if b.FirstName != "" || b.LastName != "" {
		people, err := h.PersonSearchService.SuggestPeople(b.FirstName + " " + b.LastName)
		if err != nil {
			h.Logger.Errorw("suggest publication contributor: could not suggest people", "errors", err, "request", r, "user", ctx.User.ID)
			render.InternalServerError(w, r, err)
			return
		}
		hits = make([]*models.Contributor, len(people))
		for i, person := range people {
			hits[i] = models.ContributorFromPerson(person)
		}
	}

	c := models.ContributorFromFirstLastName(b.FirstName, b.LastName)

	render.Partial(w, "publication/add_contributor_suggest", YieldAddContributorSuggest{
		Context:     ctx,
		Role:        b.Role,
		Contributor: c,
		Hits:        hits,
	})
}

func (h *Handler) ConfirmCreateContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindConfirmCreateContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("confirm create publication contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	var c *models.Contributor
	active := false
	if b.ID != "" {
		newC, newP, err := h.generateContributorFromPersonId(b.ID)
		if err != nil {
			h.ActionError(w, r, ctx.BaseContext, "confirm create publication contributor: could not generate contributor from person", err, ctx.Publication.ID)
			return
		}
		c = newC
		active = newP.Active
	} else {
		c = models.ContributorFromFirstLastName(b.FirstName, b.LastName)
	}

	render.Layout(w, "refresh_modal", "publication/confirm_create_contributor", YieldConfirmCreateContributor{
		Context:     ctx,
		Role:        b.Role,
		Contributor: c,
		Active:      active,
		Form:        confirmContributorForm(ctx, b.Role, c, nil),
	})
}

func (h *Handler) CreateContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindCreateContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("create publication contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	var c *models.Contributor
	active := false
	if b.ID != "" {
		newC, newP, err := h.generateContributorFromPersonId(b.ID)
		if err != nil {
			h.ActionError(w, r, ctx.BaseContext, "create publication contributor: could not generate contributor from person", err, ctx.Publication.ID)
			return
		}
		c = newC
		active = newP.Active
	} else {
		if b.FirstName == "" {
			b.FirstName = "[missing]"
		}
		if b.LastName == "" {
			b.LastName = "[missing]"
		}
		c = models.ContributorFromFirstLastName(b.FirstName, b.LastName)
	}
	c.CreditRole = b.CreditRole

	ctx.Publication.AddContributor(b.Role, c)

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "publication/confirm_create_contributor", YieldConfirmCreateContributor{
			Context:     ctx,
			Role:        b.Role,
			Contributor: c,
			Active:      active,
			Form:        confirmContributorForm(ctx, b.Role, c, validationErrs.(*okay.Errors)),
		})

		return
	}

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(ctx.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		h.Logger.Errorf("create publication contributor: Could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	if b.AddNext {
		c := &models.Contributor{}

		suggestURL := h.PathFor("publication_add_contributor_suggest", "id", ctx.Publication.ID, "role", b.Role).String()

		render.Partial(w, "publication/add_next_contributor", YieldAddContributor{
			Context:     ctx,
			Role:        b.Role,
			Contributor: c,
			Form:        contributorForm(ctx, c, suggestURL),
		})

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
		h.Logger.Warnw("edit publication contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	c, err := ctx.Publication.GetContributor(b.Role, b.Position)
	if err != nil {
		h.Logger.Errorw("edit publication contributor: could not get the contributor", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	active := false
	if c.Person != nil {
		active = c.Person.Active
	}

	firstName := b.FirstName
	lastName := b.LastName
	if firstName == "" && lastName == "" {
		firstName = c.FirstName()
		lastName = c.LastName()
	}

	people, err := h.PersonSearchService.SuggestPeople(firstName + " " + lastName)
	if err != nil {
		h.Logger.Errorw("suggest publication contributor: could not suggest people", "errors", err, "request", r, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	hits := make([]*models.Contributor, len(people))
	for i, person := range people {
		hits[i] = models.ContributorFromPerson(person)
	}

	// exclude the current contributor
	if c.PersonID != "" {
		for i, hit := range hits {
			if hit.PersonID == c.PersonID {
				if i == 0 {
					hits = hits[1:]
				} else {
					hits = append(hits[:i], hits[i+1:]...)
				}
				break
			}
		}
	}

	suggestURL := h.PathFor("publication_edit_contributor_suggest", "id", ctx.Publication.ID, "role", b.Role, "position", fmt.Sprintf("%d", b.Position)).String()

	render.Layout(w, "show_modal", "publication/edit_contributor", YieldEditContributor{
		Context:     ctx,
		Role:        b.Role,
		Position:    b.Position,
		Contributor: c,
		FirstName:   firstName,
		LastName:    lastName,
		Active:      active,
		Hits:        hits,
		Form:        contributorForm(ctx, c, suggestURL),
	})
}

func (h *Handler) EditContributorSuggest(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindEditContributorSuggest{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("suggest publication contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	c, err := ctx.Publication.GetContributor(b.Role, b.Position)
	if err != nil {
		h.Logger.Errorw("edit publication contributor: could not get the contributor", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	var hits []*models.Contributor

	if b.FirstName != "" || b.LastName != "" {
		people, err := h.PersonSearchService.SuggestPeople(b.FirstName + " " + b.LastName)
		if err != nil {
			h.Logger.Errorw("suggest publication contributor: could not suggest people", "errors", err, "request", r, "user", ctx.User.ID)
			render.InternalServerError(w, r, err)
			return
		}

		hits = make([]*models.Contributor, len(people))
		for i, person := range people {
			hits[i] = models.ContributorFromPerson(person)
		}

		// exclude the current contributor
		if c.PersonID != "" {
			for i, hit := range hits {
				if hit.PersonID == c.PersonID {
					hits = append(hits[:i], hits[i+1:]...)
					break
				}
			}
		}
	}

	render.Partial(w, "publication/edit_contributor_suggest", YieldEditContributorSuggest{
		Context:     ctx,
		Role:        b.Role,
		Position:    b.Position,
		Contributor: c,
		FirstName:   b.FirstName,
		LastName:    b.LastName,
		Hits:        hits,
	})
}

func (h *Handler) ConfirmUpdateContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindConfirmUpdateContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("confirm update publication contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	oldC, err := ctx.Publication.GetContributor(b.Role, b.Position)
	if err != nil {
		h.Logger.Errorw("edit publication contributor: could not get the contributor", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	var c *models.Contributor
	active := false
	if b.ID != "" {
		newC, newP, err := h.generateContributorFromPersonId(b.ID)
		if err != nil {
			h.ActionError(w, r, ctx.BaseContext, "edit publication contributor: could not generate contributor from person", err, ctx.Publication.ID)
			return
		}
		c = newC
		active = newP.Active
	} else {
		c = models.ContributorFromFirstLastName(b.FirstName, b.LastName)
	}

	c.CreditRole = oldC.CreditRole

	render.Layout(w, "refresh_modal", "publication/confirm_update_contributor", YieldConfirmUpdateContributor{
		Context:     ctx,
		Role:        b.Role,
		Position:    b.Position,
		Contributor: c,
		Active:      active,
		Form:        confirmContributorForm(ctx, b.Role, c, nil),
		EditNext:    b.Position+1 < len(ctx.Publication.Contributors(b.Role)),
	})
}

func (h *Handler) UpdateContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindUpdateContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update publication contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	var c *models.Contributor
	active := false
	if b.ID != "" {
		newC, newP, err := h.generateContributorFromPersonId(b.ID)
		if err != nil {
			h.Logger.Errorw("update publication contributor: could not fetch person", "errors", err, "personid", b.ID, "publication", ctx.Publication.ID, "user", ctx.User.ID)
			render.InternalServerError(w, r, err)
			return
		}
		c = newC
		active = newP.Active
	} else {
		if b.FirstName == "" {
			b.FirstName = "[missing]"
		}
		if b.LastName == "" {
			b.LastName = "[missing]"
		}
		c = models.ContributorFromFirstLastName(b.FirstName, b.LastName)
	}
	c.CreditRole = b.CreditRole

	if err := ctx.Publication.SetContributor(b.Role, b.Position, c); err != nil {
		h.Logger.Errorw("update publication contributor: could not set the contributor", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "publication/confirm_update_contributor", YieldConfirmUpdateContributor{
			Context:     ctx,
			Role:        b.Role,
			Position:    b.Position,
			Contributor: c,
			Active:      active,
			Form:        confirmContributorForm(ctx, b.Role, c, validationErrs.(*okay.Errors)),
			EditNext:    b.Position+1 < len(ctx.Publication.Contributors(b.Role)),
		})

		return
	}

	err := h.Repo.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(ctx.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		h.Logger.Errorf("update publication contributor: Could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	if b.EditNext && b.Position+1 < len(ctx.Publication.Contributors(b.Role)) {
		nextPos := b.Position + 1
		nextC := ctx.Publication.Contributors(b.Role)[nextPos]
		people, err := h.PersonSearchService.SuggestPeople(nextC.Name())
		if err != nil {
			h.Logger.Errorw("suggest publication contributor: could not suggest people", "errors", err, "request", r, "user", ctx.User.ID)
			render.InternalServerError(w, r, err)
			return
		}
		hits := make([]*models.Contributor, len(people))
		for i, person := range people {
			hits[i] = models.ContributorFromPerson(person)
		}

		suggestURL := h.PathFor("publication_edit_contributor_suggest", "id", ctx.Publication.ID, "role", b.Role, "position", fmt.Sprintf("%d", nextPos)).String()

		nextActive := false
		if nextC.Person != nil {
			nextActive = nextC.Person.Active
		}

		render.Partial(w, "publication/edit_next_contributor", YieldEditContributor{
			Context:     ctx,
			Role:        b.Role,
			Position:    nextPos,
			Contributor: nextC,
			FirstName:   nextC.FirstName(),
			LastName:    nextC.LastName(),
			Active:      nextActive,
			Hits:        hits,
			Form:        contributorForm(ctx, nextC, suggestURL),
		})

		return
	}

	render.View(w, "publication/refresh_contributors", YieldContributors{
		Context: ctx,
		Role:    b.Role,
	})
}

func ConfirmDeleteContributor(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	b := BindDeleteContributor{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete publication contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   fmt.Sprintf("Are you sure you want to remove this %s?", c.Loc.Get("publication.contributor.role."+b.Role)),
		DeleteUrl:  c.PathTo("publication_delete_contributor", "id", publication.ID, "role", b.Role, "position", strconv.Itoa(b.Position)),
		SnapshotID: publication.SnapshotID,
	}).Render(r.Context(), w)
}

func DeleteContributor(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindDeleteContributor{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("delete publication contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	p := ctx.GetPublication(r)

	if err := p.RemoveContributor(b.Role, b.Position); err != nil {
		c.Log.Warnw("delete publication contributor: could not remove contributor", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	if validationErrs := p.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(c.Loc, validationErrs.(*okay.Errors)))
		views.ReplaceModal(views.FormErrorsDialog("Can't delete this contributor due to the following errors", errors)).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("delete publication contributor: Could not save the publication:", "error", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(fmt.Sprintf("#contributors-%s-body", b.Role), publicationviews.ContributorsBody(c, p, b.Role)).Render(r.Context(), w)
}

func OrderContributors(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindOrderContributors{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("order publication contributors: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	contributors := p.Contributors(b.Role)
	if len(b.Positions) != len(contributors) {
		err := fmt.Errorf("positions don't match number of contributors")
		c.Log.Warnw("order publication contributors: could not order contributors", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	newContributors := make([]*models.Contributor, len(contributors))
	for i, pos := range b.Positions {
		newContributors[i] = contributors[pos]
	}
	p.SetContributors(b.Role, newContributors)

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		w.Header().Add("HX-Retarget", "#modals")
		w.Header().Add("HX-Reswap", "innerHTML")
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("order publication contributors: Could not save the publication:", "errors", err, "identifier", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(fmt.Sprintf("#contributors-%s-body", b.Role), publicationviews.ContributorsBody(c, p, b.Role)).Render(r.Context(), w)
}

func contributorForm(_ Context, c *models.Contributor, suggestURL string) *form.Form {
	return form.New().
		WithTheme("cols").
		AddSection(
			&form.Text{
				Template: "contributor_name",
				Name:     "first_name",
				Value:    c.FirstName(),
				Label:    "First name",
				Vars: struct {
					SuggestURL string
				}{
					SuggestURL: suggestURL,
				},
			},
			&form.Text{
				Template: "contributor_name",
				Name:     "last_name",
				Value:    c.LastName(),
				Label:    "Last name",
				Vars: struct {
					SuggestURL string
				}{
					SuggestURL: suggestURL,
				},
			},
		)
}

func confirmContributorForm(ctx Context, role string, c *models.Contributor, errors *okay.Errors) *form.Form {
	var fields []form.Field

	if c.PersonID != "" {
		fields = append(fields, &form.Hidden{
			Name:  "id",
			Value: c.PersonID,
		})
	} else {
		fields = append(fields,
			&form.Hidden{
				Name:  "first_name",
				Value: c.FirstName(),
			}, &form.Hidden{
				Name:  "last_name",
				Value: c.LastName(),
			})
	}

	if role == "author" {
		fields = append(fields, &form.SelectRepeat{
			Name:        "credit_role",
			Label:       "Roles",
			Options:     localize.VocabularySelectOptions(ctx.Loc, "credit_roles"),
			Values:      c.CreditRole,
			EmptyOption: true,
		})
	}

	return form.New().
		WithErrors(localize.ValidationErrors(ctx.Loc, errors)).
		WithTheme("cols").
		AddSection(fields...)
}

func (h *Handler) generateContributorFromPersonId(id string) (*models.Contributor, *models.Person, error) {
	p, err := h.PersonService.GetPerson(id)
	if err != nil {
		return nil, nil, err
	}
	c := models.ContributorFromPerson(p)
	return c, p, nil
}
