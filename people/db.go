package people

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ugent-library/crypt"
)

type Conn interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Begin(context.Context) (pgx.Tx, error)
}

const getOrganizationIDQuery = `
SELECT id
FROM organizations o
JOIN organization_identifiers oi ON o.id = oi.organization_id AND oi.kind = $1 and oi.value = $2;
`

const getOrganizationByIdentifierQuery = `
SELECT id,
       parent_id,
	   json_agg(DISTINCT jsonb_build_object('kind', ids.kind, 'value', ids.value)) FILTER (WHERE ids.organization_id IS NOT NULL) AS identifiers,
	   names,
	   ceased,
	   position,
	   created_at,
	   updated_at
FROM organizations o
JOIN organization_identifiers oi ON o.id = oi.organization_id AND oi.kind = $1 and oi.value = $2
LEFT JOIN organization_identifiers ids ON o.id = ids.organization_id
GROUP BY o.id;
`

const getAllOrganizationsQuery = `
SELECT id,
       parent_id,
	   json_agg(DISTINCT jsonb_build_object('kind', ids.kind, 'value', ids.value)) FILTER (WHERE ids.organization_id IS NOT NULL) AS identifiers,
	   names,
	   ceased,
	   created_at,
	   updated_at
FROM organizations o
LEFT JOIN organization_identifiers ids ON o.id = ids.organization_id
GROUP BY o.id;
`

const getParentOrganizations = `
WITH RECURSIVE orgs AS (
	SELECT id,
	       parent_id,
	       names,
	       ceased,
		   created_at,
		   updated_at,
		   0 AS level
	FROM organizations
	WHERE id = $1

	UNION

	SELECT o.id,
           o.parent_id,
	       o.names,
	       o.ceased,
		   o.created_at,
		   o.updated_at,	
		   orgs.level + 1
	FROM organizations o
	INNER JOIN orgs
    ON o.id = orgs.parent_id 		
)
SELECT json_agg(DISTINCT jsonb_build_object('kind', ids.kind, 'value', ids.value)) FILTER (WHERE ids.organization_id IS NOT NULL) AS identifiers,
       o.names,
       o.ceased,
	   o.created_at,
	   o.updated_at
FROM orgs o
LEFT JOIN organization_identifiers ids ON o.id = ids.organization_id
GROUP BY o.id,
         o.names,
         o.ceased,
		 o.created_at,
		 o.updated_at,
		 o.level
ORDER BY o.level;
`

const insertOrganizationQuery = `
INSERT INTO organizations (
	parent_id, 
	names,
	ceased,
	position,
	created_at,
	updated_at
)
VALUES (
	$1,
	$2,
	$3,
	(SELECT COUNT(*) FROM organizations WHERE parent_id = $1),
	COALESCE($4,CURRENT_TIMESTAMP),
	COALESCE($5,CURRENT_TIMESTAMP)
)
RETURNING id;
`

const insertOrganizationIdentifierQuery = `
INSERT INTO organization_identifiers (
	organization_id,
	kind,
	value
) VALUES ($1, $2, $3);
`

const insertAffiliationQuery = `
INSERT INTO affiliations (
	person_id,
	organization_id
) VALUES ($1, $2);
`

// SELECT p.id,
// 	   json_agg(json_build_object('kind', ids.kind, 'value', ids.value)) AS identifiers,
//   	   json_agg(DISTINCT jsonb_build_object('organization_id', a.organization_id, 'organization', a.organization)) FILTER (WHERE a.person_id IS NOT NULL) as affiliations,
//        p.name,
//        p.preferred_name,
// 	   p.given_name,
// 	   p.preferred_given_name,
// 	   p.family_name,
// 	   p.preferred_family_name,
// 	   p.honorific_prefix,
// 	   p.email,
// 	   p.active,
// 	   p.role,
// 	   p.username,
// 	   p.attributes,
// 	   p.tokens,
// 	   p.created_at,
// 	   p.updated_at
// FROM people p
// JOIN person_identifiers pi ON p.id = pi.person_id AND pi.kind = 'ugentID' and pi.value = '801001590251'
// LEFT JOIN person_identifiers ids ON p.id = ids.person_id
// LEFT JOIN (
//   SELECT a.person_id, o.id as organization_id, json_build_object('identifiers', json_agg(jsonb_build_object('kind', oi.kind, 'value', oi.value)), 'names', o.names, 'ceased', o.ceased, 'createdAt', o.created_at, 'updatedAt', o.updated_at) as organization
//   FROM affiliations a
//   INNER JOIN organizations o ON a.organization_id = o.id
//   LEFT JOIN organization_identifiers oi ON o.id = oi.organization_id
//   GROUP BY a.person_id, a.organization_id, o.id, o.names, o.ceased, o.created_at, o.updated_at
// ) a on p.id = a.person_id
// WHERE p.replaced_by_id IS NULL
// GROUP BY p.id;

