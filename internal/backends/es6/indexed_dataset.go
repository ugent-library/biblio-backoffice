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
	// index only fields
	HasMessage bool     `json:"has_message"`
	Faculty    []string `json:"faculty"`
}

func NewIndexedDataset(d *models.Dataset) *indexedDataset {
	id := &indexedDataset{
		Dataset:     *d,
		DateCreated: internal_time.FormatTimeUTC(d.DateCreated),
		DateUpdated: internal_time.FormatTimeUTC(d.DateUpdated),
		HasMessage:  len(d.Message) > 0,
	}

	if d.DateFrom != nil {
		id.DateFrom = internal_time.FormatTimeUTC(d.DateFrom)
	}
	if d.DateUntil != nil {
		id.DateUntil = internal_time.FormatTimeUTC(d.DateUntil)
	}

	// extract faculty from department trees
	for _, val := range id.Department {
		for _, dept := range val.Tree {
			// we naively assume that any 2 letter org is a faculty
			if len(dept.ID) == 2 {
				exists := false
				for _, fac := range id.Faculty {
					if fac == dept.ID {
						exists = true
						break
					}
				}

				if !exists {
					id.Faculty = append(id.Faculty, dept.ID)
				}
			}
		}
	}

	return id
}
