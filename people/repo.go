package people

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/crypt"
)

const idKind = "id"

type Repo struct {
	conn               Conn
	tokenSecret        []byte
	deactivationPeriod time.Duration
}

type RepoConfig struct {
	Conn               *pgxpool.Pool
	TokenSecret        []byte
	DeactivationPeriod time.Duration
}

func NewRepo(c RepoConfig) (*Repo, error) {
	return &Repo{
		conn:               c.Conn,
		tokenSecret:        c.TokenSecret,
		deactivationPeriod: c.DeactivationPeriod,
	}, nil
}

func (r *Repo) ImportOrganizations(ctx context.Context, iter Iter[ImportOrganizationParams]) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("repo.ImportOrganizations: %w", err)
	}
	defer tx.Rollback(ctx)

	var iterErr error
	err = iter(ctx, func(o ImportOrganizationParams) bool {
		// return error if identifier is already known
		for _, ident := range o.Identifiers {
			_, err := getOrganizationByIdentifier(ctx, tx, ident.Kind, ident.Value)
			if errors.Is(err, pgx.ErrNoRows) {
				continue
			}
			if err != nil {
				iterErr = err
				return false
			}
			iterErr = &DuplicateError{ident.String()}
			return false
		}

		if !o.Identifiers.Has(idKind) {
			o.Identifiers = append(o.Identifiers, newID())
		}

		iterErr = insertOrganization(ctx, tx, o)
		return iterErr == nil
	})
	if iterErr != nil {
		return fmt.Errorf("repo.ImportOrganizations: %w", iterErr)
	}
	if err != nil {
		return fmt.Errorf("repo.ImportOrganizations: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *Repo) ImportPerson(ctx context.Context, p ImportPersonParams) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("repo.ImportPerson: %w", err)
	}
	defer tx.Rollback(ctx)

	// return error if identifier is already known
	for _, ident := range p.Identifiers {
		_, err := getPersonByIdentifier(ctx, tx, ident.Kind, ident.Value)
		if errors.Is(err, pgx.ErrNoRows) {
			continue
		}
		if err != nil {
			return fmt.Errorf("repo.ImportPerson: %w", err)
		}
		return fmt.Errorf("repo.ImportPerson: %w", &DuplicateError{ident.String()})
	}

	if !p.Identifiers.Has(idKind) {
		p.Identifiers = append(p.Identifiers, newID())
	}

	for i, t := range p.Tokens {
		v, err := crypt.Encrypt(r.tokenSecret, t.Value)
		if err != nil {
			return fmt.Errorf("repo.ImportPerson: can't encrypt %s token: %w", t.Kind, err)
		}
		p.Tokens[i] = Token{Kind: t.Kind, Value: v}
	}

	if err := insertPerson(ctx, tx, p); err != nil {
		return fmt.Errorf("repo.ImportPerson: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *Repo) CountOrganizations(ctx context.Context) (int64, error) {
	var count int64
	err := r.conn.QueryRow(ctx, "SELECT COUNT(*) FROM organizations").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repo.CountOrganizations: %w", err)
	}
	return count, nil
}

func (r *Repo) DeleteAllOrganizations(ctx context.Context) error {
	if _, err := r.conn.Exec(ctx, "TRUNCATE organizations CASCADE"); err != nil {
		return fmt.Errorf("repo.DeleteAllOrganizations: %w", err)
	}
	return nil
}

func (r *Repo) GetOrganizationByIdentifier(ctx context.Context, kind, value string) (*Organization, error) {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo.GetOrganizationByIdentifier: identifier %q: %w", Identifier{Kind: kind, Value: value}.String(), err)
	}
	defer tx.Rollback(ctx)

	row, err := getOrganizationByIdentifier(ctx, tx, kind, value)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("repo.GetOrganizationByIdentifier: identifier %q: %w", Identifier{Kind: kind, Value: value}.String(), models.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("repo.GetOrganizationByIdentifier: identifier %q: %w", Identifier{Kind: kind, Value: value}.String(), err)
	}

	org := row.toOrganization()

	if row.ParentID.Valid {
		parentRows, err := tx.Query(ctx, getParentOrganizations, row.ParentID.Int64)
		if err != nil {
			return nil, fmt.Errorf("repo.GetOrganizationByIdentifier: identifier %q: %w", Identifier{Kind: kind, Value: value}.String(), err)
		}
		defer parentRows.Close()
		for parentRows.Next() {
			var o organizationRow
			if err := parentRows.Scan(
				&o.Identifiers,
				&o.Names,
				&o.Ceased,
				&o.CeasedOn,
				&o.CreatedAt,
				&o.UpdatedAt,
			); err != nil {
				return nil, fmt.Errorf("repo.GetOrganizationByIdentifier: identifier %q: %w", Identifier{Kind: kind, Value: value}.String(), err)
			}
			org.Parents = append(org.Parents, o.toParentOrganization())
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("repo.GetOrganizationByIdentifier: identifier %q: %w", Identifier{Kind: kind, Value: value}, err)
	}

	return org, nil
}

func (r *Repo) CountPeople(ctx context.Context) (int64, error) {
	var count int64
	err := r.conn.QueryRow(ctx, "SELECT COUNT(*) FROM people").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repo.CountPeople: %w", err)
	}
	return count, nil
}

