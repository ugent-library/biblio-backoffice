package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
)

type dbOrTx interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
}

type Store struct {
	db dbOrTx
}

func New(dsn string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Atomic(ctx context.Context, fn func(backends.Store) error) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(&Store{db: tx}); err == nil {
		tx.Commit(ctx)
	}

	return nil
}

func (s *Store) GetPublication(id string) (*models.Publication, error) {
	var data json.RawMessage
	err := s.db.QueryRow(context.Background(), "select data from publications where data_to is null and id=$1", id).Scan(&data)
	if err != nil {
		return nil, err
	}

	d := &models.Publication{}
	if err := json.Unmarshal(data, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Store) GetPublications(ids []string) ([]*models.Publication, error) {
	var publications []*models.Publication

	pgIds := &pgtype.TextArray{}
	pgIds.Set(ids)
	rows, err := s.db.Query(context.Background(), "select data from publications where data_to is null and id=any($1)", pgIds)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var data json.RawMessage
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}

		d := &models.Publication{}
		if err := json.Unmarshal(data, d); err != nil {
			return nil, err
		}

		publications = append(publications, d)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return publications, nil
}

func (s *Store) StorePublication(p *models.Publication) error {
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
		return err
	}

	ctx := context.Background()
	tx, err := s.db.Begin(ctx)
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, "update publications set data_to = $2 where id = $1 and data_to is null", p.ID, now); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, "insert into publications(id, data, data_from) values ($1, $2, $3)", p.ID, data, now); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Store) EachPublication(fn func(*models.Publication) bool) error {
	rows, err := s.db.Query(context.Background(), "select data from publications where data_to is null")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var data json.RawMessage
		if err := rows.Scan(&data); err != nil {
			return err
		}

		d := &models.Publication{}
		if err := json.Unmarshal(data, d); err != nil {
			return err
		}

		if ok := fn(d); !ok {
			return nil
		}
	}

	return rows.Err()
}

func (s *Store) GetDataset(id string) (*models.Dataset, error) {
	var data json.RawMessage
	err := s.db.QueryRow(context.Background(), "select data from datasets where data_to is null and id=$1", id).Scan(&data)
	if err != nil {
		return nil, err
	}

	d := &models.Dataset{}
	if err := json.Unmarshal(data, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Store) GetDatasets(ids []string) ([]*models.Dataset, error) {
	var datasets []*models.Dataset

	pgIds := &pgtype.TextArray{}
	pgIds.Set(ids)
	rows, err := s.db.Query(context.Background(), "select data from datasets where data_to is null and id=any($1)", pgIds)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var data json.RawMessage
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}

		d := &models.Dataset{}
		if err := json.Unmarshal(data, d); err != nil {
			return nil, err
		}

		datasets = append(datasets, d)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return datasets, nil
}

func (s *Store) StoreDataset(d *models.Dataset) error {
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
		return err
	}

	ctx := context.Background()
	tx, err := s.db.Begin(ctx)
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, "update datasets set data_to = $2 where id = $1 and data_to is null", d.ID, now); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, "insert into datasets(id, data, data_from) values ($1, $2, $3)", d.ID, data, now); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Store) EachDataset(fn func(*models.Dataset) bool) error {
	rows, err := s.db.Query(context.Background(), "select data from datasets where data_to is null")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var data json.RawMessage
		if err := rows.Scan(&data); err != nil {
			return err
		}

		d := &models.Dataset{}
		if err := json.Unmarshal(data, d); err != nil {
			return err
		}

		if ok := fn(d); !ok {
			return nil
		}
	}

	return rows.Err()
}
