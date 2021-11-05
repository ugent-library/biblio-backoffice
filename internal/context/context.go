package context

import (
	"context"

	"github.com/ugent-library/biblio-backend/internal/models"
)

type key int

const (
	userKey key = iota
	originalUserKey
	activeMenuKey
	publicationKey
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

func GetOriginalUser(c context.Context) *models.User {
	if v := c.Value(originalUserKey); v != nil {
		return v.(*models.User)
	}
	return nil
}

func WithOriginalUser(c context.Context, user *models.User) context.Context {
	return context.WithValue(c, originalUserKey, user)
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

func GetPublication(c context.Context) *models.Publication {
	if v := c.Value(publicationKey); v != nil {
		return v.(*models.Publication)
	}
	return nil
}

func WithPublication(c context.Context, pub *models.Publication) context.Context {
	return context.WithValue(c, publicationKey, pub)
}

func GetDataset(c context.Context) *models.Dataset {
	if v := c.Value(publicationKey); v != nil {
		return v.(*models.Dataset)
	}
	return nil
}

func WithDataset(c context.Context, pub *models.Dataset) context.Context {
	return context.WithValue(c, publicationKey, pub)
}
