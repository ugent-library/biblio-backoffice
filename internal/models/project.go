package models

import "time"

type Project struct {
	ID          string     `json:"_id,omitempty"`
	Title       string     `json:"title,omitempty"`
	StartDate   string     `json:"start_date,omitempty"`
	EndDate     string     `json:"end_date,omitempty"`
	DateCreated *time.Time `json:"date_created,omitempty"`
	DateUpdated *time.Time `json:"date_updated,omitempty"`
	EUProject   *EUProject `json:"eu,omitempty"`
}

type EUProject struct {
	ID                 string
	Acronym            string
	CallID             string
	FrameworkProgramme string
}

type RelatedProject struct {
	ProjectID string   `json:"project_id,omitempty"`
	Project   *Project `json:"-"`
}
