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
	mp := &models.Project{
		ID:      p.Identifiers.Get("iweto"),
		IWETOID: p.Identifiers.Get("iweto"),
		Acronym: p.Identifiers.Get("iweto"),
		GISMOID: p.Attributes.Get("gismo", "gismo_id"),
	}

	mp.EUProject = &models.EUProject{
		ID:                 p.Attributes.Get("cordis", "eu_id"),
		CallID:             p.Attributes.Get("cordis", "eu_call_id"),
		Acronym:            p.Attributes.Get("cordis", "eu_acronym"),
		FrameworkProgramme: p.Attributes.Get("cordis", "eu_framework_programme"),
	}

	if len(p.Names) > 0 {
		mp.Title = p.Names[0].Value
	}

	if len(p.Descriptions) > 0 {
		mp.Description = p.Descriptions[0].Value
	}

	mp.StartDate = p.StartDate
	mp.EndDate = p.EndDate

	return mp
}
