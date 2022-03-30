package pg

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
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
