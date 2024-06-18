// TODO all mutating methods should call Validate() before saving
package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/ugent-library/biblio-backoffice/db"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/mutate"
	"github.com/ugent-library/biblio-backoffice/snapstore"
)

type Mutation struct {
	Name string
	Args []string
	// Line is not always set
	Line int
}

type Repo struct {
	config           Config
	client           *snapstore.Client
	publicationStore *snapstore.Store
	datasetStore     *snapstore.Store
	opts             snapstore.Options
	// sqlc
	queries *db.Queries
	conn    *pgxpool.Pool
}

type Config struct {
	Conn                 *pgxpool.Pool
	PublicationListeners []PublicationListener
	DatasetListeners     []DatasetListener
	PublicationMutators  map[string]PublicationMutator
	DatasetMutators      map[string]DatasetMutator
	PublicationLoaders   []PublicationVisitor
	DatasetLoaders       []DatasetVisitor
}

type PublicationListener = func(*models.Publication)
type DatasetListener = func(*models.Dataset)
type PublicationMutator = func(*models.Publication, []string) error
type DatasetMutator = func(*models.Dataset, []string) error
type PublicationVisitor = func(*models.Publication) error
type DatasetVisitor = func(*models.Dataset) error

func New(c Config) (*Repo, error) {
	client := snapstore.New(c.Conn, []string{"publications", "datasets"},
		snapstore.WithIDGenerator(func() (string, error) {
			return ulid.Make().String(), nil
		}),
	)

	return &Repo{
		config:           c,
		client:           client,
		publicationStore: client.Store("publications"),
		datasetStore:     client.Store("datasets"),
		queries:          db.New(c.Conn),
		conn:             c.Conn,
	}, nil
}

func (s *Repo) publicationNotify(p *models.Publication) {
	for _, fn := range s.config.PublicationListeners {
		fn(p)
	}
}

func (s *Repo) datasetNotify(d *models.Dataset) {
	for _, fn := range s.config.DatasetListeners {
		fn(d)
	}
}

func (s *Repo) tx(ctx context.Context, fn func(*Repo) error) error {
	return s.client.Tx(ctx, func(opts snapstore.Options) error {
		return fn(&Repo{
			config:           s.config,
			client:           s.client,
			publicationStore: s.publicationStore,
			datasetStore:     s.datasetStore,
			opts:             opts,
			queries:          db.New(s.conn),
		})
	})
}

func (s *Repo) GetPublication(id string) (*models.Publication, error) {
	snap, err := s.publicationStore.GetCurrentSnapshot(id, s.opts)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("repo.GetPublication %s: %w", id, models.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("repo.GetPublication %s: %w", id, err)
	}
	p, err := s.snapshotToPublication(snap)
	if err != nil {
		return nil, fmt.Errorf("repo.GetPublication %s: %w", id, err)
	}
	return p, nil
}

func (s *Repo) GetPublications(ids []string) ([]*models.Publication, error) {
	c, err := s.publicationStore.GetByID(ids, s.opts)
	if err != nil {
		return nil, fmt.Errorf("repo.GetPublications: %w", err)
	}
	defer c.Close()
	var publications []*models.Publication
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return nil, fmt.Errorf("repo.GetPublications: %w", err)
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return nil, fmt.Errorf("repo.GetPublications: %w", err)
		}
		publications = append(publications, p)
	}
	if c.Err() != nil {
		return nil, fmt.Errorf("repo.GetPublications: %w", c.Err())
	}
	return publications, nil
}

func (s *Repo) ImportPublication(p *models.Publication) error {
	if p.DateCreated == nil {
		return fmt.Errorf("repo.ImportPublication %s: DateCreated not set", p.ID)
	}
	if p.DateUpdated == nil {
		return fmt.Errorf("repo.ImportPublication %s: DateUpdated not set", p.ID)
	}
	if p.DateFrom == nil {
		return fmt.Errorf("repo.ImportPublication %s: DateFrom not set", p.ID)
	}
	if p.DateUntil != nil {
		return fmt.Errorf("repo.ImportPublication %s: DateUntil not nil", p.ID)
	}

	snap, err := s.publicationToSnapshot(p)
	if err != nil {
		return fmt.Errorf("repo.ImportPublication %s: %w", p.ID, err)
	}

	if err := s.publicationStore.ImportSnapshot(snap, s.opts); err != nil {
		return fmt.Errorf("repo.ImportPublication %s: %w", p.ID, err)
	}

	for _, fn := range s.config.PublicationLoaders {
		if err := fn(p); err != nil {
			return fmt.Errorf("repo.ImportPublication %s: %w", p.ID, err)
		}
	}

	s.publicationNotify(p)

	return nil
}

