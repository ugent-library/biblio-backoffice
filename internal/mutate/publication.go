package mutate

import (
	"errors"

	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
)

func ProjectAdd(projectService backends.ProjectService) func(*models.Publication, []string) error {
	return func(p *models.Publication, args []string) error {
		if len(args) != 1 {
			return errors.New("project id is missing")
		}
		project, err := projectService.GetProject(args[0])
		if err != nil {
			return err
		}
		p.AddProject(&models.PublicationProject{
			ID:   project.ID,
			Name: project.Title,
		})
		return nil
	}
}

func ClassificationSet(p *models.Publication, args []string) error {
	if len(args) != 1 {
		return errors.New("classification is missing")
	}
	p.Classification = args[0]
	return nil
}

func KeywordAdd(p *models.Publication, args []string) error {
	for _, arg := range args {
		if !validation.InArray(p.Keyword, arg) {
			p.Keyword = append(p.Keyword, arg)
		}
	}
	return nil
}

func VABBYearAdd(p *models.Publication, args []string) error {
	for _, arg := range args {
		if !validation.InArray(p.VABBYear, arg) {
			p.VABBYear = append(p.VABBYear, arg)
		}
	}
	return nil
}
