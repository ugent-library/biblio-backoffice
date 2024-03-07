package api

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/samber/lo"
	"github.com/ugent-library/biblio-backoffice/people"
	"github.com/ugent-library/biblio-backoffice/projects"
)

type Service struct {
	peopleRepo   *people.Repo
	peopleIndex  *people.Index
	projectsRepo *projects.Repo
}

func NewService(
	peopleRepo *people.Repo,
	peopleIndex *people.Index,
	projectsRepo *projects.Repo,
) *Service {
	return &Service{
		peopleRepo:   peopleRepo,
		peopleIndex:  peopleIndex,
		projectsRepo: projectsRepo,
	}
}

func (s *Service) GetOrganization(ctx context.Context, req *GetOrganizationRequest) (GetOrganizationRes, error) {
	o, err := s.peopleIndex.GetOrganizationByIdentifier(ctx, req.Identifier.Kind, req.Identifier.Value)
	if errors.Is(err, people.ErrNotFound) {
		return nil, &ErrorStatusCode{
			StatusCode: 404,
			Response: Error{
				Code:    404,
				Message: "Organization not found",
			},
		}
	}
	if err != nil {
		return nil, err
	}

	return &GetOrganization{Organization: convertOrganization(o)}, nil
}

// TODO use index
func (s *Service) GetPerson(ctx context.Context, req *GetPersonRequest) (GetPersonRes, error) {
	p, err := s.peopleIndex.GetPersonByIdentifier(ctx, req.Identifier.Kind, req.Identifier.Value)
	if errors.Is(err, people.ErrNotFound) {
		return nil, &ErrorStatusCode{
			StatusCode: 404,
			Response: Error{
				Code:    404,
				Message: "Person not found",
			},
		}
	}
	if err != nil {
		return nil, err
	}

	return &GetPerson{Person: convertPerson(p)}, nil
}

func (s *Service) SearchOrganizations(ctx context.Context, req *SearchOrganizationsRequest) (*SearchOrganizations, error) {
	hits, err := s.peopleIndex.SearchOrganizations(ctx, req.Query.Value)
	if err != nil {
		return nil, err
	}

	return &SearchOrganizations{
		Hits: lo.Map(hits, func(v *people.Organization, _ int) Organization { return convertOrganization(v) }),
	}, nil
}

func (s *Service) SearchPeople(ctx context.Context, req *SearchPeopleRequest) (*SearchPeople, error) {
	hits, err := s.peopleIndex.SearchPeople(ctx, req.Query.Value)
	if err != nil {
		return nil, err
	}

	return &SearchPeople{
		Hits: lo.Map(hits, func(v *people.Person, _ int) Person { return convertPerson(v) }),
	}, nil
}

func (s *Service) ImportOrganizations(ctx context.Context, req *ImportOrganizationsRequest) (ImportOrganizationsRes, error) {
	iter := func(ctx context.Context, fn func(people.ImportOrganizationParams) bool) error {
		for _, params := range req.Organizations {
			if !fn(convertImportOrganizationParams(params)) {
				break
			}
		}
		return nil
	}

	err := s.peopleRepo.ImportOrganizations(ctx, iter)

	var dupErr *people.DuplicateError
	if errors.As(err, &dupErr) {
		return nil, &ErrorStatusCode{
			StatusCode: 409,
			Response: Error{
				Code:    409,
				Message: dupErr.Error(),
			},
		}
	}

	if err != nil {
		return nil, err
	}

	return &ImportOrganizationsOK{}, nil
}

