package engine

import (
	"net/http"

	"github.com/ugent-library/go-orcid/orcid"
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
	orcidClient    *orcid.MemberClient
}

func New(c Config) (*Engine, error) {
	orcidClient := orcid.NewMemberClient(orcid.Config{
		ClientID:     c.ORCIDClientID,
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