func (s *Repo) SavePublication(p *models.Publication, u *models.Person) error {
	oldPublication, err := s.GetPublication(p.ID)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return fmt.Errorf("repo.SavePublication %s: %w", p.ID, err)
	}
	if reflect.DeepEqual(oldPublication, p) {
		return nil
	}

	now := time.Now()

	if p.DateCreated == nil {
		p.DateCreated = &now
	}
	p.DateUpdated = &now

	if u != nil {
		p.UserID = u.ID
		p.User = u
		p.LastUserID = u.ID
		p.LastUser = u
	} else {
		p.UserID = ""
		p.User = nil
	}

	if err := p.Validate(); err != nil {
		return fmt.Errorf("repo.SavePublication %s: %w", p.ID, err)
	}

	if p.Status == "public" && !p.HasBeenPublic {
		p.HasBeenPublic = true
	}

	if err := s.publicationStore.Add(p.ID, p, s.opts); err != nil {
		return fmt.Errorf("repo.SavePublication %s: %w", p.ID, err)
	}

	for _, fn := range s.config.PublicationLoaders {
		if err := fn(p); err != nil {
			return fmt.Errorf("repo.SavePublication %s: %w", p.ID, err)
		}
	}

	s.publicationNotify(p)

	return nil
}

func (s *Repo) UpdatePublication(snapshotID string, p *models.Publication, u *models.Person) error {
	if oldPublication, err := s.GetPublication(p.ID); err != nil {
		return fmt.Errorf("repo.UpdatePublication %s@%s: %w", p.ID, snapshotID, err)
	} else if reflect.DeepEqual(oldPublication, p) {
		return nil
	}

	oldDateUpdated := p.DateUpdated
	now := time.Now()
	p.DateUpdated = &now

	if u != nil {
		p.UserID = u.ID
		p.User = u
		p.LastUserID = u.ID
		p.LastUser = u
	} else {
		p.UserID = ""
		p.User = nil
	}

	if p.Status == "public" && !p.HasBeenPublic {
		p.HasBeenPublic = true
	}

	snapshotID, err := s.publicationStore.AddAfter(snapshotID, p.ID, p, s.opts)
	if err != nil {
		p.DateUpdated = oldDateUpdated
		return fmt.Errorf("repo.UpdatePublication %s@%s: %w", p.ID, snapshotID, err)
	}
	p.SnapshotID = snapshotID

	for _, fn := range s.config.PublicationLoaders {
		if err := fn(p); err != nil {
			return fmt.Errorf("repo.UpdatePublication %s@%s: %w", p.ID, snapshotID, err)
		}
	}

	s.publicationNotify(p)

	return nil
}

func (s *Repo) UpdatePublicationInPlace(p *models.Publication) error {
	snap, err := s.publicationStore.Update(p.SnapshotID, p.ID, p, s.opts)
	if err != nil {
		return fmt.Errorf("repo.UpdatePublicationInPlace %s: %w", p.ID, err)
	}

	np := &models.Publication{}
	if err := snap.Scan(np); err != nil {
		return fmt.Errorf("repo.UpdatePublicationInPlace %s: %w", p.ID, err)
	}

	for _, fn := range s.config.PublicationLoaders {
		if err := fn(p); err != nil {
			return fmt.Errorf("repo.UpdatePublicationInPlace %s: %w", p.ID, err)
		}
	}

	s.publicationNotify(p)

	return nil
}

func (s *Repo) MutatePublication(id string, u *models.Person, muts ...Mutation) error {
	if len(muts) == 0 {
		return nil
	}

	p, err := s.GetPublication(id)
	if err != nil {
		return fmt.Errorf("repo.MutatePublication %s: %w", id, err)
	}

	for _, mut := range muts {
		mutator, ok := s.config.PublicationMutators[mut.Name]
		if !ok {
			return fmt.Errorf("repo.MutatePublication %s: %w", p.ID, &mutate.ArgumentError{Msg: fmt.Sprintf("unknown mutation %s at line %d", mut.Name, mut.Line)})
		}
		if err := mutator(p, mut.Args); err != nil {
			var argErr *mutate.ArgumentError
			// TODO this is a messy way of adding the line number
			if mut.Line != 0 && errors.As(err, &argErr) {
				argErr.Msg = fmt.Sprintf("%s at line %d", argErr.Msg, mut.Line)
			}
			return fmt.Errorf("repo.MutatePublication %s: mutation %s: %w", p.ID, mut.Name, err)
		}
	}

	if err = p.Validate(); err != nil {
		return fmt.Errorf("repo.MutatePublication %s: %w", p.ID, err)
	}

	if err := s.UpdatePublication(p.SnapshotID, p, u); err != nil {
		return fmt.Errorf("repo.MutatePublication %s: %w", p.ID, err)
	}

	return nil
}

