package views

import (
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
)

type ContributorData struct {
	render      *render.Render
	Publication *models.Publication
	Author      *models.PublicationContributor
	Key         string
}

func NewContributorData(render *render.Render, p *models.Publication, a *models.PublicationContributor, k string) ContributorData {
	return ContributorData{render: render, Publication: p, Author: a, Key: k}
}
