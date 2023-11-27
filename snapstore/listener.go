package snapstore

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
)

type Listener struct {
	conn *pgx.Conn
}

type Event struct {
	Name       string `json:"name"`
	RecordType string `json:"record_type"`
	RecordID   string `json:"record_id"`
	SnapshotID string `json:"snapshot_id,omitempty"`
}

func NewListener(ctx context.Context, conn *pgx.Conn) (*Listener, error) {
	if _, err := conn.Exec(ctx, `listen "events"`); err != nil {
		return nil, err
	}
	return &Listener{
		conn: conn,
	}, nil
}

func (l *Listener) Listen(ctx context.Context) (Event, error) {
	evt := Event{}

	// if ctx is done, err will be non nil
	not, err := l.conn.WaitForNotification(ctx)
	if err != nil {
		return evt, err
	}

	err = json.Unmarshal([]byte(not.Payload), &evt)

	return evt, err
}
