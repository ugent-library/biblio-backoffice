package backends

import (
	"io"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
)

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

type Transaction interface {
	Rollback() error
	Commit() error
	SavePublication(p *models.Publication) (*models.Publication, error)
	SaveDataset(d *models.Dataset) (*models.Dataset, error)
}

type StorageService interface {
	Begin() (Transaction, error)
	GetDataset(string) (*models.Dataset, error)
	GetDatasets([]string) ([]*models.Dataset, error)
	SaveDataset(*models.Dataset) (*models.Dataset, error)
	EachDataset(func(*models.Dataset) bool) error
	GetPublication(string) (*models.Publication, error)
	GetPublications([]string) ([]*models.Publication, error)
	SavePublication(*models.Publication) (*models.Publication, error)
	EachPublication(func(*models.Publication) bool) error
}

type DatasetSearchService interface {
	SearchDatasets(*models.SearchArgs) (*models.DatasetHits, error)
	IndexDataset(*models.Dataset) error
	IndexDatasets(<-chan *models.Dataset)
}

type PublicationSearchService interface {
	SearchPublications(*models.SearchArgs) (*models.PublicationHits, error)
	IndexPublication(*models.Publication) error
	IndexPublications(<-chan *models.Publication)
}

type PublicationService interface {
	GetPublication(string) (*models.Publication, error)
	GetPublicationDatasets(string) ([]*models.Dataset, error)
	GetDatasetPublications(string) ([]*models.Publication, error)
	CreateUserPublication(string, string) (*models.Publication, error)
	ImportUserPublicationByIdentifier(string, string, string) (*models.Publication, error)
	ImportUserPublications(string, string, io.Reader) (string, error)
	UpdatePublication(*models.Publication) (*models.Publication, error)
	PublishPublication(*models.Publication) (*models.Publication, error)
	BatchPublishPublications(string, *models.SearchArgs) error
	AddPublicationDataset(string, string) error
	RemovePublicationDataset(string, string) error
	ServePublicationFile(string, http.ResponseWriter, *http.Request)
	AddPublicationFile(string, *models.PublicationFile, io.Reader) error
	UpdatePublicationFile(string, *models.PublicationFile) error
	RemovePublicationFile(id, fileID string) error
	DeletePublication(string) error
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

type LicenseSearchService interface {
	SuggestLicenses(string) ([]models.Completion, error)
}

type MediaTypeSearchService interface {
	SuggestMediaTypes(string) ([]models.Completion, error)
}