func (r *Repo) DeleteAllPeople(ctx context.Context) error {
	if _, err := r.conn.Exec(ctx, "TRUNCATE people CASCADE"); err != nil {
		return fmt.Errorf("repo.DeleteAllPeople: %w", err)
	}
	return nil
}

func (r *Repo) GetPersonByIdentifier(ctx context.Context, kind, value string) (*Person, error) {
	row, err := getPersonByIdentifier(ctx, r.conn, kind, value)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("repo.GetPersonByIdentifier: identifier %q: %w", Identifier{Kind: kind, Value: value}.String(), models.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("repo.GetPersonByIdentifier: identifier %q: %w", Identifier{Kind: kind, Value: value}.String(), err)
	}
	return row.toPerson(r.tokenSecret)
}

func (r *Repo) GetActivePersonByIdentifier(ctx context.Context, kind, value string) (*Person, error) {
	row, err := getActivePersonByIdentifier(ctx, r.conn, kind, value)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("repo.GetActivePersonByIdentifier: identifier %q: %w", Identifier{Kind: kind, Value: value}.String(), models.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("repo.GetActivePersonByIdentifier: identifier %q: %w", Identifier{Kind: kind, Value: value}.String(), err)
	}
	return row.toPerson(r.tokenSecret)
}

func (r *Repo) GetActivePersonByUsername(ctx context.Context, username string) (*Person, error) {
	row, err := getActivePersonByUsername(ctx, r.conn, username)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("repo.GetActivePersonByUsername: username %q: %w", username, models.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("repo.GetActivePersonByUsername: username %q: %w", username, err)
	}
	return row.toPerson(r.tokenSecret)
}

// TODO get "conn busy" error when wrapping in tx
func (r *Repo) EachOrganization(ctx context.Context, fn func(*Organization) bool) error {
	rows, err := r.conn.Query(ctx, getAllOrganizationsQuery)
	if err != nil {
		return fmt.Errorf("repo.EachOrganization: %w", err)
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
			&row.CeasedOn,
			&row.Position,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return fmt.Errorf("repo.EachOrganization: %w", err)
		}

		org := row.toOrganization()

		if row.ParentID.Valid {
			parentRows, err := r.conn.Query(ctx, getParentOrganizations, row.ParentID.Int64)
			if err != nil {
				return fmt.Errorf("repo.EachOrganization: %w", err)
			}
			defer parentRows.Close()
			for parentRows.Next() {
				var o organizationRow
				if err := parentRows.Scan(
					&o.Identifiers,
					&o.Names,
					&o.Ceased,
					&o.CeasedOn,
					&o.CreatedAt,
					&o.UpdatedAt,
				); err != nil {
					return fmt.Errorf("repo.EachOrganization: %w", err)
				}
				org.Parents = append(org.Parents, o.toParentOrganization())
			}
		}

		if ok := fn(org); !ok {
			break
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("repo.EachOrganization: %w", err)
	}

	return nil
}

