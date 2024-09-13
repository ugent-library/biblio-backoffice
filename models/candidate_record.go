package models

import (
	"encoding/json"
	"time"
)

type CandidateRecord struct {
	ID             string          `json:"id"`
	SourceName     string          `json:"source_name"`
	SourceID       string          `json:"source_id"`
	SourceMetadata []byte          `json:"source_metadata"`
	Type           string          `json:"type"`
	Metadata       json.RawMessage `json:"metadata"`
	Status         string          `json:"status"`
	DateCreated    time.Time       `json:"date_created"`
	// this is the Metadata mapped to a publication in-memory, not stored in DB
	// TODO tightly coupled with Publication for now, refactor later
	Publication    *Publication `json:"publication"`
	StatusDate     *time.Time   `json:"status_date"`
	StatusPersonID *string      `json:"status_person_id"`
	StatusPerson   *Person      `json:"status_person"`
	ImportedID     *string      `json:"imported_id"`
}
