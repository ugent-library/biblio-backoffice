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
	publication    *Publication
}

// TODO tightly coupled with Publication for now, refactor later
// TODO handle error
func (rec *CandidateRecord) AsPublication() *Publication {
	if rec.publication != nil {
		return rec.publication
	}
	p := &Publication{}
	if err := json.Unmarshal(rec.Metadata, p); err != nil {
		panic(err)
	}
	rec.publication = p
	return p
}
