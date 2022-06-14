package mutantdb

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

type Repository[T any] struct {
	store      *Store
	entityType *entityType[T]
}

func NewRepository[T any](store *Store, t *entityType[T]) *Repository[T] {
	return &Repository[T]{
		store:      store,
		entityType: t,
	}
}

func (r *Repository[T]) Get(ctx context.Context, id string) (Projection[T], error) {
	var (
		p       Projection[T]
		rawData json.RawMessage
	)

	p.ID = id

	sql := `select mutation_id, entity_data, date_created, date_updated
	from projections
	where entity_id = $1 and entity_type = $2
	limit 1`

	err := r.store.conn.
		QueryRow(ctx, sql, id, r.entityType.name).
		Scan(&p.MutationID, &rawData, &p.DateCreated, &p.DateUpdated)

	if errors.Is(err, pgx.ErrNoRows) {
		return p, ErrNotFound
	} else if err != nil {
		return p, fmt.Errorf("mutantdb: failed to scan projection: %w", err)
	}

	if err := json.Unmarshal(rawData, &p.Data); err != nil {
		return p, fmt.Errorf("mutantdb: failed to deserialize entity data: %w", err)
	}

	return p, nil
}

func (r *Repository[T]) GetAt(ctx context.Context, id, mutationID string) (Projection[T], error) {
	// TODO: set timestamp fields
	p := Projection[T]{
		ID:         id,
		MutationID: mutationID,
		Data:       r.entityType.factory(),
	}

	sql := `select mutation_name, mutation_data
	from mutations
	where entity_id = $1 and entity_type = $2 and
	seq <= (select seq from mutations where entity_id = $1 and entity_type = $2 and mutation_id = $3)`

	rows, err := r.store.conn.Query(ctx, sql, id, r.entityType.name, mutationID)
	if err != nil {
		return p, fmt.Errorf("mutantdb: failed to query mutations: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			name string
			data json.RawMessage
		)

		if err := rows.Scan(&name, &data); err != nil {
			return p, fmt.Errorf("mutantdb: failed to scan mutation: %w", err)
		}

		h := r.store.GetMutator(r.entityType.name, name)
		if h == nil {
			return p, fmt.Errorf("mutantdb: mutator %s %s not found", r.entityType.name, name)
		}

		// TODO improve this
		d, err := h.Apply(p.Data, data)
		if err != nil {
			return p, err
		}
		if t, ok := d.(T); ok {
			p.Data = t
		} else {
			return p, fmt.Errorf("mutantdb: invalid entity data type %T", t)
		}
	}

	if err := rows.Err(); err != nil {
		return p, fmt.Errorf("mutantdb: %w", err)
	}

	return p, nil
}

func (r *Repository[T]) GetAll(ctx context.Context) (Cursor[T], error) {
	sql := `select entity_id, mutation_id, entity_data, date_created, date_updated
	from projections
	where entity_type = $1`

	var (
		c   Cursor[T]
		err error
	)

	c.rows, err = r.store.conn.Query(ctx, sql, r.entityType.name)
	if err != nil {
		return c, fmt.Errorf("mutantdb: failed to query projections: %w", err)
	}

	return c, nil
}
