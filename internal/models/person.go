package models

import "time"

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
	LastName     string         `json:"last_name"`
	ORCID        string         `json:"orcid"`
	UGentID      []string       `json:"ugent_id"`
}
