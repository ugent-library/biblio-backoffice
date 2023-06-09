// TODO all mutating methods should call Validate() before saving
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/publication"
	"github.com/ugent-library/biblio-backoffice/internal/snapstore"
)

type Repository struct {
	config           Config
	client           *snapstore.Client
	publicationStore *snapstore.Store
	datasetStore     *snapstore.Store
	opts             snapstore.Options
}

type Config struct {
	DSN                  string
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

func New(c Config) (*Repository, error) {
	db, err := pgxpool.Connect(context.Background(), c.DSN)
	if err != nil {
		return nil, err
	}

	client := snapstore.New(db, []string{"publications", "datasets"},
		snapstore.WithIDGenerator(func() (string, error) {
			return ulid.Make().String(), nil
		}),
	)

	return &Repository{
		config:           c,
		client:           client,
		publicationStore: client.Store("publications"),
		datasetStore:     client.Store("datasets"),
	}, nil
}

func (s *Repository) publicationNotify(p *models.Publication) {
	for _, fn := range s.config.PublicationListeners {
		fn(p)
	}
}

func (s *Repository) datasetNotify(d *models.Dataset) {
	for _, fn := range s.config.DatasetListeners {
		fn(d)
	}
}

func (s *Repository) tx(ctx context.Context, fn func(backends.Repository) error) error {
	return s.client.Tx(ctx, func(opts snapstore.Options) error {
		return fn(&Repository{
			config:           s.config,
			client:           s.client,
			publicationStore: s.publicationStore,
			datasetStore:     s.datasetStore,
			opts:             opts,
		})
	})
}

func (s *Repository) GetPublication(id string) (*models.Publication, error) {
	snap, err := s.publicationStore.GetCurrentSnapshot(id, s.opts)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, backends.ErrNotFound
		}
		return nil, err
	}
	p, err := s.snapshotToPublication(snap)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Repository) GetPublications(ids []string) ([]*models.Publication, error) {
	c, err := s.publicationStore.GetByID(ids, s.opts)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	var publications []*models.Publication
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return nil, err
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return nil, err
		}
		publications = append(publications, p)
	}
	if c.Err() != nil {
		return nil, c.Err()
	}
	return publications, nil
}

func (s *Repository) ImportPublication(p *models.Publication) error {
	if p.DateCreated == nil {
		return fmt.Errorf("unable to import old publication %s: date_created is not set", p.ID)
	}
	if p.DateUpdated == nil {
		return fmt.Errorf("unable to import old publication %s: date_updated is not set", p.ID)
	}
	if p.DateFrom == nil {
		return fmt.Errorf("unable to import old publication %s: date_from is not set", p.ID)
	}
	if p.DateUntil != nil {
		return fmt.Errorf("unable to import old publication %s: date_until should be nil", p.ID)
	}

	snap, err := s.publicationToSnapshot(p)
	if err != nil {
		return err
	}

	if err := s.publicationStore.ImportSnapshot(snap, s.opts); err != nil {
		return err
	}

	for _, fn := range s.config.PublicationLoaders {
		if err := fn(p); err != nil {
			return err
		}
	}

	s.publicationNotify(p)

	return nil
}

func (s *Repository) SavePublication(p *models.Publication, u *models.User) error {
	now := time.Now()

	if p.DateCreated == nil {
		p.DateCreated = &now
	}
	p.DateUpdated = &now

	if u != nil {
		p.UserID = u.Person.ID
		p.User = &u.Person
		p.LastUserID = u.Person.ID
		p.LastUser = &u.Person
	} else {
		p.UserID = ""
		p.User = nil
	}

	// TODO move outside of store?
	p = publication.DefaultPipeline.Process(p)

	if err := p.Validate(); err != nil {
		return err
	}

	if err := s.publicationStore.Add(p.ID, p, s.opts); err != nil {
		return err
	}

	s.publicationNotify(p)

	return nil
}

func (s *Repository) UpdatePublication(snapshotID string, p *models.Publication, u *models.User) error {
	// TODO move outside of store
	p = publication.DefaultPipeline.Process(p)
	oldDateUpdated := p.DateUpdated
	now := time.Now()
	p.DateUpdated = &now

	if u != nil {
		p.UserID = u.Person.ID
		p.User = &u.Person
		p.LastUserID = u.Person.ID
		p.LastUser = &u.Person
	} else {
		p.UserID = ""
		p.User = nil
	}

	snapshotID, err := s.publicationStore.AddAfter(snapshotID, p.ID, p, s.opts)
	if err != nil {
		p.DateUpdated = oldDateUpdated
		return err
	}
	p.SnapshotID = snapshotID

	s.publicationNotify(p)

	return nil
}

