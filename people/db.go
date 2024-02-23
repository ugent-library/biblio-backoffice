package people

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Conn interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Begin(context.Context) (pgx.Tx, error)
}

const getPersonByIdentifierQuery = `
SELECT id,
       identifiers,
	   name,
       preferred_name,
	   given_name,
	   preferred_given_name,
	   family_name,
	   preferred_family_name,
	   honorific_prefix,
	   email,
	   attributes,
	   created_at,
	   updated_at
FROM people
WHERE identifiers::jsonb @> jsonb_build_array(jsonb_build_object('kind', $1, 'value', $2)) AND
	  replaced_by_id IS NULL;
`

type personRow struct {
	ID int64
	Person
}

func getPersonByIdentifier(ctx context.Context, conn Conn, kind, value string) (*personRow, error) {
	var row personRow
	err := conn.QueryRow(ctx, getPersonByIdentifierQuery, kind, value).Scan(
		&row.ID,
		&row.Identifiers,
		&row.Name,
		&row.PreferredName,
		&row.GivenName,
		&row.PreferredGivenName,
		&row.FamilyName,
		&row.PreferredFamilyName,
		&row.HonorificPrefix,
		&row.Email,
		&row.Attributes,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

const insertPersonQuery = `
INSERT INTO people (
	identifiers,
	name,
	preferred_name,
	given_name,
	preferred_given_name,
	family_name,
	preferred_family_name,
	honorific_prefix,
	email,
	attributes
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id;
`

func insertPerson(ctx context.Context, conn Conn, params AddPersonParams) (int64, error) {
	var id int64
	err := conn.QueryRow(ctx, insertPersonQuery,
		params.Identifiers,
		params.Name,
		params.PreferredName,
		params.GivenName,
		params.PreferredGivenName,
		params.FamilyName,
		params.PreferredFamilyName,
		params.HonorificPrefix,
		params.Email,
		params.Attributes,
	).Scan(&id)
	return id, err
}

const updatePersonQuery = `
UPDATE people SET 
	identifiers = $2,
	name = $3,
	preferred_name = $4,
	given_name = $5,
	preferred_given_name = $5,
	family_name = $6,
	preferred_family_name = $7,
	honorific_prefix = $8,
	email = $9,
	attributes = $10,
	updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
`

func updatePerson(ctx context.Context, conn Conn, id int64, params AddPersonParams) error {
	_, err := conn.Exec(ctx, updatePersonQuery,
		id,
		params.Identifiers,
		params.Name,
		params.PreferredName,
		params.GivenName,
		params.PreferredGivenName,
		params.FamilyName,
		params.PreferredFamilyName,
		params.HonorificPrefix,
		params.Email,
		params.Attributes,
	)
	return err
}

const setPersonReplacedByQuery = `
UPDATE people
SET replaced_by_id = $2
WHERE id = $1;
`

func setPersonReplacedBy(ctx context.Context, conn Conn, id, replacedByID int64) error {
	_, err := conn.Exec(ctx, setPersonReplacedByQuery, id, replacedByID)
	return err
}
