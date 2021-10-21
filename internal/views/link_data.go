package views

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
)

type LinkData struct {
	Data
	render      *render.Render
	Publication *models.Publication
	Link        *models.PublicationLink
	Key         string
}

func NewLinkData(r *http.Request, render *render.Render, p *models.Publication, l *models.PublicationLink, k string) LinkData {
	return LinkData{Data: NewData(r), render: render, Publication: p, Link: l, Key: k}
}
