package server

import (
	"context"
	"encoding/base64"
	"log"
	"strings"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func auth(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "basic")

	if err != nil {
		return nil, status.Errorf(codes.Internal, "authentication failed: %s", err)
	}

	c, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "authentication failed: invalid base 64 in header: %s", err)
	}

	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return ctx, status.Error(codes.Unauthenticated, `invalid basic auth format`)
	}

	user, password := cs[:s], cs[s+1:]

	log.Println(user)

	if user != "foo" || password != "bar" {
		return ctx, status.Error(codes.Unauthenticated, "invalid user or password")
	}

	return ctx, nil
}
