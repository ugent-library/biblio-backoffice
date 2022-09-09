package es6

import (
	"github.com/ugent-library/biblio-backend/internal/models"
)

type IndexedDataset struct {
	models.Dataset
	// not needed anymore in es7 with date nano type
	DateCreated string `json:"date_created,omitempty"`
	DateUpdated string `json:"date_updated,omitempty"`
	DateFrom    string `json:"date_from,omitempty"`
	DateUntil   string `json:"date_until,omitempty"`
	HasMessage  bool   `json:"has_message"`
}

func NewIndexedDataset(dataset *models.Dataset) *IndexedDataset {
	idataset := &IndexedDataset{
		Dataset:     *dataset,
		DateCreated: FormatTimeUTC(dataset.DateCreated),
		DateUpdated: FormatTimeUTC(dataset.DateUpdated),
		HasMessage:  len(dataset.Message) > 0,
	}
	if dataset.DateFrom != nil {
		idataset.DateFrom = FormatTimeUTC(dataset.DateFrom)
	}
	if dataset.DateUntil != nil {
		idataset.DateUntil = FormatTimeUTC(dataset.DateUntil)
	}
	return idataset
}
