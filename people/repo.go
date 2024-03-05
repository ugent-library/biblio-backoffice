package people

import (
	"context"
	"errors"
	"slices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const idKind = "id"

var ErrNotFound = errors.New("not found")

type Repo struct {
	conn Conn
}

type RepoConfig struct {
	Conn *pgxpool.Pool
}

func NewRepo(c RepoConfig) (*Repo, error) {
	return &Repo{
		conn: c.Conn,
	}, nil
}

// TODO token encryption
// TODO make idempotent, do nothing if an ident exists
func (r *Repo) ImportOrganizations(ctx context.Context, iter Iter[ImportOrganizationParams]) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var iterErr error
	err = iter(ctx, func(o ImportOrganizationParams) bool {
		if !o.Identifiers.Has(idKind) {
			o.Identifiers = append(o.Identifiers, newID())
		}
		iterErr = insertOrganization(ctx, tx, o)
		return iterErr == nil
	})
	if iterErr != nil {
		return iterErr
	}
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// TODO make idempotent, do nothing if an ident exists
func (r *Repo) ImportPerson(ctx context.Context, p ImportPersonParams) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if !p.Identifiers.Has(idKind) {
		p.Identifiers = append(p.Identifiers, newID())
	}

	if err := insertPerson(ctx, tx, p); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repo) GetOrganizationByIdentifier(ctx context.Context, kind, value string) (*Organization, error) {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	row, err := getOrganizationByIdentifier(ctx, tx, kind, value)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	org := row.toOrganization()

	if row.ParentID.Valid {
		parentRows, err := tx.Query(ctx, getParentOrganizations, row.ParentID.Int64)
		if err != nil {
			return nil, err
		}
		defer parentRows.Close()
		for parentRows.Next() {
			var o organizationRow
			if err := parentRows.Scan(
				&o.Identifiers,
				&o.Names,
				&o.Ceased,
				&o.CreatedAt,
				&o.UpdatedAt,
			); err != nil {
				return nil, err
			}
			org.Parents = append(org.Parents, o.toParentOrganization())
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return org, nil
}

func (r *Repo) GetPersonByIdentifier(ctx context.Context, kind, value string) (*Person, error) {
	row, err := getPersonByIdentifier(ctx, r.conn, kind, value, false)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return row.toPerson(), nil
}

func (r *Repo) GetActivePersonByIdentifier(ctx context.Context, kind, value string) (*Person, error) {
	row, err := getPersonByIdentifier(ctx, r.conn, kind, value, true)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return row.toPerson(), nil
}

func (r *Repo) GetActivePersonByUsername(ctx context.Context, username string) (*Person, error) {
	row, err := getActivePersonByUsername(ctx, r.conn, username)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return row.toPerson(), nil
}

// TODO get "conn busy" error when wrapping in tx
func (r *Repo) EachOrganization(ctx context.Context, fn func(*Organization) bool) error {
	rows, err := r.conn.Query(ctx, getAllOrganizationsQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var row organizationRow
		if err := rows.Scan(
			&row.ID,
			&row.ParentID,
			&row.Identifiers,
			&row.Names,
			&row.Ceased,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return err
		}

		org := row.toOrganization()

		if row.ParentID.Valid {
			parentRows, err := r.conn.Query(ctx, getParentOrganizations, row.ParentID.Int64)
			if err != nil {
				return err
			}
			defer parentRows.Close()
			for parentRows.Next() {
				var o organizationRow
				if err := parentRows.Scan(
					&o.Identifiers,
					&o.Names,
					&o.Ceased,
					&o.CreatedAt,
					&o.UpdatedAt,
				); err != nil {
					return err
				}
				org.Parents = append(org.Parents, o.toParentOrganization())
			}
		}

		if ok := fn(org); !ok {
			break
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

// TODO see comment about tx in EachOrganization
func (r *Repo) EachPerson(ctx context.Context, fn func(*Person) bool) error {
	rows, err := r.conn.Query(ctx, getAllPeopleQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var pr personRow
		var orgIDs []int64
		if err := rows.Scan(
			&pr.ID,
			&pr.Identifiers,
			&orgIDs,
			&pr.Name,
			&pr.PreferredName,
			&pr.GivenName,
			&pr.PreferredGivenName,
			&pr.FamilyName,
			&pr.PreferredFamilyName,
			&pr.HonorificPrefix,
			&pr.Email,
			&pr.Active,
			&pr.Username,
			&pr.Attributes,
			&pr.CreatedAt,
			&pr.UpdatedAt,
		); err != nil {
			return err
		}

		for _, orgID := range orgIDs {
			var org *Organization
			rows, err := r.conn.Query(ctx, getParentOrganizations, orgID)
			if err != nil {
				return err
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
					return err
				}
				if org == nil {
					org = o.toOrganization()
				} else {
					org.Parents = append(org.Parents, o.toParentOrganization())
				}
			}

			pr.Affiliations = append(pr.Affiliations, Affiliation{Organization: org})
		}

		if ok := fn(pr.toPerson()); !ok {
			break
		}
	}

	return rows.Err()
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

	var rows []*personRow

	for _, id := range params.Identifiers {
		row, err := getPersonByIdentifier(ctx, tx, id.Kind, id.Value, false)
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
		if !slices.ContainsFunc(params.Identifiers, func(id Identifier) bool { return id.Kind == idKind }) {
			params.Identifiers = append(params.Identifiers, newID())
		}
		if _, err := createPerson(ctx, tx, params); err != nil {
			return err
		}
	case 1:
		params = transferValues(rows, params)
		if err := updatePerson(ctx, tx, rows[0].ID, params); err != nil {
			return err
		}
	default:
		params = transferValues(rows, params)
		id, err := createPerson(ctx, tx, params)
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
	return Identifier{Kind: idKind, Value: uuid.NewString()}
}

func transferValues(rows []*personRow, params AddPersonParams) AddPersonParams {
	for _, row := range rows {
		for _, rowID := range row.Identifiers {
			if rowID.Kind != idKind {
				continue
			}
			if !slices.Contains(params.Identifiers, rowID) {
				params.Identifiers = append(params.Identifiers, rowID)
			}
		}

		if params.PreferredName == "" {
			params.PreferredName = row.PreferredName.String
		}
		if params.PreferredGivenName == "" {
			params.PreferredGivenName = row.PreferredGivenName.String
		}
		if params.PreferredFamilyName == "" {
			params.PreferredFamilyName = row.PreferredFamilyName.String
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
