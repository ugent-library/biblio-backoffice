package ctx

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/engine"
)

var UserKey = &key{"User"}

type key struct {
	name string
}

func (c *key) String() string {
	return c.name
}

func HasUser(r *http.Request) bool {
	_, ok := r.Context().Value(UserKey).(*engine.User)
	return ok
}

func GetUser(r *http.Request) *engine.User {
	return r.Context().Value(UserKey).(*engine.User)
}
