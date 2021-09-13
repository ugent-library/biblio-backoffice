package engine

import (
	"fmt"
	"net/url"

	"github.com/ugent-library/biblio-backend/internal/models"
)

func (e *Engine) GetUser(id string) (*models.User, error) {
	user := &models.User{}
	if _, err := e.get(fmt.Sprintf("/user/%s", url.PathEscape(id)), nil, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (e *Engine) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	if _, err := e.get(fmt.Sprintf("/user/username/%s", url.PathEscape(username)), nil, user); err != nil {
		return nil, err
	}
	return user, nil
}
