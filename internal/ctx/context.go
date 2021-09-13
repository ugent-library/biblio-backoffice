package ctx

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
)

var UserKey = &key{"User"}

type key struct {
	name string
}

func (c *key) String() string {
	return c.name
}

func HasUser(r *http.Request) bool {
	_, ok := r.Context().Value(UserKey).(*models.User)
	return ok
}

func GetUser(r *http.Request) *models.User {
	return r.Context().Value(UserKey).(*models.User)
}