func (s *Repo) PublicationsAfter(t time.Time, limit, offset int) (int, []*models.Publication, error) {
	n, err := s.publicationStore.CountSql(
		"SELECT * FROM publications WHERE date_until IS NULL AND date_from >= $1",
		[]any{t},
		s.opts,
	)
	if err != nil {
		return 0, nil, fmt.Errorf("repo.PublicationsAfter: %w", err)
	}

	c, err := s.publicationStore.Select(
		"SELECT * FROM publications WHERE date_until IS NULL AND date_from >= $1 ORDER BY date_from ASC LIMIT $2 OFFSET $3",
		[]any{t, limit, offset},
		s.opts,
	)
	if err != nil {
		return 0, nil, fmt.Errorf("repo.PublicationsAfter: %w", err)
	}

	publications := make([]*models.Publication, 0, limit)

	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return 0, nil, fmt.Errorf("repo.PublicationsAfter: %w", err)
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return 0, nil, fmt.Errorf("repo.PublicationsAfter: %w", err)
		}

		publications = append(publications, p)
	}

	if c.Err() != nil {
		return 0, nil, fmt.Errorf("repo: publications after: %w", c.Err())
	}

	return n, publications, nil
}

func (s *Repo) PublicationsBetween(t1, t2 time.Time, fn func(*models.Publication) bool) error {
	c, err := s.publicationStore.Select(
		"SELECT * FROM publications WHERE date_until IS NULL AND date_from >= $1 AND date_from <= $2 ORDER BY date_from ASC",
		[]any{t1, t2},
		s.opts,
	)
	if err != nil {
		return fmt.Errorf("repo.PublicationsBetween: %w", err)
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return fmt.Errorf("repo.PublicationsBetween: %w", err)
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return fmt.Errorf("repo.PublicationsBetween: %w", err)
		}
		if ok := fn(p); !ok {
			break
		}
	}

	if c.Err() != nil {
		return fmt.Errorf("repo.PublicationsBetween: %w", c.Err())
	}

	return nil
}

func (s *Repo) EachPublication(fn func(*models.Publication) bool) error {
	c, err := s.publicationStore.GetAll(s.opts)
	if err != nil {
		return fmt.Errorf("repo.EachPublication: %w", err)
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return fmt.Errorf("repo.EachPublication: %w", err)
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return fmt.Errorf("repo.EachPublication: %w", err)
		}
		if ok := fn(p); !ok {
			break
		}
	}

	if c.Err() != nil {
		return fmt.Errorf("repo.EachPublication: %w", c.Err())
	}

	return nil
}

func (s *Repo) EachPublicationSnapshot(fn func(*models.Publication) bool) error {
	c, err := s.publicationStore.GetAllSnapshots(s.opts)
	if err != nil {
		return fmt.Errorf("repo.EachPublicationSnapshot: %w", err)
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return fmt.Errorf("repo.EachPublicationSnapshot: %w", err)
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return fmt.Errorf("repo.EachPublicationSnapshot: %w", err)
		}
		if ok := fn(p); !ok {
			break
		}
	}

	if c.Err() != nil {
		return fmt.Errorf("repo.EachPublicationSnapshot: %w", c.Err())
	}

	return nil
}

func (s *Repo) EachPublicationWithStatus(status string, fn func(*models.Publication) bool) error {
	sql := `SELECT * FROM publications WHERE date_until IS NULL AND data->>'status' = $1`

	c, err := s.publicationStore.Select(sql, []any{status}, s.opts)
	if err != nil {
		return fmt.Errorf("repo.EachPublicationWithStatus: %w", err)
	}
	defer c.Close()

	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return fmt.Errorf("repo.EachPublicationWithStatus: %w", err)
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return fmt.Errorf("repo.EachPublicationWithStatus: %w", err)
		}
		if ok := fn(p); !ok {
			break
		}
	}

	if c.Err() != nil {
		return fmt.Errorf("repo.EachPublicationWithStatus: %w", c.Err())
	}

	return nil
}

