package mutantdb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

type entityMutations struct {
	entityID        string
	entityName      string
	lastDateCreated time.Time
	lastMutationID  string
	mutations       []Mutation
}

// TODO DetectConflict() sets conflict detection
// TODO Get() returns projections
type Append struct {
	store          *Store
	detectConflict bool
	entities       []*entityMutations
	mutations      []Mutation
}

func NewAppend(store *Store, mutations ...Mutation) *Append {
	op := &Append{
		store:     store,
		mutations: mutations,
	}

	for _, m := range mutations {
		var found bool
		for _, e := range op.entities {
			if e.entityName == m.EntityType().Name() && e.entityID == m.EntityID() {
				e.mutations = append(e.mutations, m)
				found = true
				break
			}
		}
		if !found {
			op.entities = append(op.entities, &entityMutations{
				entityID:   m.EntityID(),
				entityName: m.EntityType().Name(),
				mutations:  []Mutation{m},
			})
		}
	}

	return op
}

// TODO bulk insert mutations https://github.com/jackc/pgx/issues/764
func (op *Append) Do(ctx context.Context) error {
	if len(op.mutations) == 0 {
		return nil
	}

	tx, err := op.store.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("mutantdb: failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// insert mutations
	for _, m := range op.mutations {
		var rawData, rawMeta json.RawMessage
		if m.Data() != nil {
			if rawData, err = json.Marshal(m.Data()); err != nil {
				return fmt.Errorf("mutantdb: failed to serialize mutation data: %w", err)
			}
		}
		if m.Meta() != nil {
			if rawMeta, err = json.Marshal(m.Meta()); err != nil {
				return fmt.Errorf("mutantdb: failed to serialize mutation meta: %w", err)
			}
		}

		// generate mutation id
		mutationID, err := op.store.idGenerator()
		if err != nil {
			return fmt.Errorf("mutantdb: failed to generate id: %w", err)
		}

		var dateCreated time.Time

		if err = tx.QueryRow(ctx,
			`insert into mutations (mutation_id, entity_id, entity_type, mutation_name, mutation_data, mutation_meta)
			values ($1, $2, $3, $4, $5, $6)
			returning date_created`,
			mutationID, m.EntityID(), m.EntityType().Name(), m.Name(), rawData, rawMeta,
		).Scan(&dateCreated); err != nil {
			return fmt.Errorf("mutantdb: failed to insert mutation: %w", err)
		}

		// remember last mutation id and date created for each entity
		for _, e := range op.entities {
			if e.entityID == m.EntityID() && e.entityName == m.EntityType().Name() {
				e.lastDateCreated = dateCreated
				e.lastMutationID = mutationID
				break
			}
		}
	}

	// upsert projections
	for _, e := range op.entities {
		// get projection data
		var (
			newEntity bool
			rawData   json.RawMessage
			d         any
			err       error
		)

		err = tx.QueryRow(ctx,
			`select entity_data from projections where entity_id = $1 and entity_type = $2 limit 1`,
			e.entityID, e.entityName,
		).Scan(&rawData)

		if err == pgx.ErrNoRows {
			newEntity = true
			d = e.mutations[0].EntityType().New()
		} else if err != nil {
			return fmt.Errorf("mutantdb: failed to get projection: %w", err)
		} else {
			d = rawData
		}

		// apply mutations
		for _, m := range e.mutations {
			if d, err = m.Apply(d); err != nil {
				return fmt.Errorf("mutantdb: failed to apply mutation: %w", err)
			}
		}

		// upsert projection
		if newEntity {
			sql := `insert into projections (entity_id, entity_type, mutation_id, entity_data, date_created, date_updated)
			values ($1, $2, $3, $4, $5, $6)`
			if _, err = tx.Exec(ctx, sql, e.entityID, e.entityName, e.lastMutationID, d, e.lastDateCreated, e.lastDateCreated); err != nil {
				return fmt.Errorf("mutantdb: failed to insert projection: %w", err)
			}
		} else {
			sql := `update projections set mutation_id = $1, entity_data = $2, date_updated = $3 where entity_id = $4 and entity_type = $5`
			if _, err = tx.Exec(ctx, sql, e.lastMutationID, d, e.lastDateCreated, e.entityID, e.entityName); err != nil {
				return fmt.Errorf("mutantdb: failed to update projection: %w", err)
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("mutantdb: failed to commit transaction: %w", err)
	}

	return nil
}
