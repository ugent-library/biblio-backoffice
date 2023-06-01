package publication

import (
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

// TODO eliminate need for this

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
	ensureHasBeenPublic,
)

var PublishPipeline = NewPipeline()
var UnpublishPipeline = NewPipeline()

func ensureHasBeenPublic(p *models.Publication) *models.Publication {
	if p.Status == "public" && !p.HasBeenPublic {
		p.HasBeenPublic = true
	}
	return p
}
