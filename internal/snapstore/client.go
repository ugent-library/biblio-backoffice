package snapstore

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Strategy int

const (
	StrategyMine Strategy = iota
	StrategyAbort
)

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
		db:             c.db,
		name:           name,
		versionsTable:  pgx.Identifier.Sanitize([]string{name + "_versions"}),
		snapshotsTable: pgx.Identifier.Sanitize([]string{name + "_snapshots"}),
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

type Store struct {
	db             DB
	name           string
	versionsTable  string
	snapshotsTable string
}

func (s *Store) Name() string {
	return s.name
}

func (s *Store) AddVersion(affinityID, id string, data interface{}, o Options) error {
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

	sql := `insert into ` + s.versionsTable + `(affinity_id, id, data)
	        values ($1, $2, $3)`

	if _, err = db.Exec(ctx, sql, affinityID, id, d); err != nil {
		return err
	}

	return nil
}

func (s *Store) AddSnapshot(affinityID, id string, strategy Strategy, o Options) error {
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

	if strategy == StrategyAbort {
		// TODO check if another affinity has already added a version after the last snapshot
	}

	sql := `
	with version as (
		select version_id, id, data 
		from ` + s.versionsTable + `
		where affinity_id = $1 and id = $2
		order by date_created desc
		limit 1
	), snapshot as (
	   insert into ` + s.snapshotsTable + `(id, data)
	   select id, data
	   from version
	   returning snapshot_id, date_from
	), old_snapshots as (
		update ` + s.snapshotsTable + `
		set date_until=snapshot.date_from
		from snapshot
		where ` + s.snapshotsTable + `.id = $2 and ` + s.snapshotsTable + `.snapshot_id != snapshot.snapshot_id
	)
	update ` + s.versionsTable + `
	set snapshot_id=snapshot.snapshot_id
	from version, snapshot
	where ` + s.versionsTable + `.version_id = version.version_id`

	if _, err := db.Exec(ctx, sql, affinityID, id); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetVersion(affinityID, id string, data interface{}, o Options) error {
	var d json.RawMessage

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

	sql := `select data from ` + s.versionsTable + `
	where affinity_id=$1 and id=$2
	order by date_created desc
	limit 1`

	if err := db.QueryRow(ctx, sql, affinityID, id).Scan(&d); err != nil {
		return err
	}

	return json.Unmarshal(d, data)
}

func (s *Store) Get(id string, data interface{}, o Options) error {
	var d json.RawMessage

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
	select data from ` + s.snapshotsTable + `
	where date_until = 'infinity'::timestamptz and id = $1
	limit 1`

	if err := db.QueryRow(ctx, sql, id).Scan(&d); err != nil {
		return err
	}

	return json.Unmarshal(d, data)
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
	sql := `
	select data from ` + s.snapshotsTable + `
	where date_until = 'infinity'::timestamptz and id = any($1)`

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

	sql := `
	select data from ` + s.snapshotsTable + `
	where date_until = 'infinity'::timestamptz`

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

type Snapshot struct {
	ID        string
	Data      json.RawMessage
	DateFrom  time.Time
	DateUntil time.Time
}
