package datasetviewing

import (
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/models"
)

type Context struct {
	handlers.BaseContext
	Dataset     *models.Dataset
	RedirectURL string
}
