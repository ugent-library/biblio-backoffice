package es6

import (
	"github.com/ugent-library/biblio-backend/internal/models"
	internal_time "github.com/ugent-library/biblio-backend/internal/time"
)

type indexedPublication struct {
	models.Publication
	// not needed anymore in es7 with date nano type
	DateCreated string `json:"date_created"`
	DateFrom    string `json:"date_from,omitempty"`
	DateUntil   string `json:"date_until,omitempty"`
	DateUpdated string `json:"date_updated"`
	// index only field
	HasMessage bool `json:"has_message"`
}

func NewIndexedPublication(publication *models.Publication) *indexedPublication {
	ipub := &indexedPublication{
		Publication: *publication,
		DateCreated: internal_time.FormatTimeUTC(publication.DateCreated),
		DateUpdated: internal_time.FormatTimeUTC(publication.DateUpdated),
		HasMessage:  len(publication.Message) > 0,
	}
	if publication.DateFrom != nil {
		ipub.DateFrom = internal_time.FormatTimeUTC(publication.DateFrom)
	}
	if publication.DateUntil != nil {
		ipub.DateUntil = internal_time.FormatTimeUTC(publication.DateUntil)
	}

	return ipub
}
