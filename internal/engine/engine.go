package engine

import (
	"io"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-orcid/orcid"
)

type DatasetService interface {
	GetDataset(string) (*models.Dataset, error)
	GetDatasetPublications(string) ([]*models.Publication, error)
	ImportUserDatasetByIdentifier(string, string, string) (*models.Dataset, error)
	UpdateDataset(*models.Dataset) (*models.Dataset, error)
	PublishDataset(*models.Dataset) (*models.Dataset, error)
	DeleteDataset(string) error
}

type DatasetSearchService interface {
	UserDatasets(string, *SearchArgs) (*models.DatasetHits, error)
}

type PublicationService interface {
	GetPublication(string) (*models.Publication, error)
	GetPublicationDatasets(string) ([]*models.Dataset, error)
	CreateUserPublication(string, string) (*models.Publication, error)
	ImportUserPublicationByIdentifier(string, string, string) (*models.Publication, error)
	ImportUserPublications(string, string, io.Reader) (string, error)
	UpdatePublication(*models.Publication) (*models.Publication, error)
	PublishPublication(*models.Publication) (*models.Publication, error)
	BatchPublishPublications(string, *SearchArgs) error
	AddPublicationDataset(string, string) error
	RemovePublicationDataset(string, string) error
	ServePublicationFile(string, http.ResponseWriter, *http.Request)
	AddPublicationFile(string, *models.PublicationFile, io.Reader) error
	UpdatePublicationFile(string, *models.PublicationFile) error
	RemovePublicationFile(id, fileID string) error
	DeletePublication(string) error
}

type PublicationSearchService interface {
	UserPublications(string, *SearchArgs) (*models.PublicationHits, error)
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

type SuggestionService interface {
	SuggestOrganizations(string) ([]models.Completion, error)
	SuggestPeople(string) ([]models.Person, error)
	SuggestProjects(string) ([]models.Completion, error)
}

type Engine struct {
	ORCIDSandbox bool
	ORCIDClient  *orcid.MemberClient
	DatasetService
	DatasetSearchService
	PersonService
	ProjectService
	PublicationService
	PublicationSearchService
	UserService
	SuggestionService
}
