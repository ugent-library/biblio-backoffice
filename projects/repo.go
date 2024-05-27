package projects

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ugent-library/biblio-backoffice/models"
)

const idKind = "id"

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
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("repo.GetProjectByIdentifier: %w", models.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	return row.toProject(), nil
}

func (r *Repo) EachProject(ctx context.Context, fn func(*Project) bool) error {
	rows, err := r.conn.Query(ctx, getAllProjectsQuery)
	if err != nil {
		return fmt.Errorf("repo.EachProject: %w", err)
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
			&r.PublicationCount,
			&r.Attributes,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return fmt.Errorf("repo.EachProject: %w", err)
		}
		if ok := fn(r.toProject()); !ok {
			break
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("repo.EachProject: %w", err)
	}

	return rows.Err()
}

func (r *Repo) CountProjects(ctx context.Context) (int64, error) {
	var count int64
	err := r.conn.QueryRow(ctx, "SELECT COUNT(*) FROM projects").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repo.CountProjects: %w", err)
	}
	return count, nil
}

func (r *Repo) ImportProject(ctx context.Context, p ImportProjectParams) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("repo.ImportProject: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, ident := range p.Identifiers {
		_, err := getProjectByIdentifier(ctx, tx, ident.Kind, ident.Value)
		if errors.Is(err, pgx.ErrNoRows) {
			continue
		}
		if err != nil {
			return fmt.Errorf("repo.ImportProject: %w", err)
		}

		return &DuplicateError{ident.String()}
	}

	if !p.Identifiers.Has(idKind) {
		p.Identifiers = append(p.Identifiers, newID())
	}

	err = insertProject(ctx, tx, p)
	if err != nil {
		return fmt.Errorf("repo.ImportProject: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("repo.ImportProject: %w", err)
	}

	return nil
}

func (r *Repo) AddProject(ctx context.Context, p AddProjectParams) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("repo.AddProject: %w", err)
	}
	defer tx.Rollback(ctx)

	var rows projectRows

	for _, ident := range p.Identifiers {
		row, err := getProjectByIdentifier(ctx, tx, ident.Kind, ident.Value)
		if errors.Is(err, pgx.ErrNoRows) {
			continue
		}
		if err != nil {
			return fmt.Errorf("repo.AddProject: %w", err)
		}

		if !rows.Has(row.ID) {
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
		if !p.Identifiers.Has(idKind) {
			p.Identifiers = append(p.Identifiers, newID())
		}
		if _, err := createProject(ctx, tx, p); err != nil {
			return fmt.Errorf("repo.AddProject: %w", err)
		}
	case 1:
		p = transferValues(rows, p)
		if err := updateProject(ctx, tx, rows[0].ID, p); err != nil {
			return fmt.Errorf("repo.AddProject: %w", err)
		}
	default:
		p = transferValues(rows, p)
		id, err := createProject(ctx, tx, p)
		if err != nil {
			return fmt.Errorf("repo.AddProject: %w", err)
		}
		for _, row := range rows {
			if err := setProjectReplacedBy(ctx, tx, row.ID, id); err != nil {
				return fmt.Errorf("repo.AddProject: %w", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("repo.AddProject: %w", err)
	}

	return nil
}

func (r *Repo) SetProjectPublicationCount(ctx context.Context, idKind, idValue string, n int) error {
	if _, err := r.conn.Exec(ctx, setProjectPublicationCountQuery, idKind, idValue, n); err != nil {
		return fmt.Errorf("repo.SetProjectPublicationCount: %w", err)
	}
	return nil
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
