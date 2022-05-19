package publications

import (
	"strings"

	"github.com/ugent-library/biblio-backend/internal/models"
)

type PipelineFunc func(*models.Publication) *models.Publication

type Pipeline []PipelineFunc

func NewPipeline(fns ...PipelineFunc) Pipeline {
	return Pipeline(fns)
}

func (pl Pipeline) Process(p *models.Publication) *models.Publication {
	for _, fn := range pl {
		p = fn(p)
	}
	return p
}

func (pl Pipeline) Func() PipelineFunc {
	return func(p *models.Publication) *models.Publication {
		return pl.Process(p)
	}
}

var DefaultPipeline = NewPipeline(
	Vacuum,
	EnsureFullName,
)

func Vacuum(p *models.Publication) *models.Publication {
	p.AdditionalInfo = strings.TrimSpace(p.AdditionalInfo)
	p.AlternativeTitle = vacuumStringSlice(p.AlternativeTitle)
	p.ArticleNumber = strings.TrimSpace(p.ArticleNumber)
	p.ArxivID = strings.TrimSpace(p.ArxivID)
	p.DefenseDate = strings.TrimSpace(p.DefenseDate)
	p.DefensePlace = strings.TrimSpace(p.DefensePlace)
	p.DefenseTime = strings.TrimSpace(p.DefenseTime)
	p.DOI = strings.TrimSpace(p.DOI)
	p.Edition = strings.TrimSpace(p.Edition)
	p.EISBN = vacuumStringSlice(p.EISBN)
	p.EISSN = vacuumStringSlice(p.EISSN)
	p.ESCIID = strings.TrimSpace(p.ESCIID)
	p.ISBN = vacuumStringSlice(p.ISBN)
	p.ISSN = vacuumStringSlice(p.ISSN)
	p.Issue = strings.TrimSpace(p.Issue)
	p.IssueTitle = strings.TrimSpace(p.IssueTitle)
	p.Keyword = vacuumStringSlice(p.Keyword)
	p.Language = vacuumStringSlice(p.Language)
	p.PageCount = strings.TrimSpace(p.PageCount)
	p.PageFirst = strings.TrimSpace(p.PageFirst)
	p.PageLast = strings.TrimSpace(p.PageLast)
	p.PlaceOfPublication = strings.TrimSpace(p.PlaceOfPublication)
	p.Publication = strings.TrimSpace(p.Publication)
	p.PublicationAbbreviation = strings.TrimSpace(p.PublicationAbbreviation)
	p.Publisher = strings.TrimSpace(p.Publisher)
	p.PubMedID = strings.TrimSpace(p.PubMedID)
	p.ReportNumber = strings.TrimSpace(p.ReportNumber)
	p.ResearchField = vacuumStringSlice(p.ResearchField)
	p.SeriesTitle = strings.TrimSpace(p.SeriesTitle)
	p.Title = strings.TrimSpace(p.Title)
	p.URL = strings.TrimSpace(p.URL)
	p.Volume = strings.TrimSpace(p.Volume)
	p.WOSID = strings.TrimSpace(p.WOSID)
	p.Year = strings.TrimSpace(p.Year)
	return p
}

func vacuumStringSlice(vals []string) []string {
	newVals := []string{}
	for _, v := range vals {
		v = strings.TrimSpace(v)
		if v != "" {
			newVals = append(newVals, v)
		}
	}
	return newVals
}

func EnsureFullName(p *models.Publication) *models.Publication {
	for _, c := range p.Author {
		if c.FullName == "" {
			c.FullName = strings.Join([]string{c.FirstName, c.LastName}, " ")
		}
	}
	for _, c := range p.Editor {
		if c.FullName == "" {
			c.FullName = strings.Join([]string{c.FirstName, c.LastName}, " ")
		}
	}
	for _, c := range p.Supervisor {
		if c.FullName == "" {
			c.FullName = strings.Join([]string{c.FirstName, c.LastName}, " ")
		}
	}
	return p
}
