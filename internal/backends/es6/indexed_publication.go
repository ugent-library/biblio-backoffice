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
	Department   []string `json:"department,omitempty"`
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

	// extract faculty and department id from department trees
	for _, rel := range ip.RelatedOrganizations {
		for _, org := range rel.Organization.Tree {
			if !validation.InArray(ip.Department, org.ID) {
				ip.Department = append(ip.Department, org.ID)
			}

			if validation.InArray(faculties, org.ID) && !validation.InArray(ip.Faculty, org.ID) {
				ip.Faculty = append(ip.Faculty, org.ID)
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
