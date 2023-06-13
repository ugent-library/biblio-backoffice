package backends

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/go-orcid/orcid"
)

type Services struct {
	ORCIDSandbox              bool
	ORCIDClient               *orcid.MemberClient
	Repository                Repository
	FileStore                 FileStore
	SearchService             SearchService
	DatasetSearchIndex        DatasetIndex
	PublicationSearchIndex    PublicationIndex
	OrganizationService       OrganizationService
	PersonService             PersonService
	ProjectService            ProjectService
	UserService               UserService
	OrganizationSearchService OrganizationSearchService
	PersonSearchService       PersonSearchService
	ProjectSearchService      ProjectSearchService
	UserSearchService         UserSearchService
	LicenseSearchService      LicenseSearchService
	MediaTypeSearchService    MediaTypeSearchService
	PublicationSources        map[string]PublicationGetter
	DatasetSources            map[string]DatasetGetter
	PublicationEncoders       map[string]PublicationEncoder
	PublicationDecoders       map[string]PublicationDecoderFactory
	PublicationListExporters  map[string]PublicationListExporterFactory
	DatasetListExporters      map[string]DatasetListExporterFactory
	HandleService             HandleService
	// Tasks                      *tasks.Hub
}

type PublicationEncoder func(*models.Publication) ([]byte, error)

type PublicationDecoderFactory func(io.Reader) PublicationDecoder

type PublicationDecoder interface {
	Decode(*models.Publication) error
}

type DatasetGetter interface {
	GetDataset(string) (*models.Dataset, error)
}

type PublicationGetter interface {
	GetPublication(string) (*models.Publication, error)
}

type Repository interface {
	GetPublication(string) (*models.Publication, error)
	GetPublications([]string) ([]*models.Publication, error)
	SavePublication(*models.Publication, *models.User) error
	ImportPublication(*models.Publication) error
	UpdatePublication(string, *models.Publication, *models.User) error
	UpdatePublicationInPlace(*models.Publication) error
	MutatePublication(string, *models.User, ...Mutation) error
	PublicationsAfter(time.Time, int, int) (int, []*models.Publication, error)
	PublicationsBetween(time.Time, time.Time, func(*models.Publication) bool) error
	EachPublication(func(*models.Publication) bool) error
	EachPublicationSnapshot(func(*models.Publication) bool) error
	EachPublicationWithoutHandle(func(*models.Publication) bool) error
	PublicationHistory(string, func(*models.Publication) bool) error
	UpdatePublicationEmbargoes() (int, error)
	PurgeAllPublications() error
	PurgePublication(string) error
	GetDataset(string) (*models.Dataset, error)
	GetDatasets([]string) ([]*models.Dataset, error)
	ImportDataset(*models.Dataset) error
	SaveDataset(*models.Dataset, *models.User) error
	UpdateDataset(string, *models.Dataset, *models.User) error
	MutateDataset(string, *models.User, ...Mutation) error
	DatasetsAfter(time.Time, int, int) (int, []*models.Dataset, error)
	DatasetsBetween(time.Time, time.Time, func(*models.Dataset) bool) error
	EachDataset(func(*models.Dataset) bool) error
	EachDatasetSnapshot(func(*models.Dataset) bool) error
	EachDatasetWithoutHandle(func(*models.Dataset) bool) error
	DatasetHistory(string, func(*models.Dataset) bool) error
	UpdateDatasetEmbargoes() (int, error)
	PurgeAllDatasets() error
	PurgeDataset(string) error
	GetPublicationDatasets(*models.Publication) ([]*models.Dataset, error)
	GetVisiblePublicationDatasets(*models.User, *models.Publication) ([]*models.Dataset, error)
	GetDatasetPublications(*models.Dataset) ([]*models.Publication, error)
	GetVisibleDatasetPublications(*models.User, *models.Dataset) ([]*models.Publication, error)
	AddPublicationDataset(*models.Publication, *models.Dataset, *models.User) error
	RemovePublicationDataset(*models.Publication, *models.Dataset, *models.User) error
}

type FileStore interface {
	Exists(context.Context, string) (bool, error)
	Get(context.Context, string) (io.ReadCloser, error)
	Add(context.Context, io.Reader, string) (string, error)
	Delete(context.Context, string) error
	DeleteAll(context.Context) error
}

