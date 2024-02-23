package api

import (
	"context"

	"github.com/ugent-library/biblio-backoffice/people"
	"github.com/ugent-library/biblio-backoffice/projects"
)

type Service struct {
	peopleRepo   *people.Repo
	projectsRepo *projects.Repo
}

func NewService(
	peopleRepo *people.Repo,
	projectsRepo *projects.Repo,
) *Service {
	return &Service{
		peopleRepo:   peopleRepo,
		projectsRepo: projectsRepo,
	}
}

func (s *Service) AddPerson(ctx context.Context, req *AddPersonRequest) error {
	p := req.Person

	identifiers := make([]people.Identifier, len(p.Identifiers))
	for i, id := range p.Identifiers {
		identifiers[i] = people.Identifier(id)
	}

	attributes := make([]people.Attribute, len(p.Attributes))
	for i, attr := range p.Attributes {
		attributes[i] = people.Attribute(attr)
	}

	return s.peopleRepo.AddPerson(ctx, people.AddPersonParams{
		Identifiers:         identifiers,
		Name:                p.Name,
		PreferredName:       p.PreferredName.Value,
		GivenName:           p.GivenName.Value,
		PreferredGivenName:  p.PreferredGivenName.Value,
		FamilyName:          p.FamilyName.Value,
		PreferredFamilyName: p.PreferredFamilyName.Value,
		HonorificPrefix:     p.HonorificPrefix.Value,
		Email:               p.Email.Value,
		// Active:              p.Active.Value,
		// Username:            p.Username.Value,
		Attributes: attributes,
	})
}

func (s *Service) AddProject(ctx context.Context, p *Project) error {
	attributes := make([]projects.Attribute, len(p.Attributes))
	for i, attr := range p.Attributes {
		attributes[i] = projects.Attribute(attr)
	}

	identifiers := make([]projects.Identifier, len(p.Identifiers))
	for i, id := range p.Identifiers {
		identifiers[i] = projects.Identifier(id)
	}

	names := make([]projects.Text, len(p.Name))
	for i, name := range p.Name {
		names[i] = projects.Text(name)
	}

	descriptions := make([]projects.Text, len(p.Description))
	for i, desc := range p.Description {
		descriptions[i] = projects.Text(desc)
	}

	foundingDate := ""
	if v, ok := p.GetFoundingDate().Get(); ok {
		foundingDate = v
	}

	dissolutionDate := ""
	if v, ok := p.GetDissolutionDate().Get(); ok {
		dissolutionDate = v
	}

	return s.projectsRepo.AddProject(ctx, &projects.Project{
		Names:           names,
		Descriptions:    descriptions,
		FoundingDate:    foundingDate,
		DissolutionDate: dissolutionDate,
		Attributes:      attributes,
		Identifiers:     identifiers,
	})
}

func (s *Service) NewError(ctx context.Context, err error) *ErrorStatusCode {
	return &ErrorStatusCode{
		StatusCode: 500,
		Response: Error{
			Code:    500,
			Message: err.Error(),
		},
	}
}