func (s *Repository) UpdatePublicationInPlace(p *models.Publication) error {
	snap, err := s.publicationStore.Update(p.SnapshotID, p.ID, p, s.opts)
	if err != nil {
		return err
	}

	np := &models.Publication{}
	if err := snap.Scan(np); err != nil {
		return err
	}

	s.publicationNotify(p)

	return nil
}

func (s *Repository) MutatePublication(id string, u *models.User, muts ...backends.Mutation) error {
	if len(muts) == 0 {
		return nil
	}

	p, err := s.GetPublication(id)
	if err != nil {
		return err
	}

	for _, mut := range muts {
		mutator, ok := s.config.PublicationMutators[mut.Op]
		if !ok {
			return fmt.Errorf("unknown mutation '%s'", mut.Op)
		}
		if err := mutator(p, mut.Args); err != nil {
			return err
		}
	}

	if err = p.Validate(); err != nil {
		return err
	}

	return s.UpdatePublication(p.SnapshotID, p, u)
}

func (s *Repository) PublicationsAfter(t time.Time, limit, offset int) (int, []*models.Publication, error) {
	n, err := s.publicationStore.CountSql(
		"SELECT * FROM publications WHERE date_until IS NULL AND date_from >= $1",
		[]any{t},
		s.opts,
	)
	if err != nil {
		return 0, nil, err
	}

	c, err := s.publicationStore.Select(
		"SELECT * FROM publications WHERE date_until IS NULL AND date_from >= $1 ORDER BY date_from ASC LIMIT $2 OFFSET $3",
		[]any{t, limit, offset},
		s.opts,
	)
	if err != nil {
		return 0, nil, err
	}

	publications := make([]*models.Publication, 0, limit)

	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return 0, nil, err
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return 0, nil, err
		}

		publications = append(publications, p)
	}

	if c.Err() != nil {
		return 0, nil, err
	}

	return n, publications, nil
}

func (s *Repository) PublicationsBetween(t1, t2 time.Time, fn func(*models.Publication) bool) error {
	c, err := s.publicationStore.Select(
		"SELECT * FROM publications WHERE date_until IS NULL AND date_from >= $1 AND date_from <= $2 ORDER BY date_from ASC",
		[]any{t1, t2},
		s.opts,
	)
	if err != nil {
		return err
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return err
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return err
		}
		if ok := fn(p); !ok {
			break
		}
	}
	return c.Err()
}

func (s *Repository) EachPublication(fn func(*models.Publication) bool) error {
	c, err := s.publicationStore.GetAll(s.opts)
	if err != nil {
		return err
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return err
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return err
		}
		if ok := fn(p); !ok {
			break
		}
	}
	return c.Err()
}

func (s *Repository) EachPublicationSnapshot(fn func(*models.Publication) bool) error {
	c, err := s.publicationStore.GetAllSnapshots(s.opts)
	if err != nil {
		return err
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return err
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return err
		}
		if ok := fn(p); !ok {
			break
		}
	}
	return c.Err()
}

// TODO add handle with a listener, then this method isn't needed anymore
func (s *Repository) EachPublicationWithoutHandle(fn func(*models.Publication) bool) error {
	sql := `
		SELECT * FROM publications WHERE date_until IS NULL AND
		data->>'status' = 'public' AND
		NOT data ? 'handle'
		`
	c, err := s.publicationStore.Select(sql, nil, s.opts)
	if err != nil {
		return err
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return err
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return err
		}
		if ok := fn(p); !ok {
			break
		}
	}
	return c.Err()
}

func (s *Repository) PublicationHistory(id string, fn func(*models.Publication) bool) error {
	c, err := s.publicationStore.GetHistory(id, s.opts)
	if err != nil {
		return err
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return err
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return err
		}
		if ok := fn(p); !ok {
			break
		}
	}
	return c.Err()
}

func (s *Repository) UpdatePublicationEmbargoes() (int, error) {
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
		return n, err
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return n, err
		}
		p, err := s.snapshotToPublication(snap)
		if err != nil {
			return n, err
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
			return n, err
		}

		n++
	}

	return n, c.Err()
}

func (s *Repository) GetDataset(id string) (*models.Dataset, error) {
	snap, err := s.datasetStore.GetCurrentSnapshot(id, s.opts)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, backends.ErrNotFound
		}
		return nil, err
	}
	d, err := s.snapshotToDataset(snap)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Repository) GetDatasets(ids []string) ([]*models.Dataset, error) {
	c, err := s.datasetStore.GetByID(ids, s.opts)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	var datasets []*models.Dataset
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return nil, err
		}
		d, err := s.snapshotToDataset(snap)
		if err != nil {
			return nil, err
		}
		datasets = append(datasets, d)
	}
	if c.Err() != nil {
		return nil, c.Err()
	}
	return datasets, nil
}

