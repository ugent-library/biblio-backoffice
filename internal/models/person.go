package models

import "time"

type PersonDepartment struct {
	ID string `json:"_id"`
}

type Person struct {
	DateCreated *time.Time         `json:"date_created" form:"-"`
	DateUpdated *time.Time         `json:"date_updated" form:"-"`
	Department  []PersonDepartment `json:"department" form:"-"`
	FirstName   string             `json:"first_name" form:"first_name"`
	FullName    string             `json:"full_name" form:"-"`
	ID          string             `json:"_id" form:"-"`
	LastName    string             `json:"last_name" form:"last_name"`
	ORCID       string             `json:"orcid" form:"-"`
	UGentID     []string           `json:"ugent_id" form:"-"`
}
