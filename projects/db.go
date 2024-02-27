package projects

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

type projectRow struct {
	ID              int64
	Identifiers     []Identifier
	Names           []Text
	Descriptions    []Text
	FoundingDate    pgtype.Text
	DissolutionDate pgtype.Text
	Deleted         bool
	Attributes      []Attribute
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (row *projectRow) toProject() *Project {
	return &Project{
		Identifiers:     row.Identifiers,
		Names:           row.Names,
		Descriptions:    row.Descriptions,
		FoundingDate:    row.FoundingDate.String,
		DissolutionDate: row.DissolutionDate.String,
		Deleted:         row.Deleted,
		Attributes:      row.Attributes,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
}

const getProjectByIdentifierQuery = `
SELECT id,
	   json_agg(json_build_object('kind', ids.kind, 'value', ids.value)) AS identifiers,
	   names,
	   descriptions,
	   founding_date,
	   dissolution_date,
	   deleted,
	   attributes,
	   created_at,
	   updated_at
FROM projects p
JOIN project_identifiers pi ON p.id = pi.project_id AND pi.kind = $1 AND pi.value = $2
LEFT JOIN project_identifiers ids on p.id = ids.project_id
WHERE p.replaced_by_id IS NULL
GROUP BY p.id;
`

func getProjectByIdentifier(ctx context.Context, conn Conn, kind, value string) (*projectRow, error) {
	var r projectRow

	err := conn.QueryRow(ctx, getProjectByIdentifierQuery, kind, value).Scan(
		&r.ID,
		&r.Identifiers,
		&r.Names,
		&r.Descriptions,
		&r.FoundingDate,
		&r.DissolutionDate,
		&r.Deleted,
		&r.Attributes,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

const getAllProjectsQuery = `
SELECT id,
	   json_agg(json_build_object('kind', ids.kind, 'value', ids.value)) AS identifiers,
	   names,
	   descriptions,
	   founding_date,
	   dissolution_date,
	   deleted,
	   attributes,
	   created_at,
	   updated_at
FROM projects p
LEFT JOIN project_identifiers ids on p.id = ids.project_id
WHERE p.replaced_by_id IS NULL
GROUP BY p.id;
`

const insertProjectQuery = `
INSERT INTO projects (
	names,
	descriptions,
	founding_date,
	dissolution_date,
	deleted,
	attributes
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
`

func createProject(ctx context.Context, conn Conn, params AddProjectParams) (int64, error) {
	var id int64
	err := conn.QueryRow(ctx, insertProjectQuery,
		params.Names,
		params.Descriptions,
		pgtype.Text{Valid: params.FoundingDate != "", String: params.FoundingDate},
		pgtype.Text{Valid: params.DissolutionDate != "", String: params.DissolutionDate},
		params.Deleted,
		params.Attributes,
	).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, replaceProjectIdentifiers(ctx, conn, id, params.Identifiers)
}

const updateProjectQuery = `
UPDATE projects SET
	names = $2,
	descriptions = $3,
	founding_date = $4,
	dissolution_date = $5,
	deleted = $6,
	attributes = $7,
	updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
`

func updateProject(ctx context.Context, conn Conn, id int64, params AddProjectParams) error {
	_, err := conn.Exec(ctx, updateProjectQuery,
		id,
		params.Names,
		params.Descriptions,
		pgtype.Text{Valid: params.FoundingDate != "", String: params.FoundingDate},
		pgtype.Text{Valid: params.DissolutionDate != "", String: params.DissolutionDate},
		params.Deleted,
		params.Attributes,
	)
	if err != nil {
		return err
	}

	return replaceProjectIdentifiers(ctx, conn, id, params.Identifiers)
}

const deleteProjectIdentifiersQuery = `
DELETE FROM project_identifiers
WHERE project_id = $1;
`

const insertProjectIdentifierQuery = `
INSERT INTO project_identifiers (
	project_id,
	kind,
	value
) VALUES ($1, $2, $3);
`

func replaceProjectIdentifiers(ctx context.Context, conn Conn, pID int64, ids []Identifier) error {
	if _, err := conn.Exec(ctx, deleteProjectIdentifiersQuery, pID); err != nil {
		return err
	}
	for _, id := range ids {
		if _, err := conn.Exec(ctx, insertProjectIdentifierQuery, pID, id.Kind, id.Value); err != nil {
			return err
		}
	}
	return nil
}

const setProjectReplacedByQuery = `
UPDATE projects
SET replaced_by_id = $2
WHERE id = $1
`

func setProjectReplacedBy(ctx context.Context, conn Conn, id, replacedByID int64) error {
	_, err := conn.Exec(ctx, setProjectReplacedByQuery, id, replacedByID)
	return err
}
