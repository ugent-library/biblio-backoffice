package mutantdb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

type Append struct {
	conn            PgConn
	idGenerator     func() (string, error)
	entityID        string
	mutations       []Mutation
	afterMutationID string
}

func (op *Append) AfterMutation(mutID string) *Append {
	op.afterMutationID = mutID
	return op
}

func (op *Append) Do(ctx context.Context) error {
	if len(op.mutations) == 0 {
		return nil
	}

	tx, err := op.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("mutantdb: failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var (
		entityName         = op.mutations[0].EntityType().Name()
		rawMutData         json.RawMessage
		rawMutMeta         json.RawMessage
		lastMutDateCreated time.Time
		lastMutID          string
		newEntity          bool
		rawEntityData      json.RawMessage
		entityData         any
	)

	//--- get current projection data

	err = tx.QueryRow(ctx,
		`select mutation_id, entity_data from projections where entity_id = $1 and entity_type = $2 limit 1`,
		op.entityID, entityName,
	).Scan(&lastMutID, &rawEntityData)

	if err == pgx.ErrNoRows {
		newEntity = true
		entityData = op.mutations[0].EntityType().New()
	} else if err != nil {
		return fmt.Errorf("mutantdb: failed to get projection: %w", err)
	} else {
		entityData = rawEntityData
	}

	//--- detect conflicts

	if op.afterMutationID != "" && op.afterMutationID != lastMutID {
		return &ErrConflict{
			CurrentMutationID:  lastMutID,
			ExpectedMutationID: op.afterMutationID,
		}
	}

	//--- insert mutations

	for _, mut := range op.mutations {
		if mut.EntityType().Name() != entityName {
			return fmt.Errorf("mutantdb: cannot apply mutations to different entities: %s != %s", mut.EntityType().Name(), entityName)
		}

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
		lastMutID, err = op.idGenerator()
		if err != nil {
			return fmt.Errorf("mutantdb: failed to generate id: %w", err)
		}

		// TODO insert all mutations in one statement
		if err = tx.QueryRow(ctx,
			`insert into mutations (mutation_id, entity_id, entity_type, mutation_name, mutation_data, mutation_meta)
		values ($1, $2, $3, $4, $5, $6)
		returning date_created`,
			lastMutID, op.entityID, entityName, mut.Name(), rawMutData, rawMutMeta,
		).Scan(&lastMutDateCreated); err != nil {
			return fmt.Errorf("mutantdb: failed to insert mutation: %w", err)
		}
	}

	//--- apply mutations

	for _, mut := range op.mutations {
		if entityData, err = mut.Apply(entityData); err != nil {
			return fmt.Errorf("mutantdb: failed to apply mutation %s: %w", mut.EntityType().Name(), err)
		}
	}

	//--- upsert projection

	rawEntityData, err = json.Marshal(entityData)
	if err != nil {
		return fmt.Errorf("mutantdb: failed to serialize projection data: %w", err)
	}

	if newEntity {
		sql := `insert into projections (entity_id, entity_type, mutation_id, entity_data, date_created, date_updated)
		values ($1, $2, $3, $4, $5, $6)`
		if _, err = tx.Exec(ctx, sql, op.entityID, entityName, lastMutID, rawEntityData,
			lastMutDateCreated, lastMutDateCreated); err != nil {
			return fmt.Errorf("mutantdb: failed to insert projection: %w", err)
		}
	} else {
		sql := `update projections set mutation_id = $1, entity_data = $2, date_updated = $3 where entity_id = $4 and entity_type = $5`
		if _, err = tx.Exec(ctx, sql, lastMutID, rawEntityData, lastMutDateCreated, op.entityID, entityName); err != nil {
			return fmt.Errorf("mutantdb: failed to update projection: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("mutantdb: failed to commit transaction: %w", err)
	}

	return nil
}
