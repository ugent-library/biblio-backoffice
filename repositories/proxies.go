package repositories

import (
	"context"
)

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
