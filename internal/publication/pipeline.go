package publication

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
	EnsureFullName,
)

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
