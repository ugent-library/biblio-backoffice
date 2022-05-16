package snapstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Snapshot struct {
	SnapshotID string
	ID         string
	Data       json.RawMessage
	DateFrom   *time.Time
	DateUntil  *time.Time
}

func (s *Snapshot) Scan(data interface{}) error {
	return json.Unmarshal(s.Data, data)
}

type Conflict struct {
}

func (e *Conflict) Error() string {
	return "version conflict"
}

type DB interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
}

type Client struct {
	db     *pgxpool.Pool
	stores map[string]*Store
}

type Store struct {
	db           DB
	name         string
	table        string
	listeners    []func(*Snapshot)
	listnenersMu sync.RWMutex
}

type Transaction struct {
	db DB
}

type Options struct {
	Context     context.Context
	Transaction *Transaction
}

func New(db *pgxpool.Pool, stores []string) *Client {
	c := &Client{db: db, stores: make(map[string]*Store)}
	for _, name := range stores {
		c.stores[name] = c.newStore(name)
	}
	return c
}

func (c *Client) newStore(name string) *Store {
	return &Store{
		db:    c.db,
		name:  name,
		table: pgx.Identifier.Sanitize([]string{name}),
	}
}

func (c *Client) Store(name string) *Store {
	return c.stores[name]
}

func (c *Client) Transaction(ctx context.Context, fn func(Options) error) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err = fn(Options{Context: ctx, Transaction: &Transaction{tx}}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Store) Name() string {
	return s.name
}

func (s *Store) Listen(fn func(*Snapshot)) {
	s.listnenersMu.Lock()
	defer s.listnenersMu.Unlock()
	s.listeners = append(s.listeners, fn)
}

func (s *Store) notify(snap *Snapshot) {
	s.listnenersMu.RLock()
	defer s.listnenersMu.RUnlock()
	// TODO do this non-blocking
	for _, fn := range s.listeners {
		fn(snap)
	}
}

func (s *Store) AddAfter(snapshotID, id string, data interface{}, o Options) error {
	if snapshotID == "" {
		return errors.New("snapshot id is empty")
	}
	if id == "" {
		return errors.New("id is empty")
	}

	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var (
		ctx context.Context
		db  DB
	)
	if o.Context == nil {
		ctx = context.Background()
	} else {
		ctx = o.Context
	}
	if o.Transaction == nil {
		db = s.db
	} else {
		db = o.Transaction.db
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	snap := Snapshot{}
	snapSql := `select id, date_from, date_until from ` + s.table + `
	where snapshot_id = $1
	limit 1`
	err = db.QueryRow(ctx, snapSql, snapshotID).Scan(&snap.ID, &snap.DateFrom, &snap.DateUntil)

	if err == pgx.ErrNoRows {
		return fmt.Errorf("unknown snapshot %s", snapshotID)
	} else if err != nil {
		return err
	}

	if snap.ID != id {
		return fmt.Errorf("id mismatch: snapshot %s belongs to %s, not %s", snapshotID, snap.ID, id)
	}

	if snap.DateUntil != nil {
		// TODO: add info needed to solve the conflict
		return &Conflict{}
	}

	now := time.Now()

	sqlUpdate := "update " + s.table + " set date_until = $1 where id = $2 and date_until is null"

	if _, err = tx.Exec(ctx, sqlUpdate, now, id); err != nil {
		return err
	}

	newSnapshotID := uuid.NewString()

	sqlInsert := `insert into ` + s.table + `(snapshot_id, id, data, date_from) values ($1, $2, $3, $4)`

	if _, err = tx.Exec(ctx, sqlInsert, newSnapshotID, id, d, now); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	s.notify(&Snapshot{
		SnapshotID: newSnapshotID,
		ID:         id,
		Data:       json.RawMessage(d),
		DateFrom:   &now,
	})

	return nil
}

func (s *Store) Add(id string, data interface{}, o Options) error {
	if id == "" {
		return errors.New("id is empty")
	}

	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var (
		ctx context.Context
		db  DB
	)
	if o.Context == nil {
		ctx = context.Background()
	} else {
		ctx = o.Context
	}
	if o.Transaction == nil {
		db = s.db
	} else {
		db = o.Transaction.db
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	now := time.Now()

	sqlUpdate := "update " + s.table + " set date_until = $1 where id = $2 and date_until is null"

	if _, err = tx.Exec(ctx, sqlUpdate, now, id); err != nil {
		return err
	}

	sqlInsert := `insert into ` + s.table + `(snapshot_id, id, data, date_from) values ($1, $2, $3, $4)`

	newSnapshotID := uuid.NewString()

	if _, err = tx.Exec(ctx, sqlInsert, newSnapshotID, id, d, now); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	s.notify(&Snapshot{
		SnapshotID: newSnapshotID,
		ID:         id,
		Data:       json.RawMessage(d),
		DateFrom:   &now,
	})

	return nil

}

func (s *Store) GetCurrentSnapshot(id string, o Options) (*Snapshot, error) {
	var (
		ctx context.Context
		db  DB
	)
	if o.Context == nil {
		ctx = context.Background()
	} else {
		ctx = o.Context
	}
	if o.Transaction == nil {
		db = s.db
	} else {
		db = o.Transaction.db
	}

	sql := `
	select snapshot_id, data from ` + s.table + `
	where date_until is null and id = $1
	limit 1`

	snap := Snapshot{}

	if err := db.QueryRow(ctx, sql, id).Scan(&snap.SnapshotID, &snap.Data); err != nil {
		return nil, err
	}

	return &snap, nil
}

func (s *Store) GetByID(ids []string, o Options) *Cursor {
	var (
		ctx context.Context
		db  DB
	)
	if o.Context == nil {
		ctx = context.Background()
	} else {
		ctx = o.Context
	}
	if o.Transaction == nil {
		db = s.db
	} else {
		db = o.Transaction.db
	}

	pgIds := &pgtype.TextArray{}
	pgIds.Set(ids)
	sql := "select data from " + s.table + " where date_until is null and id = any($1)"

	c := &Cursor{}
	c.rows, c.err = db.Query(ctx, sql, pgIds)
	return c
}

func (s *Store) GetAll(o Options) *Cursor {
	var (
		ctx context.Context
		db  DB
	)
	if o.Context == nil {
		ctx = context.Background()
	} else {
		ctx = o.Context
	}
	if o.Transaction == nil {
		db = s.db
	} else {
		db = o.Transaction.db
	}

	sql := "select data from " + s.table + " where date_until is null"

	c := &Cursor{}
	c.rows, c.err = db.Query(ctx, sql)
	return c
}

type Cursor struct {
	err  error
	rows pgx.Rows
}

func (c *Cursor) Next() bool {
	return c.err == nil && c.rows.Next()
}

func (c *Cursor) Scan(data interface{}) error {
	if c.err != nil {
		return c.err
	}
	var d json.RawMessage
	if c.err = c.rows.Scan(&d); c.err == nil {
		c.err = json.Unmarshal(d, data)
	}
	return c.err
}

func (c *Cursor) Close() {
	c.rows.Close()
}

func (c *Cursor) Err() error {
	if c.err != nil {
		return c.err
	}
	return c.rows.Err()
}
