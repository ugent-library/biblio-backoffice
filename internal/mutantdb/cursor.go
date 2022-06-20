package mutantdb

import (
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type Cursor[T any] struct {
	rows pgx.Rows
}

func (c Cursor[T]) HasNext() bool {
	return c.rows.Next()
}

func (c Cursor[T]) Next() (Projection[T], error) {
	var (
		p       Projection[T]
		rawData json.RawMessage
	)

	if err := c.rows.Scan(&p.ID, &p.MutationID, &rawData, &p.DateCreated, &p.DateUpdated); err != nil {
		return p, fmt.Errorf("mutantdb: failed to scan projection: %w", err)
	}

	if err := json.Unmarshal(rawData, &p.Data); err != nil {
		return p, fmt.Errorf("mutantdb: failed to deserialize entity data: %w", err)
	}

	return p, nil
}

func (c Cursor[T]) Close() {
	c.rows.Close()
}

func (c Cursor[T]) Error() error {
	if err := c.rows.Err(); err != nil {
		return fmt.Errorf("mutantdb: %w", err)
	}
	return nil
}
