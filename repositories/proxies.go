package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (s *Repo) IsProxyFor(proxyID string, personIDs []string) bool {
	q := `
		select exists(select 1 from proxies where proxy_person_id = $1 and person_id = any($2));
	`
	var exists bool
	if err := s.conn.QueryRow(context.TODO(), q, proxyID, personIDs).Scan(&exists); err != nil {
		// TODO log error
		return false
	}

	return exists
}

func (s *Repo) IsProxy(proxyID string) bool {
	q := `
		select exists(select 1 from proxies where proxy_person_id = $1);
	`
	var exists bool
	if err := s.conn.QueryRow(context.TODO(), q, proxyID).Scan(&exists); err != nil {
		// TODO log error
		return false
	}

	return exists
}

func (s *Repo) HasProxy(personID string) bool {
	q := `
		select exists(select 1 from proxies where person_id = $1);
	`
	var exists bool
	if err := s.conn.QueryRow(context.TODO(), q, personID).Scan(&exists); err != nil {
		// TODO log error
		return false
	}

	return exists
}

func (r *Repo) FindProxies(ctx context.Context, personIDs []string) ([][]string, error) {
	var q string
	var args []any

	if len(personIDs) > 0 {
		q = `
			select proxy_person_id, person_id from proxies
			where  proxy_person_id = any($1) or person_id = any($1)
			order by proxy_person_id, person_id;
		`
		args = []any{personIDs}
	} else {
		q = `
			select proxy_person_id, person_id from proxies
			order by proxy_person_id, person_id;
		`
	}

	rows, err := r.conn.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pairs [][]string

	for rows.Next() {
		pair := make([]string, 2)
		if err := rows.Scan(&pair[0], &pair[1]); err != nil {
			return nil, err
		}
		pairs = append(pairs, pair)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pairs, err
}

func (r *Repo) ProxyPersonIDs(ctx context.Context, proxyID string) ([]string, error) {
	q := `
		select person_id from proxies
		where proxy_person_id = $1;
	`
	rows, err := r.conn.Query(ctx, q, proxyID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

func (r *Repo) ProxyIDs(ctx context.Context, personID string) ([]string, error) {
	q := `
		select proxy_person_id from proxies
		where person_id = $1;
	`
	rows, err := r.conn.Query(ctx, q, personID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

func (r *Repo) AddProxyPerson(ctx context.Context, proxyID, personID string) error {
	q := `
		insert into proxies (proxy_person_id, person_id)
		values ($1, $2)
		on conflict (proxy_person_id, person_id) do nothing;
	`
	_, err := r.conn.Exec(ctx, q, proxyID, personID)
	return err
}

func (r *Repo) RemoveProxyPerson(ctx context.Context, proxyID, personID string) error {
	q := `
		delete from proxies
		where proxy_person_id = $1 and person_id = $2;
	`
	_, err := r.conn.Exec(ctx, q, proxyID, personID)
	return err
}
