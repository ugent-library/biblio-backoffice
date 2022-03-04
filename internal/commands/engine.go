package commands

import (
	"log"

	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/backends/ianamedia"
	"github.com/ugent-library/biblio-backend/internal/backends/librecat"
	"github.com/ugent-library/biblio-backend/internal/backends/spdxlicenses"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/go-orcid/orcid"
	"go.temporal.io/sdk/client"
)

func newEngine() *engine.Engine {
	temporal, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}

	librecatClient := librecat.New(librecat.Config{
		URL:      viper.GetString("librecat-url"),
		Username: viper.GetString("librecat-username"),
		Password: viper.GetString("librecat-password"),
	})

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
		DatasetService:            librecatClient,
		DatasetSearchService:      librecatClient,
		PublicationService:        librecatClient,
		PublicationSearchService:  librecatClient,
		PersonService:             librecatClient,
		ProjectService:            librecatClient,
		UserService:               librecatClient,
		OrganizationSearchService: librecatClient,
		PersonSearchService:       librecatClient,
		ProjectSearchService:      librecatClient,
		LicenseSearchService:      spdxlicenses.New(),
		MediaTypeSearchService:    ianamedia.New(),
	})

	if err != nil {
		log.Fatal(err)
	}

	return e
}
