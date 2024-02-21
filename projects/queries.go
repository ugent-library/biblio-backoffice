package projects

import (
	"context"
	"encoding/json"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type conn interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Begin(context.Context) (pgx.Tx, error)
}

const getProjectsByIdentifierQuery = `
SELECT p.* FROM projects p WHERE p.identifiers::JSONB @> $1
`

type getProjectsByIdentifierParams struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type getProjectsByIdentifierRow struct {
	ID              int64
	Names           []Text
	Descriptions    []Text
	FoundingDate    pgtype.Text
	DissolutionDate pgtype.Text
	Attributes      []Attribute
	Identifiers     []Identifier
	CreatedAt       pgtype.Timestamptz
	UpdatedAt       pgtype.Timestamptz
}

func getProjectsByIdentifier(ctx context.Context, conn conn, arg []getProjectsByIdentifierParams) ([]*getProjectsByIdentifierRow, error) {
	var rows []*getProjectsByIdentifierRow

	j, err := json.Marshal(arg)
	if err != nil {
		return rows, err
	}

	err = pgxscan.Select(ctx, conn, &rows, getProjectsByIdentifierQuery, string(j))
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return rows, nil
}

const createProjectQuery = `
INSERT INTO projects(
	names,
	descriptions,
	founding_date,
	dissolution_date,
	attributes,
	identifiers
)
VAlUES($1, $2, $3, $4, $5, $6)
RETURNING id
`

type createProjectParams struct {
	Names           []Text
	Descriptions    []Text
	FoundingDate    pgtype.Text
	DissolutionDate pgtype.Text
	Attributes      []Attribute
	Identifiers     []Identifier
}

func createProject(ctx context.Context, conn conn, arg *createProjectParams) (int64, error) {
	var id int64

	row := conn.QueryRow(ctx, createProjectQuery,
		arg.Names,
		arg.Descriptions,
		arg.FoundingDate,
		arg.DissolutionDate,
		arg.Attributes,
		arg.Identifiers,
	)
	err := row.Scan(&id)

	return id, err
}

const updateProjectQuery = `
UPDATE projects SET (
	names,
	descriptions,
	founding_date,
	dissolution_date,
	attributes,
	identifiers,
	updated_at
) = ($2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
WHERE id = $1
`

type updateProjectParams struct {
	ID              int64
	Names           []Text
	Descriptions    []Text
	FoundingDate    pgtype.Text
	DissolutionDate pgtype.Text
	Attributes      []Attribute
	Identifiers     []Identifier
}

func updateProject(ctx context.Context, conn conn, arg *updateProjectParams) error {
	_, err := conn.Exec(ctx, updateProjectQuery,
		arg.ID,
		arg.Names,
		arg.Descriptions,
		arg.FoundingDate,
		arg.DissolutionDate,
		arg.Attributes,
		arg.Identifiers,
	)

	return err
}
