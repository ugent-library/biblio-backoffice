package views

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
)

type AbstractData struct {
	Data
	render      *render.Render
	Publication *models.Publication
	Abstract    *models.Text
	Key         string
}

func NewAbstractData(r *http.Request, render *render.Render, p *models.Publication, a *models.Text, k string) AbstractData {
	return AbstractData{Data: NewData(r), render: render, Publication: p, Abstract: a, Key: k}
}
