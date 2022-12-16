package backends

import (
	"context"
	"errors"
	"io"

	"github.com/ugent-library/biblio-backend/internal/backends/filestore"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-orcid/orcid"
)

type Services struct {
	ORCIDSandbox               bool
	ORCIDClient                *orcid.MemberClient
	Repository                 Repository
	FileStore                  *filestore.Store
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
	CountPublications(*RepositoryQueryArgs) (int, error)
	SearchPublications(*RepositoryQueryArgs) ([]*models.Publication, error)
	SelectPublications(string, []any, func(*models.Publication) bool) error
	EachPublication(func(*models.Publication) bool) error
	EachPublicationSnapshot(func(*models.Publication) bool) error
	PublicationHistory(string, func(*models.Publication) bool) error
	PurgeAllPublications() error
	PurgePublication(string) error
	AddDatasetListener(func(*models.Dataset))
	GetDataset(string) (*models.Dataset, error)
	GetDatasets([]string) ([]*models.Dataset, error)
	ImportCurrentDataset(*models.Dataset) error
	ImportOldDataset(*models.Dataset) error
	SaveDataset(*models.Dataset, *models.User) error
	UpdateDataset(string, *models.Dataset, *models.User) error
	CountDatasets(*RepositoryQueryArgs) (int, error)
	SearchDatasets(*RepositoryQueryArgs) ([]*models.Dataset, error)
	SelectDatasets(string, []any, func(*models.Dataset) bool) error
	EachDataset(func(*models.Dataset) bool) error
	EachDatasetSnapshot(func(*models.Dataset) bool) error
	DatasetHistory(string, func(*models.Dataset) bool) error
	PurgeAllDatasets() error
	PurgeDataset(string) error
	GetPublicationDatasets(*models.Publication) ([]*models.Dataset, error)
	GetVisiblePublicationDatasets(*models.User, *models.Publication) ([]*models.Dataset, error)
	GetDatasetPublications(*models.Dataset) ([]*models.Publication, error)
	GetVisibleDatasetPublications(*models.User, *models.Dataset) ([]*models.Publication, error)
	AddPublicationDataset(*models.Publication, *models.Dataset, *models.User) error
	RemovePublicationDataset(*models.Publication, *models.Dataset, *models.User) error
}

type DatasetSearchService interface {
	Search(*models.SearchArgs) (*models.DatasetHits, error)
	Index(*models.Dataset) error
	IndexMultiple(<-chan *models.Dataset)
	Delete(id string) error
	WithScope(string, ...string) DatasetSearchService
	CreateIndex() error
	DeleteIndex() error
}

type PublicationSearchService interface {
	Search(*models.SearchArgs) (*models.PublicationHits, error)
	Index(*models.Publication) error
	IndexMultiple(<-chan *models.Publication)
	Delete(id string) error
	WithScope(string, ...string) PublicationSearchService
	CreateIndex() error
	DeleteIndex() error
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
	// GetHandle(string) (*models.Handle, error)
	UpsertHandle(string) (*models.Handle, error)
	// DeleteHandle(string) (*models.Handle, error)
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
