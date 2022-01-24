package engine

import (
	"io"
	"log"
	"net/http"

	"github.com/nats-io/nats.go"
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

type Config struct {
	NATS         *nats.Conn
	ORCIDSandbox bool
	ORCIDClient  *orcid.MemberClient
	DatasetService
	DatasetSearchService
	PersonService
	ProjectService
	PublicationService
	PublicationSearchService
	UserService
	OrganizationSearchService
	PersonSearchService
	ProjectSearchService
	LicenseSearchService
	MediaTypeSearchService
}

type Engine struct {
	Config
	js nats.JetStreamContext
}

func New(c Config) (*Engine, error) {
	e := &Engine{Config: c}

	js, err := e.NATS.JetStream()
	if err != nil {
		return e, err
	}
	if err = createWorkStream(js); err != nil {
		return e, err
	}
	e.js = js

	return e, nil
}

func createWorkStream(js nats.JetStreamContext) error {
	streamName := "WORK"
	streamSubjects := "WORK.*"

	// Check if the WORK stream already exists; if not, create it.
	stream, err := js.StreamInfo(streamName)
	if err != nil {
		log.Println(err)
	}

	// js.DeleteStream("WORK")
	if stream == nil {
		log.Printf("creating stream %q and subjects %q", streamName, streamSubjects)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:      streamName,
			Subjects:  []string{streamSubjects},
			Storage:   nats.FileStorage,
			Retention: nats.WorkQueuePolicy,
			// Discard:    nats.DiscardOld,
			// Duplicates: 1 * time.Hour,
			// MaxMsgs:    -1,
			// MaxBytes:   -1,
			// MaxAge:     -1,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
