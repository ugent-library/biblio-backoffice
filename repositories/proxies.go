package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (s *Repo) isProxy(proxyID string, personIDs []string) bool {
	q := `
		select exists(select 1 from proxies where proxy_person_id = $1 and person_id = any($2));
	`
	var exists bool
	if err := s.conn.QueryRow(context.TODO(), q, &exists); err != nil {
		// TODO log error
		return false
	}

	return exists
}

func (r *Repo) EachProxy(ctx context.Context, fn func(string, string) bool) error {
	q := `
		select proxy_person_id, person_id from proxies
		order by proxy_person_id, person_id;
	`
	rows, err := r.conn.Query(ctx, q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var proxyID string
	var personID string
	for rows.Next() {
		if err := rows.Scan(&proxyID, &personID); err != nil {
			return err
		}
		if !fn(proxyID, personID) {
			break
		}
	}

	return rows.Err()
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
