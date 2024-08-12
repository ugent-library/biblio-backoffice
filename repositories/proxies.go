package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
)

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
