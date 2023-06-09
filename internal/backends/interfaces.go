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
	ORCIDSandbox               bool
	ORCIDClient                *orcid.MemberClient
	Repository                 Repository
	FileStore                  FileStore
	DatasetSearchService       DatasetSearchService
	PublicationSearchService   PublicationSearchService
	OrganizationService        OrganizationService
	PersonService              PersonService
	ProjectService             ProjectService
	UserService                UserService
	OrganizationSearchService  OrganizationSearchService
	PersonSearchService        PersonSearchService
	ProjectSearchService       ProjectSearchService
	UserSearchService          UserSearchService
	LicenseSearchService       LicenseSearchService
	MediaTypeSearchService     MediaTypeSearchService
	PublicationSources         map[string]PublicationGetter
	DatasetSources             map[string]DatasetGetter
	PublicationEncoders        map[string]PublicationEncoder
	PublicationDecoders        map[string]PublicationDecoderFactory
	PublicationListExporters   map[string]PublicationListExporterFactory
	PublicationSearcherService PublicationSearcherService
	DatasetListExporters       map[string]DatasetListExporterFactory
	DatasetSearcherService     DatasetSearcherService
	HandleService              HandleService
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

type DatasetIDSearchService interface {
	Search(*models.SearchArgs) (*models.SearchHits, error)
	Index(*models.Dataset) error
	Delete(id string) error
	DeleteAll() error
	WithScope(string, ...string) DatasetIDSearchService
	NewBulkIndexer(BulkIndexerConfig) (BulkIndexer[*models.Dataset], error)
	NewIndexSwitcher(BulkIndexerConfig) (IndexSwitcher[*models.Dataset], error)
}

type DatasetSearchService interface {
	Search(*models.SearchArgs) (*models.DatasetHits, error)
	Index(*models.Dataset) error
	Delete(id string) error
	DeleteAll() error
	WithScope(string, ...string) DatasetSearchService
	NewBulkIndexer(BulkIndexerConfig) (BulkIndexer[*models.Dataset], error)
	NewIndexSwitcher(BulkIndexerConfig) (IndexSwitcher[*models.Dataset], error)
}

type datasetSearchService struct {
	DatasetIDSearchService
	repo Repository
}

func NewDatasetSearchService(s DatasetIDSearchService, r Repository) DatasetSearchService {
	return &datasetSearchService{
		DatasetIDSearchService: s,
		repo:                   r,
	}
}

func (s *datasetSearchService) Search(args *models.SearchArgs) (*models.DatasetHits, error) {
	h, err := s.DatasetIDSearchService.Search(args)
	if err != nil {
		return nil, err
	}
	pubs, err := s.repo.GetDatasets(h.Hits)
	if err != nil {
		return nil, err
	}
	return &models.DatasetHits{
		Pagination: h.Pagination,
		Hits:       pubs,
		Facets:     h.Facets,
	}, nil
}

func (s *datasetSearchService) WithScope(field string, terms ...string) DatasetSearchService {
	return &datasetSearchService{
		DatasetIDSearchService: s.DatasetIDSearchService.WithScope(field, terms...),
		repo:                   s.repo,
	}
}

type PublicationIDSearchService interface {
	Search(*models.SearchArgs) (*models.SearchHits, error)
	Index(*models.Publication) error
	Delete(id string) error
	DeleteAll() error
	WithScope(string, ...string) PublicationIDSearchService
	NewBulkIndexer(BulkIndexerConfig) (BulkIndexer[*models.Publication], error)
	NewIndexSwitcher(BulkIndexerConfig) (IndexSwitcher[*models.Publication], error)
}

type PublicationSearchService interface {
	Search(*models.SearchArgs) (*models.PublicationHits, error)
	Index(*models.Publication) error
	Delete(id string) error
	DeleteAll() error
	WithScope(string, ...string) PublicationSearchService
	NewBulkIndexer(BulkIndexerConfig) (BulkIndexer[*models.Publication], error)
	NewIndexSwitcher(BulkIndexerConfig) (IndexSwitcher[*models.Publication], error)
}

type publicationSearchService struct {
	PublicationIDSearchService
	repo Repository
}

func NewPublicationSearchService(s PublicationIDSearchService, r Repository) PublicationSearchService {
	return &publicationSearchService{
		PublicationIDSearchService: s,
		repo:                       r,
	}
}

func (s *publicationSearchService) Search(args *models.SearchArgs) (*models.PublicationHits, error) {
	h, err := s.PublicationIDSearchService.Search(args)
	if err != nil {
		return nil, err
	}
	pubs, err := s.repo.GetPublications(h.Hits)
	if err != nil {
		return nil, err
	}
	return &models.PublicationHits{
		Pagination: h.Pagination,
		Hits:       pubs,
		Facets:     h.Facets,
	}, nil
}

func (s *publicationSearchService) WithScope(field string, terms ...string) PublicationSearchService {
	return &publicationSearchService{
		PublicationIDSearchService: s.PublicationIDSearchService.WithScope(field, terms...),
		repo:                       s.repo,
	}
}

type PublicationIDSearcherService interface {
	GetMaxSize() int
	SetMaxSize(int)
	WithScope(string, ...string) PublicationIDSearcherService
	Searcher(*models.SearchArgs, func(string)) error
}

type PublicationSearcherService interface {
	GetMaxSize() int
	SetMaxSize(int)
	WithScope(string, ...string) PublicationSearcherService
	Searcher(*models.SearchArgs, func(*models.Publication)) error
}

type publicationSearcherService struct {
	PublicationIDSearcherService
	repo Repository
}

func NewPublicationSearcherService(s PublicationIDSearcherService, r Repository) PublicationSearcherService {
	return &publicationSearcherService{
		PublicationIDSearcherService: s,
		repo:                         r,
	}
}

func (s *publicationSearcherService) Searcher(args *models.SearchArgs, fn func(*models.Publication)) error {
	return s.PublicationIDSearcherService.Searcher(args, func(id string) {
		// TODO handle error
		pub, _ := s.repo.GetPublication(id)
		fn(pub)
	})
}

func (s *publicationSearcherService) WithScope(field string, terms ...string) PublicationSearcherService {
	return &publicationSearcherService{
		PublicationIDSearcherService: s.PublicationIDSearcherService.WithScope(field, terms...),
		repo:                         s.repo,
	}
}

type DatasetIDSearcherService interface {
	GetMaxSize() int
	SetMaxSize(int)
	WithScope(string, ...string) DatasetIDSearcherService
	Searcher(*models.SearchArgs, func(string)) error
}

type DatasetSearcherService interface {
	GetMaxSize() int
	SetMaxSize(int)
	WithScope(string, ...string) DatasetSearcherService
	Searcher(*models.SearchArgs, func(*models.Dataset)) error
}

type datasetSearcherService struct {
	DatasetIDSearcherService
	repo Repository
}

func NewDatasetSearcherService(s DatasetIDSearcherService, r Repository) DatasetSearcherService {
	return &datasetSearcherService{
		DatasetIDSearcherService: s,
		repo:                     r,
	}
}

func (s *datasetSearcherService) Searcher(args *models.SearchArgs, fn func(*models.Dataset)) error {
	return s.DatasetIDSearcherService.Searcher(args, func(id string) {
		// TODO handle error
		pub, _ := s.repo.GetDataset(id)
		fn(pub)
	})
}

func (s *datasetSearcherService) WithScope(field string, terms ...string) DatasetSearcherService {
	return &datasetSearcherService{
		DatasetIDSearcherService: s.DatasetIDSearcherService.WithScope(field, terms...),
		repo:                     s.repo,
	}
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
