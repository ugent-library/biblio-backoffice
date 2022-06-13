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
	streamType *streamType[T]
}

func NewRepository[T any](store *Store, t *streamType[T]) *Repository[T] {
	return &Repository[T]{
		store:      store,
		streamType: t,
	}
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
		QueryRow(ctx, sql, streamID, r.streamType.name).
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

func (r *Repository[T]) GetAt(ctx context.Context, streamID, eventID string) (Projection[T], error) {
	// TODO: set timestamp fields
	p := Projection[T]{
		StreamID: streamID,
		EventID:  eventID,
		Data:     r.streamType.factory(),
	}

	sql := `select name, data
	from events
	where stream_id = $1 and stream_type = $2 and
	seq <= (select seq from events where stream_id = $1 and stream_type = $2 and id = $3)`

	rows, err := r.store.conn.Query(ctx, sql, streamID, r.streamType.name, eventID)
	if err != nil {
		return p, fmt.Errorf("eventstore: failed to query events: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			name string
			data json.RawMessage
		)

		if err := rows.Scan(&name, &data); err != nil {
			return p, fmt.Errorf("eventstore: failed to scan event: %w", err)
		}

		h := r.store.GetEventHandler(r.streamType.name, name)
		if h == nil {
			return p, fmt.Errorf("eventstore: eventhandler %s %s not found", r.streamType.name, name)
		}

		// TODO improve this
		d, err := h.Apply(p.Data, data)
		if err != nil {
			return p, err
		}
		if t, ok := d.(T); ok {
			p.Data = t
		} else {
			return p, fmt.Errorf("eventstore: invalid projection data type %T", t)
		}
	}

	if err := rows.Err(); err != nil {
		return p, fmt.Errorf("eventstore: %w", err)
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

	c.rows, err = r.store.conn.Query(ctx, sql, r.streamType.name)
	if err != nil {
		return c, fmt.Errorf("eventstore: failed to query projections: %w", err)
	}

	return c, nil
}
