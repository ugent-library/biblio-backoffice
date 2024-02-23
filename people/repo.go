package people

import (
	"context"
	"errors"
	"log"
	"slices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

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

func (r *Repo) GetPersonByIdentifier(ctx context.Context, kind, value string) (*Person, error) {
	row, err := getPersonByIdentifier(ctx, r.conn, kind, value)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &row.Person, nil
}

type AddPersonParams struct {
	Identifiers         []Identifier
	Name                string
	PreferredName       string
	GivenName           string
	FamilyName          string
	PreferredGivenName  string
	PreferredFamilyName string
	HonorificPrefix     string
	Email               string
	Active              bool
	Username            string
	Attributes          []Attribute
}

func (r *Repo) AddPerson(ctx context.Context, params AddPersonParams) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	log.Print("get existing...")

	var rows []*personRow

	for _, id := range params.Identifiers {
		row, err := getPersonByIdentifier(ctx, tx, id.Kind, id.Value)
		if err != nil && err != pgx.ErrNoRows {
			return err
		}
		if err == pgx.ErrNoRows {
			continue
		}

		if !slices.ContainsFunc(rows, func(r *personRow) bool { return r.ID == row.ID }) {
			rows = append(rows, row)
		}
	}

	slices.SortFunc(rows, func(a, b *personRow) int {
		if a.UpdatedAt.Before(b.UpdatedAt) {
			return 1
		}
		return -1
	})

	switch len(rows) {
	case 0:
		log.Print("insert...")
		if !slices.ContainsFunc(params.Identifiers, func(id Identifier) bool { return id.Kind == "id" }) {
			params.Identifiers = append(params.Identifiers, newID())
		}
		if _, err := insertPerson(ctx, tx, params); err != nil {
			return err
		}
	case 1:
		log.Print("update...")
		params = transferValues(rows, params)
		if err := updatePerson(ctx, tx, rows[0].ID, params); err != nil {
			return err
		}
	default:
		log.Print("merge...")
		params = transferValues(rows, params)
		id, err := insertPerson(ctx, tx, params)
		if err != nil {
			return err
		}
		for _, row := range rows {
			if err := setPersonReplacedBy(ctx, tx, row.ID, id); err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func newID() Identifier {
	return Identifier{Kind: "id", Value: uuid.NewString()}
}

func transferValues(rows []*personRow, params AddPersonParams) AddPersonParams {
	for _, row := range rows {
		for _, rowID := range row.Identifiers {
			if rowID.Kind != "id" {
				continue
			}
			if !slices.Contains(params.Identifiers, rowID) {
				params.Identifiers = append(params.Identifiers, rowID)
			}
		}

		if params.PreferredName == "" {
			params.PreferredName = row.PreferredName
		}
		if params.PreferredGivenName == "" {
			params.PreferredGivenName = row.PreferredGivenName
		}
		if params.PreferredFamilyName == "" {
			params.PreferredFamilyName = row.PreferredFamilyName
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
