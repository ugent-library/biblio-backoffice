package eventstore

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

	if err := c.rows.Scan(&p.StreamID, &p.EventID, &rawData, &p.DateCreated, &p.DateUpdated); err != nil {
		return p, fmt.Errorf("eventstore: failed to scan projection: %w", err)
	}

	if err := json.Unmarshal(rawData, &p.Data); err != nil {
		return p, fmt.Errorf("eventstore: failed to deserialize projection data: %w", err)
	}

	return p, nil
}

func (c Cursor[T]) Close() {
	c.rows.Close()
}

func (c Cursor[T]) Error() error {
	if err := c.rows.Err(); err != nil {
		return fmt.Errorf("eventstore: %w", err)
	}
	return nil
}
