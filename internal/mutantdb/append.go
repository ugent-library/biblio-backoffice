package mutantdb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

type Append[T any] struct {
	conn        PgConn
	idGenerator func() (string, error)
	entityID    string
	entityType  *Type[T]
	mutations   []Mutation[T]

	// options
	tx            pgx.Tx
	expectedMutID string

	// data needed to return projections
	entityData          T
	firstMutDateCreated time.Time
	lastMutID           string
	lastMutDateCreated  time.Time
}

func (op *Append[T]) Tx(tx pgx.Tx) *Append[T] {
	op.tx = tx
	return op
}

func (op *Append[T]) After(mutID string) *Append[T] {
	op.expectedMutID = mutID
	return op
}

func (op *Append[T]) Do(ctx context.Context) error {
	if len(op.mutations) == 0 {
		return fmt.Errorf("mutantdb: no mutations to append")
	}

	conn := op.conn
	if op.tx != nil {
		conn = op.tx
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("mutantdb: failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var (
		rawMutData    json.RawMessage
		rawMutMeta    json.RawMessage
		newEntity     bool
		rawEntityData json.RawMessage
	)

	//--- get current projection data

	err = tx.QueryRow(ctx,
		`select mutation_id, entity_data, date_created from projections where entity_id = $1 and entity_type = $2 limit 1`,
		op.entityID, op.entityType.Name(),
	).Scan(&op.lastMutID, &rawEntityData, &op.firstMutDateCreated)

	if err == pgx.ErrNoRows {
		newEntity = true
		op.entityData = op.entityType.New()
	} else if err != nil {
		return fmt.Errorf("mutantdb: failed to get projection: %w", err)
	} else {
		if err := json.Unmarshal(rawEntityData, &op.entityData); err != nil {
			return fmt.Errorf("mutantdb: failed to deserialize entity data: %w", err)
		}
	}

	//--- detect conflicts

	if op.expectedMutID != "" && op.expectedMutID != op.lastMutID {
		return &ErrConflict{
			CurrentMutationID:  op.lastMutID,
			ExpectedMutationID: op.expectedMutID,
		}
	}

	//--- insert mutations

	for _, mut := range op.mutations {
		if mut.Data() != nil {
			if rawMutData, err = json.Marshal(mut.Data()); err != nil {
				return fmt.Errorf("mutantdb: failed to serialize mutation data: %w", err)
			}
		}
		if mut.Meta() != nil {
			if rawMutMeta, err = json.Marshal(mut.Meta()); err != nil {
				return fmt.Errorf("mutantdb: failed to serialize mutation meta: %w", err)
			}
		}

		// generate mutation id
		op.lastMutID, err = op.idGenerator()
		if err != nil {
			return fmt.Errorf("mutantdb: failed to generate id: %w", err)
		}

		// TODO insert all mutations in one statement
		if err = tx.QueryRow(ctx,
			`insert into mutations (mutation_id, entity_id, entity_type, mutation_name, mutation_data, mutation_meta)
		values ($1, $2, $3, $4, $5, $6)
		returning date_created`,
			op.lastMutID, op.entityID, op.entityType.Name(), mut.Name(), rawMutData, rawMutMeta,
		).Scan(&op.lastMutDateCreated); err != nil {
			return fmt.Errorf("mutantdb: failed to insert mutation: %w", err)
		}
	}

	//--- apply mutations

	for _, mut := range op.mutations {
		if op.entityData, err = mut.Apply(op.entityData); err != nil {
			return fmt.Errorf("mutantdb: failed to apply mutation %s: %w", mut.Name(), err)
		}
	}

	//--- upsert projection

	rawEntityData, err = json.Marshal(op.entityData)
	if err != nil {
		return fmt.Errorf("mutantdb: failed to serialize projection data: %w", err)
	}

	if newEntity {
		sql := `insert into projections (entity_id, entity_type, mutation_id, entity_data, date_created, date_updated)
		values ($1, $2, $3, $4, $5, $6)`
		if _, err = tx.Exec(ctx, sql, op.entityID, op.entityType.Name(), op.lastMutID, rawEntityData,
			op.lastMutDateCreated, op.lastMutDateCreated); err != nil {
			return fmt.Errorf("mutantdb: failed to insert projection: %w", err)
		}
	} else {
		sql := `update projections set mutation_id = $1, entity_data = $2, date_updated = $3 where entity_id = $4 and entity_type = $5`
		if _, err = tx.Exec(ctx, sql, op.lastMutID, rawEntityData, op.lastMutDateCreated,
			op.entityID, op.entityType.Name()); err != nil {
			return fmt.Errorf("mutantdb: failed to update projection: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("mutantdb: failed to commit transaction: %w", err)
	}

	return nil
}

func (op *Append[T]) Get(ctx context.Context) (Projection[T], error) {
	p := Projection[T]{
		ID: op.entityID,
	}

	if err := op.Do(ctx); err != nil {
		return p, err
	}

	p.Data = op.entityData
	p.MutationID = op.lastMutID
	p.DateCreated = op.firstMutDateCreated
	p.DateUpdated = op.lastMutDateCreated

	return p, nil
}
