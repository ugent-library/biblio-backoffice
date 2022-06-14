package mutantdb

import (
	"encoding/json"
	"time"
)

type Projection[T any] struct {
	ID          string
	Data        T
	MutationID  string
	DateCreated time.Time
	DateUpdated time.Time
}

type RawProjection struct {
	ID          string
	Type        string
	Data        json.RawMessage
	MutationID  string
	DateCreated time.Time
	DateUpdated time.Time
}
