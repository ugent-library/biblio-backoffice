package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/publications"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
)

type Store struct {
	client           *snapstore.Client
	publicationStore *snapstore.Store
	datasetStore     *snapstore.Store
	opts             snapstore.Options
}

func New(dsn string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	client := snapstore.New(db, []string{"publications", "datasets"})

	return &Store{
		client:           client,
		publicationStore: client.Store("publications"),
		datasetStore:     client.Store("datasets"),
	}, nil
}

func (s *Store) AddPublicationListener(fn func(*models.Publication)) {
	s.publicationStore.Listen(func(snap *snapstore.Snapshot) {
		p := &models.Publication{}
		if err := snap.Scan(p); err == nil {
			p.SnapshotID = snap.SnapshotID
			fn(p)
		}
	})
}

func (s *Store) AddDatasetListener(fn func(*models.Dataset)) {
	s.datasetStore.Listen(func(snap *snapstore.Snapshot) {
		d := &models.Dataset{}
		if err := snap.Scan(d); err == nil {
			d.SnapshotID = snap.SnapshotID
			fn(d)
		}
	})
}

func (s *Store) Transaction(ctx context.Context, fn func(backends.Store) error) error {
	return s.client.Transaction(ctx, func(opts snapstore.Options) error {
		return fn(&Store{client: s.client, opts: opts})
	})
}

func (s *Store) GetPublication(id string) (*models.Publication, error) {
	p := &models.Publication{}
	snap, err := s.publicationStore.GetCurrentSnapshot(id, s.opts)
	if err != nil {
		return nil, err
	}
	if err := snap.Scan(p); err != nil {
		return nil, err
	}
	p.SnapshotID = snap.SnapshotID
	return p, nil
}

func (s *Store) GetPublications(ids []string) ([]*models.Publication, error) {
	c := s.publicationStore.GetByID(ids, s.opts)
	defer c.Close()
	var publications []*models.Publication
	for c.Next() {
		d := &models.Publication{}
		if err := c.Scan(d); err == nil {
			publications = append(publications, d)
		} else {
			return nil, err
		}
	}
	if c.Err() != nil {
		return nil, c.Err()
	}
	return publications, nil
}

func (s *Store) UpdatePublication(p *models.Publication) error {
	now := time.Now()

	if p.DateCreated == nil {
		p.DateCreated = &now
	}
	p.DateUpdated = &now

	// TODO move outside of store
	publications.DefaultPipeline.Process(p)

	if err := p.Validate(); err != nil {
		return err
	}

	// TODO this needs to be a separate update action
	if p.SnapshotID != "" {
		if err := s.publicationStore.AddAfter(p.SnapshotID, p.ID, p, s.opts); err != nil {
			return err
		}
		return nil
	}

	if err := s.publicationStore.Add(p.ID, p, s.opts); err != nil {
		return err
	}

	return nil
}

func (s *Store) EachPublication(fn func(*models.Publication) bool) error {
	c := s.publicationStore.GetAll(s.opts)
	defer c.Close()
	for c.Next() {
		p := &models.Publication{}
		if err := c.Scan(p); err == nil {
			if ok := fn(p); !ok {
				break
			}
		} else {
			return err
		}
	}
	return c.Err()
}

func (s *Store) GetDataset(id string) (*models.Dataset, error) {
	d := &models.Dataset{}
	snap, err := s.datasetStore.GetCurrentSnapshot(id, s.opts)
	if err != nil {
		return nil, err
	}
	if err := snap.Scan(d); err != nil {
		return nil, err
	}
	d.SnapshotID = snap.SnapshotID
	return d, nil
}

func (s *Store) GetDatasets(ids []string) ([]*models.Dataset, error) {
	c := s.datasetStore.GetByID(ids, s.opts)
	defer c.Close()
	var datasets []*models.Dataset
	for c.Next() {
		d := &models.Dataset{}
		if err := c.Scan(d); err == nil {
			datasets = append(datasets, d)
		} else {
			return nil, err
		}
	}
	if c.Err() != nil {
		return nil, c.Err()
	}
	return datasets, nil
}

func (s *Store) UpdateDataset(d *models.Dataset) error {
	now := time.Now()

	if d.DateCreated == nil {
		d.DateCreated = &now
	}
	d.DateUpdated = &now

	d.Vacuum()

	if err := d.Validate(); err != nil {
		return err
	}

	// TODO this needs to be a separate update action
	if d.SnapshotID != "" {
		if err := s.datasetStore.AddAfter(d.SnapshotID, d.ID, d, s.opts); err != nil {
			return err
		}
		return nil
	}

	if err := s.datasetStore.Add(d.ID, d, s.opts); err != nil {
		return err
	}

	return nil
}

func (s *Store) EachDataset(fn func(*models.Dataset) bool) error {
	c := s.datasetStore.GetAll(s.opts)
	defer c.Close()
	for c.Next() {
		d := &models.Dataset{}
		if err := c.Scan(d); err == nil {
			if ok := fn(d); !ok {
				break
			}
		} else {
			return err
		}
	}
	return c.Err()
}

func (s *Store) GetPublicationDatasets(p *models.Publication) ([]*models.Dataset, error) {
	datasetIds := make([]string, len(p.RelatedDataset))
	for _, rd := range p.RelatedDataset {
		datasetIds = append(datasetIds, rd.ID)
	}
	return s.GetDatasets(datasetIds)
}

func (s *Store) GetDatasetPublications(d *models.Dataset) ([]*models.Publication, error) {
	publicationIds := make([]string, len(d.RelatedPublication))
	for _, rp := range d.RelatedPublication {
		publicationIds = append(publicationIds, rp.ID)
	}
	return s.GetPublications(publicationIds)
}

func (s *Store) AddPublicationDataset(p *models.Publication, d *models.Dataset) error {
	return s.Transaction(context.Background(), func(s backends.Store) error {
		if !p.HasRelatedDataset(d.ID) {
			p.RelatedDataset = append(p.RelatedDataset, models.RelatedDataset{ID: d.ID})
			if err := s.UpdatePublication(p); err != nil {
				return err
			}
		}
		if !d.HasRelatedPublication(p.ID) {
			d.RelatedPublication = append(d.RelatedPublication, models.RelatedPublication{ID: p.ID})
			if err := s.UpdateDataset(d); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Store) RemovePublicationDataset(p *models.Publication, d *models.Dataset) error {
	return s.Transaction(context.Background(), func(s backends.Store) error {
		if p.HasRelatedDataset(d.ID) {
			p.RemoveRelatedDataset(d.ID)
			if err := s.UpdatePublication(p); err != nil {
				return err
			}
		}
		if d.HasRelatedPublication(p.ID) {
			d.RemoveRelatedPublication(p.ID)
			if err := s.UpdateDataset(d); err != nil {
				return err
			}
		}

		return nil
	})
}
