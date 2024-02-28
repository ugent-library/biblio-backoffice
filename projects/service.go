package projects

import (
	"context"

	"github.com/ugent-library/biblio-backoffice/models"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetProject(id string) (*models.Project, error) {
	ctx := context.Background()
	rp, err := s.repo.GetProjectByIdentifier(ctx, "iweto", id)
	if err != nil {
		return nil, err
	}

	p := toProject(rp)

	return &p, nil
}

type SearchService struct {
	index *Index
}

func NewSearchService(index *Index) *SearchService {
	return &SearchService{
		index: index,
	}
}

func (s *SearchService) SuggestProjects(qs string) ([]models.Project, error) {
	ctx := context.Background()

	hits, err := s.index.SearchProjects(ctx, qs)
	if err != nil {
		return nil, err
	}

	projects := make([]models.Project, len(hits))
	for i, h := range hits {
		projects[i] = toProject(h)
	}

	return projects, nil
}

func toProject(rp *Project) models.Project {
	p := models.Project{}

	p.EUProject = &models.EUProject{}

	for _, id := range rp.Identifiers {
		if id.Kind == "iweto" {
			p.ID = id.Value
			p.IWETOID = id.Value
			p.Acronym = id.Value
		}

		if id.Kind == "gismo" {
			p.GISMOID = id.Value
		}

		if id.Kind == "cordis" {
			p.EUProject.ID = id.Value
		}
	}

	for _, name := range rp.Names {
		if name.Lang == "und" {
			p.Title = name.Value
		}
	}

	for _, desc := range rp.Descriptions {
		if desc.Lang == "und" {
			p.Description = desc.Value
		}
	}

	for _, attr := range rp.Attributes {
		if attr.Scope == "gismo" {
			if attr.Key == "eu_call_id" {
				p.EUProject.CallID = attr.Key
			}

			if attr.Key == "eu_acronym" {
				p.EUProject.Acronym = attr.Value
			}
		}
	}

	p.StartDate = rp.FoundingDate
	p.EndDate = rp.DissolutionDate

	return p
}
