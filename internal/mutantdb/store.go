package mutantdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var NotFound = errors.New("not found")

var DefaultIDGenerator = func() (string, error) {
	return uuid.NewString(), nil
}

type PgConn interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
}

type Store struct {
	conn        PgConn
	idGenerator func() (string, error)
	mutators    map[string]map[string]Mutator
	mutatorsMu  sync.RWMutex
}

func Connect(ctx context.Context, dsn string, opts ...func(*Store)) (*Store, error) {
	conn, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("mutantdb: failed to connect: %w", err)
	}

	return New(conn, opts...), nil
}

func New(conn PgConn, opts ...func(*Store)) *Store {
	s := &Store{
		conn:        conn,
		idGenerator: DefaultIDGenerator,
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

func WithMutators(mutators ...Mutator) func(*Store) {
	return func(s *Store) {
		for _, m := range mutators {
			s.AddMutator(m)
		}
	}
}

func (s *Store) AddMutator(m Mutator) *Store {
	s.mutatorsMu.Lock()
	defer s.mutatorsMu.Unlock()

	if s.mutators == nil {
		s.mutators = make(map[string]map[string]Mutator)
	}
	if s.mutators[m.EntityName()] == nil {
		s.mutators[m.EntityName()] = make(map[string]Mutator)
	}
	s.mutators[m.EntityName()][m.Name()] = m

	return s
}

func (s *Store) GetMutator(entityType, name string) Mutator {
	s.mutatorsMu.RLock()
	defer s.mutatorsMu.RUnlock()
	if m, ok := s.mutators[entityType]; ok {
		if h, ok := m[name]; ok {
			return h
		}
	}
	return nil
}

func (s *Store) Append(ctx context.Context, mutations ...Mutation) error {
	if len(mutations) == 0 {
		return nil
	}

	// TODO avoid allocating
	// TODO refactor, entity id is only unique per entity type
	mutationMap := make(map[string][]Mutation)
	for _, m := range mutations {
		mutationMap[m.EntityID()] = append(mutationMap[m.EntityID()], m)
	}

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("mutantdb: failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var lastMutationID string

	for _, e := range mutations {
		var rawData, rawMeta json.RawMessage
		if e.Data() != nil {
			if rawData, err = json.Marshal(e.Data()); err != nil {
				return fmt.Errorf("mutantdb: failed to serialize mutation data: %w", err)
			}
		}
		if e.Meta() != nil {
			if rawMeta, err = json.Marshal(e.Meta()); err != nil {
				return fmt.Errorf("mutantdb: failed to serialize mutation meta: %w", err)
			}
		}

		// generate mutation mutationID
		mutationID, err := s.idGenerator()
		if err != nil {
			return fmt.Errorf("mutantdb: failed to generate id: %w", err)
		}

		lastMutationID = mutationID

		sql := `insert into mutations (mutation_id, entity_id, entity_type, mutation_name, mutation_data, mutation_meta)
		values ($1, $2, $3, $4, $5, $6)`

		if _, err = tx.Exec(ctx, sql, mutationID, e.EntityID(), e.EntityType().Name(), e.Name(), rawData, rawMeta); err != nil {
			return fmt.Errorf("mutantdb: failed to insert mutation: %w", err)
		}
	}

	for entityID, mutations := range mutationMap {
		entityType := mutations[0].EntityType()
		snap, err := s.getProjection(ctx, tx, entityID)
		if err != nil {
			return err
		}

		var d any
		if snap == nil {
			d = entityType.New()
		} else {
			d = snap.Data
		}

		for _, e := range mutations {
			if d, err = e.Apply(d); err != nil {
				return fmt.Errorf("mutantdb: failed to apply mutation: %w", err)
			}
		}

		rawData, err := json.Marshal(d)
		if err != nil {
			return fmt.Errorf("mutantdb: failed to serialize projection data: %w", err)
		}

		// TODO set date_updated to last date_created of mutations
		now := time.Now()
		if snap == nil {
			sql := `insert into projections (entity_id, entity_type, mutation_id, entity_data, date_created, date_updated)
		values ($1, $2, $3, $4, $5, $6)`
			if _, err = tx.Exec(ctx, sql, entityID, entityType.Name(), lastMutationID, rawData, now, now); err != nil {
				return fmt.Errorf("mutantdb: failed to insert projection: %w", err)
			}
		} else {
			// TODO check row count or use one on conflict statement
			sql := `update projections set mutation_id = $1, mutation_data = $2, date_updated = $3 where entity_id = $4`
			if _, err = tx.Exec(ctx, sql, lastMutationID, rawData, now, entityID); err != nil {
				return fmt.Errorf("mutantdb: failed to update projection: %w", err)
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("mutantdb: failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Store) getProjection(ctx context.Context, tx PgConn, entityID string) (*RawProjection, error) {
	p := RawProjection{ID: entityID}

	sql := `select entity_type, mutation_id, entity_data, date_created, date_updated from projections
	where entity_id = $1
	limit 1`
	err := tx.QueryRow(ctx, sql, entityID).Scan(&p.Type, &p.MutationID, &p.Data, &p.DateCreated, &p.DateUpdated)

	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("mutantdb: failed to get projection: %w", err)
	}

	return &p, nil
}