func (s *Repository) ImportDataset(d *models.Dataset) error {
	if d.DateCreated == nil {
		return fmt.Errorf("unable to import dataset %s: date_created is not set", d.ID)
	}
	if d.DateUpdated == nil {
		return fmt.Errorf("unable to import dataset %s: date_updated is not set", d.ID)
	}
	if d.DateFrom == nil {
		return fmt.Errorf("unable to import dataset %s: date_from is not set", d.ID)
	}
	if d.DateUntil != nil {
		return fmt.Errorf("unable to import dataset %s: date_until should be nil", d.ID)
	}

	snap, err := s.datasetToSnapshot(d)
	if err != nil {
		return err
	}

	if err := s.datasetStore.ImportSnapshot(snap, s.opts); err != nil {
		return err
	}

	for _, fn := range s.config.DatasetLoaders {
		if err := fn(d); err != nil {
			return err
		}
	}

	s.datasetNotify(d)

	return nil
}

func (s *Repository) SaveDataset(d *models.Dataset, u *models.User) error {
	now := time.Now()

	if d.DateCreated == nil {
		d.DateCreated = &now
	}
	d.DateUpdated = &now

	if u != nil {
		d.UserID = u.Person.ID
		d.User = &u.Person
		d.LastUserID = u.Person.ID
		d.LastUser = &u.Person
	} else {
		d.UserID = ""
		d.User = nil
	}

	if err := d.Validate(); err != nil {
		return err
	}

	//TODO: move outside
	if d.Status == "public" && !d.HasBeenPublic {
		d.HasBeenPublic = true
	}

	if err := s.datasetStore.Add(d.ID, d, s.opts); err != nil {
		return err
	}

	s.datasetNotify(d)

	return nil
}

func (s *Repository) UpdateDataset(snapshotID string, d *models.Dataset, u *models.User) error {
	//TODO: move outside
	if d.Status == "public" && !d.HasBeenPublic {
		d.HasBeenPublic = true
	}
	oldDateUpdated := d.DateUpdated
	now := time.Now()
	d.DateUpdated = &now

	if u != nil {
		d.UserID = u.Person.ID
		d.User = &u.Person
		d.LastUserID = u.Person.ID
		d.LastUser = &u.Person
	} else {
		d.UserID = ""
		d.User = nil
	}

	snapshotID, err := s.datasetStore.AddAfter(snapshotID, d.ID, d, s.opts)
	if err != nil {
		d.DateUpdated = oldDateUpdated
		return err
	}
	d.SnapshotID = snapshotID

	s.datasetNotify(d)

	return nil
}

func (s *Repository) MutateDataset(id string, u *models.User, muts ...backends.Mutation) error {
	if len(muts) == 0 {
		return nil
	}

	d, err := s.GetDataset(id)
	if err != nil {
		return err
	}

	for _, mut := range muts {
		mutator, ok := s.config.DatasetMutators[mut.Op]
		if !ok {
			return fmt.Errorf("unknown mutation '%s'", mut.Op)
		}
		if err := mutator(d, mut.Args); err != nil {
			return err
		}
	}

	if err = d.Validate(); err != nil {
		return err
	}

	return s.UpdateDataset(d.SnapshotID, d, u)
}

func (s *Repository) DatasetsAfter(t time.Time, limit, offset int) (int, []*models.Dataset, error) {
	n, err := s.datasetStore.CountSql(
		"SELECT * FROM datasets WHERE date_until IS NULL AND date_from >= $1",
		[]any{t},
		s.opts,
	)
	if err != nil {
		return 0, nil, err
	}

	c, err := s.datasetStore.Select(
		"SELECT * FROM datasets WHERE date_until IS NULL AND date_from >= $1 ORDER BY date_from ASC LIMIT $2 OFFSET $3",
		[]any{t, limit, offset},
		s.opts,
	)
	if err != nil {
		return 0, nil, err
	}

	datasets := make([]*models.Dataset, 0, limit)

	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return 0, nil, err
		}
		p, err := s.snapshotToDataset(snap)
		if err != nil {
			return 0, nil, err
		}

		datasets = append(datasets, p)
	}

	if c.Err() != nil {
		return 0, nil, err
	}

	return n, datasets, nil
}

func (s *Repository) DatasetsBetween(t1, t2 time.Time, fn func(*models.Dataset) bool) error {
	c, err := s.datasetStore.Select(
		"SELECT * FROM datasets WHERE date_until IS NULL AND date_from >= $1 AND date_from <= $2 ORDER BY date_from ASC",
		[]any{t1, t2},
		s.opts,
	)
	if err != nil {
		return err
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return err
		}
		p, err := s.snapshotToDataset(snap)
		if err != nil {
			return err
		}
		if ok := fn(p); !ok {
			break
		}
	}
	return c.Err()
}