const getPersonByIdentifierQuery = `
SELECT p.id,
	   json_agg(DISTINCT jsonb_build_object('kind', ids.kind, 'value', ids.value)) FILTER (WHERE ids.person_id IS NOT NULL) AS identifiers,
	   array_agg(DISTINCT a.organization_id) FILTER (WHERE a.person_id IS NOT NULL) AS affiliations,
       p.name,
       p.preferred_name,
	   p.given_name,
	   p.preferred_given_name,
	   p.family_name,
	   p.preferred_family_name,
	   p.honorific_prefix,
	   p.email,
	   p.active,
	   p.role,
	   p.username,
	   p.attributes,
	   p.tokens,
	   p.created_at,
	   p.updated_at
FROM people p
JOIN person_identifiers pi ON p.id = pi.person_id AND pi.kind = $1 and pi.value = $2
LEFT JOIN person_identifiers ids ON p.id = ids.person_id
LEFT JOIN affiliations a ON p.id = a.person_id
WHERE p.replaced_by_id IS NULL
GROUP BY p.id;
`

const getAllPeopleQuery = `
SELECT id,
	   json_agg(DISTINCT jsonb_build_object('kind', ids.kind, 'value', ids.value)) FILTER (WHERE ids.person_id IS NOT NULL) AS identifiers,
  	   array_agg(DISTINCT a.organization_id) FILTER (WHERE a.person_id IS NOT NULL) AS affiliations,
	   name,
       preferred_name,
	   given_name,
	   preferred_given_name,
	   family_name,
	   preferred_family_name,
	   honorific_prefix,
	   email,
	   active,
	   username,
	   attributes,
	   created_at,
	   updated_at
FROM people p
LEFT JOIN person_identifiers ids ON p.id = ids.person_id
LEFT JOIN affiliations a ON p.id = a.person_id
WHERE p.replaced_by_id IS NULL
GROUP BY p.id;
`

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
	active,
	role,
	username,
	attributes,
	tokens,
	created_at,
	updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
	COALESCE($14,CURRENT_TIMESTAMP),
	COALESCE($15,CURRENT_TIMESTAMP))
