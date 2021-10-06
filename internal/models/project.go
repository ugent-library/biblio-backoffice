package models

import "time"

type Project struct {
	DataFormat  string     `json:"data_format,omitempty"`
	DateCreated *time.Time `json:"date_created,omitempty"`
	DateUpdated *time.Time `json:"date_updated,omitempty"`
	StartDate   string     `json:"start_date,omitempty"`
	EndDate     string     `json:"end_date,omitempty"`
	ID          string     `json:"_id,omitempty"`
	Name        string     `json:"name,omitempty"`
}
