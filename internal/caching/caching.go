// TODO use something with lower gc load and more predictable size like
// https://github.com/dgraph-io/ristretto
// https://github.com/allegro/bigcache

package caching

import (
	"errors"
	"time"

	"github.com/bluele/gcache"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
)

type userService struct {
	service backends.UserService
	cache   gcache.Cache
}

func NewUserService(service backends.UserService) backends.UserService {
	cache := gcache.New(2500).Expiration(5 * time.Minute).LRU().Build()
	return &userService{
		cache:   cache,
		service: service,
	}
}

func (s *userService) GetUser(id string) (*models.User, error) {
	v, err := s.cache.Get(id)
	if errors.Is(err, gcache.KeyNotFoundError) {
		u, err := s.service.GetUser(id)
		if err != nil {
			return nil, err
		}
		s.cache.Set(id, u)
		return u, nil
	}
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
	cache := gcache.New(5000).Expiration(30 * time.Minute).LRU().Build()
	return &personService{
		cache:   cache,
		service: service,
	}
}

func (s *personService) GetPerson(id string) (*models.Person, error) {
	v, err := s.cache.Get(id)
	if errors.Is(err, gcache.KeyNotFoundError) {
		p, err := s.service.GetPerson(id)
		if err != nil {
			return nil, err
		}
		s.cache.Set(id, p)
		return p, nil
	}
	if err != nil {
		return nil, err
	}
	return v.(*models.Person), nil
}

type organizationService struct {
	service backends.OrganizationService
	cache   gcache.Cache
}

func NewOrganzationService(service backends.OrganizationService) backends.OrganizationService {
	cache := gcache.New(1000).Expiration(30 * time.Minute).LRU().Build()
	return &organizationService{
		cache:   cache,
		service: service,
	}
}

func (s *organizationService) GetOrganization(id string) (*models.Organization, error) {
	v, err := s.cache.Get(id)
	if errors.Is(err, gcache.KeyNotFoundError) {
		o, err := s.service.GetOrganization(id)
		if err != nil {
			return nil, err
		}
		s.cache.Set(id, o)
		return o, nil
	}
	if err != nil {
		return nil, err
	}
	return v.(*models.Organization), nil
}

type projectService struct {
	service backends.ProjectService
	cache   gcache.Cache
}

func NewprojectService(service backends.ProjectService) backends.ProjectService {
	cache := gcache.New(2500).Expiration(30 * time.Minute).LRU().Build()
	return &projectService{
		cache:   cache,
		service: service,
	}
}

func (s *projectService) GetProject(id string) (*models.Project, error) {
	v, err := s.cache.Get(id)
	if errors.Is(err, gcache.KeyNotFoundError) {
		p, err := s.service.GetProject(id)
		if err != nil {
			return nil, err
		}
		s.cache.Set(id, p)
		return p, nil
	}
	if err != nil {
		return nil, err
	}
	return v.(*models.Project), nil
}
