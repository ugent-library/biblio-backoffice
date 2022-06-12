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
	// eventHandlers   map[string]EventHandler[T]
	// eventHandlersMu sync.RWMutex
}

func NewRepository[T any](store *Store, streamType string) *Repository[T] {
	r := &Repository[T]{
		store:      store,
		streamType: streamType,
		// eventHandlers: make(map[string]EventHandler[T]),
	}
	// store.AddHandler(r)
	return r
}

func (r *Repository[T]) StreamType() string {
	return r.streamType
}

// func (r *Repository[T]) AddEventHandlers(handlers ...EventHandler[T]) {
// 	r.eventHandlersMu.Lock()
// 	defer r.eventHandlersMu.Unlock()
// 	for _, h := range handlers {
// 		r.eventHandlers[h.Name()] = h
// 	}
// }

// func (r *Repository[T]) Handle(eventName string, d, eventData any) (any, error) {
// 	r.eventHandlersMu.RLock()
// 	h, ok := r.eventHandlers[eventName]
// 	r.eventHandlersMu.RUnlock()

// 	if !ok {
// 		return nil, fmt.Errorf("eventstore: no handler for event %s", eventName)
// 	}

// 	var data T

// 	switch t := d.(type) {
// 	case nil:
// 		// do nothing
// 	case T:
// 		data = t
// 	case json.RawMessage:
// 		if err := json.Unmarshal(t, data); err != nil {
// 			return data, fmt.Errorf("eventstore: failed to deserialize projection data into %T: %w", eventData, err)
// 		}
// 	default:
// 		return data, fmt.Errorf("eventstore: invalid projection data type %T", t)
// 	}

// 	return h.Apply(data, eventData)
// }

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
