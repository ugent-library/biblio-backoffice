package api

import (
	"context"

	"github.com/go-faster/errors"
)

type ApiSecurityHandler struct {
	APIKey string
}

func (s *ApiSecurityHandler) HandleApiKey(ctx context.Context, operationName string, t ApiKey) (context.Context, error) {
	if t.APIKey == s.APIKey {
		return ctx, nil
	}
	return ctx, errors.New("unauthorized")
}