type BulkIndexerConfig struct {
	OnError      func(error)
	OnIndexError func(string, error)
}

type IndexSwitcher[T any] interface {
	Index(context.Context, T) error
	Switch(context.Context) error
}

type BulkIndexer[T any] interface {
	Index(context.Context, T) error
	Close(context.Context) error
}

type DatasetIDIndex interface {
	Search(*models.SearchArgs) (*models.SearchHits, error)
	Each(searchArgs *models.SearchArgs, maxSize int, cb func(string)) error
	Delete(id string) error
	DeleteAll() error
	WithScope(string, ...string) DatasetIDIndex
}

type DatasetIndex interface {
	Search(*models.SearchArgs) (*models.DatasetHits, error)
	Each(searchArgs *models.SearchArgs, maxSize int, cb func(*models.Dataset)) error
	Delete(id string) error
	DeleteAll() error
	WithScope(string, ...string) DatasetIndex
}

type datasetIndex struct {
	DatasetIDIndex
	repo Repository
}

func NewDatasetIndex(di DatasetIDIndex, r Repository) DatasetIndex {
	return &datasetIndex{
		DatasetIDIndex: di,
		repo:           r,
	}
}

func (ds *datasetIndex) Search(args *models.SearchArgs) (*models.DatasetHits, error) {
	h, err := ds.DatasetIDIndex.Search(args)
	if err != nil {
		return nil, err
	}

	datasets, err := ds.repo.GetDatasets(h.Hits)
	if err != nil {
		return nil, err
	}

	return &models.DatasetHits{
		Pagination: h.Pagination,
		Hits:       datasets,
		Facets:     h.Facets,
	}, nil
}

func (ds *datasetIndex) Each(searchArgs *models.SearchArgs, maxSize int, cb func(*models.Dataset)) error {
	return ds.DatasetIDIndex.Each(searchArgs, maxSize, func(id string) {
		// TODO handle error
		dataset, _ := ds.repo.GetDataset(id)
		cb(dataset)
	})
}

func (ds *datasetIndex) WithScope(field string, terms ...string) DatasetIndex {
	return &datasetIndex{
		DatasetIDIndex: ds.DatasetIDIndex.WithScope(field, terms...),
		repo:           ds.repo,
	}
}

type PublicationIDIndex interface {
	Search(*models.SearchArgs) (*models.SearchHits, error)
	Each(searchArgs *models.SearchArgs, maxSize int, cb func(string)) error
	Delete(id string) error
	DeleteAll() error
	WithScope(string, ...string) PublicationIDIndex
}

type PublicationIndex interface {
	Search(*models.SearchArgs) (*models.PublicationHits, error)
	Each(searchArgs *models.SearchArgs, maxSize int, cb func(*models.Publication)) error
	Delete(id string) error
	DeleteAll() error
	WithScope(string, ...string) PublicationIndex
}

type publicationIndex struct {
	PublicationIDIndex
	repo Repository
}

func NewPublicationIndex(di PublicationIDIndex, r Repository) PublicationIndex {
	return &publicationIndex{
		PublicationIDIndex: di,
		repo:               r,
	}
}

func (pi *publicationIndex) Search(args *models.SearchArgs) (*models.PublicationHits, error) {
	h, err := pi.PublicationIDIndex.Search(args)
	if err != nil {
		return nil, err
	}

	publications, err := pi.repo.GetPublications(h.Hits)
	if err != nil {
		return nil, err
	}

	return &models.PublicationHits{
		Pagination: h.Pagination,
		Hits:       publications,
		Facets:     h.Facets,
	}, nil
}

func (pi *publicationIndex) Each(searchArgs *models.SearchArgs, maxSize int, cb func(*models.Publication)) error {
	return pi.PublicationIDIndex.Each(searchArgs, maxSize, func(id string) {
		// TODO handle error
		publication, _ := pi.repo.GetPublication(id)
		cb(publication)
	})
}

func (pi *publicationIndex) WithScope(field string, terms ...string) PublicationIndex {
	return &publicationIndex{
		PublicationIDIndex: pi.PublicationIDIndex.WithScope(field, terms...),
		repo:               pi.repo,
	}
}