// TODO add handle with a listener, then this method isn't needed anymore
func (s *Repo) EachPublicationWithoutHandle(fn func(*models.Publication) bool) error {
	sql := `
		SELECT * FROM publications WHERE date_until IS NULL AND
		data->>'status' = 'public' AND
		NOT data ? 'handle'
		`
	c, err := s.publicationStore.Select(sql, nil, s.opts)
	if err != nil {
		return fmt.Errorf("repo.EachPublicationWithoutHandle: %w", err)
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return fmt.Errorf("repo.EachPublicationWithoutHandle: %w", err)
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return fmt.Errorf("repo.EachPublicationWithoutHandle: %w", err)
		}
		if ok := fn(p); !ok {
			break
		}
	}

	if c.Err() != nil {
		return fmt.Errorf("repo.EachPublicationWithoutHandle: %w", c.Err())
	}

	return nil
}

func (s *Repo) GetPublicationSnapshotBefore(id string, dateFrom time.Time) (*models.Publication, error) {
	snap, err := s.publicationStore.GetSnapshotBefore(id, dateFrom, s.opts)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("repo.GetPublicationSnapshotBefore %s: %w", id, models.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("repo.GetPublicationSnapshotBefore %s: %w", id, err)
	}
	p, err := s.snapshotToPublication(snap)
	if err != nil {
		return nil, fmt.Errorf("repo.GetPublicationSnapshotBefore %s: %w", id, err)
	}
	return p, nil
}

func (s *Repo) PublicationHistory(id string, fn func(*models.Publication) bool) error {
	c, err := s.publicationStore.GetHistory(id, s.opts)
	if err != nil {
		return fmt.Errorf("repo.PublicationHistory %s: %w", id, err)
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return fmt.Errorf("repo.PublicationHistory %s: %w", id, err)
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return fmt.Errorf("repo.PublicationHistory %s: %w", id, err)
		}
		if ok := fn(p); !ok {
			break
		}
	}

	if c.Err() != nil {
		return fmt.Errorf("repo.PublicationHistory %s: %w", id, c.Err())
	}

	return nil
}

func (s *Repo) UpdatePublicationEmbargoes() (int, error) {
	var n int
	embargoAccessLevel := "info:eu-repo/semantics/embargoedAccess"
	now := time.Now().Format("2006-01-02")
	sql := `
		SELECT * FROM publications WHERE date_until IS NULL AND
		data->'file' IS NOT NULL AND
		EXISTS(
			SELECT 1 FROM jsonb_array_elements(data->'file') AS f
			WHERE f->>'access_level' = $1 AND
			f->>'embargo_date' <= $2
		)
		`
	c, err := s.publicationStore.Select(sql, []any{embargoAccessLevel, now}, s.opts)
	if err != nil {
		return n, fmt.Errorf("repo.UpdatePublicationEmbargoes: %w", err)
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return n, fmt.Errorf("repo.UpdatePublicationEmbargoes: %w", err)
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return n, fmt.Errorf("repo.UpdatePublicationEmbargoes: %w", err)
		}

		// clear expired embargoes
		for _, file := range p.File {
			if file.AccessLevel != embargoAccessLevel {
				continue
			}
			// TODO: what with empty embargo_date?
			if file.EmbargoDate == "" {
				continue
			}
			if file.EmbargoDate > now {
				continue
			}
			file.ClearEmbargo()
		}

		if err = s.UpdatePublication(p.SnapshotID, p, nil); err != nil {
			return n, fmt.Errorf("repo.UpdatePublicationEmbargoes: %w", err)
		}

		n++
	}

	if c.Err() != nil {
		return n, fmt.Errorf("repo.UpdatePublicationEmbargoes: %w", c.Err())
	}

	return n, nil
}

func (s *Repo) GetDataset(id string) (*models.Dataset, error) {
	snap, err := s.datasetStore.GetCurrentSnapshot(id, s.opts)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("repo.GetDataset %s: %w", id, models.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	d, err := s.snapshotToDataset(snap)
	if err != nil {
		return nil, fmt.Errorf("repo.GetDataset %s: %w", id, err)
	}
	return d, nil
}

func (s *Repo) GetDatasets(ids []string) ([]*models.Dataset, error) {
	c, err := s.datasetStore.GetByID(ids, s.opts)
	if err != nil {
		return nil, fmt.Errorf("repo.GetDatasets: %w", err)
	}
	defer c.Close()
	var datasets []*models.Dataset
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return nil, fmt.Errorf("repo.GetDatasets: %w", err)
		}
		d, err := s.snapshotToDataset(snap)
		if err != nil {
			return nil, fmt.Errorf("repo.GetDatasets: %w", err)
		}
		datasets = append(datasets, d)
	}

	if c.Err() != nil {
		return nil, fmt.Errorf("repo.GetDatasets: %w", c.Err())
	}

	return datasets, nil
}

