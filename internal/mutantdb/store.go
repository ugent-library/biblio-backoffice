package mutantdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

var ErrNotFound = errors.New("not found")

type ErrConflict struct {
	CurrentMutationID, ExpectedMutationID string
}

func (e *ErrConflict) Error() string {
	return "conflict detected"
}

type Conn interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
}

type Projection[T any] struct {
	ID          string
	Data        T
	MutationID  string
	DateCreated time.Time
	DateUpdated time.Time
}

type mutatorMap[T any] struct {
	data map[string]Mutator[T]
	mu   sync.RWMutex
}

func (mm *mutatorMap[T]) Add(m Mutator[T]) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if mm.data == nil {
		mm.data = make(map[string]Mutator[T])
	}

	mm.data[m.Name()] = m
}

func (mm *mutatorMap[T]) Get(name string) Mutator[T] {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	if m, ok := mm.data[name]; ok {
		return m
	}
	return nil
}

type store[T any] struct {
	conn        Conn
	idGenerator func() (string, error)
	entityType  *Type[T]
	mutators    *mutatorMap[T]
}

type Store[T any] struct {
	store[T]
}

func NewStore[T any](conn Conn, t *Type[T]) *Store[T] {
	return &Store[T]{store[T]{
		conn:        conn,
		entityType:  t,
		idGenerator: generateUUID,
		mutators:    &mutatorMap[T]{},
	}}
}

func (s *Store[T]) WithIDGenerator(fn func() (string, error)) *Store[T] {
	s.idGenerator = fn
	return s
}

func (s *Store[T]) WithMutators(mutators ...Mutator[T]) *Store[T] {
	for _, m := range mutators {
		s.mutators.Add(m)
	}
	return s
}

func (s *store[T]) Conn() Conn {
	return s.conn
}

func (s *store[T]) Tx(tx pgx.Tx) *store[T] {
	return &store[T]{
		conn:        tx,
		entityType:  s.entityType,
		idGenerator: s.idGenerator,
		mutators:    s.mutators,
	}
}

func (s *store[T]) AppendAfter(ctx context.Context, entityID, mutationID string, mutations ...Mutation[T]) (Projection[T], error) {
	return s.append(ctx, entityID, mutationID, mutations)
}

func (s *store[T]) Append(ctx context.Context, entityID string, mutations ...Mutation[T]) (Projection[T], error) {
	return s.append(ctx, entityID, "", mutations)
}

func (s *store[T]) Get(ctx context.Context, id string) (Projection[T], error) {
	var (
		p       Projection[T]
		rawData json.RawMessage
	)

	p.ID = id

	sql := `select mutation_id, entity_data, date_created, date_updated
	from projections
	where entity_id = $1 and entity_type = $2
	limit 1`

	err := s.conn.
		QueryRow(ctx, sql, id, s.entityType.Name()).
		Scan(&p.MutationID, &rawData, &p.DateCreated, &p.DateUpdated)

	if errors.Is(err, pgx.ErrNoRows) {
		return p, ErrNotFound
	} else if err != nil {
		return p, fmt.Errorf("mutantdb: failed to scan projection: %w", err)
	}

	if err = json.Unmarshal(rawData, &p.Data); err != nil {
		return p, fmt.Errorf("mutantdb: failed to deserialize entity data: %w", err)
	}

	return p, nil
}

func (s *store[T]) GetAt(ctx context.Context, id, mutationID string) (Projection[T], error) {
	p := Projection[T]{
		ID:         id,
		MutationID: mutationID,
		Data:       s.entityType.New(),
	}

	sql := `select mutation_name, mutation_data, date_created
	from mutations
	where entity_id = $1 and entity_type = $2 and
	seq <= (select seq from mutations where entity_id = $1 and entity_type = $2 and mutation_id = $3)`

	rows, err := s.conn.Query(ctx, sql, id, s.entityType.Name(), mutationID)
	if err != nil {
		return p, fmt.Errorf("mutantdb: failed to query mutations: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			err         error
			name        string
			mutData     json.RawMessage
			dateCreated time.Time
		)

		if err = rows.Scan(&name, &mutData, &dateCreated); err != nil {
			return p, fmt.Errorf("mutantdb: failed to scan mutation: %w", err)
		}

		if p.DateCreated.IsZero() {
			p.DateCreated = dateCreated
		}
		p.DateUpdated = dateCreated

		m := s.mutators.Get(name)
		if m == nil {
			return p, fmt.Errorf("mutantdb: mutator %s %s not found", s.entityType.Name(), name)
		}

		p.Data, err = m.Apply(p.Data, mutData)
		if err != nil {
			return p, err
		}
	}

	if err := rows.Err(); err != nil {
		return p, fmt.Errorf("mutantdb: %w", err)
	}

	return p, nil
}

