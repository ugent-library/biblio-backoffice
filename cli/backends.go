package cli

import (
	"context"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/spf13/cobra"
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
	"github.com/ugent-library/biblio-backoffice/internal/mutate"
	"github.com/ugent-library/biblio-backoffice/models"

	"github.com/ugent-library/biblio-backoffice/internal/backends/ianamedia"
	"github.com/ugent-library/biblio-backoffice/internal/backends/jsonl"
	"github.com/ugent-library/biblio-backoffice/internal/backends/pubmed"
	"github.com/ugent-library/biblio-backoffice/internal/backends/ris"
	"github.com/ugent-library/biblio-backoffice/internal/backends/spdxlicenses"
	"github.com/ugent-library/biblio-backoffice/repositories"
	"github.com/ugent-library/orcid"
)

func newServices() *backends.Services {
	authorityClient, authorityClientErr := authority.New(authority.Config{
		MongoDBURI: config.MongoDBURL,
		ESURI:      strings.Split(config.Frontend.Es6URL, ","),
	})
	if authorityClientErr != nil {
		panic(authorityClientErr)
	}

	orcidConfig := orcid.Config{
		ClientID:     config.ORCID.ClientID,
		ClientSecret: config.ORCID.ClientSecret,
		Sandbox:      config.ORCID.Sandbox,
	}
	orcidClient := orcid.NewMemberClient(orcidConfig)

	citeprocURL := config.CiteprocURL

	var handleService backends.HandleService = nil

	if config.Handle.Enabled {
		handleService = handle.NewClient(
			handle.Config{
				BaseURL:         config.Handle.URL,
				FrontEndBaseURL: fmt.Sprintf("%s/publication", config.Frontend.URL),
				Prefix:          config.Handle.Prefix,
				Username:        config.Handle.Username,
				Password:        config.Handle.Password,
			},
		)
	}

	organizationService := caching.NewOrganizationService(authorityClient)

	// always add organization info to user affiliations
	userService := &backends.UserWithOrganizationsService{
		UserService:         caching.NewUserService(authorityClient),
		OrganizationService: organizationService,
	}

	// always add organization info to person affiliations
	personService := &backends.PersonWithOrganizationsService{
		PersonService:       caching.NewPersonService(authorityClient),
		OrganizationService: organizationService,
	}

	projectService := caching.NewProjectService(authorityClient)

	repo := newRepo(personService, organizationService, projectService)

	searchService := newSearchService()

	return &backends.Services{
		FileStore:                 newFileStore(),
		ORCIDSandbox:              orcidConfig.Sandbox,
		ORCIDClient:               orcidClient,
		Repo:                      repo,
		SearchService:             searchService,
		DatasetSearchIndex:        searchService.NewDatasetIndex(repo),
		PublicationSearchIndex:    searchService.NewPublicationIndex(repo),
		OrganizationService:       organizationService,
		PersonService:             personService,
		ProjectService:            projectService,
		UserService:               userService,
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
		DatasetListExporters: map[string]backends.DatasetListExporterFactory{
			"xlsx": excel_dataset.NewExporter,
		},
		HandleService: handleService,
	}
}