func (s *Repo) ImportDataset(d *models.Dataset) error {
	if d.DateCreated == nil {
		return fmt.Errorf("repo.ImportDataset %s: DateCreated not set", d.ID)
	}
	if d.DateUpdated == nil {
		return fmt.Errorf("repo.ImportDataset %s: DateUpdated not set", d.ID)
	}
	if d.DateFrom == nil {
		return fmt.Errorf("repo.ImportDataset %s: DateFrom not set", d.ID)
	}
	if d.DateUntil != nil {
		return fmt.Errorf("repo.ImportDataset %s: DateUntil not nil", d.ID)
	}

	snap, err := s.datasetToSnapshot(d)
	if err != nil {
		return fmt.Errorf("repo.ImportDataset %s: %w", d.ID, err)
	}

	if err := s.datasetStore.ImportSnapshot(snap, s.opts); err != nil {
		return fmt.Errorf("repo.ImportDataset %s: %w", d.ID, err)
	}

	for _, fn := range s.config.DatasetLoaders {
		if err := fn(d); err != nil {
			return fmt.Errorf("repo.ImportDataset %s: %w", d.ID, err)
		}
	}

	s.datasetNotify(d)

	return nil
}

func (s *Repo) SaveDataset(d *models.Dataset, u *models.Person) error {
	oldDataset, err := s.GetDataset(d.ID)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return fmt.Errorf("repo.SaveDataset %s: %w", d.ID, err)
	}
	if reflect.DeepEqual(oldDataset, d) {
		return nil
	}

	now := time.Now()

	if d.DateCreated == nil {
		d.DateCreated = &now
	}
	d.DateUpdated = &now

	if u != nil {
		d.UserID = u.ID
		d.User = u
		d.LastUserID = u.ID
		d.LastUser = u
	} else {
		d.UserID = ""
		d.User = nil
	}

	if err := d.Validate(); err != nil {
		return fmt.Errorf("repo.SaveDataset %s: %w", d.ID, err)
	}

	if d.Status == "public" && !d.HasBeenPublic {
		d.HasBeenPublic = true
	}

	if err := s.datasetStore.Add(d.ID, d, s.opts); err != nil {
		return fmt.Errorf("repo.SaveDataset %s: %w", d.ID, err)
	}

	for _, fn := range s.config.DatasetLoaders {
		if err := fn(d); err != nil {
			return fmt.Errorf("repo.SaveDataset %s: %w", d.ID, err)
		}
	}

	s.datasetNotify(d)

	return nil
}

func (s *Repo) UpdateDataset(snapshotID string, d *models.Dataset, u *models.Person) error {
	if oldDataset, err := s.GetDataset(d.ID); err != nil {
		return fmt.Errorf("repo.UpdateDataset %s@%s: %w", d.ID, snapshotID, err)
	} else if reflect.DeepEqual(oldDataset, d) {
		return nil
	}

	oldDateUpdated := d.DateUpdated
	now := time.Now()
	d.DateUpdated = &now

	if u != nil {
		d.UserID = u.ID
		d.User = u
		d.LastUserID = u.ID
		d.LastUser = u
	} else {
		d.UserID = ""
		d.User = nil
	}

	if d.Status == "public" && !d.HasBeenPublic {
		d.HasBeenPublic = true
	}

	snapshotID, err := s.datasetStore.AddAfter(snapshotID, d.ID, d, s.opts)
	if err != nil {
		d.DateUpdated = oldDateUpdated
		return fmt.Errorf("repo.UpdateDataset %s@%s: %w", d.ID, snapshotID, err)
	}
	d.SnapshotID = snapshotID

	for _, fn := range s.config.DatasetLoaders {
		if err := fn(d); err != nil {
			return fmt.Errorf("repo.UpdateDataset %s@%s: %w", d.ID, snapshotID, err)
		}
	}

	s.datasetNotify(d)

	return nil
}

