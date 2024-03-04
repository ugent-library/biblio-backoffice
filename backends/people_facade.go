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
	ctx := context.TODO()
	p, err := f.repo.GetPersonByIdentifier(ctx, "id", id)
	if err != nil {
		return nil, err
	}

	return toPerson(p), nil
}

// TODO affiliations
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

	return mp
}
