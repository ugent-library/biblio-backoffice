package views

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/models"
)

type Data struct {
	User  *models.User
}

func NewData(r *http.Request) Data {
	return Data{User: context.User(r.Context())}
}
