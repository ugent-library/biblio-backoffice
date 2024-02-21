package projects

import (
	"context"
	"errors"
	"slices"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
)

var ErrNotFound = errors.New("not found")

type Repo struct {
	conn conn
}

type RepoConfig struct {
	Conn conn
}

func NewRepo(c RepoConfig) (*Repo, error) {
	return &Repo{
		conn: c.Conn,
	}, nil
}

func (r *Repo) AddProject(ctx context.Context, p *Project) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var existingProjects []getProjectsByIdentifierRow
	for _, id := range p.Identifiers {
		found, err := getProjectsByIdentifier(ctx, tx, getProjectsByIdentifierParams(id))
		if err != nil && err != pgx.ErrNoRows {
			return err
		}
		if err == pgx.ErrNoRows {
			continue
		}

		for _, fp := range found {
			if !slices.ContainsFunc(existingProjects, func(p getProjectsByIdentifierRow) bool { return p.ID == fp.ID }) {
				existingProjects = append(existingProjects, *fp)
			}
		}
	}

	slices.SortFunc(existingProjects, func(a, b getProjectsByIdentifierRow) int {
		if a.UpdatedAt.Time.Before(b.UpdatedAt.Time) {
			return 1
		}

		return -1
	})

	if len(existingProjects) == 0 {
		createProject(ctx, tx, &createProjectParams{
			Name:            p.Names,
			Description:     p.Descriptions,
			FoundingDate:    pgtype.Text{String: p.FoundingDate},
			DissolutionDate: pgtype.Text{String: p.DissolutionDate},
			Attributes:      p.Attributes,
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
