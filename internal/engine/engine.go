package engine

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/rabbitmq/amqp091-go"
	"github.com/ugent-library/biblio-backend/internal/message"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/task"
	"github.com/ugent-library/go-orcid/orcid"
	"go.temporal.io/sdk/client"
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
	Datasets(*SearchArgs) (*models.DatasetHits, error)
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
	Publications(*SearchArgs) (*models.PublicationHits, error)
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
	Temporal     client.Client
	MQ           *amqp091.Connection
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
	Tasks    *task.Hub
	Messages *message.Hub
}

func New(c Config) (*Engine, error) {
	e := &Engine{Config: c, Messages: message.NewHub()}

	temporal, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	// defer c.Close()
	e.Temporal = temporal

	mqCh, err := e.MQ.Channel()
	if err != nil {
		log.Fatal(err)
	}

	err = mqCh.ExchangeDeclare(
		"tasks", // exchange name
		"topic", // exchange type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, err
	}

	// receive notifications
	err = mqCh.ExchangeDeclare(
		"notifications", // exchange name
		"fanout",        // exchange type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return nil, err
	}

	q, err := mqCh.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	err = mqCh.QueueBind(
		q.Name,          // queue name
		"",              // routing key
		"notifications", // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := mqCh.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}

	e.Tasks = task.NewHub(mqCh)

	go e.Messages.Run()

	go func() {
		// dispatch message based on user_id
		for d := range msgs {
			msg := struct {
				UserID string `json:"user_id"`
			}{}
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Println(err)
			}
			e.Messages.Dispatch(msg.UserID, []byte(d.Body))
		}
	}()

	return e, nil
}
