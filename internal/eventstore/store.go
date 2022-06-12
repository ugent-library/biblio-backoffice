package eventstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// TODO add data factory for ex nihilo creation of stream and event data and remove nil check

var NotFound = errors.New("stream not found")

var DefaultIDGenerator = func() (string, error) {
	return uuid.NewString(), nil
}

type PgConn interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
}

// type StreamEventHandler interface {
// 	StreamType() string
// 	Handle(string, any, any) (any, error)
// }

type Store struct {
	conn        PgConn
	idGenerator func() (string, error)
	// handlers    map[string]StreamEventHandler
	// handlersMu  sync.RWMutex
}

type RawSnapshot struct {
	StreamID    string
	StreamType  string
	EventID     string
	Data        json.RawMessage
	DateCreated time.Time
	DateUpdated time.Time
}

func Connect(ctx context.Context, dsn string, opts ...func(*Store)) (*Store, error) {
	conn, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("eventstore: failed to connect: %w", err)
	}

	return New(conn, opts...), nil
}

func New(conn PgConn, opts ...func(*Store)) *Store {
	s := &Store{
		conn:        conn,
		idGenerator: DefaultIDGenerator,
		// handlers:    make(map[string]StreamEventHandler),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithIDGenerator(fn func() (string, error)) func(*Store) {
	return func(s *Store) {
		s.idGenerator = fn
	}
}

// func (s *Store) AddHandler(h StreamEventHandler) {
// 	s.handlersMu.Lock()
// 	defer s.handlersMu.Unlock()
// 	s.handlers[h.StreamType()] = h
// }

func (s *Store) Append(ctx context.Context, events ...Event) error {
	if len(events) == 0 {
		return nil
	}

	// TODO avoid allocating
	// TODO refactor, stream id is only unique per stream type
	eventMap := make(map[string][]Event)
	for _, event := range events {
		eventMap[event.StreamID()] = append(eventMap[event.StreamID()], event)
	}

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("eventstore: failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var lastEventID string

	for _, e := range events {
		var rawData, rawMeta json.RawMessage
		if e.Data() != nil {
			if rawData, err = json.Marshal(e.Data()); err != nil {
				return fmt.Errorf("eventstore: failed to serialize event data: %w", err)
			}
		}
		if e.Meta() != nil {
			if rawMeta, err = json.Marshal(e.Meta()); err != nil {
				return fmt.Errorf("eventstore: failed to serialize event meta: %w", err)
			}
		}

		// generate event id
		id, err := s.idGenerator()
		if err != nil {
			return fmt.Errorf("eventstore: failed to generate id: %w", err)
		}

		lastEventID = id

		sql := `insert into events(id, stream_id, stream_type, name, data, meta)
		values ($1, $2, $3, $4, $5, $6)`

		if _, err = tx.Exec(ctx, sql, id, e.StreamID(), e.StreamType(), e.Name(), rawData, rawMeta); err != nil {
			return fmt.Errorf("eventstore: failed to insert event: %w", err)
		}
	}

	for streamID, events := range eventMap {
		streamType := events[0].StreamType()
		snap, err := s.getSnapshot(ctx, tx, streamID)
		if err != nil {
			return err
		}
		// s.handlersMu.RLock()
		// p, ok := s.handlers[streamType]
		// s.handlersMu.RUnlock()
		// if !ok {
		// 	return fmt.Errorf("eventstore: no stream handler for %s", streamType)
		// }

		// TODO use factory
		var d any
		if snap != nil {
			d = snap.Data
		}

		for _, e := range events {
			if d, err = e.Apply(d); err != nil {
				return fmt.Errorf("eventstore: failed to handle event: %w", err)
			}
		}

		rawData, err := json.Marshal(d)
		if err != nil {
			return fmt.Errorf("eventstore: failed to serialize projection data: %w", err)
		}

		// TODO set date_updated to last date_created of events
		now := time.Now()
		if snap == nil {
			sql := `insert into projections(stream_id, stream_type, event_id, data, date_created, date_updated)
		values($1, $2, $3, $4, $5, $6)`
			if _, err = tx.Exec(ctx, sql, streamID, streamType, lastEventID, rawData, now, now); err != nil {
				return fmt.Errorf("eventstore: failed to insert projection: %w", err)
			}
		} else {
			// TODO check row count or use one on conflict statement
			sql := `update projections set event_id = $1, data = $2, date_updated = $3 where stream_id = $4`
			if _, err = tx.Exec(ctx, sql, lastEventID, rawData, now, streamID); err != nil {
				return fmt.Errorf("eventstore: failed to update projection: %w", err)
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("eventstore: failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Store) getSnapshot(ctx context.Context, tx PgConn, streamID string) (*RawSnapshot, error) {
	snap := RawSnapshot{StreamID: streamID}

	sql := `select stream_type, event_id, data, date_created, date_updated from projections
	where stream_id = $1
	limit 1`
	err := tx.QueryRow(ctx, sql, streamID).Scan(&snap.StreamType, &snap.EventID, &snap.Data, &snap.DateCreated, &snap.DateUpdated)

	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("eventstore: failed to get projection: %w", err)
	}

	return &snap, nil
}
