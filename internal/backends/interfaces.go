package backends

import (
	"io"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
)

type DatasetSource interface {
	GetDataset(string) (*models.Dataset, error)
}

type DatasetService interface {
	GetDataset(string) (*models.Dataset, error)
	CreateDataset(*models.Dataset) (*models.Dataset, error)
	UpdateDataset(*models.Dataset) (*models.Dataset, error)
	DeleteDataset(string) error
}

type DatasetSearchService interface {
	Datasets(*models.SearchArgs) (*models.DatasetHits, error)
	UserDatasets(string, *models.SearchArgs) (*models.DatasetHits, error)
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

type PublicationSearchService interface {
	Publications(*models.SearchArgs) (*models.PublicationHits, error)
	UserPublications(string, *models.SearchArgs) (*models.PublicationHits, error)
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
