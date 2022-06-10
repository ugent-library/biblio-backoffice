package eventstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

type Projection[T any] struct {
	ID          string
	EventID     string
	DateCreated time.Time
	DateUpdated time.Time
	Data        T
}

type Repository[T any] interface {
	Type() string
	Get(ctx context.Context, id string) (Projection[T], error)
}

type repository[T any] struct {
	store      *Store
	streamType string
}

func NewRepository[T any](store *Store, streamType string) Repository[T] {
	return &repository[T]{
		store:      store,
		streamType: streamType,
	}
}

func (r *repository[T]) Type() string {
	return r.streamType
}

func (r *repository[T]) Get(ctx context.Context, streamID string) (Projection[T], error) {
	var p Projection[T]

	sql := `select event_id, data, date_created, date_updated
	from projections
	where stream_id = $1 and stream_type = $2
	limit 1`

	var rawData json.RawMessage

	err := r.store.conn.
		QueryRow(ctx, sql, streamID, r.streamType).
		Scan(&p.EventID, &rawData, &p.DateCreated, &p.DateUpdated)

	if errors.Is(err, pgx.ErrNoRows) {
		return p, NotFound
	} else if err != nil {
		return p, fmt.Errorf("eventstore: failed to get projection: %w", err)
	}

	if err := json.Unmarshal(rawData, &p.Data); err != nil {
		return p, fmt.Errorf("eventstore: failed to deserialize projection data: %w", err)
	}

	return p, nil
}
