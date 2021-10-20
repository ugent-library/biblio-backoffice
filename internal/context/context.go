package context

import (
	"context"

	"github.com/ugent-library/biblio-backend/internal/models"
)

type key int

const (
	userKey key = iota
	activeMenuKey
)

func GetUser(c context.Context) *models.User {
	if v := c.Value(userKey); v != nil {
		return v.(*models.User)
	}
	return nil
}

func WithUser(c context.Context, user *models.User) context.Context {
	return context.WithValue(c, userKey, user)
}

func GetActiveMenu(c context.Context) string {
	if v := c.Value(activeMenuKey); v != nil {
		return v.(string)
	}
	return ""
}

func WithActiveMenu(c context.Context, menu string) context.Context {
	return context.WithValue(c, activeMenuKey, menu)
}
