package engine

import (
	"fmt"
	"net/url"

	"github.com/ugent-library/biblio-backend/internal/models"
)

func (e *Engine) GetPerson(id string) (*models.Person, error) {
	p := &models.Person{}
	if _, err := e.get(fmt.Sprintf("/person/%s", url.PathEscape(id)), nil, p); err != nil {
		return nil, err
	}
	return p, nil
}
