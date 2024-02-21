package api

import (
	"context"

	"github.com/ugent-library/biblio-backoffice/projects"
)

type Service struct {
	repo *projects.Repo
}

func NewService(repo *projects.Repo) *Service {
	return &Service{
		repo: repo,
	}
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

	return s.repo.AddProject(ctx, &projects.Project{
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
