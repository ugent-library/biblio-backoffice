package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (s *Repo) IsProxyFor(proxyIDs []string, personIDs []string) bool {
	q := `
		select exists(select 1 from proxies where proxy_person_id = any($1) and person_id = any($2));
	`
	var exists bool
	if err := s.conn.QueryRow(context.TODO(), q, proxyIDs, personIDs).Scan(&exists); err != nil {
		// TODO log error
		return false
	}

	return exists
}

func (s *Repo) IsProxy(proxyIDs []string) bool {
	q := `
		select exists(select 1 from proxies where proxy_person_id = any($1));
	`
	var exists bool
	if err := s.conn.QueryRow(context.TODO(), q, proxyIDs).Scan(&exists); err != nil {
		// TODO log error
		return false
	}

	return exists
}

func (s *Repo) HasProxy(personIDs []string) bool {
	q := `
		select exists(select 1 from proxies where person_id = any($1));
	`
	var exists bool
	if err := s.conn.QueryRow(context.TODO(), q, personIDs).Scan(&exists); err != nil {
		// TODO log error
		return false
	}

	return exists
}

func (r *Repo) FindProxies(ctx context.Context, personIDs []string, limit, offset int) (int, [][]string, error) {
	var q string
	var args []any

	if len(personIDs) > 0 {
		q = `
			select count(*) over() as total, proxy_person_id, person_id from proxies
			where proxy_person_id = any($1) or person_id = any($1)
			order by proxy_person_id, person_id
			limit $2
			offset $3;
		`
		args = []any{personIDs, limit, offset}
	} else {
		q = `
			select count(*) over() as total, proxy_person_id, person_id from proxies
			order by proxy_person_id, person_id
			limit $1
			offset $2;
		`
		args = []any{limit, offset}
	}

	rows, err := r.conn.Query(ctx, q, args...)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	var total int
	var pairs [][]string

	for rows.Next() {
		pair := make([]string, 2)
		if err := rows.Scan(&total, &pair[0], &pair[1]); err != nil {
			return 0, nil, err
		}
		pairs = append(pairs, pair)
	}

	if err := rows.Err(); err != nil {
		return 0, nil, err
	}

	return total, pairs, nil
}

func (r *Repo) ProxyPersonIDs(ctx context.Context, proxyIDs []string) ([]string, error) {
	q := `
		select person_id from proxies
		where proxy_person_id = any($1);
	`
	rows, err := r.conn.Query(ctx, q, proxyIDs)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

func (r *Repo) ProxyIDs(ctx context.Context, personIDs []string) ([]string, error) {
	q := `
		select proxy_person_id from proxies
		where person_id = any($1);
	`
	rows, err := r.conn.Query(ctx, q, personIDs)
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

func (r *Repo) RemoveProxyPerson(ctx context.Context, proxyIDs, personIDs []string) error {
	q := `
		delete from proxies
		where proxy_person_id = any($1) and person_id = any($2);
	`
	_, err := r.conn.Exec(ctx, q, proxyIDs, personIDs)
	return err
}