func (s *Repository) EachDataset(fn func(*models.Dataset) bool) error {
	c, err := s.datasetStore.GetAll(s.opts)
	if err != nil {
		return err
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return err
		}
		d, err := s.snapshotToDataset(snap)
		if err != nil {
			return err
		}
		if ok := fn(d); !ok {
			break
		}
	}
	return c.Err()
}

func (s *Repository) EachDatasetSnapshot(fn func(*models.Dataset) bool) error {
	c, err := s.datasetStore.GetAllSnapshots(s.opts)
	if err != nil {
		return err
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return err
		}
		d, err := s.snapshotToDataset(snap)
		if err != nil {
			return err
		}
		if ok := fn(d); !ok {
			break
		}
	}
	return c.Err()
}

func (s *Repository) EachDatasetWithoutHandle(fn func(*models.Dataset) bool) error {
	sql := `
		SELECT * FROM datasets WHERE date_until IS NULL AND
		data->>'status' = 'public' AND
		NOT data ? 'handle'
		`
	c, err := s.datasetStore.Select(sql, nil, s.opts)
	if err != nil {
		return err
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return err
		}
		d, err := s.snapshotToDataset(snap)
		if err != nil {
			return err
		}
		if ok := fn(d); !ok {
			break
		}
	}
	return c.Err()
}

func (s *Repository) DatasetHistory(id string, fn func(*models.Dataset) bool) error {
	c, err := s.publicationStore.GetHistory(id, s.opts)
	if err != nil {
		return err
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return err
		}
		d, err := s.snapshotToDataset(snap)
		if err != nil {
			return err
		}
		if ok := fn(d); !ok {
			break
		}
	}
	return c.Err()
}

func (s *Repository) UpdateDatasetEmbargoes() (int, error) {
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
		return n, err
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return n, err
		}
		d, err := s.snapshotToDataset(snap)
		if err != nil {
			return n, err
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
			return n, err
		}

		n++
	}

	return n, c.Err()
}

func (s *Repository) GetPublicationDatasets(p *models.Publication) ([]*models.Dataset, error) {
	datasetIds := make([]string, len(p.RelatedDataset))
	for _, rd := range p.RelatedDataset {
		datasetIds = append(datasetIds, rd.ID)
	}
	return s.GetDatasets(datasetIds)
}

func (s *Repository) GetVisiblePublicationDatasets(u *models.User, p *models.Publication) ([]*models.Dataset, error) {
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

func (s *Repository) GetDatasetPublications(d *models.Dataset) ([]*models.Publication, error) {
	publicationIds := make([]string, len(d.RelatedPublication))
	for _, rp := range d.RelatedPublication {
		publicationIds = append(publicationIds, rp.ID)
	}
	return s.GetPublications(publicationIds)
}

func (s *Repository) GetVisibleDatasetPublications(u *models.User, d *models.Dataset) ([]*models.Publication, error) {
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

func (s *Repository) AddPublicationDataset(p *models.Publication, d *models.Dataset, u *models.User) error {
	return s.tx(context.Background(), func(s backends.Repository) error {
		if !p.HasRelatedDataset(d.ID) {
			p.RelatedDataset = append(p.RelatedDataset, models.RelatedDataset{ID: d.ID})
			if err := s.SavePublication(p, u); err != nil {
				return err
			}
		}
		if !d.HasRelatedPublication(p.ID) {
			d.RelatedPublication = append(d.RelatedPublication, models.RelatedPublication{ID: p.ID})
			if err := s.SaveDataset(d, u); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Repository) RemovePublicationDataset(p *models.Publication, d *models.Dataset, u *models.User) error {
	return s.tx(context.Background(), func(s backends.Repository) error {
		if p.HasRelatedDataset(d.ID) {
			p.RemoveRelatedDataset(d.ID)
			if err := s.SavePublication(p, u); err != nil {
				return err
			}
		}
		if d.HasRelatedPublication(p.ID) {
			d.RemoveRelatedPublication(p.ID)
			if err := s.SaveDataset(d, u); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Repository) PurgeAllPublications() error {
	return s.publicationStore.PurgeAll(s.opts)
}

func (s *Repository) PurgePublication(id string) error {
	return s.publicationStore.Purge(id, s.opts)
}

func (s *Repository) PurgeAllDatasets() error {
	return s.datasetStore.PurgeAll(s.opts)
}

func (s *Repository) PurgeDataset(id string) error {
	return s.datasetStore.Purge(id, s.opts)
}

func (s *Repository) snapshotToPublication(snap *snapstore.Snapshot) (*models.Publication, error) {
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

func (s *Repository) publicationToSnapshot(p *models.Publication) (*snapstore.Snapshot, error) {
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

func (s *Repository) snapshotToDataset(snap *snapstore.Snapshot) (*models.Dataset, error) {
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

func (s *Repository) datasetToSnapshot(d *models.Dataset) (*snapstore.Snapshot, error) {
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
