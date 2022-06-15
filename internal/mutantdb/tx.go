package mutantdb

import (
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

type Tx struct {
	pgTx        pgx.Tx
	idGenerator func() (string, error)
}

func (tx *Tx) Rollback(ctx context.Context) (err error) {
	return tx.pgTx.Rollback(ctx)
}

func (tx *Tx) Commit(ctx context.Context) (err error) {
	return tx.pgTx.Commit(ctx)
}

func (tx *Tx) Append(entityID string, mutations ...Mutation) *Append {
	return &Append{
		conn:        tx.pgTx,
		idGenerator: tx.idGenerator,
		entityID:    entityID,
		mutations:   mutations,
	}
}
