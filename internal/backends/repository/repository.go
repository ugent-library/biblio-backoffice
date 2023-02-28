package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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
	client           *snapstore.Client
	publicationStore *snapstore.Store
	datasetStore     *snapstore.Store
	opts             snapstore.Options
}

func New(dsn string) (*Repository, error) {
	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	client := snapstore.New(db, []string{"publications", "datasets"},
		snapstore.WithIDGenerator(func() (string, error) {
			return ulid.Make().String(), nil
		}),
	)

	return &Repository{
		client:           client,
		publicationStore: client.Store("publications"),
		datasetStore:     client.Store("datasets"),
	}, nil
}

func (s *Repository) AddPublicationListener(fn func(*models.Publication)) {
	s.publicationStore.Listen(func(snap *snapstore.Snapshot) {
		p := &models.Publication{}
		if err := snap.Scan(p); err == nil {
			p.SnapshotID = snap.SnapshotID
			p.DateFrom = snap.DateFrom
			p.DateUntil = snap.DateUntil
			fn(p)
		}
	})
}

func (s *Repository) AddDatasetListener(fn func(*models.Dataset)) {
	s.datasetStore.Listen(func(snap *snapstore.Snapshot) {
		d := &models.Dataset{}
		if err := snap.Scan(d); err == nil {
			d.SnapshotID = snap.SnapshotID
			d.DateFrom = snap.DateFrom
			d.DateUntil = snap.DateUntil
			fn(d)
		}
	})
}

