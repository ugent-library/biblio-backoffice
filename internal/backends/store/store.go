package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
)

type Store struct {
	client *snapstore.Client
	opts   snapstore.Options
}

func New(dsn string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &Store{client: snapstore.New(db, []string{"publications", "datasets"})}, nil
}

func (s *Store) Transaction(ctx context.Context, fn func(backends.Store) error) error {
	return s.client.Transaction(ctx, func(opts snapstore.Options) error {
		return fn(&Store{client: s.client, opts: opts})
	})
}

func (s *Store) GetPublication(id string) (*models.Publication, error) {
	p := &models.Publication{}
	snap, err := s.client.Store("publications").Get(id, s.opts)
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
	c := s.client.Store("publications").GetByID(ids, s.opts)
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

func (s *Store) StorePublication(p *models.Publication) error {
	now := time.Now()

	if p.DateCreated == nil {
		p.DateCreated = &now
	}
	p.DateUpdated = &now

	if p.ID == "" {
		p.ID = uuid.NewString()
	}

	store := s.client.Store("publications")

	// TODO this needs to be a separate update action
	if p.SnapshotID != "" {
		if err := store.UpdateSnapshot(p.SnapshotID, p, s.opts); err != nil {
			return err
		}
		return nil
	}

	if err := store.Add(p.ID, p, s.opts); err != nil {
		return err
	}

	return nil
}

func (s *Store) EachPublication(fn func(*models.Publication) bool) error {
	c := s.client.Store("publications").GetAll(s.opts)
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
	snap, err := s.client.Store("datasets").Get(id, s.opts)
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
	c := s.client.Store("datasets").GetByID(ids, s.opts)
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

func (s *Store) StoreDataset(d *models.Dataset) error {

	now := time.Now()

	if d.DateCreated == nil {
		d.DateCreated = &now
	}
	d.DateUpdated = &now

	if d.ID == "" {
		d.ID = uuid.NewString()
	}

	store := s.client.Store("datasets")

	// TODO this needs to be a separate update action
	if d.SnapshotID != "" {
		if err := store.UpdateSnapshot(d.SnapshotID, d, s.opts); err != nil {
			return err
		}
		return nil
	}

	if err := store.Add(d.ID, d, s.opts); err != nil {
		return err
	}

	return nil
}

func (s *Store) EachDataset(fn func(*models.Dataset) bool) error {
	c := s.client.Store("datasets").GetAll(s.opts)
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
