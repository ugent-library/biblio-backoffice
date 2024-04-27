package models

import (
	"encoding/json"
	"time"
)

type CandidateRecord struct {
	ID             string
	SourceName     string
	SourceID       string
	SourceMetadata []byte
	Type           string
	Metadata       json.RawMessage
	Status         string
	DateCreated    time.Time
	// TODO tightly coupled with Publication for now, refactor later
	Publication *Publication
}
