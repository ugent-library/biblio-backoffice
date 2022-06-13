package eventstore

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

type Repository[T any] struct {
	store      *Store
	streamType string
}

func NewRepository[T any](store *Store, streamType string) *Repository[T] {
	return &Repository[T]{
		store:      store,
		streamType: streamType,
	}
}

func (r *Repository[T]) StreamType() string {
	return r.streamType
}

func (r *Repository[T]) Get(ctx context.Context, streamID string) (Projection[T], error) {
	var (
		p       Projection[T]
		rawData json.RawMessage
	)

	p.StreamID = streamID

	sql := `select event_id, data, date_created, date_updated
	from projections
	where stream_id = $1 and stream_type = $2
	limit 1`

	err := r.store.conn.
		QueryRow(ctx, sql, streamID, r.streamType).
		Scan(&p.EventID, &rawData, &p.DateCreated, &p.DateUpdated)

	if errors.Is(err, pgx.ErrNoRows) {
		return p, NotFound
	} else if err != nil {
		return p, fmt.Errorf("eventstore: failed to scan projection: %w", err)
	}

	if err := json.Unmarshal(rawData, &p.Data); err != nil {
		return p, fmt.Errorf("eventstore: failed to deserialize projection data: %w", err)
	}

	return p, nil
}

func (r *Repository[T]) GetAll(ctx context.Context) (Cursor[T], error) {
	sql := `select stream_id, event_id, data, date_created, date_updated
	from projections
	where stream_type = $1`

	var (
		c   Cursor[T]
		err error
	)

	c.rows, err = r.store.conn.Query(ctx, sql, r.streamType)
	if err != nil {
		return c, fmt.Errorf("eventstore: failed to query projections: %w", err)
	}

	return c, nil
}
