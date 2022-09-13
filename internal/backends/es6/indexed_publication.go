package es6

import (
	"github.com/ugent-library/biblio-backend/internal/models"
)

type IndexedPublication struct {
	models.Publication
	// not needed anymore in es7 with date nano type
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	DateFrom    string `json:"date_from,omitempty"`
	DateUntil   string `json:"date_until,omitempty"`
	HasMessage  bool   `json:"has_message"`
}

func NewIndexedPublication(publication *models.Publication) *IndexedPublication {
	ipub := &IndexedPublication{
		Publication: *publication,
		DateCreated: FormatTimeUTC(publication.DateCreated),
		DateUpdated: FormatTimeUTC(publication.DateUpdated),
		HasMessage:  len(publication.Message) > 0,
	}
	if publication.DateFrom != nil {
		ipub.DateFrom = FormatTimeUTC(publication.DateFrom)
	}
	if publication.DateUntil != nil {
		ipub.DateUntil = FormatTimeUTC(publication.DateUntil)
	}

	return ipub
}
