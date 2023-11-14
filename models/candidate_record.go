package models

import (
	"encoding/json"
	"time"
)

type CandidateRecord struct {
	SourceName     string
	SourceID       string
	SourceMetadata []byte
	Type           string
	Metadata       json.RawMessage
	AssignedUserID string
	DateCreated    time.Time
}
