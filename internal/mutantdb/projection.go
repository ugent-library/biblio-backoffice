package mutantdb

import (
	"time"
)

type Projection[T any] struct {
	ID          string
	Data        T
	MutationID  string
	DateCreated time.Time
	DateUpdated time.Time
}
