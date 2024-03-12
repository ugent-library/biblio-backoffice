package backends

import (
	"context"

	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/projects"
)

type ProjectsFacade struct {
	index *projects.Index
}

func NewProjectsFacade(index *projects.Index) *ProjectsFacade {
	return &ProjectsFacade{
		index: index,
	}
}

func (f *ProjectsFacade) GetProject(id string) (*models.Project, error) {
	p, err := f.index.GetProjectByIdentifier(context.TODO(), "iweto", id)
	if err != nil {
		return nil, err
	}

	return toProject(p), nil
}

func (f *ProjectsFacade) SuggestProjects(q string) ([]*models.Project, error) {
	results, err := f.index.SearchProjects(context.TODO(), projects.SearchParams{Limit: 20, Query: q})
	if err != nil {
		return nil, err
	}

	projects := make([]*models.Project, len(results.Hits))
	for i, p := range results.Hits {
		projects[i] = toProject(p)
	}

	return projects, nil
}

func toProject(p *projects.Project) *models.Project {
	mp := &models.Project{}

	mp.EUProject = &models.EUProject{}

	for _, id := range p.Identifiers {
		if id.Kind == "iweto" {
			mp.ID = id.Value
			mp.IWETOID = id.Value
			mp.Acronym = id.Value
		}

		if id.Kind == "gismo" {
			mp.GISMOID = id.Value
		}

		if id.Kind == "cordis" {
			mp.EUProject.ID = id.Value
		}
	}

	if len(p.Names) > 0 {
		mp.Title = p.Names[0].Value
	}

	if len(p.Descriptions) > 0 {
		mp.Description = p.Descriptions[0].Value
	}

	for _, attr := range p.Attributes {
		if attr.Scope == "gismo" {
			if attr.Key == "eu_call_id" {
				mp.EUProject.CallID = attr.Value
			}

			if attr.Key == "eu_acronym" {
				mp.EUProject.Acronym = attr.Value
			}
		}
	}

	mp.StartDate = p.StartDate
	mp.EndDate = p.EndDate

	return mp
}
