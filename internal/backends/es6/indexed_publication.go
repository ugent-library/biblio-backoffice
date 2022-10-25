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
	// index only fields
	HasMessage bool     `json:"has_message"`
	Faculty    []string `json:"faculty"`
}

func NewIndexedPublication(p *models.Publication) *indexedPublication {
	ip := &indexedPublication{
		Publication: *p,
		DateCreated: internal_time.FormatTimeUTC(p.DateCreated),
		DateUpdated: internal_time.FormatTimeUTC(p.DateUpdated),
		HasMessage:  len(p.Message) > 0,
	}

	if p.DateFrom != nil {
		ip.DateFrom = internal_time.FormatTimeUTC(p.DateFrom)
	}
	if p.DateUntil != nil {
		ip.DateUntil = internal_time.FormatTimeUTC(p.DateUntil)
	}

	// extract faculty from department trees
	for _, val := range p.Department {
		for _, dept := range val.Tree {
			// we naively assume that any 2 letter org is a faculty
			if len(dept.ID) == 2 {
				exists := false
				for _, fac := range ip.Faculty {
					if fac == dept.ID {
						exists = true
						break
					}
				}

				if !exists {
					ip.Faculty = append(ip.Faculty, dept.ID)
				}
			}
		}
	}

	return ip
}