func (s *Repository) Transaction(ctx context.Context, fn func(backends.Repository) error) error {
	return s.client.Transaction(ctx, func(opts snapstore.Options) error {
		return fn(&Repository{
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
	p, err := snapshotToPublication(snap)
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
		p, err := snapshotToPublication(snap)
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

func (s *Repository) importPublication(p *models.Publication) error {
	snap, err := publicationToSnapshot(p)
	if err != nil {
		return err
	}
	return s.publicationStore.ImportSnapshot(snap, s.opts)
}

func (s *Repository) ImportCurrentPublication(p *models.Publication) error {
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
	return s.importPublication(p)
}

func (s *Repository) ImportOldPublication(p *models.Publication) error {
	if p.DateCreated == nil {
		return fmt.Errorf("unable to import old publication %s: date_created is not set", p.ID)
	}
	if p.DateUpdated == nil {
		return fmt.Errorf("unable to import old publication %s: date_updated is not set", p.ID)
	}
	if p.DateFrom == nil {
		return fmt.Errorf("unable to import old publication %s: date_from is not set", p.ID)
	}
	if p.DateUntil == nil {
		return fmt.Errorf("unable to import old publication %s: date_until is not set", p.ID)
	}
	return s.importPublication(p)
}

func (s *Repository) SavePublication(p *models.Publication, u *models.User) error {
	now := time.Now()

	if p.DateCreated == nil {
		p.DateCreated = &now
	}
	p.DateUpdated = &now

	if u != nil {
		p.User = &models.PublicationUser{
			ID:   u.ID,
			Name: u.FullName,
		}
		p.LastUser = p.User
	} else {
		p.User = nil
	}

	// TODO move outside of store
	p = publication.DefaultPipeline.Process(p)

	if err := p.Validate(); err != nil {
		return err
	}

	return s.publicationStore.Add(p.ID, p, s.opts)
}

func (s *Repository) UpdatePublication(snapshotID string, p *models.Publication, u *models.User) error {
	// TODO move outside of store
	p = publication.DefaultPipeline.Process(p)
	oldDateUpdated := p.DateUpdated
	now := time.Now()
	p.DateUpdated = &now

	if u != nil {
		p.User = &models.PublicationUser{
			ID:   u.ID,
			Name: u.FullName,
		}
		p.LastUser = p.User
	} else {
		p.User = nil
	}

	snapshotID, err := s.publicationStore.AddAfter(snapshotID, p.ID, p, s.opts)
	if err != nil {
		p.DateUpdated = oldDateUpdated
		return err
	}
	p.SnapshotID = snapshotID
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

	return nil
}

func (s *Repository) CountPublications(args *backends.RepositoryQueryArgs) (int, error) {
	sql := "SELECT * FROM publications WHERE date_until is null"
	values := make([]any, 0)
	// TODO: how to safely quote field names?
	for _, filter := range args.Filters {
		values = append(values, filter.Value)
		sql += " AND " + filter.Field + " " + filter.Op + " $" + strconv.Itoa(len(values))
	}
	return s.publicationStore.CountSql(sql, values, s.opts)
}

func (s *Repository) CountDatasets(args *backends.RepositoryQueryArgs) (int, error) {
	sql := "SELECT * FROM datasets WHERE date_until is null"
	values := make([]any, 0)
	// TODO: how to safely quote field names?
	for _, filter := range args.Filters {
		values = append(values, filter.Value)
		sql += " AND " + filter.Field + " " + filter.Op + " $" + strconv.Itoa(len(values))
	}
	return s.datasetStore.CountSql(sql, values, s.opts)
}

func (s *Repository) SearchPublications(args *backends.RepositoryQueryArgs) ([]*models.Publication, error) {
	sql := "SELECT snapshot_id, id, data, date_from, date_until FROM publications WHERE date_until IS NULL"
	values := make([]any, 0)
	if args.Offset < 0 {
		args.Offset = 0
	}
	if args.Limit < 0 {
		args.Limit = 20
	}

	// TODO: how to safely quote field names?
	for _, filter := range args.Filters {
		values = append(values, filter.Value)
		sql += " AND " + filter.Field + " " + filter.Op + " $" + strconv.Itoa(len(values))
	}

	if args.Order != "" {
		sql += " ORDER BY " + args.Order
	}
	values = append(values, args.Limit)
	sql += " LIMIT $" + strconv.Itoa(len(values))
	values = append(values, args.Offset)
	sql += " OFFSET $" + strconv.Itoa(len(values))

	publications := make([]*models.Publication, 0, args.Limit)

	err := s.SelectPublications(sql, values, func(publication *models.Publication) bool {
		publications = append(publications, publication)
		return true
	})

	if err != nil {
		return nil, err
	}

	return publications, nil
}

func (s *Repository) SearchDatasets(args *backends.RepositoryQueryArgs) ([]*models.Dataset, error) {
	sql := "SELECT snapshot_id, id, data, date_from, date_until FROM datasets WHERE date_until IS NULL"
	values := make([]any, 0)
	if args.Offset < 0 {
		args.Offset = 0
	}
	if args.Limit < 0 {
		args.Limit = 20
	}

	// TODO: how to safely quote field names?
	for _, filter := range args.Filters {
		values = append(values, filter.Value)
		sql += " AND " + filter.Field + " " + filter.Op + " $" + strconv.Itoa(len(values))
	}

	if args.Order != "" {
		sql += " ORDER BY " + args.Order
	}
	values = append(values, args.Limit)
	sql += " LIMIT $" + strconv.Itoa(len(values))
	values = append(values, args.Offset)
	sql += " OFFSET $" + strconv.Itoa(len(values))
	datasets := make([]*models.Dataset, 0, args.Limit)

	err := s.SelectDatasets(sql, values, func(dataset *models.Dataset) bool {
		datasets = append(datasets, dataset)
		return true
	})

	if err != nil {
		return nil, err
	}

	return datasets, nil
}

func (s *Repository) SelectPublications(sql string, values []any, fn func(*models.Publication) bool) error {
	c, err := s.publicationStore.Select(sql, values, s.opts)
	if err != nil {
		return err
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return err
		}
		p, err := snapshotToPublication(snap)
		if err != nil {
			return err
		}
		if ok := fn(p); !ok {
			break
		}
	}
	return c.Err()
}

func (s *Repository) PublicationsBetween(t1, t2 time.Time, fn func(*models.Publication) bool) error {
	c, err := s.publicationStore.Select(
		"SELECT * FROM publications WHERE date_until IS NULL AND date_from >= $1 AND date_from <= $2",
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
		p, err := snapshotToPublication(snap)
		if err != nil {
			return err
		}
		if ok := fn(p); !ok {
			break
		}
	}
	return c.Err()
}

func (s *Repository) EachPublication(ctx context.Context, fn func(*models.Publication) error) error {
	c, err := s.publicationStore.GetAll(s.opts)
	if err != nil {
		return err
	}
	defer c.Close()

	for c.HasNext() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			snap, err := c.Next()
			if err != nil {
				return err
			}
			p, err := snapshotToPublication(snap)
			if err != nil {
				return err
			}
			if err := fn(p); err != nil {
				return err
			}
		}
	}

	return c.Err()
}

func (s *Repository) EachPublicationSnapshot(ctx context.Context, fn func(*models.Publication) bool) error {
	c, err := s.publicationStore.GetAllSnapshots(s.opts)
	if err != nil {
		return err
	}
	defer c.Close()

loop:
	for c.HasNext() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			snap, err := c.Next()
			if err != nil {
				return err
			}
			p, err := snapshotToPublication(snap)
			if err != nil {
				return err
			}
			if ok := fn(p); !ok {
				break loop
			}
		}
	}
	return c.Err()
}

func (s *Repository) PublicationHistory(ctx context.Context, id string, fn func(*models.Publication) bool) error {
	c, err := s.publicationStore.GetHistory(id, s.opts)
	if err != nil {
		return err
	}
	defer c.Close()

loop:
	for c.HasNext() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			snap, err := c.Next()
			if err != nil {
				return err
			}
			p, err := snapshotToPublication(snap)
			if err != nil {
				return err
			}
			if ok := fn(p); !ok {
				break loop
			}
		}
	}
	return c.Err()
}