// TODO see comment about tx in EachOrganization
func (r *Repo) EachPerson(ctx context.Context, fn func(*Person) bool) error {
	rows, err := r.conn.Query(ctx, getAllPeopleQuery)
	if err != nil {
		return fmt.Errorf("repo.EachPerson: %w", err)
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
			&pr.PublicationCount,
			&pr.Attributes,
			&pr.CreatedAt,
			&pr.UpdatedAt,
		); err != nil {
			return fmt.Errorf("repo.EachPerson: %w", err)
		}

		for _, orgID := range orgIDs {
			var org *Organization
			rows, err := r.conn.Query(ctx, getParentOrganizations, orgID)
			if err != nil {
				return fmt.Errorf("repo.EachPerson: %w", err)
			}
			defer rows.Close()
			for rows.Next() {
				var o organizationRow
				if err := rows.Scan(
					&o.Identifiers,
					&o.Names,
					&o.Ceased,
					&o.CeasedOn,
					&o.CreatedAt,
					&o.UpdatedAt,
				); err != nil {
					return fmt.Errorf("repo.EachPerson: %w", err)
				}
				if org == nil {
					org = o.toOrganization()
				} else {
					org.Parents = append(org.Parents, o.toParentOrganization())
				}
			}

			pr.Affiliations = append(pr.Affiliations, Affiliation{Organization: org})
		}

		p, err := pr.toPerson(r.tokenSecret)
		if err != nil {
			return fmt.Errorf("repo.EachPerson: %w", err)
		}

		if ok := fn(p); !ok {
			break
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("repo.EachPerson: %w", err)
	}

	return nil
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
	Affiliations        []AffiliationParams
}

func (r *Repo) AddPerson(ctx context.Context, params AddPersonParams) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("repo.AddPerson: %w", err)
	}
	defer tx.Rollback(ctx)

	var rows []*personRow

	for _, id := range params.Identifiers {
		row, err := getPersonByIdentifier(ctx, tx, id.Kind, id.Value)
		if errors.Is(err, pgx.ErrNoRows) {
			continue
		}
		if err != nil {
			return fmt.Errorf("repo.AddPerson: %w", err)
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
			return fmt.Errorf("repo.AddPerson: %w", err)
		}
	case 1:
		params = transferValues(rows, params)
		if err := updatePerson(ctx, tx, rows[0].ID, params); err != nil {
			return fmt.Errorf("repo.AddPerson: %w", err)
		}
	default:
		params = transferValues(rows, params)
		id, err := createPerson(ctx, tx, params)
		if err != nil {
			return fmt.Errorf("repo.AddPerson: %w", err)
		}
		for _, row := range rows {
			if err := setPersonReplacedBy(ctx, tx, row.ID, id); err != nil {
				return fmt.Errorf("repo.AddPerson: %w", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("repo.AddPerson: %w", err)
	}

	return nil
}

func (r *Repo) DeactivatePeople(ctx context.Context) error {
	if r.deactivationPeriod == 0 {
		return nil
	}
	t := time.Now().Add(-r.deactivationPeriod)
	if _, err := r.conn.Exec(ctx, deactivatePeopleQuery, t); err != nil {
		return fmt.Errorf("repo.DeactivatePeople: %w", err)
	}
	return nil
}

func (r *Repo) SetPersonPublicationCount(ctx context.Context, idKind, idValue string, n int) error {
	if _, err := r.conn.Exec(ctx, setPersonPublicationCountQuery, idKind, idValue, n); err != nil {
		return fmt.Errorf("repo.SetPersonPublicationCount: %w", err)
	}
	return nil
}

func (r *Repo) SetPersonPreferredName(ctx context.Context, idKind, idValue, givenName, familyName string) error {
	if _, err := r.conn.Exec(ctx, setPersonPreferredNameQuery, idKind, idValue, givenName, familyName); err != nil {
		return fmt.Errorf("repo.SetPersonPreferredName: %w", err)
	}
	return nil
}

func newID() Identifier {
	return Identifier{Kind: idKind, Value: uuid.NewString()}
}

// TODO transfer tokens?
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
