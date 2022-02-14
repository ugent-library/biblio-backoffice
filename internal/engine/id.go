package engine

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// TODO cache rand source (sync pool or in goroutine)
// do we need ulid.Monotonic?
func NewCorrelationID() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}