func (s *Repository) GetDataset(id string) (*models.Dataset, error) {
	snap, err := s.datasetStore.GetCurrentSnapshot(id, s.opts)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, backends.ErrNotFound
		}
		return nil, err
	}
	d, err := snapshotToDataset(snap)
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
		d, err := snapshotToDataset(snap)
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

func (s *Repository) importDataset(d *models.Dataset) error {
	snap, err := datasetToSnapshot(d)
	if err != nil {
		return err
	}
	return s.datasetStore.ImportSnapshot(snap, s.opts)
}

func (s *Repository) ImportCurrentDataset(d *models.Dataset) error {
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
	return s.importDataset(d)
}

func (s *Repository) ImportOldDataset(d *models.Dataset) error {
	if d.DateCreated == nil {
		return fmt.Errorf("unable to import dataset %s: date_created is not set", d.ID)
	}
	if d.DateUpdated == nil {
		return fmt.Errorf("unable to import dataset %s: date_updated is not set", d.ID)
	}
	if d.DateFrom == nil {
		return fmt.Errorf("unable to import dataset %s: date_from is not set", d.ID)
	}
	if d.DateUntil == nil {
		return fmt.Errorf("unable to import dataset %s: date_until is not set", d.ID)
	}
	return s.importDataset(d)
}

func (s *Repository) SaveDataset(d *models.Dataset, u *models.User) error {
	now := time.Now()

	if d.DateCreated == nil {
		d.DateCreated = &now
	}
	d.DateUpdated = &now

	if u != nil {
		d.User = &models.DatasetUser{
			ID:   u.ID,
			Name: u.FullName,
		}
		d.LastUser = d.User
	} else {
		d.User = nil
	}

	if err := d.Validate(); err != nil {
		return err
	}

	//TODO: move outside
	if d.Status == "public" && !d.HasBeenPublic {
		d.HasBeenPublic = true
	}

	return s.datasetStore.Add(d.ID, d, s.opts)
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
		d.User = &models.DatasetUser{
			ID:   u.ID,
			Name: u.FullName,
		}
		d.LastUser = d.User
	} else {
		d.User = nil
	}

	snapshotID, err := s.datasetStore.AddAfter(snapshotID, d.ID, d, s.opts)
	if err != nil {
		d.DateUpdated = oldDateUpdated
		return err
	}
	d.SnapshotID = snapshotID
	return nil
}

func (s *Repository) SelectDatasets(sql string, values []any, fn func(*models.Dataset) bool) error {
	c, err := s.datasetStore.Select(sql, values, s.opts)
	if err != nil {
		return err
	}
	defer c.Close()
	for c.HasNext() {
		snap, err := c.Next()
		if err != nil {
			return err
		}
		d, err := snapshotToDataset(snap)
		if err != nil {
			return err
		}
		if ok := fn(d); !ok {
			break
		}
	}
	return c.Err()
}

func (s *Repository) DatasetsBetween(t1, t2 time.Time, fn func(*models.Dataset) bool) error {
	c, err := s.datasetStore.Select(
		"SELECT * FROM datasets WHERE date_until IS NULL AND date_from >= $1 AND date_from <= $2",
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
		p, err := snapshotToDataset(snap)
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
		d, err := snapshotToDataset(snap)
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
		d, err := snapshotToDataset(snap)
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
		d, err := snapshotToDataset(snap)
		if err != nil {
			return err
		}
		if ok := fn(d); !ok {
			break
		}
	}
	return c.Err()
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
		if u.CanDeletePublication(publication) {
			filteredPublications = append(filteredPublications, publication)
		}
	}
	return filteredPublications, nil
}

func (s *Repository) AddPublicationDataset(p *models.Publication, d *models.Dataset, u *models.User) error {
	return s.Transaction(context.Background(), func(s backends.Repository) error {
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
	return s.Transaction(context.Background(), func(s backends.Repository) error {
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

func snapshotToPublication(snap *snapstore.Snapshot) (*models.Publication, error) {
	p := &models.Publication{}
	if err := snap.Scan(p); err != nil {
		return nil, err
	}
	p.SnapshotID = snap.SnapshotID
	p.DateFrom = snap.DateFrom
	p.DateUntil = snap.DateUntil
	return p, nil
}

func publicationToSnapshot(p *models.Publication) (*snapstore.Snapshot, error) {
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

func snapshotToDataset(snap *snapstore.Snapshot) (*models.Dataset, error) {
	d := &models.Dataset{}
	if err := snap.Scan(d); err != nil {
		return nil, err
	}
	d.SnapshotID = snap.SnapshotID
	d.DateFrom = snap.DateFrom
	d.DateUntil = snap.DateUntil
	return d, nil
}

func datasetToSnapshot(d *models.Dataset) (*snapstore.Snapshot, error) {
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
