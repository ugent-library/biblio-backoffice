package es6

import (
	"time"

	"github.com/ugent-library/biblio-backend/internal/models"
)

type IndexedPublication struct {
	models.Publication
	// not needed anymore in es7 with date nano type
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	HasMessage  bool   `json:"has_message"`
}

func NewIndexedPublication(publication *models.Publication) *IndexedPublication {
	ipub := &IndexedPublication{
		Publication: *publication,
		DateCreated: publication.DateCreated.Format(time.RFC3339),
		DateUpdated: publication.DateUpdated.Format(time.RFC3339),
		HasMessage:  len(publication.Message) > 0,
	}
	return ipub
}
