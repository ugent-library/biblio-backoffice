package mutate

import (
	"errors"
	"fmt"

	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/validation"
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

func JournalTitleSet(p *models.Publication, args []string) error {
	if len(args) != 1 {
		return errors.New("journal title is missing")
	}
	if p.Type != "journal_article" {
		return errors.New("record is not of type journal_article")
	}
	p.Publication = args[0]
	return nil
}

func JournalAbbreviationSet(p *models.Publication, args []string) error {
	if len(args) != 1 {
		return errors.New("journal abbreviation is missing")
	}
	if p.Type != "journal_article" {
		return errors.New("record is not of type journal_article")
	}
	p.PublicationAbbreviation = args[0]
	return nil
}

func ISBNAdd(p *models.Publication, args []string) error {
	for _, arg := range args {
		if !validation.InArray(p.ISBN, arg) {
			p.ISBN = append(p.ISBN, arg)
		}
	}
	return nil
}

func ISBNRemove(p *models.Publication, args []string) error {
	var vals []string
	for _, val := range p.ISBN {
		if !validation.InArray(args, val) {
			vals = append(vals, val)
		}
	}
	p.ISBN = vals
	return nil
}

func EISBNAdd(p *models.Publication, args []string) error {
	for _, arg := range args {
		if !validation.InArray(p.EISBN, arg) {
			p.EISBN = append(p.EISBN, arg)
		}
	}
	return nil
}

func EISBNRemove(p *models.Publication, args []string) error {
	var vals []string
	for _, val := range p.EISBN {
		if !validation.InArray(args, val) {
			vals = append(vals, val)
		}
	}
	p.EISBN = vals
	return nil
}

func ISSNAdd(p *models.Publication, args []string) error {
	for _, arg := range args {
		if !validation.InArray(p.ISSN, arg) {
			p.ISSN = append(p.ISSN, arg)
		}
	}
	return nil
}

func ISSNRemove(p *models.Publication, args []string) error {
	var vals []string
	for _, val := range p.ISSN {
		if !validation.InArray(args, val) {
			vals = append(vals, val)
		}
	}
	p.ISSN = vals
	return nil
}

func EISSNAdd(p *models.Publication, args []string) error {
	for _, arg := range args {
		if !validation.InArray(p.EISSN, arg) {
			p.EISSN = append(p.EISSN, arg)
		}
	}
	return nil
}

func EISSNRemove(p *models.Publication, args []string) error {
	var vals []string
	for _, val := range p.EISSN {
		if !validation.InArray(args, val) {
			vals = append(vals, val)
		}
	}
	p.EISSN = vals
	return nil
}

func ExternalFieldsSet(p *models.Publication, args []string) error {
	if len(args) < 1 {
		return errors.New("no key supplied")
	}
	if len(args) < 2 {
		return fmt.Errorf("no values supplied for %s", args[0])
	}
	if p.ExternalFields == nil {
		p.ExternalFields = models.ExternalFields{}
	}
	p.ExternalFields.Set(args[0], args[1:]...)
	return nil
}
