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

type YieldConfirmUpdateContributor struct {
	Context
	Role        string
	Position    int
	Contributor *models.Contributor
	Active      bool
	Form        *form.Form
	EditNext    bool
}

func AddContributor(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindAddContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("add publication contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	var hits []*models.Contributor

	if b.FirstName != "" || b.LastName != "" {
		people, err := c.PersonSearchService.SuggestPeople(b.FirstName + " " + b.LastName)
		if err != nil {
			c.Log.Errorw("suggest publication contributor: could not suggest people", "errors", err, "request", r, "user", c.User.ID)
			c.HandleError(w, r, httperror.InternalServerError)
			return
		}
		hits = make([]*models.Contributor, len(people))
		for i, person := range people {
			hits[i] = models.ContributorFromPerson(person)
		}
	}

	contributor := models.ContributorFromFirstLastName(b.FirstName, b.LastName)

	// suggestURL := c.PathFor("publication_add_contributor_suggest", "id", ctx.Publication.ID, "role", b.Role).String()

	views.ShowModal(
		publicationviews.AddContributor(c, publicationviews.AddContributorArgs{
			Publication: ctx.GetPublication(r),
			Role:        b.Role,
			Contributor: contributor,
			Hits:        hits,
		}),
	).Render(r.Context(), w)
}

func AddContributorSuggest(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindAddContributorSuggest{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("suggest publication contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	var hits []*models.Contributor

	if b.FirstName != "" || b.LastName != "" {
		people, err := c.PersonSearchService.SuggestPeople(b.FirstName + " " + b.LastName)
		if err != nil {
			c.Log.Errorw("suggest publication contributor: could not suggest people", "errors", err, "request", r, "user", c.User.ID)
			c.HandleError(w, r, httperror.InternalServerError)
			return
		}
		hits = make([]*models.Contributor, len(people))
		for i, person := range people {
			hits[i] = models.ContributorFromPerson(person)
		}
	}

	contributor := models.ContributorFromFirstLastName(b.FirstName, b.LastName)

	publicationviews.AddContributorSuggest(c, publicationviews.AddContributorSuggestArgs{
		Publication: ctx.GetPublication(r),
		Role:        b.Role,
		Contributor: contributor,
		Hits:        hits,
	}).Render(r.Context(), w)
}

func ConfirmCreateContributor(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindConfirmCreateContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("confirm create publication contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	var contributor *models.Contributor
	if b.ID != "" {
		newC, _, err := generateContributorFromPersonId(c, b.ID)
		if err != nil {
			c.HandleError(w, r, err)
			return
		}
		contributor = newC
	} else {
		contributor = models.ContributorFromFirstLastName(b.FirstName, b.LastName)
	}

	views.ReplaceModal(publicationviews.ConfirmCreateContributor(c, publicationviews.ConfirmCreateContributorArgs{
		Publication: p,
		Contributor: contributor,
		Role:        b.Role,
	})).Render(r.Context(), w)
}

func CreateContributor(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindCreateContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("create publication contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	var contributor *models.Contributor
	if b.ID != "" {
		newC, _, err := generateContributorFromPersonId(c, b.ID)
		if err != nil {
			c.HandleError(w, r, err)
			return
		}
		contributor = newC
	} else {
		if b.FirstName == "" {
			b.FirstName = "[missing]"
		}
		if b.LastName == "" {
			b.LastName = "[missing]"
		}
		contributor = models.ContributorFromFirstLastName(b.FirstName, b.LastName)
	}
	contributor.CreditRole = b.CreditRole

	p.AddContributor(b.Role, contributor)

	if validationErrs := p.Validate(); validationErrs != nil {
		views.ReplaceModal(publicationviews.ConfirmCreateContributor(c, publicationviews.ConfirmCreateContributorArgs{
			Publication: p,
			Contributor: contributor,
			Role:        b.Role,
			Errors:      validationErrs.(*okay.Errors),
		})).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("create publication contributor: Could not save the publication:", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	if b.AddNext {
		views.Cat(
			views.Replace(fmt.Sprintf("#contributors-%s-body", b.Role), publicationviews.ContributorsBody(
				c, p, b.Role,
			)),
			views.ReplaceModal(publicationviews.AddContributor(c, publicationviews.AddContributorArgs{
				Publication: p,
				Role:        b.Role,
				Contributor: &models.Contributor{},
			})),
		).Render(r.Context(), w)
		return
	}

	views.CloseModalAndReplace(fmt.Sprintf("#contributors-%s-body", b.Role), publicationviews.ContributorsBody(c, p, b.Role)).Render(r.Context(), w)
}

func EditContributor(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindEditContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("edit publication contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	contributor, err := p.GetContributor(b.Role, b.Position)
	if err != nil {
		c.Log.Errorw("edit publication contributor: could not get the contributor", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	firstName := b.FirstName
	lastName := b.LastName
	if firstName == "" && lastName == "" {
		firstName = contributor.FirstName()
		lastName = contributor.LastName()
	}

	people, err := c.PersonSearchService.SuggestPeople(firstName + " " + lastName)
	if err != nil {
		c.Log.Errorw("suggest publication contributor: could not suggest people", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	hits := make([]*models.Contributor, len(people))
	for i, person := range people {
		hits[i] = models.ContributorFromPerson(person)
	}

	// exclude the current contributor
	if contributor.PersonID != "" {
		for i, hit := range hits {
			if hit.PersonID == contributor.PersonID {
				if i == 0 {
					hits = hits[1:]
				} else {
					hits = append(hits[:i], hits[i+1:]...)
				}
				break
			}
		}
	}

	views.ShowModal(publicationviews.EditContributor(c, publicationviews.EditContributorArgs{
		Publication: p,
		Role:        b.Role,
		Position:    b.Position,
		Contributor: contributor,
		FirstName:   firstName,
		LastName:    lastName,
		Hits:        hits,
	})).Render(r.Context(), w)
}

func EditContributorSuggest(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindEditContributorSuggest{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("suggest publication contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	contributor, err := p.GetContributor(b.Role, b.Position)
	if err != nil {
		c.Log.Errorw("edit publication contributor: could not get the contributor", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	var hits []*models.Contributor

	if b.FirstName != "" || b.LastName != "" {
		people, err := c.PersonSearchService.SuggestPeople(b.FirstName + " " + b.LastName)
		if err != nil {
			c.Log.Errorw("suggest publication contributor: could not suggest people", "errors", err, "request", r, "user", c.User.ID)
			c.HandleError(w, r, httperror.InternalServerError)
			return
		}

		hits = make([]*models.Contributor, len(people))
		for i, person := range people {
			hits[i] = models.ContributorFromPerson(person)
		}

		// exclude the current contributor
		if contributor.PersonID != "" {
			for i, hit := range hits {
				if hit.PersonID == contributor.PersonID {
					hits = append(hits[:i], hits[i+1:]...)
					break
				}
			}
		}
	}

	publicationviews.EditContributorSuggest(c, publicationviews.EditContributorSuggestArgs{
		Publication: p,
		Role:        b.Role,
		Position:    b.Position,
		Contributor: contributor,
		FirstName:   b.FirstName,
		LastName:    b.LastName,
		Hits:        hits,
	}).Render(r.Context(), w)
}

func ConfirmUpdateContributor(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindConfirmUpdateContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("confirm update publication contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	oldC, err := p.GetContributor(b.Role, b.Position)
	if err != nil {
		c.Log.Errorw("edit publication contributor: could not get the contributor", "errors", err, "publication", p.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	var contributor *models.Contributor
	if b.ID != "" {
		newC, _, err := generateContributorFromPersonId(c, b.ID)
		if err != nil {
			c.HandleError(w, r, err)
			return
		}
		contributor = newC
	} else {
		contributor = models.ContributorFromFirstLastName(b.FirstName, b.LastName)
	}

	contributor.CreditRole = oldC.CreditRole

	views.ReplaceModal(publicationviews.ConfirmUpdateContributor(c, publicationviews.ConfirmUpdateContributorArgs{
		Publication: p,
		Role:        b.Role,
		Position:    b.Position,
		Contributor: contributor,
		EditNext:    b.Position+1 < len(p.Contributors(b.Role)),
	})).Render(r.Context(), w)
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

func generateContributorFromPersonId(c *ctx.Ctx, id string) (*models.Contributor, *models.Person, error) {
	person, err := c.PersonService.GetPerson(id)
	if err != nil {
		return nil, nil, err
	}

	contributor := models.ContributorFromPerson(person)
	return contributor, person, nil
}
