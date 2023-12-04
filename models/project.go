package models

import "time"

type Project struct {
	ID          string
	Title       string
	Description string
	Acronym     string
	StartDate   string
	EndDate     string
	DateCreated *time.Time
	DateUpdated *time.Time
	EUProject   *EUProject
	GISMOID     string
	IWETOID     string
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
