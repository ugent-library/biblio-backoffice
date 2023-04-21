// TODO use something with lower gc load and more predictable size like
// https://github.com/dgraph-io/ristretto
// https://github.com/allegro/bigcache
// or memcached, redis

package caching

import (
	"time"

	"github.com/bluele/gcache"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

type userService struct {
	service backends.UserService
	cache   gcache.Cache
}

func NewUserService(service backends.UserService) backends.UserService {
	cache := gcache.New(2500).
		Expiration(5 * time.Minute).
		LRU().
		LoaderFunc(func(key any) (any, error) {
			return service.GetUser(key.(string))
		}).
		Build()
	return &userService{
		cache:   cache,
		service: service,
	}
}

func (s *userService) GetUser(id string) (*models.User, error) {
	v, err := s.cache.Get(id)
	if err != nil {
		return nil, err
	}
	return v.(*models.User), nil
}

func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	return s.service.GetUserByUsername(username)
}

type personService struct {
	service backends.PersonService
	cache   gcache.Cache
}

func NewPersonService(service backends.PersonService) backends.PersonService {
	cache := gcache.New(5000).
		Expiration(30 * time.Minute).
		LoaderFunc(func(key any) (any, error) {
			return service.GetPerson(key.(string))
		}).
		LRU().
		Build()
	return &personService{
		cache:   cache,
		service: service,
	}
}

func (s *personService) GetPersons(ids []string) ([]*models.Person, error) {
	persons := make([]*models.Person, 0)
	for _, id := range ids {
		person, err := s.GetPerson(id)
		if err != nil && err == backends.ErrNotFound {
			continue
		} else if err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}
	return persons, nil
}

func (s *personService) GetPerson(id string) (*models.Person, error) {
	v, err := s.cache.Get(id)
	if err != nil {
		return nil, err
	}
	return v.(*models.Person), nil
}

type organizationService struct {
	service backends.OrganizationService
	cache   gcache.Cache
}

func NewOrganizationService(service backends.OrganizationService) backends.OrganizationService {
	cache := gcache.New(1000).
		Expiration(30 * time.Minute).
		LoaderFunc(func(key any) (any, error) {
			return service.GetOrganization(key.(string))
		}).
		LRU().
		Build()
	return &organizationService{
		cache:   cache,
		service: service,
	}
}

func (s *organizationService) GetOrganization(id string) (*models.Organization, error) {
	v, err := s.cache.Get(id)
	if err != nil {
		return nil, err
	}
	return v.(*models.Organization), nil
}

type projectService struct {
	service backends.ProjectService
	cache   gcache.Cache
}

func NewProjectService(service backends.ProjectService) backends.ProjectService {
	cache := gcache.New(2500).
		Expiration(30 * time.Minute).
		LoaderFunc(func(key any) (any, error) {
			return service.GetProject(key.(string))
		}).
		LRU().
		Build()
	return &projectService{
		cache:   cache,
		service: service,
	}
}

func (s *projectService) GetProject(id string) (*models.Project, error) {
	v, err := s.cache.Get(id)
	if err != nil {
		return nil, err
	}
	return v.(*models.Project), nil
}
