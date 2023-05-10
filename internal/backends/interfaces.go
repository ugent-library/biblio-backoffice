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
	Transaction(context.Context, func(Repository) error) error
	AddPublicationListener(func(*models.Publication))
	GetPublication(string) (*models.Publication, error)
	GetPublications([]string) ([]*models.Publication, error)
	SavePublication(*models.Publication, *models.User) error
	ImportCurrentPublication(*models.Publication) error
	ImportOldPublication(*models.Publication) error
	UpdatePublication(string, *models.Publication, *models.User) error
	UpdatePublicationInPlace(p *models.Publication) error
	PublicationsAfter(t time.Time, limit, offset int) (int, []*models.Publication, error)
	PublicationsBetween(time.Time, time.Time, func(*models.Publication) bool) error
	EachPublication(func(*models.Publication) bool) error
	EachPublicationSnapshot(func(*models.Publication) bool) error
	EachPublicationWithoutHandle(func(*models.Publication) bool) error
	PublicationHistory(string, func(*models.Publication) bool) error
	UpdatePublicationEmbargoes() (int, error)
	PurgeAllPublications() error
	PurgePublication(string) error
	AddDatasetListener(func(*models.Dataset))
	GetDataset(string) (*models.Dataset, error)
	GetDatasets([]string) ([]*models.Dataset, error)
	ImportCurrentDataset(*models.Dataset) error
	ImportOldDataset(*models.Dataset) error
	SaveDataset(*models.Dataset, *models.User) error
	UpdateDataset(string, *models.Dataset, *models.User) error
	DatasetsAfter(t time.Time, limit, offset int) (int, []*models.Dataset, error)
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

type DatasetSearchService interface {
	Search(*models.SearchArgs) (*models.DatasetHits, error)
	Index(*models.Dataset) error
	Delete(id string) error
	DeleteAll() error
	WithScope(string, ...string) DatasetSearchService
	NewBulkIndexer(BulkIndexerConfig) (BulkIndexer[*models.Dataset], error)
	NewIndexSwitcher(BulkIndexerConfig) (IndexSwitcher[*models.Dataset], error)
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

type PublicationSearcherService interface {
	GetMaxSize() int
	SetMaxSize(int)
	WithScope(string, ...string) PublicationSearcherService
	Searcher(*models.SearchArgs, func(*models.Publication)) error
}

type DatasetSearcherService interface {
	GetMaxSize() int
	SetMaxSize(int)
	WithScope(string, ...string) DatasetSearcherService
	Searcher(*models.SearchArgs, func(*models.Dataset)) error
}

type OrganizationService interface {
	GetOrganization(string) (*models.Organization, error)
}

type PersonService interface {
	GetPerson(string) (*models.Person, error)
	GetPersons([]string) ([]*models.Person, error)
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

type RepositoryFilter struct {
	Field string
	Op    string
	Value string
}

type RepositoryQueryArgs struct {
	Limit   int
	Offset  int
	Order   string
	Filters []*RepositoryFilter
}
