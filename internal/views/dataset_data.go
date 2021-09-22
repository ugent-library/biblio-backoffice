package views

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/unrolled/render"
)

type DatasetData struct {
	Data
	render  *render.Render
	Dataset *models.Publication
}

func NewDatasetData(r *http.Request, render *render.Render, p *models.Publication) DatasetData {
	return DatasetData{Data: NewData(r), render: render, Dataset: p}
}
