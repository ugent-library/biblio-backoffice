package engine

import (
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/go-orcid/orcid"
	"go.temporal.io/sdk/client"
)

type Config struct {
	Temporal     client.Client
	ORCIDSandbox bool
	ORCIDClient  *orcid.MemberClient
	backends.DatasetService
	backends.DatasetSearchService
	backends.PersonService
	backends.ProjectService
	backends.PublicationService
	backends.PublicationSearchService
	backends.UserService
	backends.OrganizationSearchService
	backends.PersonSearchService
	backends.ProjectSearchService
	backends.LicenseSearchService
	backends.MediaTypeSearchService
	DatasetSources map[string]backends.DatasetSource
}

type Engine struct {
	Config
}

func New(c Config) (*Engine, error) {
	e := &Engine{Config: c}
	return e, nil
}
