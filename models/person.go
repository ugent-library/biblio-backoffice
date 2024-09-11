package models

import (
	"time"
)

type Affiliation struct {
	OrganizationID string        `json:"organization_id"`
	Organization   *Organization `json:"-"`
}

type Person struct {
	Active       bool           `json:"active"`
	DateCreated  *time.Time     `json:"date_created"`
	DateUpdated  *time.Time     `json:"date_updated"`
	Affiliations []*Affiliation `json:"affiliations"`
	Email        string         `json:"email"`
	FirstName    string         `json:"first_name"`
	FullName     string         `json:"full_name"`
	ID           string         `json:"id"`
	IDs          []string       `json:"ids"`
	LastName     string         `json:"last_name"`
	ORCID        string         `json:"orcid"`
	UGentID      []string       `json:"ugent_id"`
	// fields below are only relevant for Active people (users)
	Username   string `json:"username"`
	Role       string `json:"role"`
	ORCIDToken string `json:"orcid_token"`
}

func (p *Person) AffiliatedWith(orgID string) bool {
	for _, aff := range p.Affiliations {
		if aff.OrganizationID == orgID {
			return true
		}
		if aff.Organization != nil {
			for _, org := range aff.Organization.Tree {
				if org.ID == orgID {
					return true
				}
			}
		}
	}
	return false
}