func (s *store[T]) GetAll(ctx context.Context) (Cursor[T], error) {
	sql := `select entity_id, mutation_id, entity_data, date_created, date_updated
	from projections
	where entity_type = $1`

	var (
		c   Cursor[T]
		err error
	)

	c.rows, err = s.conn.Query(ctx, sql, s.entityType.Name())
	if err != nil {
		return c, fmt.Errorf("mutantdb: failed to query projections: %w", err)
	}

	return c, nil
}

func (s *store[T]) append(ctx context.Context, entityID, expectedMutID string, mutations []Mutation[T]) (Projection[T], error) {
	//--- projection to return

	p := Projection[T]{
		ID: entityID,
	}

	//--- check args

	if len(mutations) == 0 {
		return p, fmt.Errorf("mutantdb: no mutations to append")
	}

	//--- start tx

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return p, fmt.Errorf("mutantdb: failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	//--- tmp vars

	var (
		newEntity           bool
		rawEntityData       json.RawMessage
		entityData          T
		firstMutDateCreated time.Time
		lastMutID           string
		lastMutDateCreated  time.Time
	)

	//--- get current projection data

	err = tx.QueryRow(ctx,
		`select mutation_id, entity_data, date_created from projections where entity_id = $1 and entity_type = $2 limit 1`,
		entityID, s.entityType.Name(),
	).Scan(&lastMutID, &rawEntityData, &firstMutDateCreated)

	if err == pgx.ErrNoRows {
		newEntity = true
		entityData = s.entityType.New()
	} else if err != nil {
		return p, fmt.Errorf("mutantdb: failed to get projection: %w", err)
	} else {
		if err := json.Unmarshal(rawEntityData, &entityData); err != nil {
			return p, fmt.Errorf("mutantdb: failed to deserialize entity data: %w", err)
		}
	}

	//--- detect conflicts

	if expectedMutID != "" && expectedMutID != lastMutID {
		return p, &ErrConflict{
			CurrentMutationID:  lastMutID,
			ExpectedMutationID: expectedMutID,
		}
	}

	//--- insert mutations

	for _, mut := range mutations {
		var rawMutData, rawMutMeta json.RawMessage

		if mut.Data() != nil {
			if rawMutData, err = json.Marshal(mut.Data()); err != nil {
				return p, fmt.Errorf("mutantdb: failed to serialize mutation data: %w", err)
			}
		}
		if mut.Meta() != nil {
			if rawMutMeta, err = json.Marshal(mut.Meta()); err != nil {
				return p, fmt.Errorf("mutantdb: failed to serialize mutation meta: %w", err)
			}
		}

		// generate mutation id
		lastMutID, err = s.idGenerator()
		if err != nil {
			return p, fmt.Errorf("mutantdb: failed to generate id: %w", err)
		}

		// TODO insert all mutations in one statement
		if err = tx.QueryRow(ctx,
			`insert into mutations (mutation_id, entity_id, entity_type, mutation_name, mutation_data, mutation_meta)
		values ($1, $2, $3, $4, $5, $6)
		returning date_created`,
			lastMutID, entityID, s.entityType.Name(), mut.Name(), rawMutData, rawMutMeta,
		).Scan(&lastMutDateCreated); err != nil {
			return p, fmt.Errorf("mutantdb: failed to insert mutation: %w", err)
		}
	}

	//--- apply mutations

	for _, mut := range mutations {
		if entityData, err = mut.Apply(entityData); err != nil {
			return p, fmt.Errorf("mutantdb: failed to apply mutation %s: %w", mut.Name(), err)
		}
	}

	//--- upsert projection

	rawEntityData, err = json.Marshal(entityData)
	if err != nil {
		return p, fmt.Errorf("mutantdb: failed to serialize projection data: %w", err)
	}

	if newEntity {
		sql := `insert into projections (entity_id, entity_type, mutation_id, entity_data, date_created, date_updated)
		values ($1, $2, $3, $4, $5, $6)`
		if _, err = tx.Exec(ctx, sql, entityID, s.entityType.Name(), lastMutID, rawEntityData,
			lastMutDateCreated, lastMutDateCreated); err != nil {
			return p, fmt.Errorf("mutantdb: failed to insert projection: %w", err)
		}
	} else {
		sql := `update projections set mutation_id = $1, entity_data = $2, date_updated = $3 where entity_id = $4 and entity_type = $5`
		if _, err = tx.Exec(ctx, sql, lastMutID, rawEntityData, lastMutDateCreated,
			entityID, s.entityType.Name()); err != nil {
			return p, fmt.Errorf("mutantdb: failed to update projection: %w", err)
		}
	}

	//--- commit tx

	if err = tx.Commit(ctx); err != nil {
		return p, fmt.Errorf("mutantdb: failed to commit transaction: %w", err)
	}

	//--- return projection

	p.Data = entityData
	p.MutationID = lastMutID
	p.DateCreated = firstMutDateCreated
	p.DateUpdated = lastMutDateCreated

	return p, nil
}

func generateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