func newRepo(personService backends.PersonService, organizationService backends.OrganizationService, projectService backends.ProjectService) *repositories.Repo {
	ctx := context.Background()

	bp := newPublicationBulkIndexerService()
	bd := newDatasetBulkIndexerService()

	repo, err := repositories.New(repositories.Config{
		DSN: config.PgConn,

		PublicationListeners: []repositories.PublicationListener{
			func(p *models.Publication) {
				if p.DateUntil == nil {
					if err := bp.Index(ctx, p); err != nil {
						logger.Errorf("error indexing publication %s: %w", p.ID, err)
					}
				}
			},
		},

		DatasetListeners: []repositories.DatasetListener{
			func(d *models.Dataset) {
				if d.DateUntil == nil {
					if err := bd.Index(ctx, d); err != nil {
						logger.Errorf("error indexing dataset %s: %w", d.ID, err)
					}
				}
			},
		},

		PublicationLoaders: []repositories.PublicationVisitor{
			func(p *models.Publication) error {
				if p.CreatorID != "" {
					person, err := personService.GetPerson(p.CreatorID)
					if err != nil {
						logger.Warnf("error loading person %s in publication %s:, %w", p.CreatorID, p.ID, err)
						p.Creator = backends.NewDummyPerson(p.CreatorID)
					} else {
						p.Creator = person
					}
				}
				if p.UserID != "" {
					person, err := personService.GetPerson(p.UserID)
					if err != nil {
						logger.Warnf("error loading person %s in publication %s:, %w", p.UserID, p.ID, err)
						p.User = backends.NewDummyPerson(p.UserID)
					} else {
						p.User = person
					}
				}
				if p.LastUserID != "" {
					person, err := personService.GetPerson(p.LastUserID)
					if err != nil {
						logger.Warnf("error loading person %s in publication %s:, %w", p.LastUserID, p.ID, err)
						p.LastUser = backends.NewDummyPerson(p.LastUserID)
					} else {
						p.LastUser = person
					}
				}

				for _, role := range []string{"author", "editor", "supervisor"} {
					for _, c := range p.Contributors(role) {
						if c.PersonID == "" {
							continue
						}
						person, err := personService.GetPerson(c.PersonID)
						if err != nil {
							logger.Warnf("error loading person %s in publication %s:, %w", c.PersonID, p.ID, err)
							c.Person = backends.NewDummyPerson(c.PersonID)
						} else {
							c.Person = person
						}
					}
				}
				return nil
			},
			func(p *models.Publication) error {
				for _, rel := range p.RelatedOrganizations {
					org, err := organizationService.GetOrganization(rel.OrganizationID)
					if err != nil {
						rel.Organization = backends.NewDummyOrganization(rel.OrganizationID)
					} else {
						rel.Organization = org
					}
				}
				return nil
			},
			func(p *models.Publication) error {
				for _, rel := range p.RelatedProjects {
					project, err := projectService.GetProject(rel.ProjectID)
					if err != nil {
						logger.Warnf("error loading project %s in publication %s:, %w", rel.ProjectID, p.ID, err)
						rel.Project = backends.NewDummyProject(rel.ProjectID)
					} else {
						rel.Project = project
					}
				}
				return nil
			},
		},

		DatasetLoaders: []repositories.DatasetVisitor{
			func(d *models.Dataset) error {
				if d.CreatorID != "" {
					person, err := personService.GetPerson(d.CreatorID)
					if err != nil {
						logger.Warnf("error loading person %s in dataset %s:, %w", d.CreatorID, d.ID, err)
						d.Creator = backends.NewDummyPerson(d.CreatorID)
					} else {
						d.Creator = person
					}
				}
				if d.UserID != "" {
					person, err := personService.GetPerson(d.UserID)
					if err != nil {
						logger.Warnf("error loading person %s in dataset %s:, %w", d.UserID, d.ID, err)
						d.User = backends.NewDummyPerson(d.UserID)
					} else {
						d.User = person
					}
				}
				if d.LastUserID != "" {
					person, err := personService.GetPerson(d.LastUserID)
					if err != nil {
						logger.Warnf("error loading person %s in dataset %s:, %w", d.LastUserID, d.ID, err)
						d.LastUser = backends.NewDummyPerson(d.LastUserID)
					} else {
						d.LastUser = person
					}
				}

				for _, role := range []string{"author", "contributor"} {
					for _, c := range d.Contributors(role) {
						if c.PersonID == "" {
							continue
						}
						person, err := personService.GetPerson(c.PersonID)
						if err != nil {
							logger.Warnf("error loading person %s in dataset %s:, %w", c.PersonID, d.ID, err)
							c.Person = backends.NewDummyPerson(c.PersonID)
						} else {
							c.Person = person
						}
					}
				}
				return nil
			},
			func(d *models.Dataset) error {
				for _, rel := range d.RelatedOrganizations {
					org, err := organizationService.GetOrganization(rel.OrganizationID)
					if err != nil {
						rel.Organization = backends.NewDummyOrganization(rel.OrganizationID)
					} else {
						rel.Organization = org
					}
				}
				return nil
			},
			func(d *models.Dataset) error {
				for _, rel := range d.RelatedProjects {
					project, err := projectService.GetProject(rel.ProjectID)
					if err != nil {
						logger.Warnf("error loading project %s in dataset %s:, %w", rel.ProjectID, d.ID, err)
						rel.Project = backends.NewDummyProject(rel.ProjectID)
					} else {
						rel.Project = project
					}
				}
				return nil
			},
		},

		PublicationMutators: map[string]repositories.PublicationMutator{
			"project.add":              mutate.ProjectAdd(projectService),
			"classification.set":       mutate.ClassificationSet,
			"keyword.add":              mutate.KeywordAdd,
			"keyword.remove":           mutate.KeywordRemove,
			"vabb_id.set":              mutate.VABBIDSet,
			"vabb_type.set":            mutate.VABBTypeSet,
			"vabb_approved.set":        mutate.VABBApprovedSet,
			"vabb_year.add":            mutate.VABBYearAdd,
			"reviewer_tag.add":         mutate.ReviewerTagAdd,
			"reviewer_tag.remove":      mutate.ReviewerTagRemove,
			"journal_title.set":        mutate.JournalTitleSet,
			"journal_abbreviation.set": mutate.JournalAbbreviationSet,
			"isbn.add":                 mutate.ISBNAdd,
			"isbn.remove":              mutate.ISBNRemove,
			"eisbn.add":                mutate.EISBNAdd,
			"eisbn.remove":             mutate.EISBNRemove,
			"issn.add":                 mutate.ISSNAdd,
			"issn.remove":              mutate.ISSNRemove,
			"eissn.add":                mutate.EISSNAdd,
			"eissn.remove":             mutate.EISSNRemove,
		},
	})

	if err != nil {
		log.Fatalln("unable to create store", err)
	}

	return repo
}

