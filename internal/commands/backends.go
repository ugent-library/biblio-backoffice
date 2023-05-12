package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/backends/arxiv"
	"github.com/ugent-library/biblio-backoffice/internal/backends/authority"
	"github.com/ugent-library/biblio-backoffice/internal/backends/bibtex"
	"github.com/ugent-library/biblio-backoffice/internal/backends/citeproc"
	"github.com/ugent-library/biblio-backoffice/internal/backends/crossref"
	"github.com/ugent-library/biblio-backoffice/internal/backends/datacite"
	"github.com/ugent-library/biblio-backoffice/internal/backends/es6"
	excel_dataset "github.com/ugent-library/biblio-backoffice/internal/backends/excel/dataset"
	excel_publication "github.com/ugent-library/biblio-backoffice/internal/backends/excel/publication"
	"github.com/ugent-library/biblio-backoffice/internal/backends/fsstore"
	"github.com/ugent-library/biblio-backoffice/internal/backends/handle"
	"github.com/ugent-library/biblio-backoffice/internal/backends/s3store"
	"github.com/ugent-library/biblio-backoffice/internal/caching"
	"github.com/ugent-library/biblio-backoffice/internal/models"

	"github.com/ugent-library/biblio-backoffice/internal/backends/ianamedia"
	"github.com/ugent-library/biblio-backoffice/internal/backends/jsonl"
	"github.com/ugent-library/biblio-backoffice/internal/backends/pubmed"
	"github.com/ugent-library/biblio-backoffice/internal/backends/repository"
	"github.com/ugent-library/biblio-backoffice/internal/backends/ris"
	"github.com/ugent-library/biblio-backoffice/internal/backends/spdxlicenses"
	"go.uber.org/zap"

	// "github.com/ugent-library/biblio-backoffice/internal/tasks"
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
	authorityClient, authorityClientErr := authority.New(authority.Config{
		MongoDBURI: viper.GetString("mongodb-url"),
		ES6Config: es6.Config{
			ClientConfig: elasticsearch.Config{
				Addresses: strings.Split(viper.GetString("frontend-es6-url"), ","),
			},
		},
	})
	if authorityClientErr != nil {
		panic(authorityClientErr)
	}

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

	logger := newLogger()

	return &backends.Services{
		FileStore:                 newFileStore(),
		ORCIDSandbox:              orcidConfig.Sandbox,
		ORCIDClient:               orcidClient,
		Repository:                newRepository(logger),
		DatasetSearchService:      newDatasetSearchService(),
		PublicationSearchService:  newPublicationSearchService(),
		OrganizationService:       caching.NewOrganizationService(authorityClient),
		PersonService:             caching.NewPersonService(authorityClient),
		ProjectService:            caching.NewProjectService(authorityClient),
		UserService:               caching.NewUserService(authorityClient),
		OrganizationSearchService: authorityClient,
		PersonSearchService:       authorityClient,
		ProjectSearchService:      authorityClient,
		UserSearchService:         authorityClient,
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

func newRepository(logger *zap.SugaredLogger) backends.Repository {
	ctx := context.Background()

	bp := newPublicationBulkIndexerService(logger)
	bd := newDatasetBulkIndexerService(logger)

	repo, err := repository.New(viper.GetString("pg-conn"))
	if err != nil {
		log.Fatalln("unable to create store", err)
	}

	repo.AddPublicationListener(func(p *models.Publication) {
		if p.DateUntil == nil {
			if err := bp.Index(ctx, p); err != nil {
				logger.Errorf("error indexing publication %s: %w", p.ID, err)
			}
		}
	})

	repo.AddDatasetListener(func(d *models.Dataset) {
		if d.DateUntil == nil {
			if err := bd.Index(ctx, d); err != nil {
				logger.Errorf("error indexing dataset %s: %w", d.ID, err)
			}
		}
	})

	return repo
}

func newFileStore() backends.FileStore {
	if baseDir := viper.GetString("file-dir"); baseDir != "" {
		store, err := fsstore.New(fsstore.Config{
			Dir:     path.Join(baseDir, "root"),
			TempDir: path.Join(baseDir, "tmp"),
		})
		if err != nil {
			log.Fatalln("Unable to initialize filestore", err)
		}
		return store
	}

	store, err := s3store.New(s3store.Config{
		Endpoint:   viper.GetString("s3-endpoint"),
		Region:     viper.GetString("s3-region"),
		ID:         viper.GetString("s3-id"),
		Secret:     viper.GetString("s3-secret"),
		Bucket:     viper.GetString("s3-bucket"),
		TempBucket: viper.GetString("s3-temp-bucket"),
	})

	if err != nil {
		log.Fatalln("Unable to initialize filestore", err)
	}
	return store
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
		Index:          viper.GetString(t + "-index"),
		Settings:       string(settings),
		IndexRetention: viper.GetInt("index-retention"),
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

func newPublicationBulkIndexerService(logger *zap.SugaredLogger) backends.BulkIndexer[*models.Publication] {
	bp, err := newPublicationSearchService().NewBulkIndexer(backends.BulkIndexerConfig{
		OnError: func(err error) {
			logger.Errorf("Indexing failed : %s", err)
		},
		OnIndexError: func(id string, err error) {
			logger.Errorf("Indexing failed for dataset [id: %s] : %s", id, err)
		},
	})

	if err != nil {
		logger.Fatalln("unable to create publication bulk indexer", err)
	}

	return bp
}

func newDatasetBulkIndexerService(logger *zap.SugaredLogger) backends.BulkIndexer[*models.Dataset] {
	bd, err := newDatasetSearchService().NewBulkIndexer(backends.BulkIndexerConfig{
		OnError: func(err error) {
			logger.Errorf("Indexing failed : %s", err)
		},
		OnIndexError: func(id string, err error) {
			logger.Errorf("Indexing failed for dataset [id: %s] : %s", id, err)
		},
	})

	if err != nil {
		logger.Fatalln("unable to create publication bulk indexer", err)
	}

	return bd
}