func (s *Service) ImportPerson(ctx context.Context, req *ImportPersonRequest) (ImportPersonRes, error) {
	err := s.peopleRepo.ImportPerson(ctx, convertImportPersonParams(req.Person))

	var dupErr *people.DuplicateError
	if errors.As(err, &dupErr) {
		return nil, &ErrorStatusCode{
			StatusCode: 409,
			Response: Error{
				Code:    409,
				Message: dupErr.Error(),
			},
		}
	}

	if err != nil {
		return nil, err
	}

	return &ImportPersonOK{}, nil
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

func (s *Service) ImportProject(ctx context.Context, req *ImportProjectRequest) (ImportProjectRes, error) {
	err := s.projectsRepo.ImportProject(ctx, convertImportProjectParams(req.Project))

	var dupErr *projects.DuplicateProjectError
	if errors.As(err, &dupErr) {
		return nil, &ErrorStatusCode{
			StatusCode: 409,
			Response: Error{
				Code:    409,
				Message: dupErr.Error(),
			},
		}
	}

	if err != nil {
		return nil, err
	}

	return &ImportProjectOK{}, nil
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
		Identifiers:  identifiers,
		Names:        names,
		Descriptions: descriptions,
		StartDate:    startDate,
		EndDate:      endDate,
		Deleted:      p.Deleted.Value,
		Attributes:   attributes,
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

func convertOrganization(from *people.Organization) Organization {
	return Organization{
		Identifiers: lo.Map(from.Identifiers, func(v people.Identifier, _ int) Identifier { return Identifier(v) }),
		Names:       lo.Map(from.Names, func(v people.Text, _ int) Text { return Text(v) }),
		Ceased:      from.Ceased,
		Parents:     lo.Map(from.Parents, func(v people.ParentOrganization, _ int) ParentOrganization { return convertParentOrganization(v) }),
		CreatedAt:   from.CreatedAt,
		UpdatedAt:   from.UpdatedAt,
	}
}

func convertPerson(from *people.Person) Person {
	return Person{
		Identifiers:         lo.Map(from.Identifiers, func(v people.Identifier, _ int) Identifier { return Identifier(v) }),
		Name:                from.Name,
		PreferredName:       OptString{Set: from.PreferredName != "", Value: from.PreferredName},
		GivenName:           OptString{Set: from.GivenName != "", Value: from.GivenName},
		PreferredGivenName:  OptString{Set: from.PreferredGivenName != "", Value: from.PreferredGivenName},
		FamilyName:          OptString{Set: from.FamilyName != "", Value: from.FamilyName},
		PreferredFamilyName: OptString{Set: from.PreferredFamilyName != "", Value: from.PreferredFamilyName},
		HonorificPrefix:     OptString{Set: from.HonorificPrefix != "", Value: from.HonorificPrefix},
		Email:               OptString{Set: from.Email != "", Value: from.Email},
		Active:              from.Active,
		Role:                OptString{Set: from.Role != "", Value: from.Role},
		Username:            OptString{Set: from.Username != "", Value: from.Username},
		Attributes:          lo.Map(from.Attributes, func(v people.Attribute, _ int) Attribute { return Attribute(v) }),
		Tokens:              lo.Map(from.Tokens, func(v people.Token, _ int) Token { return Token(v) }),
		Affiliations: lo.Map(from.Affiliations, func(a people.Affiliation, _ int) Affiliation {
			return Affiliation{Organization: convertOrganization(a.Organization)}
		}),
		CreatedAt: from.CreatedAt,
		UpdatedAt: from.UpdatedAt,
	}
}

func convertParentOrganization(from people.ParentOrganization) ParentOrganization {
	return ParentOrganization{
		Identifiers: lo.Map(from.Identifiers, func(v people.Identifier, _ int) Identifier { return Identifier(v) }),
		Names:       lo.Map(from.Names, func(v people.Text, _ int) Text { return Text(v) }),
		Ceased:      from.Ceased,
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

func convertImportProjectParams(from ImportProjectParams) projects.ImportProjectParams {
	return projects.ImportProjectParams{
		Identifiers:  lo.Map(from.Identifiers, func(v Identifier, _ int) projects.Identifier { return projects.Identifier(v) }),
		Names:        lo.Map(from.Names, func(v Text, _ int) projects.Text { return projects.Text(v) }),
		Descriptions: lo.Map(from.Descriptions, func(v Text, _ int) projects.Text { return projects.Text(v) }),
		StartDate:    lo.Ternary(from.StartDate.Set, from.StartDate.Value, ""),
		EndDate:      lo.Ternary(from.EndDate.Set, from.EndDate.Value, ""),
		Deleted:      lo.Ternary(from.Deleted.Set, from.Deleted.Value, false),
		Attributes:   lo.Map(from.Attributes, func(v Attribute, _ int) projects.Attribute { return projects.Attribute(v) }),
		CreatedAt:    lo.Ternary(from.CreatedAt.Set, &from.CreatedAt.Value, nil),
		UpdatedAt:    lo.Ternary(from.UpdatedAt.Set, &from.UpdatedAt.Value, nil),
	}
}

func convertImportPersonParams(from ImportPersonParams) people.ImportPersonParams {
	return people.ImportPersonParams{
		Identifiers:         lo.Map(from.Identifiers, func(v Identifier, _ int) people.Identifier { return people.Identifier(v) }),
		Name:                from.Name,
		PreferredName:       from.PreferredName.Value,
		GivenName:           from.GivenName.Value,
		PreferredGivenName:  from.PreferredGivenName.Value,
		FamilyName:          from.FamilyName.Value,
		PreferredFamilyName: from.PreferredFamilyName.Value,
		HonorificPrefix:     from.HonorificPrefix.Value,
		Email:               from.Email.Value,
		Active:              from.Active.Value,
		Role:                from.Role.Value,
		Username:            from.Username.Value,
		Attributes:          lo.Map(from.Attributes, func(v Attribute, _ int) people.Attribute { return people.Attribute(v) }),
		Tokens:              lo.Map(from.Tokens, func(v Token, _ int) people.Token { return people.Token(v) }),
		Affiliations: lo.Map(from.Affiliations, func(v AffiliationParams, _ int) people.AffiliationParams {
			return people.AffiliationParams{OrganizationIdentifier: people.Identifier(v.OrganizationIdentifier)}
		}),
		CreatedAt: lo.Ternary(from.CreatedAt.Set, &from.CreatedAt.Value, nil),
		UpdatedAt: lo.Ternary(from.UpdatedAt.Set, &from.UpdatedAt.Value, nil),
	}
}
