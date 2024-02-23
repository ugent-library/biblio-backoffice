// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// AddPerson implements addPerson operation.
//
// Upsert a person.
//
// POST /add-person
func (UnimplementedHandler) AddPerson(ctx context.Context, req *AddPersonRequest) error {
	return ht.ErrNotImplemented
}

// AddProject implements addProject operation.
//
// Upsert a project.
//
// POST /add-project
func (UnimplementedHandler) AddProject(ctx context.Context, req *Project) error {
	return ht.ErrNotImplemented
}

// NewError creates *ErrorStatusCode from error returned by handler.
//
// Used for common default response.
func (UnimplementedHandler) NewError(ctx context.Context, err error) (r *ErrorStatusCode) {
	r = new(ErrorStatusCode)
	return r
}
