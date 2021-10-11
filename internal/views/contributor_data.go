package views

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
)

type ContributorData struct {
	Data
	render      *render.Render
	Publication *models.Publication
}

func NewContributorData(r *http.Request, render *render.Render, p *models.Publication) ContributorData {
	return ContributorData{Data: NewData(r), render: render, Publication: p}
}
