package pg

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ugent-library/biblio-backend/internal/models"
)

type Client struct {
	db *pgxpool.Pool
}

func New(dsn string) (*Client, error) {
	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &Client{db: db}, nil
}

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

func (c *Client) CreateDataset(d *models.Dataset) (*models.Dataset, error) {
	now := time.Now()
	d.ID = uuid.NewString()
	d.DateUpdated = &now
	d.DateCreated = &now

	data, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	_, err = c.db.Exec(ctx, "insert into datasets(id, data, data_from) values ($1, $2, $3)", d.ID, data, now)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (c *Client) UpdateDataset(d *models.Dataset) (*models.Dataset, error) {
	now := time.Now()
	d.DateUpdated = &now

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
