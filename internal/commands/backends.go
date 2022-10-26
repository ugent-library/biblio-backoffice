package commands

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/backends/arxiv"
	"github.com/ugent-library/biblio-backend/internal/backends/biblio"
	"github.com/ugent-library/biblio-backend/internal/backends/bibtex"
	"github.com/ugent-library/biblio-backend/internal/backends/citeproc"
	"github.com/ugent-library/biblio-backend/internal/backends/crossref"
	"github.com/ugent-library/biblio-backend/internal/backends/datacite"
	"github.com/ugent-library/biblio-backend/internal/backends/es6"
	excel_dataset "github.com/ugent-library/biblio-backend/internal/backends/excel/dataset"
	excel_publication "github.com/ugent-library/biblio-backend/internal/backends/excel/publication"
	"github.com/ugent-library/biblio-backend/internal/backends/handle"

	"github.com/ugent-library/biblio-backend/internal/backends/filestore"
	"github.com/ugent-library/biblio-backend/internal/backends/ianamedia"
	"github.com/ugent-library/biblio-backend/internal/backends/jsonl"
	"github.com/ugent-library/biblio-backend/internal/backends/pubmed"
	"github.com/ugent-library/biblio-backend/internal/backends/repository"
	"github.com/ugent-library/biblio-backend/internal/backends/ris"
	"github.com/ugent-library/biblio-backend/internal/backends/spdxlicenses"
	"go.uber.org/zap"

	// "github.com/ugent-library/biblio-backend/internal/tasks"
	"github.com/ugent-library/go-orcid/orcid"
)

var (
	_services     *backends.Services
	_servicesOnce sync.Once
)

func Services() *backends.Services {
	_servicesOnce.Do(func() {
		_services = newServices()
	})
	return _services
}

func newServices() *backends.Services {
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

	citeprocURL := viper.GetString("citeproc-url")

	var handleService backends.HandleService = nil

	if viper.GetBool("hdl-srv-enabled") {
		handleService = handle.NewClient(
			handle.Config{
				BaseURL:         viper.GetString("hdl-srv-url"),
				FrontEndBaseURL: fmt.Sprintf("%s/publication", viper.GetString("frontend-url")),
				Prefix:          viper.GetString("hdl-srv-prefix"),
				Username:        viper.GetString("hdl-srv-username"),
				Password:        viper.GetString("hdl-srv-password"),
			},
		)
	}

	return &backends.Services{
		FileStore:                 newFileStore(),
		ORCIDSandbox:              orcidConfig.Sandbox,
		ORCIDClient:               orcidClient,
		Repository:                newRepository(),
		DatasetSearchService:      newDatasetSearchService(),
		PublicationSearchService:  newPublicationSearchService(),
		OrganizationService:       biblioClient,
		PersonService:             biblioClient,
		ProjectService:            biblioClient,
		UserService:               biblioClient,
		OrganizationSearchService: biblioClient,
		PersonSearchService:       biblioClient,
		ProjectSearchService:      biblioClient,
		UserSearchService:         biblioClient,
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
			"cite-mla":                 citeproc.New(citeprocURL, "mla").EncodePublication,
			"cite-apa":                 citeproc.New(citeprocURL, "apa").EncodePublication,
			"cite-chicago-author-date": citeproc.New(citeprocURL, "chicago-author-date").EncodePublication,
			"cite-fwo":                 citeproc.New(citeprocURL, "fwo").EncodePublication,
			"cite-vancouver":           citeproc.New(citeprocURL, "vancouver").EncodePublication,
			"cite-ieee":                citeproc.New(citeprocURL, "ieee").EncodePublication,
		},
		PublicationDecoders: map[string]backends.PublicationDecoderFactory{
			"jsonl":  jsonl.NewDecoder,
			"ris":    ris.NewDecoder,
			"wos":    ris.NewDecoder,
			"bibtex": bibtex.NewDecoder,
		},
		// Tasks: tasks.NewHub(),
		PublicationListExporters: map[string]backends.PublicationListExporterFactory{
			"xlsx": excel_publication.NewExporter,
		},
		PublicationSearcherService: newPublicationSearcherService(),
		DatasetListExporters: map[string]backends.DatasetListExporterFactory{
			"xlsx": excel_dataset.NewExporter,
		},
		DatasetSearcherService: newDatasetSearcherService(),
		HandleService:          handleService,
	}
}

func newLogger() *zap.SugaredLogger {
	logEnv := viper.GetString("mode")

	var logger *zap.Logger
	var err error

	if logEnv == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatalln("Unable to initialize logger", err)
	}

	sugar := logger.Sugar()

	return sugar
}

func newRepository() backends.Repository {
	s, err := repository.New(viper.GetString("pg-conn"))
	if err != nil {
		log.Fatalln("unable to create store", err)
	}
	return s
}

func newFileStore() *filestore.Store {
	fs, err := filestore.New(viper.GetString("file-dir"))
	if err != nil {
		log.Fatalln("Unable to initialize filestore", err)
	}
	return fs
}

func newEs6Client(t string) *es6.Client {
	settings, err := os.ReadFile("etc/es6/" + t + ".json")
	if err != nil {
		log.Fatal(err)
	}
	client, err := es6.New(es6.Config{
		ClientConfig: elasticsearch.Config{
			Addresses: strings.Split(viper.GetString("es6-url"), ","),
		},
		Index:    viper.GetString(t + "-index"),
		Settings: string(settings),
	})
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func newPublicationSearchService() backends.PublicationSearchService {
	es6Client := newEs6Client("publication")
	return es6.NewPublications(*es6Client)

}

func newDatasetSearchService() backends.DatasetSearchService {
	es6Client := newEs6Client("dataset")
	return es6.NewDatasets(*es6Client)

}

func newPublicationSearcherService() backends.PublicationSearcherService {
	es6Client := newEs6Client("publication")
	//max size of exportable records is now 10K. Make configurable
	return es6.NewPublicationSearcher(*es6Client, 10000)
}

func newDatasetSearcherService() backends.DatasetSearcherService {
	es6Client := newEs6Client("dataset")
	//max size of exportable records is now 10K. Make configurable
	return es6.NewDatasetSearcher(*es6Client, 10000)
}
