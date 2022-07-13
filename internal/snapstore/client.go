package snapstore

// IDEAS:
// - add a method to get all snapshots for a given id
// - add a method to view a store at a given date
// - add date_created column to the table
// - store canonical json and make snapshot id a hash of that
// - improve the api:

// tx := client.BeginTx(ctx)
// tx.Store("publications").Add(ctx, id, &models.Publication{})
// tx.Commit()
// or builder style
// client.Store("publications").WithTx(tx).Get(id).Scan(&models.Publication{})
// client.Store("publications").WithTx(tx).Get(id).Snapshot()

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

func (s *Snapshot) Scan(data any) error {
	return json.Unmarshal(s.Data, data)
}

type Conflict struct {
}

func (e *Conflict) Error() string {
	return "version conflict"
}

type DB interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

type Client struct {
	db         *pgxpool.Pool
	stores     map[string]*Store
	generateID func() (string, error)
}

type Store struct {
	db           DB
	name         string
	table        string
	listeners    []func(*Snapshot)
	listnenersMu sync.RWMutex
	generateID   func() (string, error)
}

type Transaction struct {
	db DB
}

type Options struct {
	Context     context.Context
	Transaction *Transaction
}

func WithIDGenerator(fn func() (string, error)) func(*Client) {
	return func(c *Client) {
		c.generateID = fn
	}
}

func New(db *pgxpool.Pool, stores []string, opts ...func(*Client)) *Client {
	c := &Client{db: db, stores: make(map[string]*Store)}
	for _, opt := range opts {
		opt(c)
	}
	for _, name := range stores {
		c.stores[name] = c.newStore(name)
	}

	if c.generateID == nil {
		c.generateID = generateUUID
	}

	return c
}

