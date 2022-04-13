package engine

import (
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/go-orcid/orcid"
	"go.temporal.io/sdk/client"
)

type Config struct {
	Temporal                 client.Client
	ORCIDSandbox             bool
	ORCIDClient              *orcid.MemberClient
	StorageService           backends.StorageService
	DatasetSearchService     backends.DatasetSearchService
	PublicationSearchService backends.PublicationSearchService
	backends.PersonService
	backends.ProjectService
	backends.UserService
	backends.OrganizationSearchService
	backends.PersonSearchService
	backends.ProjectSearchService
	backends.LicenseSearchService
	backends.MediaTypeSearchService
	DatasetSources map[string]backends.DatasetGetter
}

type Engine struct {
	Config
}

func New(c Config) (*Engine, error) {
	e := &Engine{
		Config: c,
	}
	return e, nil
}
