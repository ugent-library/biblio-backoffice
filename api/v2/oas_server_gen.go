// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// AddPerson implements addPerson operation.
	//
	// Upsert a person.
	//
	// POST /add-person
	AddPerson(ctx context.Context, req *AddPersonRequest) error
	// AddProject implements addProject operation.
	//
	// Upsert a project.
	//
	// POST /add-project
	AddProject(ctx context.Context, req *AddProjectRequest) error
	// GetOrganization implements getOrganization operation.
	//
	// Get organization by identifier.
	//
	// POST /get-organization
	GetOrganization(ctx context.Context, req *GetOrganizationRequest) (GetOrganizationRes, error)
	// ImportOrganizations implements importOrganizations operation.
	//
	// Import organization hierarchy.
	//
	// POST /import-organizations
	ImportOrganizations(ctx context.Context, req *ImportOrganizationsRequest) error
	// ImportPerson implements importPerson operation.
	//
	// Import a person.
	//
	// POST /import-person
	ImportPerson(ctx context.Context, req *ImportPersonRequest) error
	// ImportProject implements importProject operation.
	//
	// Import a project.
	//
	// POST /import-project
	ImportProject(ctx context.Context, req *ImportProjectRequest) (ImportProjectRes, error)
	// NewError creates *ErrorStatusCode from error returned by handler.
	//
	// Used for common default response.
	NewError(ctx context.Context, err error) *ErrorStatusCode
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h   Handler
	sec SecurityHandler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, sec SecurityHandler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		sec:        sec,
		baseServer: s,
	}, nil
}
