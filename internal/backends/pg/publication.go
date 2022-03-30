package pg

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ugent-library/biblio-backend/internal/models"
)

func (c *Client) GetPublication(id string) (*models.Publication, error) {
	var data json.RawMessage
	err := c.db.QueryRow(context.Background(), "select data from publications where data_to is null and id=$1", id).Scan(&data)
	if err != nil {
		return nil, err
	}

	d := &models.Publication{}
	if err := json.Unmarshal(data, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (c *Client) CreatePublication(d *models.Publication) (*models.Publication, error) {
	now := time.Now()
	d.ID = uuid.NewString()
	d.DateUpdated = &now
	d.DateCreated = &now

	data, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	_, err = c.db.Exec(ctx, "insert into publications(id, data, data_from) values ($1, $2, $3)", d.ID, data, now)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (c *Client) UpdatePublication(d *models.Publication) (*models.Publication, error) {
	now := time.Now()
	d.DateUpdated = &now

	data, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	tx, err := c.db.Begin(ctx)
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, "update publications set data_to = $2 where id = $1 and data_to is null", d.ID, now); err != nil {
		return nil, err
	}
	if _, err = tx.Exec(ctx, "insert into publications(id, data, data_from) values ($1, $2, $3)", d.ID, data, now); err != nil {
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return d, nil
}

func (c *Client) EachPublication(fn func(*models.Publication) bool) error {
	rows, err := c.db.Query(context.Background(), "select data from publications where data_to is null")
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
