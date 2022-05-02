package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/ugent-library/biblio-backend/internal/models"
)

func (c *Client) GetDataset(id string) (*models.Dataset, error) {
	var data json.RawMessage
	err := c.db.QueryRow(context.Background(), "select data from datasets where data_to is null and id=$1", id).Scan(&data)
	if err != nil {
		return nil, err
	}

	d := &models.Dataset{}
	if err := json.Unmarshal(data, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (c *Client) GetDatasets(ids []string) ([]*models.Dataset, error) {
	var datasets []*models.Dataset

	pgIds := &pgtype.TextArray{}
	pgIds.Set(ids)
	rows, err := c.db.Query(context.Background(), "select data from datasets where data_to is null and id=any($1)", pgIds)
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

func (c *Client) SaveDataset(d *models.Dataset) (*models.Dataset, error) {
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
	tx, err := c.db.Begin(ctx)
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, "update datasets set data_to = $2 where id = $1 and data_to is null", d.ID, now); err != nil {
		return nil, err
	}
	if _, err = tx.Exec(ctx, "insert into datasets(id, data, data_from) values ($1, $2, $3)", d.ID, data, now); err != nil {
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return d, nil
}

func (c *Client) EachDataset(fn func(*models.Dataset) bool) error {
	rows, err := c.db.Query(context.Background(), "select data from datasets where data_to is null")
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
