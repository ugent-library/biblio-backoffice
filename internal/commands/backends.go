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
	"github.com/ugent-library/biblio-backend/internal/backends/ris"
	"github.com/ugent-library/biblio-backend/internal/backends/spdxlicenses"
	"github.com/ugent-library/biblio-backend/internal/backends/store"
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
	fs, err := filestore.New(viper.GetString("file-dir"))
	if err != nil {
		log.Fatalln("Unable to initialize filestore", err)
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

	return &backends.Services{
		FileStore:                 fs,
		ORCIDSandbox:              orcidConfig.Sandbox,
		ORCIDClient:               orcidClient,
		Store:                     newStore(),
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
		Tasks: tasks.NewHub(),
	}
}

func newStore() backends.Store {
	s, err := store.New(viper.GetString("pg-conn"))
	if err != nil {
		log.Fatalln("unable to create store", err)
	}
	return s
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
