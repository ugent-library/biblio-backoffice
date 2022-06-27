package repository

import (
	"context"
	"time"

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

func (s *Repository) SavePublication(p *models.Publication) error {
	now := time.Now()

	if p.DateCreated == nil {
		p.DateCreated = &now
	}
	p.DateUpdated = &now

	// TODO move outside of store
	p = publication.DefaultPipeline.Process(p)

	if err := p.Validate(); err != nil {
		return err
	}

	return s.publicationStore.Add(p.ID, p, s.opts)
}

func (s *Repository) UpdatePublication(snapshotID string, d *models.Publication) error {
	now := time.Now()
	d.DateUpdated = &now
	return s.publicationStore.AddAfter(snapshotID, d.ID, d, s.opts)
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

func (s *Repository) SaveDataset(d *models.Dataset) error {
	now := time.Now()

	if d.DateCreated == nil {
		d.DateCreated = &now
	}
	d.DateUpdated = &now

	if err := d.Validate(); err != nil {
		return err
	}

	return s.datasetStore.Add(d.ID, d, s.opts)
}

func (s *Repository) UpdateDataset(snapshotID string, d *models.Dataset) error {
	now := time.Now()
	d.DateUpdated = &now
	return s.datasetStore.AddAfter(snapshotID, d.ID, d, s.opts)
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

func (s *Repository) GetDatasetPublications(d *models.Dataset) ([]*models.Publication, error) {
	publicationIds := make([]string, len(d.RelatedPublication))
	for _, rp := range d.RelatedPublication {
		publicationIds = append(publicationIds, rp.ID)
	}
	return s.GetPublications(publicationIds)
}

func (s *Repository) AddPublicationDataset(p *models.Publication, d *models.Dataset) error {
	return s.Transaction(context.Background(), func(s backends.Repository) error {
		if !p.HasRelatedDataset(d.ID) {
			p.RelatedDataset = append(p.RelatedDataset, models.RelatedDataset{ID: d.ID})
			if err := s.SavePublication(p); err != nil {
				return err
			}
		}
		if !d.HasRelatedPublication(p.ID) {
			d.RelatedPublication = append(d.RelatedPublication, models.RelatedPublication{ID: p.ID})
			if err := s.SaveDataset(d); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Repository) RemovePublicationDataset(p *models.Publication, d *models.Dataset) error {
	return s.Transaction(context.Background(), func(s backends.Repository) error {
		if p.HasRelatedDataset(d.ID) {
			p.RemoveRelatedDataset(d.ID)
			if err := s.SavePublication(p); err != nil {
				return err
			}
		}
		if d.HasRelatedPublication(p.ID) {
			d.RemoveRelatedPublication(p.ID)
			if err := s.SaveDataset(d); err != nil {
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
