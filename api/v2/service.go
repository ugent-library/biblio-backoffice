package api

import (
	"context"

	"github.com/samber/lo"
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

func (s *Service) ImportOrganizations(ctx context.Context, req *ImportOrganizationsRequest) error {
	iter := func(ctx context.Context, fn func(people.ImportOrganizationParams) bool) error {
		for _, params := range req.Organizations {
			if !fn(convertImportOrganizationParams(params)) {
				break
			}
		}
		return nil
	}

	return s.peopleRepo.ImportOrganizations(ctx, iter)
}

func (s *Service) ImportPerson(ctx context.Context, req *ImportPersonRequest) error {
	return nil
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
		Active:              p.Active.Value,
		Username:            p.Username.Value,
		Attributes:          attributes,
	})
}

func (s *Service) AddProject(ctx context.Context, req *AddProjectRequest) error {
	p := req.Project

	identifiers := make([]projects.Identifier, len(p.Identifiers))
	for i, id := range p.Identifiers {
		identifiers[i] = projects.Identifier(id)
	}

	attributes := make([]projects.Attribute, len(p.Attributes))
	for i, attr := range p.Attributes {
		attributes[i] = projects.Attribute(attr)
	}

	names := make([]projects.Text, len(p.Names))
	for i, name := range p.Names {
		names[i] = projects.Text(name)
	}

	descriptions := make([]projects.Text, len(p.Descriptions))
	for i, desc := range p.Descriptions {
		descriptions[i] = projects.Text(desc)
	}

	startDate := ""
	if v, ok := p.GetStartDate().Get(); ok {
		startDate = v
	}

	endDate := ""
	if v, ok := p.GetEndDate().Get(); ok {
		endDate = v
	}

	return s.projectsRepo.AddProject(ctx, projects.AddProjectParams{
		Names:        names,
		Descriptions: descriptions,
		StartDate:    startDate,
		EndDate:      endDate,
		Attributes:   attributes,
		Deleted:      p.Deleted.Value,
		Identifiers:  identifiers,
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

func convertImportOrganizationParams(from ImportOrganizationParams) people.ImportOrganizationParams {
	return people.ImportOrganizationParams{
		Identifiers:      lo.Map(from.Identifiers, func(v Identifier, _ int) people.Identifier { return people.Identifier(v) }),
		ParentIdentifier: lo.Ternary(from.ParentIdentifier.Set, lo.ToPtr(people.Identifier(from.ParentIdentifier.Value)), nil),
		Names:            lo.Map(from.Names, func(v Text, _ int) people.Text { return people.Text(v) }),
		Ceased:           from.Ceased.Value,
		CreatedAt:        lo.Ternary(from.CreatedAt.Set, &from.CreatedAt.Value, nil),
		UpdatedAt:        lo.Ternary(from.UpdatedAt.Set, &from.UpdatedAt.Value, nil),
	}
}