func (s *Repo) MutateDataset(id string, u *models.Person, muts ...Mutation) error {
	if len(muts) == 0 {
		return nil
	}

	d, err := s.GetDataset(id)
	if err != nil {
		return err
	}

	for _, mut := range muts {
		mutator, ok := s.config.DatasetMutators[mut.Name]
		if !ok {
			return fmt.Errorf("repo.MutateDataset %s: %w", d.ID, &mutate.ArgumentError{Msg: fmt.Sprintf("unknown mutation %s at line %d", mut.Name, mut.Line)})
		}
		if err := mutator(d, mut.Args); err != nil {
			var argErr *mutate.ArgumentError
			// TODO this is a messy way of adding the line number
			if mut.Line != 0 && errors.As(err, &argErr) {
				argErr.Msg = fmt.Sprintf("%s at line %d", argErr.Msg, mut.Line)
			}
			return fmt.Errorf("repo.MutateDataset %s: mutation %s: %w", id, mut.Name, err)
		}
	}

	if err = d.Validate(); err != nil {
		return fmt.Errorf("repo.MutateDataset %s: %w", id, err)
	}

	if err := s.UpdateDataset(d.SnapshotID, d, u); err != nil {
		return fmt.Errorf("repo.MutateDataset %s: %w", id, err)
	}

	return nil
}

func (s *Repo) DatasetsAfter(t time.Time, limit, offset int) (int, []*models.Dataset, error) {
	n, err := s.datasetStore.CountSql(
		"SELECT * FROM datasets WHERE date_until IS NULL AND date_from >= $1",
		[]any{t},
		s.opts,
	)
	if err != nil {
		return 0, nil, fmt.Errorf("repo.DatasetsAfter: %w", err)
	}

	c, err := s.datasetStore.Select(
		"SELECT * FROM datasets WHERE date_until IS NULL AND date_from >= $1 ORDER BY date_from ASC LIMIT $2 OFFSET $3",
		[]any{t, limit, offset},
		s.opts,
	)
	if err != nil {
		return 0, nil, fmt.Errorf("repo.DatasetsAfter: %w", err)
	}

	datasets := make([]*models.Dataset, 0, limit)

	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return 0, nil, fmt.Errorf("repo.DatasetsAfter: %w", err)
		}
		p, err := s.snapshotToDataset(snap)
		if err != nil {
			return 0, nil, fmt.Errorf("repo.DatasetsAfter: %w", err)
		}

		datasets = append(datasets, p)
	}

	if c.Err() != nil {
		return 0, nil, fmt.Errorf("repo.DatasetsAfter: %w", c.Err())
	}

	return n, datasets, nil
}

func (s *Repo) DatasetsBetween(t1, t2 time.Time, fn func(*models.Dataset) bool) error {
	c, err := s.datasetStore.Select(
		"SELECT * FROM datasets WHERE date_until IS NULL AND date_from >= $1 AND date_from <= $2 ORDER BY date_from ASC",
		[]any{t1, t2},
		s.opts,
	)
	if err != nil {
		return fmt.Errorf("repo.DatasetsBetween: %w", err)
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return fmt.Errorf("repo.DatasetsBetween: %w", err)
		}
		p, err := s.snapshotToDataset(snap)
		if err != nil {
			return fmt.Errorf("repo.DatasetsBetween: %w", err)
		}
		if ok := fn(p); !ok {
			break
		}
	}

	if c.Err() != nil {
		return fmt.Errorf("repo.DatasetsBetween: %w", c.Err())
	}

	return nil
}

func (s *Repo) EachDataset(fn func(*models.Dataset) bool) error {
	c, err := s.datasetStore.GetAll(s.opts)
	if err != nil {
		return fmt.Errorf("repo.EachDataset: %w", err)
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return fmt.Errorf("repo.EachDataset: %w", err)
		}
		d, err := s.snapshotToDataset(snap)
		if err != nil {
			return fmt.Errorf("repo.EachDataset: %w", err)
		}
		if ok := fn(d); !ok {
			break
		}
	}

	if c.Err() != nil {
		return fmt.Errorf("repo.EachDataset: %w", c.Err())
	}

	return nil
}

func (s *Repo) EachDatasetSnapshot(fn func(*models.Dataset) bool) error {
	c, err := s.datasetStore.GetAllSnapshots(s.opts)
	if err != nil {
		return fmt.Errorf("repo.EachDatasetSnapshot: %w", err)
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return fmt.Errorf("repo.EachDatasetSnapshot: %w", err)
		}
		d, err := s.snapshotToDataset(snap)
		if err != nil {
			return fmt.Errorf("repo.EachDatasetSnapshot: %w", err)
		}
		if ok := fn(d); !ok {
			break
		}
	}

	if c.Err() != nil {
		return fmt.Errorf("repo.EachDatasetSnapshot: %w", c.Err())
	}

	return nil
}

