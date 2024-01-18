package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Users []*User

type Permissions map[string][]string

type User struct {
	Username string
	Password string
	Role     string
}

type BasicAuthInterceptor struct {
	users       Users
	permissions Permissions
}

func NewBasicAuthInterceptor(u Users, p Permissions) *BasicAuthInterceptor {
	return &BasicAuthInterceptor{
		users:       u,
		permissions: p,
	}
}

func (a *BasicAuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// Decode
		cs, err := a.decode(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "authentication failed: %s", err)
		}

		s := strings.IndexByte(cs, ':')
		user, password := cs[:s], cs[s+1:]

		// Authenticate
		u := a.authenticate(user, password)
		if u == nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: invalid username or password")
		}

		// Authorize
		err = a.authorize(u, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (a *BasicAuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Decode
		cs, err := a.decode(stream.Context())
		if err != nil {
			return status.Errorf(codes.Internal, "authentication failed: %s", err)
		}

		s := strings.IndexByte(cs, ':')
		user, password := cs[:s], cs[s+1:]

		// Authenticate
		u := a.authenticate(user, password)
		if u == nil {
			return status.Errorf(codes.Unauthenticated, "authentication failed: invalid username or password")
		}

		// Authorize
		err = a.authorize(u, info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func (a *BasicAuthInterceptor) authorize(u *User, method string) error {
	roles, ok := a.permissions[method]
	if !ok {
		// Everyone has access to this method, no roles defined
		return nil
	}

	for _, role := range roles {
		if role == u.Role {
			return nil
		}
	}

	return status.Error(codes.PermissionDenied, "authorization failed: no permission to access this RPC method")
}

func (a *BasicAuthInterceptor) authenticate(username string, password string) *User {
	if username == "" || password == "" {
		return nil
	}

	for _, u := range a.users {
		if u.Username == username && u.Password == password {
			return u
		}
	}

	return nil
}

func (a *BasicAuthInterceptor) decode(ctx context.Context) (string, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "basic")

	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	if token == "" {
		return "", fmt.Errorf("empty basic authentication header")
	}

	c, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", fmt.Errorf("invalid base 64 formatting in basic authentication header: %w", err)
	}

	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return "", fmt.Errorf("invalid basic authentication format (not user:password)")
	}

	return cs, nil
}
