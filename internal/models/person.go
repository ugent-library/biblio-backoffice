package models

import "time"

type PersonDepartment struct {
	ID string `json:"_id"`
}

type Person struct {
	DateCreated *time.Time         `json:"date_created"`
	DateUpdated *time.Time         `json:"date_updated"`
	Department  []PersonDepartment `json:"department"`
	FirstName   string             `json:"first_name" form:"first_name"`
	FullName    string             `json:"full_name"`
	ID          string             `json:"_id"`
	LastName    string             `json:"last_name" form:"last_name"`
	ORCID       string             `json:"orcid"`
	UGentID     []string           `json:"ugent_id" form:"ugent_id"`
}
