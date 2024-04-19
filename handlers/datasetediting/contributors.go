package datasetediting

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views/dataset"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/okay"
)

type BindAddContributor struct {
	Role      string `path:"role"`
	FirstName string `query:"first_name"`
	LastName  string `query:"last_name"`
}

type BindAddContributorSuggest struct {
	Role      string `path:"role"`
	FirstName string `query:"first_name"`
	LastName  string `query:"last_name"`
}

type BindConfirmCreateContributor struct {
	Role      string `path:"role"`
	ID        string `query:"id"`
	FirstName string `query:"first_name"`
	LastName  string `query:"last_name"`
}

type BindCreateContributor struct {
	Role      string `path:"role"`
	ID        string `form:"id"`
	FirstName string `form:"first_name"`
	LastName  string `form:"last_name"`
	AddNext   bool   `form:"add_next"`
}

type BindEditContributor struct {
	Role      string `path:"role"`
	Position  int    `path:"position"`
	FirstName string `query:"first_name"`
	LastName  string `query:"last_name"`
}

type BindEditContributorSuggest struct {
	Role      string `path:"role"`
	Position  int    `path:"position"`
	FirstName string `query:"first_name"`
	LastName  string `query:"last_name"`
}

type BindConfirmUpdateContributor struct {
	Role      string `path:"role"`
	Position  int    `path:"position"`
	ID        string `query:"id"`
	FirstName string `query:"first_name"`
	LastName  string `query:"last_name"`
}

