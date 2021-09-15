package views

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/models"
)

type Context struct {
	Data interface{}
	User *models.User
}

func newContext(r *http.Request, data interface{}) Context {
	var c Context
	if d, ok := data.(Context); ok {
		c = d
	} else {
		c = Context{Data: data}
	}

	c.User = context.User(r.Context())

	return c
}
