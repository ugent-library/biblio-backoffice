package commands

import (
	"log"
	"sync"

	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/backends/datacite"
	"github.com/ugent-library/biblio-backend/internal/backends/fc6"
	"github.com/ugent-library/biblio-backend/internal/backends/ianamedia"
	"github.com/ugent-library/biblio-backend/internal/backends/librecat"
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
	temporal, err := client.NewClient(client.Options{
		HostPort: viper.GetString("temporal-host-port"),
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}

	librecatClient := librecat.New(librecat.Config{
		URL:      viper.GetString("librecat-url"),
		Username: viper.GetString("librecat-username"),
		Password: viper.GetString("librecat-password"),
	})

	fc6Client := fc6.New(fc6.Config{})

	orcidConfig := orcid.Config{
		ClientID:     viper.GetString("orcid-client-id"),
		ClientSecret: viper.GetString("orcid-client-secret"),
		Sandbox:      viper.GetBool("orcid-sandbox"),
	}
	orcidClient := orcid.NewMemberClient(orcidConfig)

	e, err := engine.New(engine.Config{
		Temporal:                  temporal,
		ORCIDSandbox:              orcidConfig.Sandbox,
		ORCIDClient:               orcidClient,
		DatasetService:            fc6Client,
		DatasetSearchService:      librecatClient,
		PublicationService:        librecatClient,
		PublicationSearchService:  librecatClient,
		PersonService:             librecatClient,
		ProjectService:            librecatClient,
		UserService:               librecatClient,
		OrganizationSearchService: librecatClient,
		PersonSearchService:       librecatClient,
		ProjectSearchService:      librecatClient,
		DatasetSources: map[string]backends.DatasetSource{
			"datacite": datacite.New(),
		},
		LicenseSearchService:   spdxlicenses.New(),
		MediaTypeSearchService: ianamedia.New(),
	})

	if err != nil {
		log.Fatal(err)
	}

	return e
}
