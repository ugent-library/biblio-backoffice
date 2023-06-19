package snapstore

// IDEAS:
// - snapshots table with type column
// - versioning strategies: update in place unless user changed, status changed, abort if no changes, …
// - compaction method
// - introduce internal and external relations
// - add a method to get all snapshots for a given id
// - add a method to view a store at a given date
// - add date_created column to the table
// - draft versions with affinity_id: https://github.com/ugent-library/biblio-backoffice/commit/419b5ccd5de83b1010a2b629d72a526d2e33ae67
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
	db         DB
	name       string
	table      string
	generateID func() (string, error)
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

func (c *Client) Tx(ctx context.Context, fn func(Options) error) error {
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

	sqlSnap := `SELECT id, date_from, date_until FROM ` + s.table + `
	WHERE snapshot_id = $1
	LIMIT 1`

	err = tx.QueryRow(ctx, sqlSnap, snapshotID).Scan(&snap.ID, &snap.DateFrom, &snap.DateUntil)

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

	sqlUpdate := `UPDATE ` + s.table + ` SET date_until = $1
	WHERE id = $2 AND date_until IS NULL`

	if _, err := tx.Exec(ctx, sqlUpdate, now, id); err != nil {
		return "", err
	}

	newSnapshotID, err := s.generateID()
	if err != nil {
		return "", err
	}

	sqlInsert := `INSERT INTO ` + s.table + `(snapshot_id, id, data, date_from) VALUES ($1, $2, $3, $4)`

	if _, err = tx.Exec(ctx, sqlInsert, newSnapshotID, id, d, now); err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return newSnapshotID, nil
}

func (s *Store) Update(snapshotID, id string, data any, o Options) (*Snapshot, error) {
	ctx, db := s.ctxAndDb(o)

	if snapshotID == "" {
		return nil, errors.New("snapshot id is empty")
	}

	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	snap := Snapshot{
		SnapshotID: snapshotID,
		ID:         id,
		Data:       d,
	}

	sql := `UPDATE ` + s.table + ` SET data = $1
	WHERE id = $2 AND snapshot_id = $3
	RETURNING date_from,date_until`

	if err = db.QueryRow(ctx, sql, d, id, snapshotID).Scan(&snap.DateFrom, &snap.DateUntil); err != nil {
		return nil, err
	}

	return &snap, nil
}

func (store *Store) ImportSnapshot(snapshot *Snapshot, options Options) error {
	if snapshot.ID == "" {
		return errors.New("id is empty")
	}
	if snapshot.SnapshotID == "" {
		return errors.New("snapshot_id is empty")
	}
	if snapshot.DateFrom == nil {
		return errors.New("date_from is nil")
	}
	if snapshot.Data == nil {
		return errors.New("data is nil")
	}

	ctx, db := store.ctxAndDb(options)

	sql := `INSERT INTO ` + store.table +
		`(snapshot_id, id, data, date_from, date_until) VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(ctx, sql,
		snapshot.SnapshotID,
		snapshot.ID,
		snapshot.Data,
		snapshot.DateFrom,
		snapshot.DateUntil,
	)

	return err
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

	sqlUpdate := `UPDATE ` + s.table + ` SET date_until = $1
	WHERE id = $2 AND date_until IS null`

	if _, err := tx.Exec(ctx, sqlUpdate, now, id); err != nil {
		return err
	}

	sqlInsert := `INSERT INTO ` + s.table + `(snapshot_id, id, data, date_from) VALUES ($1, $2, $3, $4)`

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

	return nil

}

func (s *Store) GetCurrentSnapshot(id string, o Options) (*Snapshot, error) {
	ctx, db := s.ctxAndDb(o)

	sql := `SELECT snapshot_id, data, date_from FROM ` + s.table + `
	WHERE date_until IS NULL AND id = $1
	LIMIT 1`

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
	sql := "SELECT snapshot_id, id, data, date_from, date_until FROM " + s.table +
		" WHERE date_until IS NULL AND id = any($1) order by array_position($1, id)"

	rows, err := db.Query(ctx, sql, pgIds)
	if err != nil {
		return nil, err
	}
	return &Cursor{rows}, nil
}

func (s *Store) GetAll(o Options) (*Cursor, error) {
	ctx, db := s.ctxAndDb(o)

	sql := "SELECT snapshot_id, id, data, date_from, date_until FROM " + s.table +
		" WHERE date_until IS NULL"

	rows, err := db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	return &Cursor{rows}, nil
}

func (s *Store) Purge(id string, o Options) error {
	ctx, db := s.ctxAndDb(o)

	sql := "DELETE FROM " + s.table + " WHERE id = $1"

	_, err := db.Exec(ctx, sql, id)

	return err
}

func (s *Store) PurgeAll(o Options) error {
	ctx, db := s.ctxAndDb(o)

	sql := "TRUNCATE " + s.table

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

func (s *Store) Select(sql string, values []any, options Options) (*Cursor, error) {
	var (
		ctx context.Context
		db  DB
	)
	if options.Context == nil {
		ctx = context.Background()
	} else {
		ctx = options.Context
	}
	if options.Transaction == nil {
		db = s.db
	} else {
		db = options.Transaction.db
	}

	rows, err := db.Query(ctx, sql, values...)
	if err != nil {
		return nil, err
	}
	return &Cursor{rows}, nil
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

	sql := "SELECT snapshot_id, id, data, date_from, date_until FROM " + s.table

	rows, err := db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	return &Cursor{rows}, nil
}

func (s *Store) CountSql(sql string, values []any, o Options) (int, error) {
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

	countSql := "SELECT COUNT(*) count FROM (" + sql + ") t"
	rows, rowsErr := db.Query(ctx, countSql, values...)
	if rowsErr != nil {
		return 0, rowsErr
	}
	defer rows.Close()

	if !rows.Next() {
		// TODO: it is an error to use rows.Scan without calling Next()
		// but what to return here?
		return 0, nil
	}

	var count int = 0
	scanErr := rows.Scan(&count)
	if scanErr != nil {
		return 0, scanErr
	}
	return count, nil
}

func (s *Store) GetHistory(id string, o Options) (*Cursor, error) {
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

	sql := "SELECT snapshot_id, id, data, date_from, date_until FROM " + s.table +
		" WHERE id = $1 ORDER BY date_from DESC"

	rows, err := db.Query(ctx, sql, id)
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