type SearchService interface {
	NewDatasetIndex(Repository) DatasetIndex
	NewDatasetBulkIndexer(BulkIndexerConfig) (BulkIndexer[*models.Dataset], error)
	NewDatasetIndexSwitcher(BulkIndexerConfig) (IndexSwitcher[*models.Dataset], error)
	NewPublicationIndex(Repository) PublicationIndex
	NewPublicationBulkIndexer(BulkIndexerConfig) (BulkIndexer[*models.Publication], error)
	NewPublicationIndexSwitcher(BulkIndexerConfig) (IndexSwitcher[*models.Publication], error)
}

type OrganizationService interface {
	GetOrganization(string) (*models.Organization, error)
}

type PersonService interface {
	GetPerson(string) (*models.Person, error)
}

type ProjectService interface {
	GetProject(string) (*models.Project, error)
}

type UserService interface {
	GetUser(string) (*models.User, error)
	GetUserByUsername(string) (*models.User, error)
}

type OrganizationSearchService interface {
	SuggestOrganizations(string) ([]models.Completion, error)
}

type PersonSearchService interface {
	SuggestPeople(string) ([]models.Person, error)
}

type ProjectSearchService interface {
	SuggestProjects(string) ([]models.Completion, error)
}

type UserSearchService interface {
	SuggestUsers(string) ([]models.Person, error)
}

type LicenseSearchService interface {
	IndexAll() error
	SuggestLicenses(string) ([]models.Completion, error)
}

type MediaTypeSearchService interface {
	IndexAll() error
	SuggestMediaTypes(string) ([]models.Completion, error)
}

type PublicationListExporter interface {
	GetContentType() string
	Add(*models.Publication)
	Flush() error
}

type PublicationListExporterFactory func(io.Writer) PublicationListExporter

type DatasetListExporter interface {
	GetContentType() string
	Add(*models.Dataset)
	Flush() error
}

type DatasetListExporterFactory func(io.Writer) DatasetListExporter

type HandleService interface {
	UpsertHandle(string) (*models.Handle, error)
}

var ErrNotFound = errors.New("record not found")

type Mutation struct {
	Op   string
	Args []string
}

type PersonWithOrganizationsService struct {
	PersonService       PersonService
	OrganizationService OrganizationService
}

func (s *PersonWithOrganizationsService) GetPerson(id string) (*models.Person, error) {
	p, err := s.PersonService.GetPerson(id)
	if err != nil {
		return nil, err
	}
	for _, a := range p.Affiliations {
		o, err := s.OrganizationService.GetOrganization(a.OrganizationID)
		if err == ErrNotFound {
			a.Organization = NewDummyOrganization(a.OrganizationID)
		} else if err != nil {
			return nil, err
		} else {
			a.Organization = o
		}
	}
	return p, nil
}

type UserWithOrganizationsService struct {
	UserService         UserService
	OrganizationService OrganizationService
}

func (s *UserWithOrganizationsService) GetUser(id string) (*models.User, error) {
	u, err := s.UserService.GetUser(id)
	if err != nil {
		return nil, err
	}
	for _, a := range u.Affiliations {
		o, err := s.OrganizationService.GetOrganization(a.OrganizationID)
		if err == ErrNotFound {
			a.Organization = NewDummyOrganization(a.OrganizationID)
		} else if err != nil {
			return nil, err
		} else {
			a.Organization = o
		}
	}
	return u, nil
}

func (s *UserWithOrganizationsService) GetUserByUsername(username string) (*models.User, error) {
	u, err := s.UserService.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	for _, a := range u.Affiliations {
		o, err := s.OrganizationService.GetOrganization(a.OrganizationID)
		if err == ErrNotFound {
			a.Organization = NewDummyOrganization(a.OrganizationID)
		} else if err != nil {
			return nil, err
		} else {
			a.Organization = o
		}
	}
	return u, nil
}

// TODO remove this when we always have an uptodate db of organizations
func NewDummyOrganization(id string) *models.Organization {
	return &models.Organization{
		ID:   id,
		Name: id,
		Tree: []struct {
			ID string `json:"id,omitempty"`
		}{
			{ID: id},
		},
	}
}

func NewDummyPerson(id string) *models.Person {
	return &models.Person{
		ID:        id,
		FullName:  "[missing]",
		FirstName: "[missing]",
		LastName:  "[missing]",
	}
}

func NewDummyProject(id string) *models.Project {
	return &models.Project{
		ID:    id,
		Title: "[missing]",
	}
}
