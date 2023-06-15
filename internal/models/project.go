package models

import "time"

type Project struct {
	DateCreated *time.Time `json:"date_created,omitempty"`
	DateUpdated *time.Time `json:"date_updated,omitempty"`
	EndDate     string     `json:"end_date,omitempty"`
	ID          string     `json:"_id,omitempty"`
	StartDate   string     `json:"start_date,omitempty"`
	Title       string     `json:"title,omitempty"`
}

type RelatedProject struct {
	ProjectID string   `json:"project_id,omitempty"`
	Project   *Project `json:"-"`
}
