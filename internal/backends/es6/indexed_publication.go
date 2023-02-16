package es6

import (
	"regexp"
	"strings"
  
  "github.com/ugent-library/biblio-backoffice/internal/models"
	internal_time "github.com/ugent-library/biblio-backoffice/internal/time"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
	"github.com/ugent-library/biblio-backoffice/internal/vocabularies"
)

type indexedPublication struct {
	models.Publication
	// not needed anymore in es7 with date nano type
	DateCreated string `json:"date_created"`
	DateFrom    string `json:"date_from,omitempty"`
	DateUntil   string `json:"date_until,omitempty"`
	DateUpdated string `json:"date_updated"`
	// index only fields
	HasMessage   bool     `json:"has_message"`
	Faculty      []string `json:"faculty,omitempty"`
	FacetWOSType []string `json:"facet_wos_type,omitempty"`
}

var reSplitWOS *regexp.Regexp = regexp.MustCompile("[,;]")

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

	faculties := vocabularies.Map["faculties"]

	// extract faculty from department trees
	for _, val := range p.Department {
		for _, dept := range val.Tree {
			if validation.InArray(faculties, dept.ID) {
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

	if ip.WOSType != "" {

		wos_types := reSplitWOS.Split(ip.WOSType, -1)
		for _, wos_type := range wos_types {
			wt := strings.TrimSpace(wos_type)
			ip.FacetWOSType = append(ip.FacetWOSType, wt)
		}

	}

	return ip
}
