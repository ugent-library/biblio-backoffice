package eventstore

import "time"

type Projection[T any] struct {
	StreamID    string
	EventID     string
	DateCreated time.Time
	DateUpdated time.Time
	Data        T
}
