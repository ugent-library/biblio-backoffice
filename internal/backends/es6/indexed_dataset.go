package es6

import (
	"github.com/ugent-library/biblio-backend/internal/models"
	internal_time "github.com/ugent-library/biblio-backend/internal/time"
)

type indexedDataset struct {
	models.Dataset
	// not needed anymore in es7 with date nano type
	DateCreated string `json:"date_created,omitempty"`
	DateUpdated string `json:"date_updated,omitempty"`
	DateFrom    string `json:"date_from,omitempty"`
	DateUntil   string `json:"date_until,omitempty"`
	HasMessage  bool   `json:"has_message"`
}

func NewIndexedDataset(dataset *models.Dataset) *indexedDataset {
	idataset := &indexedDataset{
		Dataset:     *dataset,
		DateCreated: internal_time.FormatTimeUTC(dataset.DateCreated),
		DateUpdated: internal_time.FormatTimeUTC(dataset.DateUpdated),
		HasMessage:  len(dataset.Message) > 0,
	}
	if dataset.DateFrom != nil {
		idataset.DateFrom = internal_time.FormatTimeUTC(dataset.DateFrom)
	}
	if dataset.DateUntil != nil {
		idataset.DateUntil = internal_time.FormatTimeUTC(dataset.DateUntil)
	}
	return idataset
}
