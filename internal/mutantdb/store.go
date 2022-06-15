package mutantdb

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrNotFound = errors.New("not found")

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

func (s *Store) Append(mutations ...Mutation) *Append {
	return NewAppend(s, mutations...)
}
