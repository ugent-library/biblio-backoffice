package es6

import (
	"time"

	"github.com/ugent-library/biblio-backend/internal/models"
)

type IndexedDataset struct {
	models.Dataset
	// not needed anymore in es7 with date nano type
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	HasMessage  bool   `json:"has_message"`
}

func NewIndexedDataset(dataset *models.Dataset) *IndexedDataset {
	idataset := &IndexedDataset{
		Dataset:     *dataset,
		DateCreated: dataset.DateCreated.Format(time.RFC3339),
		DateUpdated: dataset.DateUpdated.Format(time.RFC3339),
		HasMessage:  len(dataset.Message) > 0,
	}
	return idataset
}
