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
		p.AddProject(project)
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

func KeywordRemove(p *models.Publication, args []string) error {
	var vals []string
	for _, val := range p.Keyword {
		if !validation.InArray(args, val) {
			vals = append(vals, val)
		}
	}
	p.Keyword = vals
	return nil
}

func VABBTypeSet(p *models.Publication, args []string) error {
	if len(args) != 1 {
		return errors.New("vabb type is missing")
	}
	p.VABBType = args[0]
	return nil
}

func VABBIDSet(p *models.Publication, args []string) error {
	if len(args) != 1 {
		return errors.New("vabb id is missing")
	}
	p.VABBID = args[0]
	return nil
}

func VABBApprovedSet(p *models.Publication, args []string) error {
	if len(args) != 1 {
		return errors.New("vabb approved value must be 'true' or 'false'")
	}
	switch args[0] {
	case "true":
		p.VABBApproved = true
	case "false":
		p.VABBApproved = false
	default:
		return errors.New("vabb approved value must be 'true' or 'false'")
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

func ReviewerTagAdd(p *models.Publication, args []string) error {
	for _, arg := range args {
		if !validation.InArray(p.ReviewerTags, arg) {
			p.ReviewerTags = append(p.ReviewerTags, arg)
		}
	}
	return nil
}

func ReviewerTagRemove(p *models.Publication, args []string) error {
	var vals []string
	for _, val := range p.ReviewerTags {
		if !validation.InArray(args, val) {
			vals = append(vals, val)
		}
	}
	p.ReviewerTags = vals
	return nil
}

func PublicationSet(p *models.Publication, args []string) error {
	if len(args) != 1 {
		return errors.New("journal title is missing")
	}
	p.Publication = args[0]
	return nil
}
