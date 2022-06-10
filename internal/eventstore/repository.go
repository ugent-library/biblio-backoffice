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
	StreamID    string
	EventID     string
	DateCreated time.Time
	DateUpdated time.Time
	Data        T
}

type Repository[T any] interface {
	Type() string
	Get(ctx context.Context, id string) (Projection[T], error)
	GetAll(ctx context.Context) (Cursor[T], error)
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

func (r *repository[T]) GetAll(ctx context.Context) (Cursor[T], error) {
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

type Cursor[T any] struct {
	rows pgx.Rows
}

func (c Cursor[T]) HasNext() bool {
	return c.rows.Next()
}

func (c Cursor[T]) Next() (Projection[T], error) {
	var (
		p       Projection[T]
		rawData json.RawMessage
	)

	if err := c.rows.Scan(&p.StreamID, &p.EventID, &rawData, &p.DateCreated, &p.DateUpdated); err != nil {
		return p, fmt.Errorf("eventstore: failed to scan projection: %w", err)
	}

	if err := json.Unmarshal(rawData, &p.Data); err != nil {
		return p, fmt.Errorf("eventstore: failed to deserialize projection data: %w", err)
	}

	return p, nil
}

func (c Cursor[T]) Close() {
	c.rows.Close()
}

func (c Cursor[T]) Error() error {
	return c.rows.Err()
}
