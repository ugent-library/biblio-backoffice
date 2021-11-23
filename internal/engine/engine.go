package engine

import (
	"net/http"

	"github.com/nics/orcid-go/orcid"
)

type Config struct {
	LibreCatURL       string
	LibreCatUsername  string
	LibreCatPassword  string
	ORCIDClientID     string
	ORCIDClientSecret string
	ORCIDSandbox      bool
}

type Engine struct {
	Config         Config
	librecatClient *http.Client
	orcidClient    *orcid.Client
}

func New(c Config) (*Engine, error) {
	orcidClient := orcid.NewClient(orcid.Config{
		ClientId:     c.ORCIDClientID,
		ClientSecret: c.ORCIDClientSecret,
		// Scopes:       []string{"/read-public"},
		Sandbox: c.ORCIDSandbox,
	})

	e := &Engine{
		Config:         c,
		librecatClient: http.DefaultClient,
		orcidClient:    orcidClient,
	}

	return e, nil
}
