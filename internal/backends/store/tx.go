package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
)

type Tx struct {
	db *pgxpool.Pool
	tx pgx.Tx
}

func (c *Client) Begin() (backends.Transaction, error) {
	tx, err := c.db.Begin(context.Background())
	if err != nil {
		return &Tx{}, err
	}
	return &Tx{db: c.db, tx: tx}, nil
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback(context.Background())
}

func (t *Tx) Commit() error {
	return t.tx.Commit(context.Background())
}

func (t *Tx) SavePublication(p *models.Publication) (*models.Publication, error) {
	now := time.Now()

	if p.DateCreated == nil {
		p.DateCreated = &now
	}
	p.DateUpdated = &now

	if p.ID == "" {
		p.ID = uuid.NewString()
	}

	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	if _, err = t.tx.Exec(ctx, "update publications set data_to = $2 where id = $1 and data_to is null", p.ID, now); err != nil {
		return nil, err
	}
	if _, err = t.tx.Exec(ctx, "insert into publications(id, data, data_from) values ($1, $2, $3)", p.ID, data, now); err != nil {
		return nil, err
	}

	return p, nil
}

func (t *Tx) SaveDataset(d *models.Dataset) (*models.Dataset, error) {
	now := time.Now()

	if d.DateCreated == nil {
		d.DateCreated = &now
	}
	d.DateUpdated = &now

	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	data, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	if _, err = t.tx.Exec(ctx, "update datasets set data_to = $2 where id = $1 and data_to is null", d.ID, now); err != nil {
		return nil, err
	}
	if _, err = t.tx.Exec(ctx, "insert into datasets(id, data, data_from) values ($1, $2, $3)", d.ID, data, now); err != nil {
		return nil, err
	}

	return d, nil
}
