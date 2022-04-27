package commands

import (
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/backends/arxiv"
	"github.com/ugent-library/biblio-backend/internal/backends/biblio"
	"github.com/ugent-library/biblio-backend/internal/backends/citeproc"
	"github.com/ugent-library/biblio-backend/internal/backends/crossref"
	"github.com/ugent-library/biblio-backend/internal/backends/datacite"
	"github.com/ugent-library/biblio-backend/internal/backends/es6"
	"github.com/ugent-library/biblio-backend/internal/backends/filestore"
	"github.com/ugent-library/biblio-backend/internal/backends/ianamedia"
	"github.com/ugent-library/biblio-backend/internal/backends/pg"
	"github.com/ugent-library/biblio-backend/internal/backends/pg/bibtex"
	"github.com/ugent-library/biblio-backend/internal/backends/pg/jsonl"
	"github.com/ugent-library/biblio-backend/internal/backends/pg/ris"
	"github.com/ugent-library/biblio-backend/internal/backends/pubmed"
	"github.com/ugent-library/biblio-backend/internal/backends/spdxlicenses"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/go-orcid/orcid"
	"go.temporal.io/sdk/client"
)

var (
	_engine     *engine.Engine
	_engineOnce sync.Once
)

func Engine() *engine.Engine {
	_engineOnce.Do(func() {
		_engine = newEngine()
	})
	return _engine
}

func newEngine() *engine.Engine {
	fs, err := filestore.New(viper.GetString("file-dir"))
	if err != nil {
		log.Fatalln("Unable to initialize filestore", err)
	}

	temporal, err := client.NewClient(client.Options{
		HostPort: viper.GetString("temporal-host-port"),
		Logger:   &temporalLogger{},
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}

	es6Client := newEs6Client()

	biblioClient := biblio.New(biblio.Config{
		URL:      viper.GetString("frontend-url"),
		Username: viper.GetString("frontend-username"),
		Password: viper.GetString("frontend-password"),
	})

	orcidConfig := orcid.Config{
		ClientID:     viper.GetString("orcid-client-id"),
		ClientSecret: viper.GetString("orcid-client-secret"),
		Sandbox:      viper.GetBool("orcid-sandbox"),
	}
	orcidClient := orcid.NewMemberClient(orcidConfig)

	e, err := engine.New(engine.Config{
		FileStore:                 fs,
		Temporal:                  temporal,
		ORCIDSandbox:              orcidConfig.Sandbox,
		ORCIDClient:               orcidClient,
		StorageService:            newStorageService(),
		DatasetSearchService:      es6Client,
		PublicationSearchService:  es6Client,
		PersonService:             biblioClient,
		ProjectService:            biblioClient,
		UserService:               biblioClient,
		OrganizationSearchService: biblioClient,
		PersonSearchService:       biblioClient,
		ProjectSearchService:      biblioClient,
		LicenseSearchService:      spdxlicenses.New(),
		MediaTypeSearchService:    ianamedia.New(),
		DatasetSources: map[string]backends.DatasetGetter{
			"datacite": datacite.New(),
		},
		PublicationSources: map[string]backends.PublicationGetter{
			"crossref": crossref.New(),
			"pubmed":   pubmed.New(),
			"arxiv":    arxiv.New(),
		},
		PublicationEncoders: map[string]backends.PublicationEncoder{
			"cite-mla":                 citeproc.New("mla").EncodePublication,
			"cite-apa":                 citeproc.New("apa").EncodePublication,
			"cite-chicago-author-date": citeproc.New("chicago-author-date").EncodePublication,
			"cite-fwo":                 citeproc.New("fwo").EncodePublication,
			"cite-vancouver":           citeproc.New("vancouver").EncodePublication,
			"cite-ieee":                citeproc.New("ieee").EncodePublication,
		},
		PublicationDecoders: map[string]backends.PublicationDecoderFactory{
			"jsonl":  jsonl.NewDecoder,
			"ris":    ris.NewDecoder,
			"wos":    ris.NewDecoder,
			"bibtex": bibtex.NewDecoder,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	return e
}

func newStorageService() backends.StorageService {
	client, err := pg.New(viper.GetString("pg-conn"))
	if err != nil {
		log.Fatalln("unable to create pg dataset service", err)
	}
	return client
}

func newEs6Client() *es6.Client {
	datasetSettings, err := ioutil.ReadFile("etc/es6/dataset.json")
	if err != nil {
		log.Fatal(err)
	}
	publicationSettings, err := ioutil.ReadFile("etc/es6/publication.json")
	if err != nil {
		log.Fatal(err)
	}
	client, err := es6.New(es6.Config{
		ClientConfig: elasticsearch.Config{
			Addresses: strings.Split(viper.GetString("es6-url"), ","),
		},
		DatasetIndex:        viper.GetString("dataset-index"),
		DatasetSettings:     string(datasetSettings),
		PublicationIndex:    viper.GetString("publication-index"),
		PublicationSettings: string(publicationSettings),
	})
	if err != nil {
		log.Fatal(err)
	}
	return client
}

type temporalLogger struct{}

func (l *temporalLogger) Debug(msg string, keyvals ...interface{}) {
	log.Println(append([]interface{}{"DEBUG", "TEMPORAL", msg}, keyvals...)...)
}

func (l *temporalLogger) Info(msg string, keyvals ...interface{}) {
	log.Println(append([]interface{}{"INFO", "TEMPORAL", msg}, keyvals...)...)
}

func (l *temporalLogger) Warn(msg string, keyvals ...interface{}) {
	log.Println(append([]interface{}{"WARN", "TEMPORAL", msg}, keyvals...)...)
}

func (l *temporalLogger) Error(msg string, keyvals ...interface{}) {
	log.Println(append([]interface{}{"ERROR", "TEMPORAL", msg}, keyvals...)...)
}