type BindUpdateContributor struct {
	Role      string `path:"role"`
	Position  int    `path:"position"`
	ID        string `form:"id"`
	FirstName string `form:"first_name"`
	LastName  string `form:"last_name"`
	EditNext  bool   `form:"edit_next"`
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

type YieldDeleteContributor struct {
	Context
	Role     string
	Position int
}

func (h *Handler) AddContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAddContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("add dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	var hits []*models.Contributor

	if b.FirstName != "" || b.LastName != "" {
		people, err := h.PersonSearchService.SuggestPeople(b.FirstName + " " + b.LastName)
		if err != nil {
			h.Logger.Errorw("suggest dataset contributor: could not suggest people", "errors", err, "request", r, "user", ctx.User.ID)
			render.InternalServerError(w, r, err)
			return
		}

		hits = make([]*models.Contributor, len(people))
		for i, person := range people {
			hits[i] = models.ContributorFromPerson(person)
		}
	}

	c := models.ContributorFromFirstLastName(b.FirstName, b.LastName)

	suggestURL := h.PathFor("dataset_add_contributor_suggest", "id", ctx.Dataset.ID, "role", b.Role).String()

	render.Layout(w, "show_modal", "dataset/add_contributor", YieldAddContributor{
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
		h.Logger.Warnw("suggest dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	var hits []*models.Contributor

	if b.FirstName != "" || b.LastName != "" {
		people, err := h.PersonSearchService.SuggestPeople(b.FirstName + " " + b.LastName)
		if err != nil {
			h.Logger.Errorw("suggest dataset contributor: could not suggest people", "errors", err, "request", r, "user", ctx.User.ID)
			render.InternalServerError(w, r, err)
			return
		}

		hits = make([]*models.Contributor, len(people))
		for i, person := range people {
			hits[i] = models.ContributorFromPerson(person)
		}
	}

	c := models.ContributorFromFirstLastName(b.FirstName, b.LastName)

	render.Partial(w, "dataset/add_contributor_suggest", YieldAddContributorSuggest{
		Context:     ctx,
		Role:        b.Role,
		Contributor: c,
		Hits:        hits,
	})
}

func (h *Handler) ConfirmCreateContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindConfirmCreateContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("confirm create dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	var c *models.Contributor
	active := false
	if b.ID != "" {
		newC, newP, err := h.generateContributorFromPersonId(b.ID)
		if err != nil {
			h.ActionError(w, r, ctx.BaseContext, "confirm create dataset contributor: could not generate contributor from person", err, ctx.Dataset.ID)
			return
		}
		c = newC
		active = newP.Active
	} else {
		c = models.ContributorFromFirstLastName(b.FirstName, b.LastName)
	}

	render.Layout(w, "refresh_modal", "dataset/confirm_create_contributor", YieldConfirmCreateContributor{
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
		h.Logger.Warnw("create dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	var c *models.Contributor
	active := false
	if b.ID != "" {
		newC, newP, err := h.generateContributorFromPersonId(b.ID)
		if err != nil {
			h.ActionError(w, r, ctx.BaseContext, "create dataset contributor: could not generate contributor from person", err, ctx.Dataset.ID)
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

	ctx.Dataset.AddContributor(b.Role, c)

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "dataset/confirm_create_contributor", YieldConfirmCreateContributor{
			Context:     ctx,
			Role:        b.Role,
			Contributor: c,
			Active:      active,
			Form:        confirmContributorForm(ctx, b.Role, c, validationErrs.(*okay.Errors)),
		})

		return
	}

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("dataset.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("create dataset contributor: Could not save the dataset:", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	if b.AddNext {
		c := &models.Contributor{}

		suggestURL := h.PathFor("dataset_add_contributor_suggest", "id", ctx.Dataset.ID, "role", b.Role).String()

		render.Partial(w, "dataset/add_next_contributor", YieldAddContributor{
			Context:     ctx,
			Role:        b.Role,
			Contributor: c,
			Form:        contributorForm(ctx, c, suggestURL),
		})

		return
	}

	render.View(w, "dataset/refresh_contributors", YieldContributors{
		Context: ctx,
		Role:    b.Role,
	})
}

func (h *Handler) EditContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindEditContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("edit dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	c, err := ctx.Dataset.GetContributor(b.Role, b.Position)
	if err != nil {
		h.Logger.Errorw("edit dataset contributor: could not get the contributor", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
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
		h.Logger.Errorw("suggest dataset contributor: could not suggest people", "errors", err, "request", r, "user", ctx.User.ID)
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

	suggestURL := h.PathFor("dataset_edit_contributor_suggest", "id", ctx.Dataset.ID, "role", b.Role, "position", fmt.Sprintf("%d", b.Position)).String()

	render.Layout(w, "show_modal", "dataset/edit_contributor", YieldEditContributor{
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
		h.Logger.Warnw("suggest dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	c, err := ctx.Dataset.GetContributor(b.Role, b.Position)
	if err != nil {
		h.Logger.Errorw("edit dataset contributor: could not get the contributor", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	var hits []*models.Contributor

	if b.FirstName != "" || b.LastName != "" {
		people, err := h.PersonSearchService.SuggestPeople(b.FirstName + " " + b.LastName)
		if err != nil {
			h.Logger.Errorw("suggest dataset contributor: could not suggest people", "errors", err, "request", r, "user", ctx.User.ID)
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

	render.Partial(w, "dataset/edit_contributor_suggest", YieldEditContributorSuggest{
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
		h.Logger.Warnw("confirm update dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	var c *models.Contributor
	active := false
	if b.ID != "" {
		newC, newP, err := h.generateContributorFromPersonId(b.ID)
		if err != nil {
			h.ActionError(w, r, ctx.BaseContext, "confirm update dataset contributor: could not generate contributor from person", err, ctx.Dataset.ID)
			return
		}
		c = newC
		active = newP.Active
	} else {
		c = models.ContributorFromFirstLastName(b.FirstName, b.LastName)
	}

	render.Layout(w, "refresh_modal", "dataset/confirm_update_contributor", YieldConfirmUpdateContributor{
		Context:     ctx,
		Role:        b.Role,
		Position:    b.Position,
		Contributor: c,
		Active:      active,
		Form:        confirmContributorForm(ctx, b.Role, c, nil),
		EditNext:    b.Position+1 < len(ctx.Dataset.Contributors(b.Role)),
	})
}

func (h *Handler) UpdateContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindUpdateContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	var c *models.Contributor
	active := false
	if b.ID != "" {
		newC, newP, err := h.generateContributorFromPersonId(b.ID)
		if err != nil {
			h.Logger.Errorw("update dataset contributor: could not fetch person", "errors", err, "personid", b.ID, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
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

	if err := ctx.Dataset.SetContributor(b.Role, b.Position, c); err != nil {
		h.Logger.Errorw("update dataset contributor: could not set the contributor", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "dataset/confirm_update_contributor", YieldConfirmUpdateContributor{
			Context:     ctx,
			Role:        b.Role,
			Position:    b.Position,
			Contributor: c,
			Active:      active,
			Form:        confirmContributorForm(ctx, b.Role, c, validationErrs.(*okay.Errors)),
			EditNext:    b.Position+1 < len(ctx.Dataset.Contributors(b.Role)),
		})

		return
	}

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("dataset.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("update dataset contributor: Could not save the dataset:", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	if b.EditNext && b.Position+1 < len(ctx.Dataset.Contributors(b.Role)) {
		nextPos := b.Position + 1
		nextC := ctx.Dataset.Contributors(b.Role)[nextPos]
		people, err := h.PersonSearchService.SuggestPeople(nextC.Name())
		if err != nil {
			h.Logger.Errorw("suggest dataset contributor: could not suggest people", "errors", err, "request", r, "user", ctx.User.ID)
			render.InternalServerError(w, r, err)
			return
		}
		hits := make([]*models.Contributor, len(people))
		for i, person := range people {
			hits[i] = models.ContributorFromPerson(person)
		}

		suggestURL := h.PathFor("dataset_edit_contributor_suggest", "id", ctx.Dataset.ID, "role", b.Role, "position", fmt.Sprintf("%d", nextPos)).String()

		nextActive := false
		if nextC.Person != nil {
			nextActive = nextC.Person.Active
		}

		render.Partial(w, "dataset/edit_next_contributor", YieldEditContributor{
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

	render.View(w, "dataset/refresh_contributors", YieldContributors{
		Context: ctx,
		Role:    b.Role,
	})
}

func (h *Handler) ConfirmDeleteContributor(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindDeleteContributor{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("confirm delete dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	dataset.ConfirmDeleteContributor(c, ctx.GetDataset(r), b.Role, b.Position).Render(r.Context(), w)
}

func (h *Handler) DeleteContributor(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDeleteContributor{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("delete dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	if err := ctx.Dataset.RemoveContributor(b.Role, b.Position); err != nil {
		h.Logger.Warnw("delete dataset contributor: could not remove contributor", "errors", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	if validationErrs := ctx.Dataset.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(ctx.Loc, validationErrs.(*okay.Errors)))
		render.Layout(w, "refresh_modal", "form_errors_dialog", struct {
			Title  string
			Errors form.Errors
		}{
			Title:  "Can't delete this contributor due to the following errors",
			Errors: errors,
		})

		return
	}

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("dataset.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("delete dataset contributor: Could not save the dataset:", "error", err, "dataset", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_contributors", YieldContributors{
		Context: ctx,
		Role:    b.Role,
	})
}

func (h *Handler) OrderContributors(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindOrderContributors{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("order dataset contributors: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	contributors := ctx.Dataset.Contributors(b.Role)
	if len(b.Positions) != len(contributors) {
		err := fmt.Errorf("positions don't match number of contributors")
		h.Logger.Warnw("order dataset contributors: could not order contributors", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}
	newContributors := make([]*models.Contributor, len(contributors))
	for i, pos := range b.Positions {
		newContributors[i] = contributors[pos]
	}
	ctx.Dataset.SetContributors(b.Role, newContributors)

	err := h.Repo.UpdateDataset(r.Header.Get("If-Match"), ctx.Dataset, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Loc.Get("dataset.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("order dataset contributors: Could not save the dataset:", "errors", err, "identifier", ctx.Dataset.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "dataset/refresh_contributors", YieldContributors{
		Context: ctx,
		Role:    b.Role,
	})
}

func contributorForm(ctx Context, c *models.Contributor, suggestURL string) *form.Form {
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
