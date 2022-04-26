package biblio

import (
	"fmt"
	"net/url"

	"github.com/ugent-library/biblio-backend/internal/models"
)

func (c *Client) GetUser(id string) (*models.User, error) {
	user := &models.User{}
	if _, err := c.get(fmt.Sprintf("/restricted/user/%s", url.PathEscape(id)), nil, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c *Client) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	if _, err := c.get(fmt.Sprintf("/restricted/user/username/%s", url.PathEscape(username)), nil, user); err != nil {
		return nil, err
	}
	return user, nil
}