func (s *Repo) EachDatasetWithoutHandle(fn func(*models.Dataset) bool) error {
	sql := `
		SELECT * FROM datasets WHERE date_until IS NULL AND
		data->>'status' = 'public' AND
		NOT data ? 'handle'
		`
	c, err := s.datasetStore.Select(sql, nil, s.opts)
	if err != nil {
		return fmt.Errorf("repo.EachDatasetWithoutHandle: %w", err)
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return fmt.Errorf("repo.EachDatasetWithoutHandle: %w", err)
		}
		d, err := s.snapshotToDataset(snap)
		if err != nil {
			return fmt.Errorf("repo.EachDatasetWithoutHandle: %w", err)
		}
		if ok := fn(d); !ok {
			break
		}
	}

	if c.Err() != nil {
		return fmt.Errorf("repo.EachDatasetWithoutHandle: %w", c.Err())
	}

	return nil
}

func (s *Repo) GetDatasetSnapshotBefore(id string, dateFrom time.Time) (*models.Dataset, error) {
	snap, err := s.datasetStore.GetSnapshotBefore(id, dateFrom, s.opts)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("repo.GetDatasetSnapshotBefore %s: %w", id, models.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("repo.GetDatasetSnapshotBefore %s: %w", id, err)
	}
	d, err := s.snapshotToDataset(snap)
	if err != nil {
		return nil, fmt.Errorf("repo.GetDatasetSnapshotBefore %s: %w", id, err)
	}
	return d, nil
}

func (s *Repo) DatasetHistory(id string, fn func(*models.Dataset) bool) error {
	c, err := s.datasetStore.GetHistory(id, s.opts)
	if err != nil {
		return fmt.Errorf("repo.GetDatasetHistory %s: %w", id, err)
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return fmt.Errorf("repo.GetDatasetHistory %s: %w", id, err)
		}
		d, err := s.snapshotToDataset(snap)
		if err != nil {
			return fmt.Errorf("repo.GetDatasetHistory %s: %w", id, err)
		}
		if ok := fn(d); !ok {
			break
		}
	}

	if c.Err() != nil {
		return fmt.Errorf("repo.GetDatasetHistory %s: %w", id, c.Err())
	}

	return nil
}

func (s *Repo) UpdateDatasetEmbargoes() (int, error) {
	var n int
	embargoAccessLevel := "info:eu-repo/semantics/embargoedAccess"
	now := time.Now().Format("2006-01-02")
	sql := `
		SELECT * FROM datasets
		WHERE date_until is null AND
		data->>'access_level' = $1 AND
		data->>'embargo_date' <> '' AND
		data->>'embargo_date' <= $2
		`
	c, err := s.datasetStore.Select(sql, []any{embargoAccessLevel, now}, s.opts)
	if err != nil {
		return n, fmt.Errorf("repo.UpdateDatasetEmbargoes: %w", err)
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return n, fmt.Errorf("repo.UpdateDatasetEmbargoes: %w", err)
		}
		d, err := s.snapshotToDataset(snap)
		if err != nil {
			return n, fmt.Errorf("repo.UpdateDatasetEmbargoes: %w", err)
		}

		// clear expired embargo
		// TODO: what with empty embargo_date?
		if d.EmbargoDate == "" {
			continue
		}
		if d.EmbargoDate > now {
			continue
		}
		d.ClearEmbargo()

		if err = s.UpdateDataset(d.SnapshotID, d, nil); err != nil {
			return n, fmt.Errorf("repo.UpdateDatasetEmbargoes: %w", err)
		}

		n++
	}

	if c.Err() != nil {
		return n, fmt.Errorf("repo.UpdateDatasetEmbargoes: %w", c.Err())
	}

	return n, nil
}

func (s *Repo) GetPublicationDatasets(p *models.Publication) ([]*models.Dataset, error) {
	datasetIds := make([]string, len(p.RelatedDataset))
	for _, rd := range p.RelatedDataset {
		datasetIds = append(datasetIds, rd.ID)
	}

	datasets, err := s.GetDatasets(datasetIds)
	if err != nil {
		return nil, fmt.Errorf("repo.GetPublicationDatasets %s: %w", p.ID, err)
	}

	return datasets, nil
}

func (s *Repo) GetVisiblePublicationDatasets(u *models.Person, p *models.Publication) ([]*models.Dataset, error) {
	datasets, err := s.GetPublicationDatasets(p)
	if err != nil {
		return nil, err
	}
	filteredDatasets := make([]*models.Dataset, 0, len(datasets))
	for _, dataset := range datasets {
		if u.CanViewDataset(dataset) {
			filteredDatasets = append(filteredDatasets, dataset)
		}
	}
	return filteredDatasets, nil
}

