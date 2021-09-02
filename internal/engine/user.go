package engine

import (
	"fmt"
	"net/url"
)

type User map[string]interface{}

func (e *Engine) GetUserByUsername(username string) (User, error) {
	var user User
	if _, err := e.get(fmt.Sprintf("/user/username/%s", url.PathEscape(username)), nil, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (e *Engine) GetUserByEmail(email string) (User, error) {
	user := make(User)
	if _, err := e.get(fmt.Sprintf("/user/email/%s", url.PathEscape(email)), nil, &user); err != nil {
		return nil, err
	}
	return user, nil
}
