package datasetediting

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	datasetviews "github.com/ugent-library/biblio-backoffice/views/dataset"
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

func AddContributor(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindAddContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("add dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	var hits []*models.Contributor

	if b.FirstName != "" || b.LastName != "" {
		people, err := c.PersonSearchService.SuggestPeople(b.FirstName + " " + b.LastName)
		if err != nil {
			c.Log.Errorw("suggest dataset contributor: could not suggest people", "errors", err, "request", r, "user", c.User.ID)
			c.HandleError(w, r, httperror.InternalServerError)
			return
		}

		hits = make([]*models.Contributor, len(people))
		for i, person := range people {
			hits[i] = models.ContributorFromPerson(person)
		}
	}

	views.ShowModal(datasetviews.AddContributor(c, datasetviews.AddContributorArgs{
		Dataset:     dataset,
		Contributor: models.ContributorFromFirstLastName(b.FirstName, b.LastName),
		Role:        b.Role,
		Hits:        hits,
	})).Render(r.Context(), w)
}

func AddContributorSuggest(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindAddContributorSuggest{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("suggest dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	var hits []*models.Contributor

	if b.FirstName != "" || b.LastName != "" {
		people, err := c.PersonSearchService.SuggestPeople(b.FirstName + " " + b.LastName)
		if err != nil {
			c.Log.Errorw("suggest dataset contributor: could not suggest people", "errors", err, "request", r, "user", c.User.ID)
			c.HandleError(w, r, httperror.InternalServerError)
			return
		}

		hits = make([]*models.Contributor, len(people))
		for i, person := range people {
			hits[i] = models.ContributorFromPerson(person)
		}
	}

	datasetviews.AddContributorSuggest(
		c, datasetviews.AddContributorSuggestArgs{
			Dataset:     dataset,
			Contributor: models.ContributorFromFirstLastName(b.FirstName, b.LastName),
			Role:        b.Role,
			Hits:        hits,
		},
	).Render(r.Context(), w)
}

func ConfirmCreateContributor(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindConfirmCreateContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("confirm create dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	var cn *models.Contributor
	if b.ID != "" {
		newC, err := generateContributorFromPersonId(c, b.ID)
		if err != nil {
			c.HandleError(w, r, err)
			return
		}
		cn = newC
	} else {
		cn = models.ContributorFromFirstLastName(b.FirstName, b.LastName)
	}

	views.ReplaceModal(datasetviews.ConfirmCreateContributor(
		c, datasetviews.ConfirmCreateContributorArgs{
			Dataset:     dataset,
			Contributor: cn,
			Role:        b.Role,
		},
	)).Render(r.Context(), w)
}

func CreateContributor(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindCreateContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("create dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	var cn *models.Contributor
	if b.ID != "" {
		newC, err := generateContributorFromPersonId(c, b.ID)
		if err != nil {
			c.HandleError(w, r, err)
			return
		}
		cn = newC
	} else {
		if b.FirstName == "" {
			b.FirstName = "[missing]"
		}
		if b.LastName == "" {
			b.LastName = "[missing]"
		}
		cn = models.ContributorFromFirstLastName(b.FirstName, b.LastName)
	}

	dataset.AddContributor(b.Role, cn)

	if validationErrs := dataset.Validate(); validationErrs != nil {
		views.ReplaceModal(datasetviews.ConfirmCreateContributor(
			c, datasetviews.ConfirmCreateContributorArgs{
				Dataset:     dataset,
				Contributor: cn,
				Role:        b.Role,
				Errors:      validationErrs.(*okay.Errors),
			},
		)).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), dataset, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("create dataset contributor: Could not save the dataset:", "errors", err, "dataset", dataset.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	if b.AddNext {
		views.Cat(
			views.Replace(fmt.Sprintf("#contributors-%s-body", b.Role), datasetviews.ContributorsBody(
				c, dataset, b.Role,
			)),
			views.ReplaceModal(datasetviews.AddContributor(c, datasetviews.AddContributorArgs{
				Dataset:     dataset,
				Role:        b.Role,
				Contributor: &models.Contributor{},
			})),
		).Render(r.Context(), w)
		return
	}

	views.CloseModalAndReplace(
		fmt.Sprintf("#contributors-%s-body", b.Role),
		datasetviews.ContributorsBody(c, dataset, b.Role),
	).Render(r.Context(), w)
}

func EditContributor(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindEditContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("edit dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	cn, err := dataset.GetContributor(b.Role, b.Position)
	if err != nil {
		c.Log.Errorw("edit dataset contributor: could not get the contributor", "errors", err, "dataset", dataset.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	firstName := b.FirstName
	lastName := b.LastName
	if firstName == "" && lastName == "" {
		firstName = cn.FirstName()
		lastName = cn.LastName()
	}

	people, err := c.PersonSearchService.SuggestPeople(firstName + " " + lastName)
	if err != nil {
		c.Log.Errorw("suggest dataset contributor: could not suggest people", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	hits := make([]*models.Contributor, len(people))
	for i, person := range people {
		hits[i] = models.ContributorFromPerson(person)
	}

	// exclude the current contributor
	if cn.PersonID != "" {
		for i, hit := range hits {
			if hit.PersonID == cn.PersonID {
				if i == 0 {
					hits = hits[1:]
				} else {
					hits = append(hits[:i], hits[i+1:]...)
				}
				break
			}
		}
	}

	views.ShowModal(datasetviews.EditContributor(c, datasetviews.EditContributorArgs{
		Dataset:     dataset,
		Role:        b.Role,
		Position:    b.Position,
		Contributor: cn,
		FirstName:   firstName,
		LastName:    lastName,
		Hits:        hits,
	})).Render(r.Context(), w)
}

func EditContributorSuggest(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindEditContributorSuggest{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("suggest dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	cn, err := dataset.GetContributor(b.Role, b.Position)
	if err != nil {
		c.Log.Errorw("edit dataset contributor: could not get the contributor", "errors", err, "dataset", dataset.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	var hits []*models.Contributor

	if b.FirstName != "" || b.LastName != "" {
		people, err := c.PersonSearchService.SuggestPeople(b.FirstName + " " + b.LastName)
		if err != nil {
			c.Log.Errorw("suggest dataset contributor: could not suggest people", "errors", err, "request", r, "user", c.User.ID)
			c.HandleError(w, r, httperror.InternalServerError)
			return
		}

		hits = make([]*models.Contributor, len(people))
		for i, person := range people {
			hits[i] = models.ContributorFromPerson(person)
		}

		// exclude the current contributor
		if cn.PersonID != "" {
			for i, hit := range hits {
				if hit.PersonID == cn.PersonID {
					hits = append(hits[:i], hits[i+1:]...)
					break
				}
			}
		}
	}

	datasetviews.EditContributorSuggest(c, datasetviews.EditContributorSuggestArgs{
		Dataset:     dataset,
		Role:        b.Role,
		Position:    b.Position,
		Contributor: cn,
		FirstName:   b.FirstName,
		LastName:    b.LastName,
		Hits:        hits,
	}).Render(r.Context(), w)
}

func ConfirmUpdateContributor(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindConfirmUpdateContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("confirm update dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	var cn *models.Contributor
	if b.ID != "" {
		newC, err := generateContributorFromPersonId(c, b.ID)
		if err != nil {
			c.HandleError(w, r, err)
			return
		}
		cn = newC
	} else {
		cn = models.ContributorFromFirstLastName(b.FirstName, b.LastName)
	}

	views.ReplaceModal(datasetviews.ConfirmUpdateContributor(c, datasetviews.ConfirmUpdateContributorArgs{
		Dataset:     dataset,
		Role:        b.Role,
		Position:    b.Position,
		Contributor: cn,
		EditNext:    b.Position+1 < len(dataset.Contributors(b.Role)),
	})).Render(r.Context(), w)
}

func UpdateContributor(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindUpdateContributor{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("update dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	var cn *models.Contributor
	if b.ID != "" {
		newC, err := generateContributorFromPersonId(c, b.ID)
		if err != nil {
			c.Log.Errorw("update dataset contributor: could not fetch person", "errors", err, "personid", b.ID, "dataset", dataset.ID, "user", c.User.ID)
			c.HandleError(w, r, httperror.InternalServerError)
			return
		}
		cn = newC
	} else {
		if b.FirstName == "" {
			b.FirstName = "[missing]"
		}
		if b.LastName == "" {
			b.LastName = "[missing]"
		}
		cn = models.ContributorFromFirstLastName(b.FirstName, b.LastName)
	}

	if err := dataset.SetContributor(b.Role, b.Position, cn); err != nil {
		c.Log.Errorw("update dataset contributor: could not set the contributor", "errors", err, "dataset", dataset.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	if validationErrs := dataset.Validate(); validationErrs != nil {
		views.ReplaceModal(datasetviews.ConfirmUpdateContributor(c, datasetviews.ConfirmUpdateContributorArgs{
			Dataset:     dataset,
			Role:        b.Role,
			Position:    b.Position,
			Contributor: cn,
			Errors:      validationErrs.(*okay.Errors),
			EditNext:    b.Position+1 < len(dataset.Contributors(b.Role)),
		})).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), dataset, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("update dataset contributor: Could not save the dataset:", "errors", err, "dataset", dataset.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	if b.EditNext && b.Position+1 < len(dataset.Contributors(b.Role)) {
		nextPos := b.Position + 1
		nextC := dataset.Contributors(b.Role)[nextPos]
		people, err := c.PersonSearchService.SuggestPeople(nextC.Name())
		if err != nil {
			c.Log.Errorw("suggest dataset contributor: could not suggest people", "errors", err, "request", r, "user", c.User.ID)
			c.HandleError(w, r, httperror.InternalServerError)
			return
		}
		hits := make([]*models.Contributor, len(people))
		for i, person := range people {
			hits[i] = models.ContributorFromPerson(person)
		}

		views.Cat(
			views.Replace(fmt.Sprintf("#contributors-%s-body", b.Role), datasetviews.ContributorsBody(c, dataset, b.Role)),
			views.ReplaceModal(datasetviews.EditContributor(c, datasetviews.EditContributorArgs{
				Dataset:     dataset,
				Role:        b.Role,
				Position:    nextPos,
				Contributor: nextC,
				FirstName:   nextC.FirstName(),
				LastName:    nextC.LastName(),
				Hits:        hits,
			})),
		).Render(r.Context(), w)
		return
	}

	views.CloseModalAndReplace(fmt.Sprintf("#contributors-%s-body", b.Role), datasetviews.ContributorsBody(
		c, dataset, b.Role,
	)).Render(r.Context(), w)
}

func ConfirmDeleteContributor(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindDeleteContributor{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   fmt.Sprintf("Are you sure you want to remove this %s?", c.Loc.Get("dataset.contributor.role."+b.Role)),
		DeleteUrl:  c.PathTo("dataset_delete_contributor", "id", dataset.ID, "role", b.Role, "position", strconv.Itoa(b.Position)),
		SnapshotID: dataset.SnapshotID,
	}).Render(r.Context(), w)
}

func DeleteContributor(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindDeleteContributor{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("delete dataset contributor: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	if err := dataset.RemoveContributor(b.Role, b.Position); err != nil {
		c.Log.Warnw("delete dataset contributor: could not remove contributor", "errors", err, "dataset", dataset.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	if validationErrs := dataset.Validate(); validationErrs != nil {
		errors := localize.ValidationErrors(c.Loc, validationErrs.(*okay.Errors))
		views.ReplaceModal(views.FormErrorsDialog("Can't delete this contributor due to the following errors", errors)).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), dataset, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("delete dataset contributor: Could not save the dataset:", "error", err, "dataset", dataset.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(
		fmt.Sprintf("#contributors-%s-body", b.Role),
		datasetviews.ContributorsBody(c, dataset, b.Role),
	).Render(r.Context(), w)
}

func OrderContributors(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindOrderContributors{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("order dataset contributors: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	contributors := dataset.Contributors(b.Role)
	if len(b.Positions) != len(contributors) {
		err := fmt.Errorf("positions don't match number of contributors")
		c.Log.Warnw("order dataset contributors: could not order contributors", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}
	newContributors := make([]*models.Contributor, len(contributors))
	for i, pos := range b.Positions {
		newContributors[i] = contributors[pos]
	}
	dataset.SetContributors(b.Role, newContributors)

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), dataset, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("order dataset contributors: Could not save the dataset:", "errors", err, "identifier", dataset.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(fmt.Sprintf("#contributors-%s-body", b.Role), datasetviews.ContributorsBody(c, dataset, b.Role)).Render(r.Context(), w)
}

func generateContributorFromPersonId(c *ctx.Ctx, id string) (*models.Contributor, error) {
	p, err := c.PersonService.GetPerson(id)
	if err != nil {
		return nil, err
	}
	cn := models.ContributorFromPerson(p)
	return cn, nil
}
