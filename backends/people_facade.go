package backends

import (
	"context"

	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/people"
)

type PeopleFacade struct {
	repo  *people.Repo
	index *people.Index
}

func NewPeopleFacade(repo *people.Repo, index *people.Index) *PeopleFacade {
	return &PeopleFacade{
		repo:  repo,
		index: index,
	}
}

func (f *PeopleFacade) GetPerson(id string) (*models.Person, error) {
	p, err := f.repo.GetPersonByIdentifier(context.TODO(), "id", id)
	if err != nil {
		return nil, err
	}

	return toPerson(p), nil
}

func (f *PeopleFacade) GetUserByUsername(username string) (*models.Person, error) {
	p, err := f.repo.GetActivePersonByUsername(context.TODO(), username)
	if err != nil {
		return nil, err
	}

	return toPerson(p), nil
}

func (f *PeopleFacade) GetUser(id string) (*models.Person, error) {
	p, err := f.repo.GetActivePersonByIdentifier(context.TODO(), "id", id)
	if err != nil {
		return nil, err
	}

	return toPerson(p), nil
}

func (f *PeopleFacade) GetOrganization(id string) (*models.Organization, error) {
	o, err := f.repo.GetOrganizationByIdentifier(context.TODO(), "biblio", id)
	if err != nil {
		return nil, err
	}
	return toOrganization(o), nil
}

func (f *PeopleFacade) SuggestPeople(qs string) ([]*models.Person, error) {
	hits, err := f.index.SearchPeople(context.TODO(), qs)
	if err != nil {
		return nil, err
	}

	people := make([]*models.Person, len(hits))
	for i, p := range hits {
		people[i] = toPerson(p)
	}

	return people, nil
}

// TODO filter out inactive people in the index
func (f *PeopleFacade) SuggestUsers(qs string) ([]*models.Person, error) {
	hits, err := f.index.SearchPeople(context.TODO(), qs)
	if err != nil {
		return nil, err
	}

	people := make([]*models.Person, 0, len(hits))
	for _, p := range hits {
		if p.Active {
			people = append(people, toPerson(p))
		}
	}

	return people, nil
}

func (f *PeopleFacade) SuggestOrganizations(qs string) ([]models.Completion, error) {
	hits, err := f.index.SearchOrganizations(context.TODO(), qs)
	if err != nil {
		return nil, err
	}

	orgs := make([]models.Completion, len(hits))
	for i, o := range hits {
		orgs[i] = models.Completion{
			ID:      o.Identifiers.Get("biblio"),
			Heading: o.Names.Get("eng"),
		}
	}

	return orgs, nil
}

func toPerson(p *people.Person) *models.Person {
	mp := &models.Person{
		ID:          p.Identifiers.Get("id"),
		Active:      p.Active,
		DateCreated: &p.CreatedAt,
		DateUpdated: &p.UpdatedAt,
		Email:       p.Email,
		FullName:    p.Name,
		FirstName:   p.GivenName,
		LastName:    p.FamilyName,
		ORCID:       p.Identifiers.Get("orcid"),
		UGentID:     p.Identifiers.GetAll("ugentID"),
		Username:    p.Username,
		Role:        p.Role,
	}

	if p.PreferredName != "" {
		mp.FullName = p.PreferredName
	}
	if p.PreferredGivenName != "" {
		mp.FirstName = p.PreferredGivenName
	}
	if p.PreferredFamilyName != "" {
		mp.LastName = p.PreferredFamilyName
	}

	for _, token := range p.Tokens {
		if token.Kind == "orcid" {
			mp.ORCIDToken = string(token.Value)
		}
	}

	for _, a := range p.Affiliations {
		mo := toOrganization(a.Organization)
		mp.Affiliations = append(mp.Affiliations, &models.Affiliation{
			OrganizationID: mo.ID,
			Organization:   mo,
		})
	}

	return mp
}

func toOrganization(o *people.Organization) *models.Organization {
	id := o.Identifiers.Get("biblio")
	tree := []models.OrganizationTreeElement{{ID: id}}
	for _, po := range o.Parents {
		tree = append(tree, models.OrganizationTreeElement{ID: po.Identifiers.Get("biblio")})
	}
	return &models.Organization{
		ID:   id,
		Name: o.Names.Get("eng"),
		Tree: tree,
	}
}
