package es6

import (
	"github.com/ugent-library/biblio-backoffice/internal/models"
	internal_time "github.com/ugent-library/biblio-backoffice/internal/time"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
	"github.com/ugent-library/biblio-backoffice/internal/vocabularies"
)

type indexedDataset struct {
	models.Dataset
	// not needed anymore in es7 with date nano type
	DateCreated string `json:"date_created,omitempty"`
	DateUpdated string `json:"date_updated,omitempty"`
	DateFrom    string `json:"date_from,omitempty"`
	DateUntil   string `json:"date_until,omitempty"`
	// index only fields
	HasMessage       bool     `json:"has_message"`
	Department       []string `json:"department,omitempty"`
	Faculty          []string `json:"faculty,omitempty"`
	IdentifierTypes  []string `json:"identifier_types,omitempty"`
	IdentifierValues []string `json:"identifier_values,omitempty"`
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

	faculties := vocabularies.Map["faculties"]

	// extract faculty and department id from department trees
	for _, rel := range id.RelatedOrganizations {
		for _, org := range rel.Organization.Tree {
			if !validation.InArray(id.Department, org.ID) {
				id.Department = append(id.Department, org.ID)
			}

			if validation.InArray(faculties, org.ID) && !validation.InArray(id.Faculty, org.ID) {
				id.Faculty = append(id.Faculty, org.ID)
			}
		}
	}

	for k, vals := range d.Identifiers {
		id.IdentifierTypes = append(id.IdentifierTypes, k)
		id.IdentifierValues = append(id.IdentifierValues, vals...)
	}

	return id
}