func (s *Repo) GetDatasetPublications(d *models.Dataset) ([]*models.Publication, error) {
	publicationIds := make([]string, len(d.RelatedPublication))
	for _, rp := range d.RelatedPublication {
		publicationIds = append(publicationIds, rp.ID)
	}

	publications, err := s.GetPublications(publicationIds)
	if err != nil {
		return nil, fmt.Errorf("repo.GetDatasetPublications %s: %w", d.ID, err)
	}

	return publications, nil
}

func (s *Repo) GetVisibleDatasetPublications(u *models.Person, d *models.Dataset) ([]*models.Publication, error) {
	publications, err := s.GetDatasetPublications(d)
	if err != nil {
		return nil, err
	}
	filteredPublications := make([]*models.Publication, 0, len(publications))
	for _, publication := range publications {
		if u.CanViewPublication(publication) {
			filteredPublications = append(filteredPublications, publication)
		}
	}
	return filteredPublications, nil
}

func (s *Repo) AddPublicationDataset(p *models.Publication, d *models.Dataset, u *models.Person) error {
	return s.tx(context.Background(), func(s *Repo) error {
		if !p.HasRelatedDataset(d.ID) {
			p.RelatedDataset = append(p.RelatedDataset, models.RelatedDataset{ID: d.ID})
			if err := s.SavePublication(p, u); err != nil {
				return fmt.Errorf("repo.AddPublicationDataset %s %s: %w", p.ID, d.ID, err)
			}
		}
		if !d.HasRelatedPublication(p.ID) {
			d.RelatedPublication = append(d.RelatedPublication, models.RelatedPublication{ID: p.ID})
			if err := s.SaveDataset(d, u); err != nil {
				return fmt.Errorf("repo.AddPublicationDataset %s %s: %w", p.ID, d.ID, err)
			}
		}

		return nil
	})
}

func (s *Repo) RemovePublicationDataset(p *models.Publication, d *models.Dataset, u *models.Person) error {
	return s.tx(context.Background(), func(s *Repo) error {
		if p.HasRelatedDataset(d.ID) {
			p.RemoveRelatedDataset(d.ID)
			if err := s.SavePublication(p, u); err != nil {
				return fmt.Errorf("repo.RemovePublicationDataset %s %s: %w", p.ID, d.ID, err)
			}
		}
		if d.HasRelatedPublication(p.ID) {
			d.RemoveRelatedPublication(p.ID)
			if err := s.SaveDataset(d, u); err != nil {
				return fmt.Errorf("repo.RemovePublicationDataset %s %s: %w", p.ID, d.ID, err)
			}
		}

		return nil
	})
}

func (s *Repo) PurgeAllPublications() error {
	return s.publicationStore.PurgeAll(s.opts)
}

func (s *Repo) PurgePublication(id string) error {
	return s.publicationStore.Purge(id, s.opts)
}

func (s *Repo) PurgeAllDatasets() error {
	return s.datasetStore.PurgeAll(s.opts)
}

func (s *Repo) PurgeDataset(id string) error {
	return s.datasetStore.Purge(id, s.opts)
}

func (s *Repo) snapshotToPublication(snap *snapstore.Snapshot) (*models.Publication, error) {
	p := &models.Publication{}
	if err := snap.Scan(p); err != nil {
		return nil, err
	}
	p.SnapshotID = snap.SnapshotID
	p.DateFrom = snap.DateFrom
	p.DateUntil = snap.DateUntil
	for _, fn := range s.config.PublicationLoaders {
		if err := fn(p); err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (s *Repo) publicationToSnapshot(p *models.Publication) (*snapstore.Snapshot, error) {
	snap := &snapstore.Snapshot{}
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	snap.Data = data
	snap.DateFrom = p.DateFrom
	snap.DateUntil = p.DateUntil
	snap.ID = p.ID
	snap.SnapshotID = p.SnapshotID
	return snap, nil
}

func (s *Repo) snapshotToDataset(snap *snapstore.Snapshot) (*models.Dataset, error) {
	d := &models.Dataset{}
	if err := snap.Scan(d); err != nil {
		return nil, err
	}
	d.SnapshotID = snap.SnapshotID
	d.DateFrom = snap.DateFrom
	d.DateUntil = snap.DateUntil
	for _, fn := range s.config.DatasetLoaders {
		if err := fn(d); err != nil {
			return nil, err
		}
	}
	return d, nil
}

func (s *Repo) datasetToSnapshot(d *models.Dataset) (*snapstore.Snapshot, error) {
	snap := &snapstore.Snapshot{}
	data, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	snap.Data = data
	snap.DateFrom = d.DateFrom
	snap.DateUntil = d.DateUntil
	snap.ID = d.ID
	snap.SnapshotID = d.SnapshotID

	return snap, nil
}
