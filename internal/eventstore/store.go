package eventstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var NotFound = errors.New("stream not found")

type PgConn interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
}

type RawProcessor interface {
	RawApply(json.RawMessage, []Event) (json.RawMessage, error)
}

type Store struct {
	conn         PgConn
	processors   map[string]RawProcessor
	processorsMu sync.RWMutex
}

type Event struct {
	ID          string
	StreamID    string
	StreamType  string
	Type        string
	Data        any
	Meta        map[string]string
	DateCreated time.Time
}

type RawSnapshot struct {
	StreamID    string
	StreamType  string
	EventID     string
	Data        json.RawMessage
	DateCreated time.Time
	DateUpdated time.Time
}

func Connect(ctx context.Context, dsn string) (*Store, error) {
	conn, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("eventstore: failed to connect: %w", err)
	}

	return New(conn), nil
}

func New(conn PgConn) *Store {
	return &Store{
		conn: conn,
	}
}

func (s *Store) AddProcessor(streamType string, p RawProcessor) {
	s.processorsMu.Lock()
	defer s.processorsMu.Unlock()
	if s.processors == nil {
		s.processors = make(map[string]RawProcessor)
	}
	s.processors[streamType] = p
}

func (s *Store) Append(ctx context.Context, events ...Event) error {
	if len(events) == 0 {
		return nil
	}

	// TODO avoid allocating
	eventMap := make(map[string][]Event)
	for _, event := range events {
		eventMap[event.StreamID] = append(eventMap[event.StreamID], event)
	}

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("eventstore: failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, e := range events {
		var rawData, rawMeta json.RawMessage
		if e.Data != nil {
			if rawData, err = json.Marshal(e.Data); err != nil {
				return fmt.Errorf("eventstore: failed to serialize event data: %w", err)
			}
		}
		if e.Meta != nil {
			if rawMeta, err = json.Marshal(e.Meta); err != nil {
				return fmt.Errorf("eventstore: failed to serialize event meta: %w", err)
			}
		}

		sql := `insert into events(id, stream_id, stream_type, type, data, meta)
		values ($1, $2, $3, $4, $5, $6)`

		if _, err = tx.Exec(ctx, sql, e.ID, e.StreamID, e.StreamType, e.Type, rawData, rawMeta); err != nil {
			return fmt.Errorf("eventstore: failed to insert event: %w", err)
		}
	}

	for streamID, events := range eventMap {
		// TODO check stream types all match
		streamType := events[0].StreamType
		lastEventID := events[len(events)-1].ID
		snap, err := s.getSnapshot(ctx, tx, streamID)
		if err != nil {
			return err
		}
		s.processorsMu.RLock()
		p, ok := s.processors[streamType]
		s.processorsMu.RUnlock()
		if !ok {
			return fmt.Errorf("eventstore: no processor for %s", streamType)
		}
		// TODO pass context
		var rawData json.RawMessage
		if snap != nil {
			rawData = snap.Data
		}
		if rawData, err = p.RawApply(rawData, events); err != nil {
			return err
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
