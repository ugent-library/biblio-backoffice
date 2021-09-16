package views

import (
	"net/http"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/models"
)

type Data struct {
	User    *models.User
	request *http.Request
}

func NewData(r *http.Request) Data {
	return Data{
		User:    context.User(r.Context()),
		request: r,
	}
}

func (d Data) OnHTMXFragment() bool {
	return strings.Contains(d.request.RequestURI, "htmx")
}