func (c *Client) newStore(name string) *Store {
	return &Store{
		db:         c.db,
		generateID: c.generateID,
		name:       name,
		table:      pgx.Identifier.Sanitize([]string{name}),
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

func (s *Store) AddAfter(snapshotID, id string, data any, o Options) (string, error) {
	if snapshotID == "" {
		return "", errors.New("snapshot id is empty")
	}
	if id == "" {
		return "", errors.New("id is empty")
	}

	d, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	ctx, db := s.ctxAndDb(o)

	tx, err := db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	snap := Snapshot{}
	snapSql := `select id, date_from, date_until from ` + s.table + `
	where snapshot_id = $1
	limit 1`
	err = tx.QueryRow(ctx, snapSql, snapshotID).Scan(&snap.ID, &snap.DateFrom, &snap.DateUntil)

	if err == pgx.ErrNoRows {
		return "", fmt.Errorf("unknown snapshot %s", snapshotID)
	} else if err != nil {
		return "", err
	}

	if snap.ID != id {
		return "", fmt.Errorf("id mismatch: snapshot %s belongs to %s, not %s", snapshotID, snap.ID, id)
	}

	if snap.DateUntil != nil {
		// TODO: add info needed to solve the conflict
		return "", &Conflict{}
	}

	now := time.Now()

	sqlUpdate := `update ` + s.table + ` set date_until = $1
	where id = $2 and date_until is null
	returning snapshot_id,id,data,date_until,date_from`

	updatedRows, err := tx.Query(ctx, sqlUpdate, now, id)
	if err != nil {
		return "", err
	}
	cursorUpdatedRows := &Cursor{updatedRows}
	oldSnapshots := []*Snapshot{}
	for cursorUpdatedRows.HasNext() {
		snap, e := cursorUpdatedRows.Next()
		if e != nil {
			return "", e
		}
		oldSnapshots = append(oldSnapshots, snap)
	}

	newSnapshotID, err := s.generateID()
	if err != nil {
		return "", err
	}

	sqlInsert := `insert into ` + s.table + `(snapshot_id, id, data, date_from) values ($1, $2, $3, $4)`

	if _, err = tx.Exec(ctx, sqlInsert, newSnapshotID, id, d, now); err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	for _, snap := range oldSnapshots {
		s.notify(snap)
	}
	s.notify(&Snapshot{
		SnapshotID: newSnapshotID,
		ID:         id,
		Data:       json.RawMessage(d),
		DateFrom:   &now,
	})

	return newSnapshotID, nil
}

func (s *Store) Add(id string, data any, o Options) error {
	if id == "" {
		return errors.New("id is empty")
	}

	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ctx, db := s.ctxAndDb(o)

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	now := time.Now()

	sqlUpdate := `update ` + s.table + ` set date_until = $1
	where id = $2 and date_until is null
	returning snapshot_id,id,data,date_until,date_from`

	updatedRows, err := tx.Query(ctx, sqlUpdate, now, id)
	if err != nil {
		return err
	}
	cursorUpdatedRows := &Cursor{updatedRows}
	oldSnapshots := []*Snapshot{}
	for cursorUpdatedRows.HasNext() {
		snap, e := cursorUpdatedRows.Next()
		if e != nil {
			return e
		}
		oldSnapshots = append(oldSnapshots, snap)
	}

	sqlInsert := `insert into ` + s.table + `(snapshot_id, id, data, date_from) values ($1, $2, $3, $4)`

	newSnapshotID, err := s.generateID()
	if err != nil {
		return err
	}

	if _, err = tx.Exec(ctx, sqlInsert, newSnapshotID, id, d, now); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	for _, snap := range oldSnapshots {
		s.notify(snap)
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
	ctx, db := s.ctxAndDb(o)

	sql := `
	select snapshot_id, data, date_from from ` + s.table + `
	where date_until is null and id = $1
	limit 1`

	snap := Snapshot{}

	if err := db.QueryRow(ctx, sql, id).Scan(&snap.SnapshotID, &snap.Data, &snap.DateFrom); err != nil {
		return nil, err
	}

	return &snap, nil
}

func (s *Store) GetByID(ids []string, o Options) (*Cursor, error) {
	ctx, db := s.ctxAndDb(o)

	pgIds := &pgtype.TextArray{}
	pgIds.Set(ids)
	sql := "select snapshot_id, id, data, date_from, date_until from " + s.table +
		" where date_until is null and id = any($1)"

	rows, err := db.Query(ctx, sql, pgIds)
	if err != nil {
		return nil, err
	}
	return &Cursor{rows}, nil
}

func (s *Store) GetAll(o Options) (*Cursor, error) {
	ctx, db := s.ctxAndDb(o)

	sql := "select snapshot_id, id, data, date_from, date_until from " + s.table +
		" where date_until is null"

	rows, err := db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	return &Cursor{rows}, nil
}

func (s *Store) Purge(id string, o Options) error {
	ctx, db := s.ctxAndDb(o)

	sql := "delete from " + s.table + " where id = $1"

	_, err := db.Exec(ctx, sql, id)

	return err
}

func (s *Store) PurgeAll(o Options) error {
	ctx, db := s.ctxAndDb(o)

	sql := "truncate " + s.table

	_, err := db.Exec(ctx, sql)

	return err
}

func (s *Store) ctxAndDb(o Options) (context.Context, DB) {
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
	return ctx, db
}

func (s *Store) GetAllSnapshots(o Options) (*Cursor, error) {
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

	sql := "select snapshot_id, id, data, date_from, date_until from " + s.table

	rows, err := db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	return &Cursor{rows}, nil
}

type Cursor struct {
	rows pgx.Rows
}

func (c *Cursor) HasNext() bool {
	return c.rows.Next()
}

func (c *Cursor) Next() (*Snapshot, error) {
	s := Snapshot{}
	err := c.rows.Scan(&s.SnapshotID, &s.ID, &s.Data, &s.DateFrom, &s.DateUntil)
	return &s, err
}

func (c *Cursor) Close() {
	c.rows.Close()
}

func (c *Cursor) Err() error {
	return c.rows.Err()
}

func generateUUID() (string, error) {
	return uuid.NewString(), nil
}
