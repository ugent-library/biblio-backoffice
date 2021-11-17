package views

import (
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
)

type ContributorData struct {
	render      *render.Render
	Publication *models.Publication
	Author      *models.Contributor
	Key         string
}

func NewContributorData(render *render.Render, p *models.Publication, a *models.Contributor, k string) ContributorData {
	return ContributorData{render: render, Publication: p, Author: a, Key: k}
}
