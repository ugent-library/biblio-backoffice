package backends

import (
	"context"
	"io"

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

type Store interface {
	Transaction(context.Context, func(Store) error) error
	GetPublication(string) (*models.Publication, error)
	GetPublications([]string) ([]*models.Publication, error)
	StorePublication(*models.Publication) error
	EachPublication(func(*models.Publication) bool) error
	GetDataset(string) (*models.Dataset, error)
	GetDatasets([]string) ([]*models.Dataset, error)
	StoreDataset(*models.Dataset) error
	EachDataset(func(*models.Dataset) bool) error
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
