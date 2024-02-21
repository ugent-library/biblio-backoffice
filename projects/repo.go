package projects

import (
	"context"
	"errors"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
		args := []getProjectsByIdentifierParams{getProjectsByIdentifierParams(id)}
		found, err := getProjectsByIdentifier(ctx, tx, args)
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
		_, err = createProject(ctx, tx, &createProjectParams{
			Names:           p.Names,
			Descriptions:    p.Descriptions,
			FoundingDate:    pgtype.Text{String: p.FoundingDate},
			DissolutionDate: pgtype.Text{String: p.DissolutionDate},
			Attributes:      p.Attributes,
			Identifiers:     p.Identifiers,
		})
		if err != nil {
			return err
		}

		return tx.Commit(ctx)
	}

	projectID := existingProjects[0].ID

	err = updateProject(ctx, tx, &updateProjectParams{
		ID:              projectID,
		Names:           p.Names,
		Descriptions:    p.Descriptions,
		FoundingDate:    pgtype.Text{String: p.FoundingDate},
		DissolutionDate: pgtype.Text{String: p.DissolutionDate},
		Attributes:      p.Attributes,
		Identifiers:     p.Identifiers,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
