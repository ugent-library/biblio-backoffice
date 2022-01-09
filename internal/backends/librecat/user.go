package librecat

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/models"
)

func (c *Client) GetUser(id string) (*models.User, error) {
	user := &models.User{}
	if _, err := c.get(fmt.Sprintf("/user/%s", url.PathEscape(id)), nil, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c *Client) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	if _, err := c.get(fmt.Sprintf("/user/username/%s", url.PathEscape(username)), nil, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c *Client) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if _, err := c.get(fmt.Sprintf("/user/email/%s", url.PathEscape(strings.ToLower(email))), nil, user); err != nil {
		return nil, err
	}
	return user, nil
}
