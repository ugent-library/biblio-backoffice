package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/publication"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/ulid"
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
		snapstore.WithIDGenerator(ulid.Generate),
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
	p := &models.Publication{}
	snap, err := s.publicationStore.GetCurrentSnapshot(id, s.opts)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, backends.ErrNotFound
		}
		return nil, err
	}
	if err := snap.Scan(p); err != nil {
		return nil, err
	}
	p.SnapshotID = snap.SnapshotID
	p.DateFrom = snap.DateFrom
	p.DateUntil = snap.DateUntil
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
		p := &models.Publication{}
		if err := snap.Scan(p); err != nil {
			return nil, err
		}
		p.SnapshotID = snap.SnapshotID
		p.DateFrom = snap.DateFrom
		p.DateUntil = snap.DateUntil
		publications = append(publications, p)
	}
	if c.Err() != nil {
		return nil, c.Err()
	}
	return publications, nil
}

func (s *Repository) publicationToSnapshot(publication *models.Publication) (*snapstore.Snapshot, error) {
	var snapshot *snapstore.Snapshot = &snapstore.Snapshot{}

	var data []byte
	var dataErr error
	data, dataErr = json.Marshal(publication)
	if dataErr != nil {
		return nil, dataErr
	}

	snapshot.Data = data
	snapshot.DateFrom = publication.DateFrom
	snapshot.DateUntil = publication.DateUntil
	snapshot.ID = publication.ID
	snapshot.SnapshotID = publication.SnapshotID

	return snapshot, nil
}

func (s *Repository) importPublication(p *models.Publication) error {
	snapshot, snapshotErr := s.publicationToSnapshot(p)
	if snapshotErr != nil {
		return snapshotErr
	}
	return s.publicationStore.ImportSnapshot(snapshot, s.opts)
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
	p.User = &models.PublicationUser{
		ID:   u.ID,
		Name: u.FullName,
	}
	snapshotID, err := s.publicationStore.AddAfter(snapshotID, p.ID, p, s.opts)
	if err != nil {
		p.DateUpdated = oldDateUpdated
		return err
	}
	p.SnapshotID = snapshotID
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
	cursor, cursorErr := s.publicationStore.Select(sql, values, s.opts)
	if cursorErr != nil {
		return cursorErr
	}
	defer cursor.Close()
	for cursor.HasNext() {
		snapshot, snapshotErr := cursor.Next()
		if snapshotErr != nil {
			return snapshotErr
		}
		publication := &models.Publication{}
		if err := snapshot.Scan(publication); err != nil {
			return err
		}
		publication.SnapshotID = snapshot.SnapshotID
		publication.DateFrom = snapshot.DateFrom
		publication.DateUntil = snapshot.DateUntil
		// TODO catch errors from fn() and pass them upstream
		if ok := fn(publication); !ok {
			break
		}
	}
	return cursor.Err()
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
		p := &models.Publication{}
		if err := snap.Scan(p); err != nil {
			return err
		}
		p.SnapshotID = snap.SnapshotID
		p.DateFrom = snap.DateFrom
		p.DateUntil = snap.DateUntil
		// TODO catch errors from fn() and pass them upstream
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
		p := &models.Publication{}
		if err := snap.Scan(p); err != nil {
			return err
		}
		p.SnapshotID = snap.SnapshotID
		p.DateFrom = snap.DateFrom
		p.DateUntil = snap.DateUntil
		// TODO catch errors from fn() and pass them upstream
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
		p := &models.Publication{}
		if err := snap.Scan(p); err != nil {
			return err
		}
		p.SnapshotID = snap.SnapshotID
		p.DateFrom = snap.DateFrom
		p.DateUntil = snap.DateUntil
		// TODO catch errors from fn() and pass them upstream
		if ok := fn(p); !ok {
			break
		}
	}
	return c.Err()
}

func (s *Repository) GetDataset(id string) (*models.Dataset, error) {
	d := &models.Dataset{}
	snap, err := s.datasetStore.GetCurrentSnapshot(id, s.opts)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, backends.ErrNotFound
		}
		return nil, err
	}
	if err := snap.Scan(d); err != nil {
		return nil, err
	}
	d.SnapshotID = snap.SnapshotID
	d.DateFrom = snap.DateFrom
	d.DateUntil = snap.DateUntil
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
		d := &models.Dataset{}
		if err := snap.Scan(d); err != nil {
			return nil, err
		}
		d.SnapshotID = snap.SnapshotID
		d.DateFrom = snap.DateFrom
		d.DateUntil = snap.DateUntil
		datasets = append(datasets, d)
	}
	if c.Err() != nil {
		return nil, c.Err()
	}
	return datasets, nil
}

func (s *Repository) datasetToSnapshot(dataset *models.Dataset) (*snapstore.Snapshot, error) {
	var snapshot *snapstore.Snapshot = &snapstore.Snapshot{}

	var data []byte
	var dataErr error
	data, dataErr = json.Marshal(dataset)
	if dataErr != nil {
		return nil, dataErr
	}

	snapshot.Data = data
	snapshot.DateFrom = dataset.DateFrom
	snapshot.DateUntil = dataset.DateUntil
	snapshot.ID = dataset.ID
	snapshot.SnapshotID = dataset.SnapshotID

	return snapshot, nil
}
func (s *Repository) importDataset(d *models.Dataset) error {
	snapshot, snapshotErr := s.datasetToSnapshot(d)
	if snapshotErr != nil {
		return snapshotErr
	}
	return s.datasetStore.ImportSnapshot(snapshot, s.opts)
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
	d.User = &models.DatasetUser{
		ID:   u.ID,
		Name: u.FullName,
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
	cursor, cursorErr := s.datasetStore.Select(sql, values, s.opts)
	if cursorErr != nil {
		return cursorErr
	}
	defer cursor.Close()
	for cursor.HasNext() {
		snapshot, snapshotErr := cursor.Next()
		if snapshotErr != nil {
			return snapshotErr
		}
		dataset := &models.Dataset{}
		if err := snapshot.Scan(dataset); err != nil {
			return err
		}
		dataset.SnapshotID = snapshot.SnapshotID
		dataset.DateFrom = snapshot.DateFrom
		dataset.DateUntil = snapshot.DateUntil
		// TODO catch errors from fn() and pass them upstream
		if ok := fn(dataset); !ok {
			break
		}
	}
	return cursor.Err()
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
		d := &models.Dataset{}
		if err := snap.Scan(d); err != nil {

			return err
		}
		d.SnapshotID = snap.SnapshotID
		d.DateFrom = snap.DateFrom
		d.DateUntil = snap.DateUntil
		// TODO catch errors from fn() and pass them upstream
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
		d := &models.Dataset{}
		if err := snap.Scan(d); err != nil {

			return err
		}
		d.SnapshotID = snap.SnapshotID
		d.DateFrom = snap.DateFrom
		d.DateUntil = snap.DateUntil
		// TODO catch errors from fn() and pass them upstream
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
		d := &models.Dataset{}
		if err := snap.Scan(d); err != nil {
			return err
		}
		d.SnapshotID = snap.SnapshotID
		d.DateFrom = snap.DateFrom
		d.DateUntil = snap.DateUntil
		// TODO catch errors from fn() and pass them upstream
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