RETURNING id;
`

const createPersonQuery = `
INSERT INTO people (
	name,
	preferred_name,
	given_name,
	preferred_given_name,
	family_name,
	preferred_family_name,
	honorific_prefix,
	email,
	active,
	username,
	attributes
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id;
`

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
	active = $10,
	username = $11,
	attributes = $12,
	updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
`

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

const setPersonReplacedByQuery = `
UPDATE people
SET replaced_by_id = $2, active = FALSE
WHERE id = $1;
`

const deactivatePeopleQuery = `
UPDATE people SET active = FALSE
WHERE updated_at < $1;
`

type organizationRow struct {
	ID          int64
	ParentID    pgtype.Int8
	Identifiers []Identifier
	Names       []Text
	Ceased      bool
	Position    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (row *organizationRow) toOrganization() *Organization {
	return &Organization{
		Identifiers: row.Identifiers,
		Names:       row.Names,
		Ceased:      row.Ceased,
		Position:    row.Position,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func (row *organizationRow) toParentOrganization() ParentOrganization {
	return ParentOrganization{
		Identifiers: row.Identifiers,
		Names:       row.Names,
		Ceased:      row.Ceased,
	}
}

type personRow struct {
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
	Active              bool
	Role                pgtype.Text
	Username            pgtype.Text
	Attributes          []Attribute
	Tokens              []Token
	Affiliations        []Affiliation
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (row *personRow) toPerson(tokenSecret []byte) (*Person, error) {
	p := &Person{
		Identifiers:         row.Identifiers,
		Name:                row.Name,
		PreferredName:       row.PreferredName.String,
		GivenName:           row.GivenName.String,
		PreferredGivenName:  row.PreferredGivenName.String,
		FamilyName:          row.FamilyName.String,
		PreferredFamilyName: row.PreferredFamilyName.String,
		HonorificPrefix:     row.HonorificPrefix.String,
		Email:               row.Email.String,
		Active:              row.Active,
		Role:                row.Role.String,
		Username:            row.Username.String,
		Attributes:          row.Attributes,
		Tokens:              row.Tokens,
		Affiliations:        row.Affiliations,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
	}
	for i, t := range p.Tokens {
		v, err := crypt.Decrypt(tokenSecret, t.Value)
		if err != nil {
			return nil, err
		}
		p.Tokens[i] = Token{Kind: t.Kind, Value: v}
	}
	return p, nil
}

func getOrganizationByIdentifier(ctx context.Context, conn Conn, kind, value string) (*organizationRow, error) {
	var r organizationRow

	err := conn.QueryRow(ctx, getOrganizationByIdentifierQuery, kind, value).Scan(
		&r.ID,
		&r.ParentID,
		&r.Identifiers,
		&r.Names,
		&r.Ceased,
		&r.Position,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func insertOrganization(ctx context.Context, conn Conn, o ImportOrganizationParams) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var id int64
	var parentID pgtype.Int8

	if ident := o.ParentIdentifier; ident != nil {
		err := tx.QueryRow(ctx, getOrganizationIDQuery, ident.Kind, ident.Value).Scan(&parentID.Int64)
		if err == pgx.ErrNoRows {
			return fmt.Errorf("organization %s not found", ident.String())
		}
		if err != nil {
			return err
		}
		parentID.Valid = true
	}

	err = tx.QueryRow(ctx, insertOrganizationQuery,
		parentID,
		o.Names,
		o.Ceased,
		o.CreatedAt,
		o.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return err
	}

	for _, ident := range o.Identifiers {
		if _, err := tx.Exec(ctx, insertOrganizationIdentifierQuery, id, ident.Kind, ident.Value); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func insertPerson(ctx context.Context, conn Conn, p ImportPersonParams) error {
	var id int64

	err := conn.QueryRow(ctx, insertPersonQuery,
		p.Name,
		pgtype.Text{Valid: p.PreferredName != "", String: p.PreferredName},
		pgtype.Text{Valid: p.GivenName != "", String: p.GivenName},
		pgtype.Text{Valid: p.PreferredGivenName != "", String: p.PreferredGivenName},
		pgtype.Text{Valid: p.FamilyName != "", String: p.FamilyName},
		pgtype.Text{Valid: p.PreferredFamilyName != "", String: p.PreferredFamilyName},
		pgtype.Text{Valid: p.HonorificPrefix != "", String: p.HonorificPrefix},
		pgtype.Text{Valid: p.Email != "", String: p.Email},
		p.Active,
		pgtype.Text{Valid: p.Role != "", String: p.Role},
		pgtype.Text{Valid: p.Username != "", String: p.Username},
		p.Attributes,
		p.Tokens,
		p.CreatedAt,
		p.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return err
	}

	for _, ident := range p.Identifiers {
		if _, err := conn.Exec(ctx, insertPersonIdentifierQuery, id, ident.Kind, ident.Value); err != nil {
			return err
		}
	}

	for _, a := range p.Affiliations {
		var organizationID int64
		err := conn.QueryRow(ctx, getOrganizationIDQuery, a.OrganizationIdentifier.Kind, a.OrganizationIdentifier.Value).Scan(&organizationID)
		if err == pgx.ErrNoRows {
			return fmt.Errorf("organization %s not found", a.OrganizationIdentifier.String())
		}
		if err != nil {
			return err
		}
		if _, err := conn.Exec(ctx, insertAffiliationQuery, id, organizationID); err != nil {
			return err
		}
	}

	return nil
}

func getPersonByIdentifier(ctx context.Context, conn Conn, kind, value string) (*personRow, error) {
	var r personRow

	var orgIDs []int64

	q := getPersonByIdentifierQuery

	err := conn.QueryRow(ctx, q, kind, value).Scan(
		&r.ID,
		&r.Identifiers,
		&orgIDs,
		&r.Name,
		&r.PreferredName,
		&r.GivenName,
		&r.PreferredGivenName,
		&r.FamilyName,
		&r.PreferredFamilyName,
		&r.HonorificPrefix,
		&r.Email,
		&r.Active,
		&r.Role,
		&r.Username,
		&r.Attributes,
		&r.Tokens,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	for _, orgID := range orgIDs {
		var org *Organization
		rows, err := conn.Query(ctx, getParentOrganizations, orgID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var o organizationRow
			if err := rows.Scan(
				&o.Identifiers,
				&o.Names,
				&o.Ceased,
				&o.CreatedAt,
				&o.UpdatedAt,
			); err != nil {
				return nil, err
			}
			if org == nil {
				org = o.toOrganization()
			} else {
				org.Parents = append(org.Parents, o.toParentOrganization())
			}
		}

		r.Affiliations = append(r.Affiliations, Affiliation{Organization: org})
	}

	return &r, nil
}

func createPerson(ctx context.Context, conn Conn, params AddPersonParams) (int64, error) {
	var id int64
	err := conn.QueryRow(ctx, createPersonQuery,
		params.Name,
		pgtype.Text{Valid: params.PreferredName != "", String: params.PreferredName},
		pgtype.Text{Valid: params.GivenName != "", String: params.GivenName},
		pgtype.Text{Valid: params.PreferredGivenName != "", String: params.PreferredGivenName},
		pgtype.Text{Valid: params.FamilyName != "", String: params.FamilyName},
		pgtype.Text{Valid: params.PreferredFamilyName != "", String: params.PreferredFamilyName},
		pgtype.Text{Valid: params.HonorificPrefix != "", String: params.HonorificPrefix},
		pgtype.Text{Valid: params.Email != "", String: params.Email},
		params.Active,
		pgtype.Text{Valid: params.Username != "", String: params.Username},
		params.Attributes,
	).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, replacePersonIdentifiers(ctx, conn, id, params.Identifiers)
}

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
		params.Active,
		pgtype.Text{Valid: params.Username != "", String: params.Username},
		params.Attributes,
	)
	if err != nil {
		return err
	}
	return replacePersonIdentifiers(ctx, conn, id, params.Identifiers)
}

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

func setPersonReplacedBy(ctx context.Context, conn Conn, id, replacedByID int64) error {
	_, err := conn.Exec(ctx, setPersonReplacedByQuery, id, replacedByID)
	return err
}
