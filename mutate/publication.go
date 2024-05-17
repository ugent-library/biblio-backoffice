package mutate

import (
	"slices"

	"github.com/ugent-library/biblio-backoffice/models"
)

type ArgumentError struct {
	Msg string
}

func (e *ArgumentError) Error() string {
	return e.Msg
}

type ProjectGetter func(string) (*models.Project, error)

func AddProject(projectGetter ProjectGetter) func(*models.Publication, []string) error {
	return func(p *models.Publication, args []string) error {
		if !p.UsesProject() {
			return &ArgumentError{"project not used for this publication type"}
		}
		if len(args) != 1 {
			return &ArgumentError{"project id is missing"}
		}
		project, err := projectGetter(args[0])
		if err != nil {
			return err
		}
		p.AddProject(project)
		return nil
	}
}

func RemoveProject(p *models.Publication, args []string) error {
	if !p.UsesProject() {
		return &ArgumentError{"project not used for this publication type"}
	}
	if len(args) != 1 {
		return &ArgumentError{"project id is missing"}
	}
	p.RemoveProject(args[0])
	return nil
}

func SetClassification(p *models.Publication, args []string) error {
	if len(args) != 1 {
		return &ArgumentError{"classification is missing"}
	}
	p.Classification = args[0]
	return nil
}

func AddKeyword(p *models.Publication, args []string) error {
	for _, arg := range args {
		if !slices.Contains(p.Keyword, arg) {
			p.Keyword = append(p.Keyword, arg)
		}
	}
	return nil
}

func RemoveKeyword(p *models.Publication, args []string) error {
	var vals []string
	for _, val := range p.Keyword {
		if !slices.Contains(args, val) {
			vals = append(vals, val)
		}
	}
	p.Keyword = vals
	return nil
}

func SetVABBType(p *models.Publication, args []string) error {
	if len(args) != 1 {
		return &ArgumentError{"vabb type is missing"}
	}
	p.VABBType = args[0]
	return nil
}

func SetVABBID(p *models.Publication, args []string) error {
	if len(args) != 1 {
		return &ArgumentError{"vabb id is missing"}
	}
	p.VABBID = args[0]
	return nil
}

func SetVABBApproved(p *models.Publication, args []string) error {
	if len(args) != 1 {
		return &ArgumentError{"value must be 'true' or 'false'"}
	}
	switch args[0] {
	case "true":
		p.VABBApproved = true
	case "false":
		p.VABBApproved = false
	default:
		return &ArgumentError{"value must be 'true' or 'false'"}
	}
	return nil
}

func AddVABBYear(p *models.Publication, args []string) error {
	for _, arg := range args {
		if !slices.Contains(p.VABBYear, arg) {
			p.VABBYear = append(p.VABBYear, arg)
		}
	}
	return nil
}

func AddReviewerTag(p *models.Publication, args []string) error {
	for _, arg := range args {
		if !slices.Contains(p.ReviewerTags, arg) {
			p.ReviewerTags = append(p.ReviewerTags, arg)
		}
	}
	return nil
}

func RemoveReviewerTag(p *models.Publication, args []string) error {
	var vals []string
	for _, val := range p.ReviewerTags {
		if !slices.Contains(args, val) {
			vals = append(vals, val)
		}
	}
	p.ReviewerTags = vals
	return nil
}

func SetJournalTitle(p *models.Publication, args []string) error {
	if len(args) != 1 {
		return &ArgumentError{"journal title is missing"}
	}
	if p.Type != "journal_article" {
		return &ArgumentError{"record is not of type journal_article"}
	}
	p.Publication = args[0]
	return nil
}

func SetJournalAbbreviation(p *models.Publication, args []string) error {
	if len(args) != 1 {
		return &ArgumentError{"journal abbreviation is missing"}
	}
	if p.Type != "journal_article" {
		return &ArgumentError{"record is not of type journal_article"}
	}
	p.PublicationAbbreviation = args[0]
	return nil
}

func AddISBN(p *models.Publication, args []string) error {
	for _, arg := range args {
		if !slices.Contains(p.ISBN, arg) {
			p.ISBN = append(p.ISBN, arg)
		}
	}
	return nil
}

func RemoveISBN(p *models.Publication, args []string) error {
	var vals []string
	for _, val := range p.ISBN {
		if !slices.Contains(args, val) {
			vals = append(vals, val)
		}
	}
	p.ISBN = vals
	return nil
}

func AddEISBN(p *models.Publication, args []string) error {
	for _, arg := range args {
		if !slices.Contains(p.EISBN, arg) {
			p.EISBN = append(p.EISBN, arg)
		}
	}
	return nil
}

func RemoveEISBN(p *models.Publication, args []string) error {
	var vals []string
	for _, val := range p.EISBN {
		if !slices.Contains(args, val) {
			vals = append(vals, val)
		}
	}
	p.EISBN = vals
	return nil
}

func AddISSN(p *models.Publication, args []string) error {
	for _, arg := range args {
		if !slices.Contains(p.ISSN, arg) {
			p.ISSN = append(p.ISSN, arg)
		}
	}
	return nil
}

func RemoveISSN(p *models.Publication, args []string) error {
	var vals []string
	for _, val := range p.ISSN {
		if !slices.Contains(args, val) {
			vals = append(vals, val)
		}
	}
	p.ISSN = vals
	return nil
}

func AddEISSN(p *models.Publication, args []string) error {
	for _, arg := range args {
		if !slices.Contains(p.EISSN, arg) {
			p.EISSN = append(p.EISSN, arg)
		}
	}
	return nil
}

func RemoveEISSN(p *models.Publication, args []string) error {
	var vals []string
	for _, val := range p.EISSN {
		if !slices.Contains(args, val) {
			vals = append(vals, val)
		}
	}
	p.EISSN = vals
	return nil
}

func SetExternalField(p *models.Publication, args []string) error {
	if len(args) < 1 {
		return &ArgumentError{"key is missing"}
	}
	if len(args) < 2 {
		return &ArgumentError{"values are missing"}
	}
	if p.ExternalFields == nil {
		p.ExternalFields = models.Values{}
	}
	p.ExternalFields.SetAll(args[0], args[1:]...)
	return nil
}

func SetStatus(p *models.Publication, args []string) error {
	if len(args) != 1 {
		return &ArgumentError{"status is missing"}
	}
	p.Status = args[0]
	return nil
}

func SetLocked(p *models.Publication, args []string) error {
	if len(args) != 1 {
		return &ArgumentError{"value must be 'true' or 'false'"}
	}
	switch args[0] {
	case "true":
		p.Locked = true
	case "false":
		p.Locked = false
	default:
		return &ArgumentError{"value must be 'true' or 'false'"}
	}
	return nil
}