func newFileStore() backends.FileStore {
	if baseDir := config.FileDir; baseDir != "" {
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
		Endpoint:   config.S3.Endpoint,
		Region:     config.S3.Region,
		ID:         config.S3.ID,
		Secret:     config.S3.Secret,
		Bucket:     config.S3.Bucket,
		TempBucket: config.S3.TempBucket,
	})

	if err != nil {
		log.Fatalln("Unable to initialize filestore", err)
	}
	return store
}

func newSearchService() backends.SearchService {
	config := es6.SearchServiceConfig{
		Addresses:        config.Es6URL,
		PublicationIndex: config.PublicationIndex,
		DatasetIndex:     config.DatasetIndex,
		IndexRetention:   config.IndexRetention,
	}

	s, err := es6.NewSearchService(config)

	if err != nil {
		logger.Fatalln("unable to create search service", err)
	}

	return s
}

func newPublicationBulkIndexerService() backends.BulkIndexer[*models.Publication] {
	ctx := context.Background()

	bp, err := newSearchService().NewPublicationBulkIndexer(backends.BulkIndexerConfig{
		OnError: func(err error) {
			logger.Errorf("Indexing failed : %s", err)
		},
		OnIndexError: func(id string, err error) {
			logger.Errorf("Indexing failed for publication [id: %s] : %s", id, err)
		},
	})

	if err != nil {
		logger.Fatalln("unable to create publication bulk indexer", err)
	}

	cobra.OnFinalize(func() {
		err := bp.Close(ctx)
		if err != nil {
			panic(err)
		}
	})

	return bp
}

func newDatasetBulkIndexerService() backends.BulkIndexer[*models.Dataset] {
	ctx := context.Background()

	bd, err := newSearchService().NewDatasetBulkIndexer(backends.BulkIndexerConfig{
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

	cobra.OnFinalize(func() {
		err := bd.Close(ctx)
		if err != nil {
			panic(err)
		}
	})

	return bd
}
