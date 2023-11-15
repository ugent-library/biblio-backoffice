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

func PublicationFromCandidateRecord(rec *CandidateRecord) (*Publication, error) {
	p := &Publication{}
	if err := json.Unmarshal(rec.Metadata, p); err != nil {
		return nil, err
	}
	return p, nil
}
