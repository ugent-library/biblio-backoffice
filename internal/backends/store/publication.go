package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
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

func (c *Client) GetPublications(ids []string) ([]*models.Publication, error) {
	var publications []*models.Publication

	pgIds := &pgtype.TextArray{}
	pgIds.Set(ids)
	rows, err := c.db.Query(context.Background(), "select data from publications where data_to is null and id=any($1)", pgIds)
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

func (c *Client) SavePublication(p *models.Publication) (*models.Publication, error) {
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
	tx, err := c.db.Begin(ctx)
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, "update publications set data_to = $2 where id = $1 and data_to is null", p.ID, now); err != nil {
		return nil, err
	}
	if _, err = tx.Exec(ctx, "insert into publications(id, data, data_from) values ($1, $2, $3)", p.ID, data, now); err != nil {
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return p, nil
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
