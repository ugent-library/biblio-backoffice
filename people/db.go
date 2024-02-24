package people

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type Conn interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Begin(context.Context) (pgx.Tx, error)
}

const getPersonByIdentifierQuery = `
WITH ids AS (
	SELECT ids1.*
	FROM person_identifiers ids1
	LEFT JOIN person_identifiers ids2 ON ids1.person_id = ids2.person_id
	WHERE ids2.kind = $1 AND ids2.value = $2	
)
SELECT id,
	   json_agg(json_build_object('kind', ids.kind, 'value', ids.value)) AS identifiers,
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
FROM people p, ids WHERE p.id = ids.person_id
GROUP BY p.id;
`

type getPersonRow struct {
	ID                  int64
	Identifiers         []Identifier
	Name                string
	PreferredName       pgtype.Text
	GivenName           pgtype.Text
	PreferredGivenName  pgtype.Text
	FamilyName          pgtype.Text
	PreferredFamilyName pgtype.Text
	HonorificPrefix     pgtype.Text
	Email               pgtype.Text
	Attributes          []Attribute
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (row *getPersonRow) toPerson() *Person {
	return &Person{
		Identifiers:         row.Identifiers,
		Name:                row.Name,
		PreferredName:       row.PreferredName.String,
		GivenName:           row.GivenName.String,
		PreferredGivenName:  row.PreferredGivenName.String,
		FamilyName:          row.FamilyName.String,
		PreferredFamilyName: row.PreferredFamilyName.String,
		HonorificPrefix:     row.HonorificPrefix.String,
		Email:               row.Email.String,
		Attributes:          row.Attributes,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
	}
}

func getPersonByIdentifier(ctx context.Context, conn Conn, kind, value string) (*getPersonRow, error) {
	var row getPersonRow

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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;
`

func createPerson(ctx context.Context, conn Conn, params AddPersonParams) (int64, error) {
	var id int64
	err := conn.QueryRow(ctx, insertPersonQuery,
		params.Name,
		pgtype.Text{Valid: params.PreferredName != "", String: params.PreferredName},
		pgtype.Text{Valid: params.GivenName != "", String: params.GivenName},
		pgtype.Text{Valid: params.PreferredGivenName != "", String: params.PreferredGivenName},
		pgtype.Text{Valid: params.FamilyName != "", String: params.FamilyName},
		pgtype.Text{Valid: params.PreferredFamilyName != "", String: params.PreferredFamilyName},
		pgtype.Text{Valid: params.HonorificPrefix != "", String: params.HonorificPrefix},
		pgtype.Text{Valid: params.Email != "", String: params.Email},
		params.Attributes,
	).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, replacePersonIdentifiers(ctx, conn, id, params.Identifiers)
}

const updatePersonQuery = `
UPDATE people SET 
	name = $2,
	preferred_name = $3,
	given_name = $4,
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
		params.Name,
		pgtype.Text{Valid: params.PreferredName != "", String: params.PreferredName},
		pgtype.Text{Valid: params.GivenName != "", String: params.GivenName},
		pgtype.Text{Valid: params.PreferredGivenName != "", String: params.PreferredGivenName},
		pgtype.Text{Valid: params.FamilyName != "", String: params.FamilyName},
		pgtype.Text{Valid: params.PreferredFamilyName != "", String: params.PreferredFamilyName},
		pgtype.Text{Valid: params.HonorificPrefix != "", String: params.HonorificPrefix},
		pgtype.Text{Valid: params.Email != "", String: params.Email},
		params.Attributes,
	)
	if err != nil {
		return err
	}
	return replacePersonIdentifiers(ctx, conn, id, params.Identifiers)
}

const deletePersonIdentifiersQuery = `
DELETE FROM person_identifiers
WHERE person_id = $1;
`

const insertPersonIdentifierQuery = `
INSERT INTO person_identifiers (
	person_id,
	kind,
	value
) VALUES ($1, $2, $3);
`

func replacePersonIdentifiers(ctx context.Context, conn Conn, pID int64, ids []Identifier) error {
	if _, err := conn.Exec(ctx, deletePersonIdentifiersQuery, pID); err != nil {
		return err
	}
	for _, id := range ids {
		if _, err := conn.Exec(ctx, insertPersonIdentifierQuery, pID, id.Kind, id.Value); err != nil {
			return err
		}
	}
	return nil
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
