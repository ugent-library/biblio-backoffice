package context

import (
	"context"

	"github.com/ugent-library/biblio-backend/internal/models"
)

var userKey = &key{"User"}

type key struct {
	name string
}

func (c *key) String() string {
	return c.name
}

func HasUser(c context.Context) bool {
	_, ok := c.Value(userKey).(*models.User)
	return ok
}

func User(c context.Context) *models.User {
	return c.Value(userKey).(*models.User)
}

func WithUser(c context.Context, user *models.User) context.Context {
	return context.WithValue(c, userKey, user)
}
