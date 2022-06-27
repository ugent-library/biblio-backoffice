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
	"github.com/ugent-library/biblio-backend/internal/backends/bibtex"
	"github.com/ugent-library/biblio-backend/internal/backends/citeproc"
	"github.com/ugent-library/biblio-backend/internal/backends/crossref"
	"github.com/ugent-library/biblio-backend/internal/backends/datacite"
	"github.com/ugent-library/biblio-backend/internal/backends/es6"
	"github.com/ugent-library/biblio-backend/internal/backends/filestore"
	"github.com/ugent-library/biblio-backend/internal/backends/ianamedia"
	"github.com/ugent-library/biblio-backend/internal/backends/jsonl"
	"github.com/ugent-library/biblio-backend/internal/backends/pubmed"
	"github.com/ugent-library/biblio-backend/internal/backends/repository"
	"github.com/ugent-library/biblio-backend/internal/backends/ris"

	// "github.com/ugent-library/biblio-backend/internal/backends/spdxlicenses"
	"github.com/ugent-library/biblio-backend/internal/tasks"
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
		// LicenseSearchService:      spdxlicenses.New(),
		MediaTypeSearchService: ianamedia.New(),
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
		Tasks: tasks.NewHub(),
	}
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
	settings, err := ioutil.ReadFile("etc/es6/" + t + ".json")
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
