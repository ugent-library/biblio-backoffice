package projects

import (
	"context"
	"errors"
	"slices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const idKind = "id"

var ErrNotFound = errors.New("not found")

type Repo struct {
	conn Conn
}

type RepoConfig struct {
	Conn Conn
}

func NewRepo(c RepoConfig) (*Repo, error) {
	return &Repo{
		conn: c.Conn,
	}, nil
}

func (r *Repo) GetProjectByIdentifier(ctx context.Context, kind, value string) (*Project, error) {
	row, err := getProjectByIdentifier(ctx, r.conn, kind, value)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return row.toProject(), nil
}

func (r *Repo) EachProject(ctx context.Context, fn func(*Project) bool) error {
	rows, err := r.conn.Query(ctx, getAllProjectsQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var r projectRow
		if err := rows.Scan(
			&r.ID,
			&r.Identifiers,
			&r.Names,
			&r.Descriptions,
			&r.StartDate,
			&r.EndDate,
			&r.Deleted,
			&r.Attributes,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return nil
		}
		if ok := fn(r.toProject()); !ok {
			break
		}
	}

	return rows.Err()
}

type AddProjectParams struct {
	Identifiers  []Identifier
	Names        []Text
	Descriptions []Text
	StartDate    string
	EndDate      string
	Deleted      bool
	Attributes   []Attribute
}

func (r *Repo) AddProject(ctx context.Context, params AddProjectParams) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var rows []*projectRow

	for _, id := range params.Identifiers {
		row, err := getProjectByIdentifier(ctx, tx, id.Kind, id.Value)
		if err != nil && err != pgx.ErrNoRows {
			return err
		}

		if err == pgx.ErrNoRows {
			continue
		}

		if !slices.ContainsFunc(rows, func(r *projectRow) bool { return r.ID == row.ID }) {
			rows = append(rows, row)
		}
	}

	slices.SortFunc(rows, func(a, b *projectRow) int {
		if a.UpdatedAt.Before(b.UpdatedAt) {
			return 1
		}
		return -1
	})

	switch len(rows) {
	case 0:
		if !slices.ContainsFunc(params.Identifiers, func(id Identifier) bool { return id.Kind == idKind }) {
			params.Identifiers = append(params.Identifiers, newID())
		}
		if _, err := createProject(ctx, tx, params); err != nil {
			return err
		}
	case 1:
		params = transferValues(rows, params)
		if err := updateProject(ctx, tx, rows[0].ID, params); err != nil {
			return err
		}
	default:
		params = transferValues(rows, params)
		id, err := createProject(ctx, tx, params)
		if err != nil {
			return err
		}
		for _, row := range rows {
			if err := setProjectReplacedBy(ctx, tx, row.ID, id); err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func newID() Identifier {
	return Identifier{Kind: idKind, Value: uuid.NewString()}
}

func transferValues(rows []*projectRow, params AddProjectParams) AddProjectParams {
	for _, row := range rows {
		for _, rowID := range row.Identifiers {
			if rowID.Kind != idKind {
				continue
			}
			if !slices.Contains(params.Identifiers, rowID) {
				params.Identifiers = append(params.Identifiers, rowID)
			}
		}

		for _, name := range row.Names {
			if !slices.Contains(params.Names, name) {
				params.Names = append(params.Names, name)
			}
		}

		for _, desc := range row.Descriptions {
			if !slices.Contains(params.Descriptions, desc) {
				params.Descriptions = append(params.Descriptions, desc)
			}
		}

		if params.StartDate == "" {
			params.EndDate = row.StartDate.String
		}

		if params.EndDate == "" {
			params.EndDate = row.EndDate.String
		}

		var attrs []Attribute
		for _, attr := range row.Attributes {
			if !slices.ContainsFunc(params.Attributes, func(a Attribute) bool { return a.Scope == attr.Scope }) {
				attrs = append(attrs, attr)
			}
		}
		if len(attrs) > 0 {
			params.Attributes = append(params.Attributes, attrs...)
		}
	}

	return params
}
